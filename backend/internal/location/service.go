package location

import (
	"context"
	"fmt"
	"log"

	"github.com/example/location-demo/internal/domain"
)

// Service implements the business logic for location operations.
// It depends on interfaces (ports), not concrete implementations.
type Service struct {
	repo     domain.LocationRepository
	external domain.ExternalLocationProvider // can be nil if no external API configured
}

// NewService creates a new location service.
func NewService(repo domain.LocationRepository, external domain.ExternalLocationProvider) *Service {
	return &Service{
		repo:     repo,
		external: external,
	}
}

// ──────────────────────────────────────────────
// Location Operations
// ──────────────────────────────────────────────

// GetByID retrieves a single location detail by ID (now includes stats).
func (s *Service) GetByID(ctx context.Context, id int64, lang string) (*domain.LocationDetail, error) {
	if lang == "" {
		lang = "en"
	}

	detail, err := s.repo.GetByID(ctx, id, lang)
	if err != nil {
		return nil, fmt.Errorf("service.GetByID: %w", err)
	}

	return detail, nil
}

// Search implements the Waterfall Strategy:
//  1. Search aliases (e.g., "Sai Gon" → Ho Chi Minh City)
//  2. Search translations (e.g., "Hồ Chí Minh" in Vietnamese)
//  3. Fallback to external API (OpenStreetMap)
//  4. Save external results to local DB for future searches
func (s *Service) Search(ctx context.Context, query string, lang string) ([]domain.SearchResult, error) {
	if lang == "" {
		lang = "en"
	}

	// Step 1: Search alias table
	results, err := s.repo.SearchByAlias(ctx, query)
	if err != nil {
		log.Printf("WARN: alias search failed: %v", err)
	}
	if len(results) > 0 {
		return results, nil
	}

	// Step 2: Search translations
	results, err = s.repo.SearchByTranslation(ctx, query, lang)
	if err != nil {
		log.Printf("WARN: translation search failed: %v", err)
	}
	if len(results) > 0 {
		return results, nil
	}

	// Step 3: Fallback to external API (if configured)
	if s.external == nil {
		return []domain.SearchResult{}, nil
	}

	locations, translations, err := s.external.Search(ctx, query, lang)
	if err != nil {
		return nil, fmt.Errorf("service.Search: external fallback: %w", err)
	}

	// Step 4: Hydrate — save results to local DB for next time
	var hydrated []domain.SearchResult
	for i, loc := range locations {
		// Collect translations for this specific location
		var locTranslations []domain.LocationTranslation
		for _, t := range translations {
			if t.LocationID == 0 || t.LocationID == loc.ID {
				locTranslations = append(locTranslations, domain.LocationTranslation{
					LangCode: t.LangCode,
					Name:     t.Name,
				})
			}
		}

		id, insertErr := s.repo.InsertLocation(ctx, &locations[i], locTranslations, nil)
		if insertErr != nil {
			log.Printf("WARN: failed to hydrate location: %v", insertErr)
			continue
		}

		name := loc.ExternalID
		if len(locTranslations) > 0 {
			name = locTranslations[0].Name
		}

		hydrated = append(hydrated, domain.SearchResult{
			ID:   id,
			Name: name,
			Type: loc.Type,
		})
	}

	return hydrated, nil
}

// GetChildren returns sub-locations of a parent (e.g., cities in a country).
func (s *Service) GetChildren(ctx context.Context, parentID int64, lang string) ([]domain.SearchResult, error) {
	if lang == "" {
		lang = "en"
	}
	return s.repo.GetChildren(ctx, parentID, lang)
}

// ──────────────────────────────────────────────
// Post Operations
// ──────────────────────────────────────────────

// CreatePost validates and creates a new post, updating location stats.
func (s *Service) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	if post.Content == "" {
		return nil, fmt.Errorf("service.CreatePost: content is required")
	}
	if post.LocationID == 0 {
		return nil, fmt.Errorf("service.CreatePost: location_id is required")
	}

	// Validate media type
	switch post.MediaType {
	case "text", "photo", "video":
		// valid
	case "":
		post.MediaType = "text"
	default:
		return nil, fmt.Errorf("service.CreatePost: invalid media_type '%s' (must be text/photo/video)", post.MediaType)
	}

	if post.UserID == 0 {
		post.UserID = 1 // Default demo user
	}

	return s.repo.CreatePost(ctx, post)
}

// GetPostsByLocation retrieves posts tagged at a location.
func (s *Service) GetPostsByLocation(ctx context.Context, locationID int64, lang string, limit, offset int) ([]domain.PostWithLocation, error) {
	if lang == "" {
		lang = "en"
	}
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetPostsByLocation(ctx, locationID, lang, limit, offset)
}

// ──────────────────────────────────────────────
// Stats & Trending
// ──────────────────────────────────────────────

// GetStats returns the pre-aggregated stats for a location.
func (s *Service) GetStats(ctx context.Context, locationID int64) (*domain.LocationStats, error) {
	return s.repo.GetStats(ctx, locationID)
}

// GetTrending returns the top trending locations.
func (s *Service) GetTrending(ctx context.Context, lang string, limit int) ([]domain.TrendingLocation, error) {
	if lang == "" {
		lang = "en"
	}
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	return s.repo.GetTrending(ctx, lang, limit)
}
