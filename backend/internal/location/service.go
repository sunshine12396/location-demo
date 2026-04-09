package location

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/example/location-demo/internal/domain"
)

// ──────────────────────────────────────────────
// Async Executor (extensible for worker pool later)
// ──────────────────────────────────────────────

type AsyncExecutor interface {
	Go(fn func())
}

type goExecutor struct{}

func (g goExecutor) Go(fn func()) {
	go fn()
}

// ──────────────────────────────────────────────
// Service
// ──────────────────────────────────────────────

type Service struct {
	repo     domain.LocationRepository
	external domain.ExternalLocationProvider
	syncDays int
	async    AsyncExecutor
}

func NewService(
	repo domain.LocationRepository,
	external domain.ExternalLocationProvider,
	syncDays int,
) *Service {
	return &Service{
		repo:     repo,
		external: external,
		syncDays: syncDays,
		async:    goExecutor{},
	}
}

// ──────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────

func normalizeLang(lang string) string {
	if lang == "" {
		return "en"
	}
	return lang
}

func (s *Service) isOutdated(updatedAt time.Time) bool {
	if s.syncDays <= 0 {
		return false
	}
	threshold := time.Duration(s.syncDays) * 24 * time.Hour
	return time.Since(updatedAt) > threshold
}

// Collect missing translation (non-blocking)
func (s *Service) maybeCollectTranslation(ctx context.Context, detail *domain.LocationDetail, lang string) {
	if detail.Name != "" || s.external == nil || detail.ExternalID == "" {
		return
	}

	_, trans, _, err := s.external.FetchByExternalID(ctx, detail.ExternalID, lang)
	if err != nil || trans == nil {
		return
	}

	log.Printf("INFO: Collecting translation for %s [%s]", detail.ExternalID, lang)

	detail.Name = trans.Name

	s.async.Go(func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.repo.AddTranslation(bgCtx, trans)
	})
}

// Trigger async sync if outdated
func (s *Service) maybeTriggerBackgroundSync(detail *domain.LocationDetail, lang string) {
	if s.external == nil || detail.ExternalID == "" {
		return
	}
	if !s.isOutdated(detail.UpdatedAt) {
		return
	}

	log.Printf("INFO: Location %d outdated → trigger sync", detail.ID)

	s.async.Go(func() {
		s.syncFromExternalBackground(detail.ID, detail.ExternalID, lang)
	})
}

// ──────────────────────────────────────────────
// Location Operations
// ──────────────────────────────────────────────

func (s *Service) GetByID(ctx context.Context, id int64, lang string) (*domain.LocationDetail, error) {
	lang = normalizeLang(lang)

	detail, err := s.repo.GetByID(ctx, id, lang)
	if err != nil {
		return nil, fmt.Errorf("GetByID: %w", err)
	}

	s.maybeCollectTranslation(ctx, detail, lang)
	s.maybeTriggerBackgroundSync(detail, lang)

	return detail, nil
}

// Background sync
func (s *Service) syncFromExternalBackground(id int64, externalID, lang string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	loc, trans, aliases, err := s.external.FetchByExternalID(ctx, externalID, lang)
	if err != nil {
		log.Printf("WARN: sync failed %s: %v", externalID, err)
		return
	}

	loc.ID = id

	if err := s.repo.UpdateLocation(ctx, loc); err != nil {
		log.Printf("WARN: update failed %d: %v", id, err)
		return
	}

	if trans != nil {
		trans.LocationID = id
		_ = s.repo.AddTranslation(ctx, trans)
	}

	if len(aliases) > 0 {
		_ = s.repo.ReplaceAliases(ctx, id, aliases)
	}

	log.Printf("INFO: synced location %d", id)
}

// ──────────────────────────────────────────────
// Search (Waterfall)
// ──────────────────────────────────────────────

func (s *Service) Search(ctx context.Context, query, lang, locType string) ([]domain.SearchResult, error) {
	lang = normalizeLang(lang)

	var lType *domain.LocationType
	if locType != "" {
		t := domain.LocationType(locType)
		lType = &t
	}

	// 1. alias
	results, err := s.repo.SearchByAlias(ctx, query, lType)
	if err == nil && len(results) > 0 {
		return results, nil
	}

	// 2. translation
	results, err = s.repo.SearchByTranslation(ctx, query, lang, lType)
	if err == nil && len(results) > 0 {
		return results, nil
	}

	// 3. external fallback
	if s.external == nil {
		return []domain.SearchResult{}, nil
	}

	locations, translations, err := s.external.Search(ctx, query, lang)
	if err != nil {
		return nil, fmt.Errorf("Search external: %w", err)
	}

	results = make([]domain.SearchResult, 0, len(locations))
	for i := range locations {
		results = append(results, domain.SearchResult{
			ID:         0,
			ExternalID: locations[i].ExternalID,
			Name:       translations[i].Name,
			Type:       locations[i].Type,
		})
	}

	return results, nil
}

func (s *Service) LocalSearch(ctx context.Context, query, lang, locType string) ([]domain.SearchResult, error) {
	lang = normalizeLang(lang)

	var lType *domain.LocationType
	if locType != "" {
		t := domain.LocationType(locType)
		lType = &t
	}

	results, err := s.repo.SearchByAlias(ctx, query, lType)
	if err == nil && len(results) > 0 {
		return results, nil
	}

	results, err = s.repo.SearchByTranslation(ctx, query, lang, lType)
	if results == nil {
		results = []domain.SearchResult{}
	}
	return results, err
}

func (s *Service) Autocomplete(ctx context.Context, query, lang string) ([]domain.SearchResult, error) {
	if query == "" || s.external == nil {
		return []domain.SearchResult{}, nil
	}
	res, err := s.external.Autocomplete(ctx, query, normalizeLang(lang))
	if res == nil {
		res = []domain.SearchResult{}
	}
	return res, err
}

