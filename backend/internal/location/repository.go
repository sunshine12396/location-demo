package location

import (
	"context"
	"fmt"
	"strings"

	"github.com/example/location-demo/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PostgresRepository implements domain.LocationRepository using GORM.
type PostgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new repository instance.
func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// GetByID retrieves a location with its translated name, parent info, and stats.
func (r *PostgresRepository) GetByExternalID(ctx context.Context, externalID string) (*domain.Location, error) {
	var loc domain.Location
	err := r.db.WithContext(ctx).Table("locations").Where("external_id = ?", externalID).First(&loc).Error
	if err != nil {
		return nil, err
	}
	return &loc, nil
}

// GetByID retrieves a location with its translated name, parent info, and stats.
func (r *PostgresRepository) GetByID(ctx context.Context, id int64, lang string) (*domain.LocationDetail, error) {
	var loc domain.Location
	err := r.db.WithContext(ctx).
		Preload("Translations", "lang_code = ?", lang).
		Preload("Parent").
		Preload("Parent.Translations", "lang_code = ?", lang).
		Preload("Stats").
		First(&loc, id).Error

	if err != nil {
		return nil, fmt.Errorf("repository.GetByID: %w", err)
	}

	// Fallback Strategy: If requested translation missing, try 'en' then any available node.
	if len(loc.Translations) == 0 {
		r.db.WithContext(ctx).Table("location_translations").
			Where("location_id = ?", id).
			Order("CASE WHEN lang_code = 'en' THEN 0 ELSE 1 END").
			Limit(1).Find(&loc.Translations)
	}

	if loc.Parent != nil && len(loc.Parent.Translations) == 0 {
		r.db.WithContext(ctx).Table("location_translations").
			Where("location_id = ?", loc.Parent.ID).
			Order("CASE WHEN lang_code = 'en' THEN 0 ELSE 1 END").
			Limit(1).Find(&loc.Parent.Translations)
	}

	// Map domain.Location to domain.LocationDetail
	name := ""
	formattedAddress := ""
	shortFormattedAddress := ""
	if len(loc.Translations) > 0 {
		name = loc.Translations[0].Name
		formattedAddress = loc.Translations[0].FormattedAddress
		shortFormattedAddress = loc.Translations[0].ShortFormattedAddress
	}

	detail := &domain.LocationDetail{
		ID:                    loc.ID,
		ExternalID:            loc.ExternalID,
		Name:                  name,
		FormattedAddress:      formattedAddress,
		ShortFormattedAddress: shortFormattedAddress,
		Type:                  loc.Type,
		Lat:                   loc.Lat,
		Lng:                   loc.Lng,
		Provider:              loc.Provider,
		UpdatedAt:             loc.UpdatedAt,
		Stats:                 loc.Stats,
	}

	if loc.Parent != nil {
		parentName := ""
		if len(loc.Parent.Translations) > 0 {
			parentName = loc.Parent.Translations[0].Name
		}
		detail.Parent = &domain.LocationSummary{
			ID:   loc.Parent.ID,
			Name: parentName,
			Type: loc.Parent.Type,
		}
	}

	return detail, nil
}

// SearchByAlias searches the alias table using case-insensitive match and optional type filter.
func (r *PostgresRepository) SearchByAlias(ctx context.Context, query string, locType *domain.LocationType) ([]domain.SearchResult, error) {
	var results []domain.SearchResult
	q := r.db.WithContext(ctx).Table("location_aliases la").
		Select("l.id, COALESCE(lt.name, MAX(la.alias)) AS name, l.type").
		Joins("JOIN locations l ON l.id = la.location_id").
		Joins("LEFT JOIN location_translations lt ON lt.location_id = l.id AND lt.lang_code = ?", "en").
		Where("LOWER(la.alias) LIKE ?", "%"+strings.ToLower(query)+"%").
		Group("l.id, l.type, lt.name")

	if locType != nil {
		q = q.Where("l.type = ?", *locType)
	}

	err := q.Limit(10).Scan(&results).Error
	return results, err
}

// SearchByTranslation searches translations using case-insensitive match and optional type filter.
func (r *PostgresRepository) SearchByTranslation(ctx context.Context, query string, lang string, locType *domain.LocationType) ([]domain.SearchResult, error) {
	var results []domain.SearchResult
	q := r.db.WithContext(ctx).Table("location_translations lt").
		Select("l.id, lt.name, l.type").
		Joins("JOIN locations l ON l.id = lt.location_id").
		Where("lt.lang_code = ? AND LOWER(lt.name) LIKE ?", lang, "%"+strings.ToLower(query)+"%")

	if locType != nil {
		q = q.Where("l.type = ?", *locType)
	}

	err := q.Limit(10).Scan(&results).Error
	return results, err
}

// InsertLocation creates a new location with its associations automatically.
func (r *PostgresRepository) InsertLocation(ctx context.Context, loc *domain.Location, translations []domain.LocationTranslation, aliases []domain.LocationAlias) (int64, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Assign associations
		loc.Translations = translations
		loc.Aliases = aliases

		// 2. Insert location with Upsert logic for high-concurrency safety
		err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "external_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"lat":        loc.Lat,
				"lng":        loc.Lng,
				"updated_at": gorm.Expr("NOW()"),
			}),
		}).Create(loc).Error

		if err != nil {
			return err
		}

		// 3. Populate ID if not returned (depends on driver behavior during conflict updates)
		if loc.ID == 0 && loc.ExternalID != "" {
			if err := tx.Select("id").Table("locations").Where("external_id = ?", loc.ExternalID).First(&loc.ID).Error; err != nil {
				return err
			}
		}

		// Update path using hierarchy logic
		var path string
		if loc.ParentID != nil {
			var parent domain.Location
			tx.Select("path").First(&parent, *loc.ParentID)
			path = fmt.Sprintf("%s.%d", parent.Path, loc.ID)
		} else {
			path = fmt.Sprintf("%d", loc.ID)
		}

		return tx.Model(loc).Update("path", path).Error
	})
	return loc.ID, err
}

