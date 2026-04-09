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
	LocationTypeWard     LocationType = "ward"
	LocationTypeLandmark LocationType = "landmark"
	LocationTypeVenue    LocationType = "venue"
	LocationTypeAddress  LocationType = "address"
	LocationTypeRegion   LocationType = "region"
	LocationTypePlace    LocationType = "place"
)

type ExternalProvider string

const (
	ProviderOSM    ExternalProvider = "osm"
	ProviderGoogle ExternalProvider = "google"
	ProviderLocal  ExternalProvider = "local"
)

// --- Core Models ---

// Location is the central entity of the domain.
// It represents a single geographical point with hierarchy support.
type Location struct {
	ID                int64              `json:"id" gorm:"primaryKey"`
	ExternalID        string             `json:"external_id" gorm:"index:idx_locations_external_id,unique"`
	Type              LocationType       `json:"type"`
	ExternalType      string             `json:"external_type,omitempty"`
	Lat               float64            `json:"lat"`
	Lng               float64            `json:"lng"`
	ParentID          *int64             `json:"parent_id,omitempty"`
	Path              string             `json:"path,omitempty"`
	Provider          ExternalProvider   `json:"provider" gorm:"default:'local'"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	AddressComponents []AddressComponent `json:"-" gorm:"-"`

	// Associations for preloading
	Translations []LocationTranslation `json:"translations,omitempty" gorm:"foreignKey:LocationID"`
	Aliases      []LocationAlias       `json:"aliases,omitempty" gorm:"foreignKey:LocationID"`
	Stats        *LocationStats        `json:"stats,omitempty" gorm:"foreignKey:LocationID"`
	Parent       *Location             `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
}

// AddressComponent represents a part of a location's hierarchy (e.g., city, country).
type AddressComponent struct {
	LongName     string
	ShortName    string
	ExternalType string
	Type         LocationType
	LanguageCode string
}

// LocationTranslation holds the localized name for a location.
type LocationTranslation struct {
	LocationID            int64  `json:"location_id" gorm:"primaryKey"`
	LangCode              string `json:"lang_code" gorm:"primaryKey"`
	Name                  string `json:"name"`
	FormattedAddress      string `json:"formatted_address,omitempty"`
	ShortFormattedAddress string `json:"short_formatted_address,omitempty"`
}

// LocationAlias stores alternative names for a location (e.g., "Sai Gon" for Ho Chi Minh City).
type LocationAlias struct {
	ID         int64  `json:"id" gorm:"primaryKey"`
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
	LocationID    int64      `json:"location_id"`
	TotalPosts    int64      `json:"total_posts"`
	TotalPhotos   int64      `json:"total_photos"`
	TotalVideos   int64      `json:"total_videos"`
	LastPostAt    *time.Time `json:"last_post_at,omitempty"`
	TrendingScore float64    `json:"trending_score"`
}

// TrendingLocation is a daily snapshot of a location's trending score.
type TrendingLocation struct {
	LocationID int64        `json:"location_id"`
	Name       string       `json:"name"`
	Type       LocationType `json:"type"`
	Score      float64      `json:"score"`
	Date       string       `json:"date"`
}

// --- API Response Models ---

// LocationDetail is the rich response returned to the client for a single location.
type LocationDetail struct {
	ID                    int64                 `json:"id" gorm:"primaryKey"`
	ExternalID            string                `json:"external_id"`
	Name                  string                `json:"name" gorm:"column:name"`
	FormattedAddress      string                `json:"formatted_address,omitempty"`
	ShortFormattedAddress string                `json:"short_formatted_address,omitempty"`
	Type                  LocationType          `json:"type"`
	Lat                   float64               `json:"lat"`
	Lng                   float64               `json:"lng"`
	Provider              ExternalProvider      `json:"provider"`
	UpdatedAt             time.Time             `json:"updated_at"`
	Parent                *LocationSummary      `json:"parent,omitempty" gorm:"-"`
	Stats                 *LocationStats        `json:"stats,omitempty" gorm:"-"`
	Translations          []LocationTranslation `json:"translations,omitempty" gorm:"-"`
}

// LocationSummary is a lightweight representation used for parent references and search results.
type LocationSummary struct {
	ID   int64        `json:"id"`
	Name string       `json:"name"`
	Type LocationType `json:"type"`
}