func (s *Service) GetChildren(ctx context.Context, parentID int64, lang string) ([]domain.SearchResult, error) {
	res, err := s.repo.GetChildren(ctx, parentID, normalizeLang(lang))
	if res == nil {
		res = []domain.SearchResult{}
	}
	return res, err
}

// ──────────────────────────────────────────────
// Hydration
// ──────────────────────────────────────────────

func (s *Service) EnsureLocationByExternalID(ctx context.Context, externalID, lang string) (int64, error) {
	lang = normalizeLang(lang)

	// 1. check local
	loc, err := s.repo.GetByExternalID(ctx, externalID)
	if err == nil && loc != nil {
		if s.isOutdated(loc.UpdatedAt) {
			s.async.Go(func() {
				s.syncFromExternalBackground(loc.ID, externalID, lang)
			})
		}
		return loc.ID, nil
	}

	// 2. fetch external
	newLoc, trans, aliases, err := s.external.FetchByExternalID(ctx, externalID, lang)
	if err != nil {
		return 0, fmt.Errorf("fetch external: %w", err)
	}

	// 3. resolve parents
	var parentID *int64
	for _, comp := range newLoc.AddressComponents {
		// Avoid circular reference or self-parenting
		if comp.LongName == trans.Name && comp.Type == newLoc.Type {
			continue
		}
		parentID = s.ensureParent(ctx, comp, newLoc.Lat, newLoc.Lng, lang, parentID)
	}
	newLoc.ParentID = parentID

	// 4. save
	id, err := s.repo.InsertLocation(ctx, newLoc, []domain.LocationTranslation{*trans}, aliases)
	if err != nil {
		return 0, fmt.Errorf("insert location: %w", err)
	}

	return id, nil
}

// ensure parent hierarchy
func (s *Service) ensureParent(
	ctx context.Context,
	comp domain.AddressComponent,
	lat, lng float64,
	lang string,
	lastParentID *int64,
) *int64 {

	if comp.LongName == "" {
		return lastParentID
	}

	loc, err := s.repo.FindByParentAndName(ctx, lastParentID, comp.LongName, comp.Type, lang)
	if err == nil && loc != nil {
		id := loc.ID
		return &id
	}

	newLoc := domain.Location{
		ExternalType: comp.ExternalType,
		Type:         comp.Type,
		ParentID:     lastParentID,
		Provider:     domain.ProviderGoogle,
	}

	var translations []domain.LocationTranslation

	if s.external != nil {
		input := domain.ResolveNameToIDInput{
			Name:         comp.LongName,
			Lang:         lang,
			IncludedType: comp.ExternalType,
			Lat:          lat,
			Lng:          lng,
		}

		out, err := s.external.ResolveNameToID(ctx, input)
		if err == nil {
			newLoc.ExternalID = out.PlaceID
			newLoc.Lat = out.Lat
			newLoc.Lng = out.Lng

			translations = append(translations, domain.LocationTranslation{
				LangCode:              lang,
				Name:                  comp.LongName,
				FormattedAddress:      out.FormattedAddress,
				ShortFormattedAddress: out.ShortFormattedAddress,
			})
		}
	}

	aliases := []domain.LocationAlias{{Alias: comp.LongName}}
	if comp.ShortName != "" {
		aliases = append(aliases, domain.LocationAlias{Alias: comp.ShortName})
	}

	id, err := s.repo.InsertLocation(ctx, &newLoc, translations, aliases)
	if err != nil {
		log.Printf("ERROR: save parent %s failed: %v", comp.LongName, err)
		return lastParentID
	}

	return &id
}

// ──────────────────────────────────────────────
// Post
// ──────────────────────────────────────────────

func (s *Service) CreatePost(ctx context.Context, post *domain.Post, externalID, lang string) (*domain.Post, error) {
	if post.Content == "" {
		return nil, fmt.Errorf("content required")
	}

	if externalID != "" {
		id, err := s.EnsureLocationByExternalID(ctx, externalID, lang)
		if err != nil {
			return nil, fmt.Errorf("hydration failed: %w", err)
		}
		post.LocationID = id
	}

	if post.LocationID == 0 {
		return nil, fmt.Errorf("location required")
	}

	switch post.MediaType {
	case "", "text":
		post.MediaType = "text"
	case "photo", "video":
	default:
		return nil, fmt.Errorf("invalid media_type")
	}

	if post.UserID == 0 {
		post.UserID = 1
	}

	return s.repo.CreatePost(ctx, post)
}
func (s *Service) GetPostsByLocation(ctx context.Context, locationID int64, lang string, limit, offset int) ([]domain.PostWithLocation, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	res, err := s.repo.GetPostsByLocation(ctx, locationID, normalizeLang(lang), limit, offset)
	if res == nil {
		res = []domain.PostWithLocation{}
	}
	return res, err
}

func (s *Service) GetPosts(ctx context.Context, locationID *int64, lang string, limit, offset int) ([]domain.PostWithLocation, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	res, err := s.repo.GetPosts(ctx, locationID, normalizeLang(lang), limit, offset)
	if res == nil {
		res = []domain.PostWithLocation{}
	}
	return res, err
}

// ──────────────────────────────────────────────
// Stats
// ──────────────────────────────────────────────

func (s *Service) GetStats(ctx context.Context, locationID int64) (*domain.LocationStats, error) {
	return s.repo.GetStats(ctx, locationID)
}

func (s *Service) GetTrending(ctx context.Context, lang string, limit int) ([]domain.TrendingLocation, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	return s.repo.GetTrending(ctx, normalizeLang(lang), limit)
}