// AddTranslation adds a single translation for a location.
func (r *PostgresRepository) AddTranslation(ctx context.Context, trans *domain.LocationTranslation) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "location_id"}, {Name: "lang_code"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"formatted_address",
			"short_formatted_address",
		}),
	}).Create(trans).Error
}

// GetChildren returns direct children of a parent location.
func (r *PostgresRepository) GetChildren(ctx context.Context, parentID int64, lang string) ([]domain.SearchResult, error) {
	var results []domain.SearchResult
	err := r.db.WithContext(ctx).Table("locations l").
		Select("l.id, COALESCE(lt.name, '') AS name, l.type").
		Joins("LEFT JOIN LATERAL (SELECT name FROM location_translations WHERE location_id = l.id ORDER BY (CASE WHEN lang_code = ? THEN 0 WHEN lang_code = 'en' THEN 1 ELSE 2 END) LIMIT 1) lt ON true", lang).
		Where("l.parent_id = ?", parentID).
		Order("lt.name").
		Scan(&results).Error

	return results, err
}

// UpdateLocation updates core attributes of a location and bumps updated_at.
func (r *PostgresRepository) UpdateLocation(ctx context.Context, loc *domain.Location) error {
	return r.db.WithContext(ctx).Table("locations").Where("id = ?", loc.ID).Updates(map[string]interface{}{
		"type":          loc.Type,
		"external_type": loc.ExternalType,
		"lat":           loc.Lat,
		"lng":           loc.Lng,
		"updated_at":    gorm.Expr("NOW()"),
	}).Error
}

// ReplaceAliases clears existing aliases and replaces them with new ones.
func (r *PostgresRepository) ReplaceAliases(ctx context.Context, locationID int64, aliases []domain.LocationAlias) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("location_aliases").Where("location_id = ?", locationID).Delete(&domain.LocationAlias{}).Error; err != nil {
			return err
		}
		for i := range aliases {
			aliases[i].LocationID = locationID
			if err := tx.Table("location_aliases").Create(&aliases[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ──────────────────────────────────────────────
// Post Queries
// ──────────────────────────────────────────────

// CreatePost inserts a new post and updates location_stats using GORM builder.
func (r *PostgresRepository) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Insert the post
		if err := tx.Create(post).Error; err != nil {
			return err
		}

		// 2. Upsert location_stats using builder
		scoreInc := 1.0
		if post.MediaType == "photo" {
			scoreInc = 2.5
		} else if post.MediaType == "video" {
			scoreInc = 3.0
		}

		photoInc := 0
		if post.MediaType == "photo" {
			photoInc = 1
		}
		videoInc := 0
		if post.MediaType == "video" {
			videoInc = 1
		}

		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "location_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"total_posts":    gorm.Expr("location_stats.total_posts + ?", 1),
				"total_photos":   gorm.Expr("location_stats.total_photos + ?", photoInc),
				"total_videos":   gorm.Expr("location_stats.total_videos + ?", videoInc),
				"last_post_at":   gorm.Expr("NOW()"),
				"trending_score": gorm.Expr("location_stats.trending_score + ?", scoreInc),
			}),
		}).Create(&domain.LocationStats{
			LocationID:    post.LocationID,
			TotalPosts:    1,
			TotalPhotos:   int64(photoInc),
			TotalVideos:   int64(videoInc),
			LastPostAt:    &post.CreatedAt,
			TrendingScore: scoreInc,
		}).Error
	})

	return post, err
}

