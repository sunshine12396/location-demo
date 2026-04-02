package location

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/example/location-demo/internal/domain"
)

// PostgresRepository implements domain.LocationRepository using PostgreSQL.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new repository instance.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// ──────────────────────────────────────────────
// Location Queries
// ──────────────────────────────────────────────

// GetByID retrieves a location with its translated name, parent info, and stats.
func (r *PostgresRepository) GetByID(ctx context.Context, id int64, lang string) (*domain.LocationDetail, error) {
	query := `
		SELECT
			l.id, l.type, l.lat, l.lng,
			COALESCE(l.slug, '') AS slug,
			COALESCE(l.is_verified, FALSE) AS is_verified,
			COALESCE(lt.name, '') AS name,
			l.parent_id,
			COALESCE(pt.name, '') AS parent_name,
			COALESCE(p.type, '') AS parent_type
		FROM locations l
		LEFT JOIN location_translations lt ON lt.location_id = l.id AND lt.lang_code = $2
		LEFT JOIN locations p ON p.id = l.parent_id
		LEFT JOIN location_translations pt ON pt.location_id = p.id AND pt.lang_code = $2
		WHERE l.id = $1
	`

	detail := &domain.LocationDetail{}
	var parentID sql.NullInt64
	var parentName, parentType string

	err := r.db.QueryRowContext(ctx, query, id, lang).Scan(
		&detail.ID, &detail.Type, &detail.Lat, &detail.Lng,
		&detail.Slug, &detail.IsVerified,
		&detail.Name,
		&parentID, &parentName, &parentType,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("location %d not found", id)
		}
		return nil, fmt.Errorf("repository.GetByID: %w", err)
	}

	if parentID.Valid {
		detail.Parent = &domain.LocationSummary{
			ID:   parentID.Int64,
			Name: parentName,
			Type: domain.LocationType(parentType),
		}
	}

	// Attach stats (non-blocking — stats may not exist yet)
	stats, _ := r.GetStats(ctx, id)
	detail.Stats = stats

	return detail, nil
}

// SearchByAlias searches the alias table using case-insensitive LIKE matching.
func (r *PostgresRepository) SearchByAlias(ctx context.Context, query string) ([]domain.SearchResult, error) {
	sqlQuery := `
		SELECT l.id, COALESCE(lt.name, la.alias) AS name, l.type
		FROM location_alias la
		JOIN locations l ON l.id = la.location_id
		LEFT JOIN location_translations lt ON lt.location_id = l.id AND lt.lang_code = 'en'
		WHERE LOWER(la.alias) LIKE LOWER($1)
		LIMIT 10
	`

	return r.querySearchResults(ctx, sqlQuery, "%"+query+"%")
}

// SearchByTranslation searches translations using case-insensitive LIKE.
func (r *PostgresRepository) SearchByTranslation(ctx context.Context, query string, lang string) ([]domain.SearchResult, error) {
	sqlQuery := `
		SELECT l.id, lt.name, l.type
		FROM location_translations lt
		JOIN locations l ON l.id = lt.location_id
		WHERE lt.lang_code = $1 AND LOWER(lt.name) LIKE LOWER($2)
		LIMIT 10
	`

	return r.querySearchResults(ctx, sqlQuery, lang, "%"+query+"%")
}

