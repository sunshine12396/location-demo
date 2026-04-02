package domain

import (
	"context"
	"time"
)

// --- Location Types ---

// LocationType defines the category of a location.
type LocationType string

const (
	LocationTypeCountry  LocationType = "country"
	LocationTypeCity     LocationType = "city"
	LocationTypeDistrict LocationType = "district"
	LocationTypeLandmark LocationType = "landmark"
	LocationTypeVenue    LocationType = "venue"
)

// --- Core Models ---

// Location is the central entity of the domain.
// It represents a single geographical point with hierarchy support.
type Location struct {
	ID         int64        `json:"id"`
	ExternalID string       `json:"external_id,omitempty"`
	Type       LocationType `json:"type"`
	Lat        float64      `json:"lat"`
	Lng        float64      `json:"lng"`
	ParentID   *int64       `json:"parent_id,omitempty"`
	Path       string       `json:"path,omitempty"`
	Slug       string       `json:"slug,omitempty"`
	IsVerified bool         `json:"is_verified"`
	CreatedAt  time.Time    `json:"created_at"`
}

// LocationTranslation holds the localized name for a location.
type LocationTranslation struct {
	LocationID int64  `json:"location_id"`
	LangCode   string `json:"lang_code"`
	Name       string `json:"name"`
}

// LocationAlias stores alternative names for a location (e.g., "Sai Gon" for Ho Chi Minh City).
type LocationAlias struct {
	ID         int64  `json:"id"`
	LocationID int64  `json:"location_id"`
	Alias      string `json:"alias"`
}

// Post represents user-generated content attached to a location.
type Post struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Content    string    `json:"content"`
	MediaType  string    `json:"media_type"` // text / photo / video
	LocationID int64     `json:"location_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// LocationStats holds pre-aggregated counters for a location.
type LocationStats struct {
	LocationID    int64     `json:"location_id"`
	TotalPosts    int64     `json:"total_posts"`
	TotalPhotos   int64     `json:"total_photos"`
	TotalVideos   int64     `json:"total_videos"`
	LastPostAt    *time.Time `json:"last_post_at,omitempty"`
	TrendingScore float64   `json:"trending_score"`
}

// TrendingLocation is a daily snapshot of a location's trending score.
type TrendingLocation struct {
	LocationID int64   `json:"location_id"`
	Name       string  `json:"name"`
	Type       LocationType `json:"type"`
	Score      float64 `json:"score"`
	Date       string  `json:"date"`
}

// --- API Response Models ---

// LocationDetail is the rich response returned to the client for a single location.
type LocationDetail struct {
	ID           int64                `json:"id"`
	Name         string               `json:"name"`
	Type         LocationType         `json:"type"`
	Lat          float64              `json:"lat"`
	Lng          float64              `json:"lng"`
	Slug         string               `json:"slug,omitempty"`
	IsVerified   bool                 `json:"is_verified"`
	Parent       *LocationSummary     `json:"parent,omitempty"`
	Stats        *LocationStats       `json:"stats,omitempty"`
	Translations []LocationTranslation `json:"translations,omitempty"`
}

// LocationSummary is a lightweight representation used for parent references and search results.
type LocationSummary struct {
	ID   int64        `json:"id"`
	Name string       `json:"name"`
	Type LocationType `json:"type"`
}

// SearchResult is what the search endpoint returns.
type SearchResult struct {
	ID      int64        `json:"id"`
	Name    string       `json:"name"`
	Type    LocationType `json:"type"`
	Country string       `json:"country,omitempty"`
}

// PostWithLocation is a post enriched with its location name for display.
type PostWithLocation struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Content      string    `json:"content"`
	MediaType    string    `json:"media_type"`
	LocationID   int64     `json:"location_id"`
	LocationName string    `json:"location_name"`
	LocationType LocationType `json:"location_type"`
	CreatedAt    time.Time `json:"created_at"`
}

// --- Repository Interface (Port) ---
// This is the contract the data layer must fulfill.
// The service layer depends on this interface, NOT on a concrete implementation.

type LocationRepository interface {
	// GetByID retrieves a single location with its translation for the given language.
	GetByID(ctx context.Context, id int64, lang string) (*LocationDetail, error)

	// SearchByAlias searches the alias table for a matching term.
	SearchByAlias(ctx context.Context, query string) ([]SearchResult, error)

	// SearchByTranslation searches the translations table for a matching term in a given language.
	SearchByTranslation(ctx context.Context, query string, lang string) ([]SearchResult, error)

	// InsertLocation creates a new location and its translations/aliases.
	InsertLocation(ctx context.Context, loc *Location, translations []LocationTranslation, aliases []LocationAlias) (int64, error)

	// GetChildren returns all direct children of a location.
	GetChildren(ctx context.Context, parentID int64, lang string) ([]SearchResult, error)

	// --- Post Operations ---

	// CreatePost inserts a new post and updates location stats.
	CreatePost(ctx context.Context, post *Post) (*Post, error)

	// GetPostsByLocation returns posts for a given location (and optionally all its descendants via path).
	GetPostsByLocation(ctx context.Context, locationID int64, lang string, limit, offset int) ([]PostWithLocation, error)

	// --- Stats & Trending ---

	// GetStats returns the pre-aggregated stats for a location.
	GetStats(ctx context.Context, locationID int64) (*LocationStats, error)

	// GetTrending returns the top trending locations for today.
	GetTrending(ctx context.Context, lang string, limit int) ([]TrendingLocation, error)
}

// --- External API Interface (Port) ---
// This defines how the service interacts with external geocoding APIs.

type ExternalLocationProvider interface {
	// Search queries an external API (OpenStreetMap, Google Places) for locations.
	Search(ctx context.Context, query string, lang string) ([]Location, []LocationTranslation, error)
}