// SearchResult is what the search endpoint returns.
type SearchResult struct {
	ID         int64        `json:"id"`
	ExternalID string       `json:"external_id,omitempty"`
	Name       string       `json:"name"`
	Type       LocationType `json:"type"`
	Country    string       `json:"country,omitempty"`
}

// PostWithLocation is a post enriched with its location name for display.
type PostWithLocation struct {
	ID           int64        `json:"id"`
	UserID       int64        `json:"user_id"`
	Content      string       `json:"content"`
	MediaType    string       `json:"media_type"`
	LocationID   int64        `json:"location_id"`
	LocationName string       `json:"location_name"`
	LocationType LocationType `json:"location_type"`
	CreatedAt    time.Time    `json:"created_at"`
}

// --- Repository Interface (Port) ---
// This is the contract the data layer must fulfill.
// The service layer depends on this interface, NOT on a concrete implementation.

type LocationRepository interface {
	// GetByID retrieves a single location with its translation for the given language.
	GetByID(ctx context.Context, id int64, lang string) (*LocationDetail, error)

	// GetByExternalID retrieves a location record by its external provider ID.
	GetByExternalID(ctx context.Context, externalID string) (*Location, error)

	// SearchByAlias searches the alias table for a matching term with optional type filtering.
	SearchByAlias(ctx context.Context, query string, locType *LocationType) ([]SearchResult, error)

	// SearchByTranslation searches the translations table for a matching term with optional type filtering.
	SearchByTranslation(ctx context.Context, query string, lang string, locType *LocationType) ([]SearchResult, error)

	InsertLocation(ctx context.Context, loc *Location, translations []LocationTranslation, aliases []LocationAlias) (int64, error)

	// AddTranslation adds a single translation for a location.
	AddTranslation(ctx context.Context, translation *LocationTranslation) error

	// GetChildren returns all direct children of a location.
	GetChildren(ctx context.Context, parentID int64, lang string) ([]SearchResult, error)

	// UpdateLocation updates core attributes of a location and bumps updated_at.
	UpdateLocation(ctx context.Context, loc *Location) error

	// ReplaceAliases clears existing aliases and replaces them with new ones.
	ReplaceAliases(ctx context.Context, locationID int64, aliases []LocationAlias) error

	// --- Post Operations ---

	// CreatePost inserts a new post and updates location stats.
	CreatePost(ctx context.Context, post *Post) (*Post, error)

	// GetPostsByLocation returns posts for a given location (and optionally all its descendants via path).
	GetPostsByLocation(ctx context.Context, locationID int64, lang string, limit, offset int) ([]PostWithLocation, error)

	// GetPosts returns recent posts across all locations, with optional location filtering.
	GetPosts(ctx context.Context, locationID *int64, lang string, limit, offset int) ([]PostWithLocation, error)

	// --- Stats & Trending ---

	// GetStats returns the pre-aggregated stats for a location.
	GetStats(ctx context.Context, locationID int64) (*LocationStats, error)

	// GetTrending returns the top trending locations for today.
	GetTrending(ctx context.Context, lang string, limit int) ([]TrendingLocation, error)

	// FindByParentAndName searches for a location by its name, type, and specific parent.
	FindByParentAndName(ctx context.Context, parentID *int64, name string, locType LocationType, lang string) (*Location, error)
}

// --- External API Interface (Port) ---
// This defines how the service interacts with external geocoding APIs.

type ResolveNameToIDInput struct {
	Name         string
	Lang         string
	IncludedType string
	Lat          float64
	Lng          float64
}

type ResolveNameToIDOutput struct {
	PlaceID               string
	Lat                   float64
	Lng                   float64
	FormattedAddress      string
	ShortFormattedAddress string
}

type ExternalLocationProvider interface {
	// Search queries an external API (OpenStreetMap, Google Places) for locations.
	Search(ctx context.Context, query string, lang string) ([]Location, []LocationTranslation, error)

	// Autocomplete provides suggestions as the user types.
	Autocomplete(ctx context.Context, query string, lang string) ([]SearchResult, error)

	// FetchByExternalID retrieves a single location's data from the external provider.
	FetchByExternalID(ctx context.Context, externalID string, lang string) (*Location, *LocationTranslation, []LocationAlias, error)

	// ResolveNameToID translates a location name (e.g., "District 1") into a Google Place ID,
	// using coordinates for bias and a specific type to ensure geographical and categorical accuracy.
	ResolveNameToID(ctx context.Context, input ResolveNameToIDInput) (*ResolveNameToIDOutput, error)
}