// GetPostsByLocation returns posts for a location and descendants using builder.
func (r *PostgresRepository) GetPostsByLocation(ctx context.Context, locationID int64, lang string, limit, offset int) ([]domain.PostWithLocation, error) {
	var loc domain.Location
	r.db.WithContext(ctx).Select("path").First(&loc, locationID)

	var posts []domain.PostWithLocation
	query := r.db.WithContext(ctx).Table("posts p").
		Select("p.id, p.user_id, p.content, p.media_type, p.location_id, COALESCE(lt.name, '') AS location_name, l.type AS location_type, p.created_at").
		Joins("JOIN locations l ON l.id = p.location_id").
		Joins("LEFT JOIN LATERAL (SELECT name FROM location_translations WHERE location_id = l.id ORDER BY (CASE WHEN lang_code = ? THEN 0 WHEN lang_code = 'en' THEN 1 ELSE 2 END) LIMIT 1) lt ON true", lang)

	if loc.Path != "" {
		query = query.Where("l.path <@ ?", loc.Path)
	} else {
		query = query.Where("p.location_id = ?", locationID)
	}

	err := query.Order("p.created_at DESC").Limit(limit).Offset(offset).Scan(&posts).Error
	return posts, err
}

// GetPosts returns all recent posts across all locations, with optional filtering.
func (r *PostgresRepository) GetPosts(ctx context.Context, locationID *int64, lang string, limit, offset int) ([]domain.PostWithLocation, error) {
	var posts []domain.PostWithLocation
	query := r.db.WithContext(ctx).Table("posts p").
		Select("p.id, p.user_id, p.content, p.media_type, p.location_id, COALESCE(lt.name, '') AS location_name, l.type AS location_type, p.created_at").
		Joins("JOIN locations l ON l.id = p.location_id").
		Joins("LEFT JOIN LATERAL (SELECT name FROM location_translations WHERE location_id = l.id ORDER BY (CASE WHEN lang_code = ? THEN 0 WHEN lang_code = 'en' THEN 1 ELSE 2 END) LIMIT 1) lt ON true", lang)

	if locationID != nil && *locationID > 0 {
		var loc domain.Location
		if err := r.db.WithContext(ctx).Select("path").First(&loc, *locationID).Error; err == nil && loc.Path != "" {
			query = query.Where("l.path <@ ?", loc.Path)
		} else {
			query = query.Where("p.location_id = ?", *locationID)
		}
	}

	err := query.Order("p.created_at DESC").Limit(limit).Offset(offset).Scan(&posts).Error
	return posts, err
}

// ──────────────────────────────────────────────
// Stats & Trending Queries
// ──────────────────────────────────────────────

// GetStats returns the pre-aggregated stats for a location.
func (r *PostgresRepository) GetStats(ctx context.Context, locationID int64) (*domain.LocationStats, error) {
	var stats domain.LocationStats
	err := r.db.WithContext(ctx).Table("location_stats").Where("location_id = ?", locationID).First(&stats).Error
	if err != nil {
		return nil, nil // Silently return nil for no stats
	}
	return &stats, nil
}

// GetTrending returns top trending locations ordered by score.
func (r *PostgresRepository) GetTrending(ctx context.Context, lang string, limit int) ([]domain.TrendingLocation, error) {
	var results []domain.TrendingLocation
	err := r.db.WithContext(ctx).Table("location_stats ls").
		Select("ls.location_id, COALESCE(lt.name, '') AS name, l.type, ls.trending_score AS score").
		Joins("JOIN locations l ON l.id = ls.location_id").
		Joins("LEFT JOIN LATERAL (SELECT name FROM location_translations WHERE location_id = l.id ORDER BY (CASE WHEN lang_code = ? THEN 0 WHEN lang_code = 'en' THEN 1 ELSE 2 END) LIMIT 1) lt ON true", lang).
		Where("ls.trending_score > 0").
		Order("ls.trending_score DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// FindByParentAndName searches for a location by name, type, and specific parent_id scoping.
func (r *PostgresRepository) FindByParentAndName(ctx context.Context, parentID *int64, name string, locType domain.LocationType, lang string) (*domain.Location, error) {
	var loc domain.Location
	query := r.db.WithContext(ctx).Table("locations l").
		Select("l.id, l.type, l.external_type, l.lat, l.lng, l.parent_id, l.path, l.provider, l.created_at, l.updated_at").
		Joins("JOIN location_translations lt ON lt.location_id = l.id").
		Where("lt.name = ? AND lt.lang_code = ? AND l.type = ?", name, lang, locType)

	if parentID != nil {
		query = query.Where("l.parent_id = ?", *parentID)
	} else {
		query = query.Where("l.parent_id IS NULL")
	}

	err := query.First(&loc).Error
	if err != nil {
		return nil, err
	}
	return &loc, nil
}