// InsertLocation creates a new location along with its translations and aliases in a transaction.
func (r *PostgresRepository) InsertLocation(ctx context.Context, loc *domain.Location, translations []domain.LocationTranslation, aliases []domain.LocationAlias) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("repository.InsertLocation: begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	// Insert location
	var id int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO locations (external_id, type, lat, lng, parent_id, path)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (external_id) DO UPDATE SET external_id = EXCLUDED.external_id
		RETURNING id
	`, loc.ExternalID, loc.Type, loc.Lat, loc.Lng, loc.ParentID, loc.Path).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repository.InsertLocation: insert location: %w", err)
	}

	// Insert translations
	for _, t := range translations {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO location_translations (location_id, lang_code, name)
			VALUES ($1, $2, $3)
			ON CONFLICT (location_id, lang_code) DO UPDATE SET name = EXCLUDED.name
		`, id, t.LangCode, t.Name)
		if err != nil {
			return 0, fmt.Errorf("repository.InsertLocation: insert translation: %w", err)
		}
	}

	// Insert aliases
	for _, a := range aliases {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO location_alias (location_id, alias)
			VALUES ($1, $2)
		`, id, strings.ToLower(a.Alias))
		if err != nil {
			return 0, fmt.Errorf("repository.InsertLocation: insert alias: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("repository.InsertLocation: commit: %w", err)
	}

	return id, nil
}

// GetChildren returns direct children of a parent location.
func (r *PostgresRepository) GetChildren(ctx context.Context, parentID int64, lang string) ([]domain.SearchResult, error) {
	sqlQuery := `
		SELECT l.id, COALESCE(lt.name, '') AS name, l.type
		FROM locations l
		LEFT JOIN location_translations lt ON lt.location_id = l.id AND lt.lang_code = $2
		WHERE l.parent_id = $1
		ORDER BY lt.name
	`

	return r.querySearchResults(ctx, sqlQuery, parentID, lang)
}

// ──────────────────────────────────────────────
// Post Queries
// ──────────────────────────────────────────────

// CreatePost inserts a new post and updates location_stats atomically.
func (r *PostgresRepository) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("repository.CreatePost: begin tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	// 1. Insert the post
	err = tx.QueryRowContext(ctx, `
		INSERT INTO posts (user_id, content, media_type, location_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`, post.UserID, post.Content, post.MediaType, post.LocationID).Scan(&post.ID, &post.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.CreatePost: insert post: %w", err)
	}

	// 2. Upsert location_stats — increment the correct counter
	_, err = tx.ExecContext(ctx, `
		INSERT INTO location_stats (location_id, total_posts, total_photos, total_videos, last_post_at, trending_score)
		VALUES ($1, 1,
			CASE WHEN $2 = 'photo' THEN 1 ELSE 0 END,
			CASE WHEN $2 = 'video' THEN 1 ELSE 0 END,
			NOW(),
			1.0 + CASE WHEN $2 = 'photo' THEN 1.5 WHEN $2 = 'video' THEN 2.0 ELSE 0 END
		)
		ON CONFLICT (location_id) DO UPDATE SET
			total_posts    = location_stats.total_posts + 1,
			total_photos   = location_stats.total_photos + CASE WHEN $2 = 'photo' THEN 1 ELSE 0 END,
			total_videos   = location_stats.total_videos + CASE WHEN $2 = 'video' THEN 1 ELSE 0 END,
			last_post_at   = NOW(),
			trending_score = location_stats.trending_score + 1.0
				+ CASE WHEN $2 = 'photo' THEN 1.5 WHEN $2 = 'video' THEN 2.0 ELSE 0 END
	`, post.LocationID, post.MediaType)
	if err != nil {
		return nil, fmt.Errorf("repository.CreatePost: update stats: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("repository.CreatePost: commit: %w", err)
	}

	return post, nil
}

// GetPostsByLocation returns posts for a location and all its descendants, joined with location name.
func (r *PostgresRepository) GetPostsByLocation(ctx context.Context, locationID int64, lang string, limit, offset int) ([]domain.PostWithLocation, error) {
	// First, get the path of the target location
	var path string
	err := r.db.QueryRowContext(ctx, "SELECT path FROM locations WHERE id = $1", locationID).Scan(&path)
	if err != nil {
		if err == sql.ErrNoRows {
			// Location might exist but have no path, or it might just have no posts.
			// Revert to direct match if we can't find a path
			path = "" 
		} else {
			return nil, fmt.Errorf("repository.GetPostsByLocation: fetching path: %w", err)
		}
	}

	var query string
	var rows *sql.Rows

	if path != "" {
		query = `
			SELECT
				p.id, p.user_id, p.content, p.media_type,
				p.location_id, COALESCE(lt.name, '') AS location_name, l.type AS location_type,
				p.created_at
			FROM posts p
			JOIN locations l ON l.id = p.location_id
			LEFT JOIN location_translations lt ON lt.location_id = l.id AND lt.lang_code = $2
			WHERE l.path LIKE $1 || '%'
			ORDER BY p.created_at DESC
			LIMIT $3 OFFSET $4
		`
		rows, err = r.db.QueryContext(ctx, query, path, lang, limit, offset)
	} else {
		// Fallback to exact match
		query = `
			SELECT
				p.id, p.user_id, p.content, p.media_type,
				p.location_id, COALESCE(lt.name, '') AS location_name, l.type AS location_type,
				p.created_at
			FROM posts p
			JOIN locations l ON l.id = p.location_id
			LEFT JOIN location_translations lt ON lt.location_id = l.id AND lt.lang_code = $2
			WHERE p.location_id = $1
			ORDER BY p.created_at DESC
			LIMIT $3 OFFSET $4
		`
		rows, err = r.db.QueryContext(ctx, query, locationID, lang, limit, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("repository.GetPostsByLocation: %w", err)
	}
	defer rows.Close()

	var posts []domain.PostWithLocation
	for rows.Next() {
		var p domain.PostWithLocation
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.Content, &p.MediaType,
			&p.LocationID, &p.LocationName, &p.LocationType,
			&p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("repository.GetPostsByLocation: scan: %w", err)
		}
		posts = append(posts, p)
	}

	return posts, rows.Err()
}

// ──────────────────────────────────────────────
// Stats & Trending Queries
// ──────────────────────────────────────────────

// GetStats returns the pre-aggregated stats for a location.
func (r *PostgresRepository) GetStats(ctx context.Context, locationID int64) (*domain.LocationStats, error) {
	query := `
		SELECT location_id, total_posts, total_photos, total_videos, last_post_at, trending_score
		FROM location_stats
		WHERE location_id = $1
	`

	stats := &domain.LocationStats{}
	var lastPostAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, locationID).Scan(
		&stats.LocationID, &stats.TotalPosts, &stats.TotalPhotos,
		&stats.TotalVideos, &lastPostAt, &stats.TrendingScore,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No stats yet — not an error
		}
		return nil, fmt.Errorf("repository.GetStats: %w", err)
	}

	if lastPostAt.Valid {
		stats.LastPostAt = &lastPostAt.Time
	}

	return stats, nil
}

// GetTrending returns top trending locations ordered by score.
func (r *PostgresRepository) GetTrending(ctx context.Context, lang string, limit int) ([]domain.TrendingLocation, error) {
	query := `
		SELECT
			ls.location_id,
			COALESCE(lt.name, '') AS name,
			l.type,
			ls.trending_score AS score
		FROM location_stats ls
		JOIN locations l ON l.id = ls.location_id
		LEFT JOIN location_translations lt ON lt.location_id = l.id AND lt.lang_code = $1
		WHERE ls.trending_score > 0
		ORDER BY ls.trending_score DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, lang, limit)
	if err != nil {
		return nil, fmt.Errorf("repository.GetTrending: %w", err)
	}
	defer rows.Close()

	var results []domain.TrendingLocation
	for rows.Next() {
		var t domain.TrendingLocation
		if err := rows.Scan(&t.LocationID, &t.Name, &t.Type, &t.Score); err != nil {
			return nil, fmt.Errorf("repository.GetTrending: scan: %w", err)
		}
		results = append(results, t)
	}

	return results, rows.Err()
}

// ──────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────

// querySearchResults is a helper to avoid duplicating scan logic.
func (r *PostgresRepository) querySearchResults(ctx context.Context, query string, args ...interface{}) ([]domain.SearchResult, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("repository.querySearchResults: %w", err)
	}
	defer rows.Close()

	var results []domain.SearchResult
	for rows.Next() {
		var sr domain.SearchResult
		if err := rows.Scan(&sr.ID, &sr.Name, &sr.Type); err != nil {
			return nil, fmt.Errorf("repository.querySearchResults: scan: %w", err)
		}
		results = append(results, sr)
	}

	return results, rows.Err()
}
