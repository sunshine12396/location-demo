package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/example/location-demo/internal/domain"
)

// GooglePlacesClient is a real implementation of domain.ExternalLocationProvider
// that calls the Google Places Text Search (New) API.
type GooglePlacesClient struct {
	apiKey     string
	httpClient *http.Client
}

// NewGooglePlacesClient initializes the external provider.
func NewGooglePlacesClient(apiKey string) *GooglePlacesClient {
	return &GooglePlacesClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Google Places API Response Structs
type GooglePlacesResponse struct {
	Places []GooglePlace `json:"places"`
}

type GooglePlace struct {
	Id          string `json:"id"`
	DisplayName struct {
		Text         string `json:"text"`
		LanguageCode string `json:"languageCode"`
	} `json:"displayName"`
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Types []string `json:"types"`
}

// Search queries the Google Places API for locations.
func (c *GooglePlacesClient) Search(ctx context.Context, query string, lang string) ([]domain.Location, []domain.LocationTranslation, error) {
	if c.apiKey == "" {
		// If no API key is configured, fail gracefully or log it.
		return []domain.Location{}, []domain.LocationTranslation{}, nil
	}

	// We use the Google Places Text Search (New) API
	// Documentation: https://developers.google.com/maps/documentation/places/web-service/text-search
	reqURL := "https://places.googleapis.com/v1/places:searchText"

	// Prepare payload
	payload := map[string]interface{}{
		"textQuery":    query,
		"languageCode": lang,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("google_client: failed to marshal payload: %w", err)
	}
	
	importBytes := string(bodyBytes)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(importBytes))
	if err != nil {
		return nil, nil, fmt.Errorf("google_client: failed to create request: %w", err)
	}

	// Google Places API requires these headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", c.apiKey)
	req.Header.Set("X-Goog-FieldMask", "places.id,places.displayName,places.location,places.types")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("google_client: http request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("google_client: received non-200 status code: %d", res.StatusCode)
	}

	var gResponse GooglePlacesResponse
	if err := json.NewDecoder(res.Body).Decode(&gResponse); err != nil {
		return nil, nil, fmt.Errorf("google_client: failed to decode response: %w", err)
	}

	var locations []domain.Location
	var translations []domain.LocationTranslation

	// Map Google Places results to our domain models
	for _, place := range gResponse.Places {
		// Determine location type roughly based on Google's types
		locType := domain.LocationTypeCity
		for _, t := range place.Types {
			if t == "administrative_area_level_1" || t == "country" {
				locType = domain.LocationTypeCountry
				break
			}
			if t == "point_of_interest" || t == "tourist_attraction" || t == "landmark" {
				locType = domain.LocationTypeLandmark
				break
			}
		}

		locations = append(locations, domain.Location{
			ExternalID: place.Id, // Using the Google Places ID as our ExternalID
			Type:       locType,
			Lat:        place.Location.Latitude,
			Lng:        place.Location.Longitude,
		})

		translations = append(translations, domain.LocationTranslation{
			// Setting LocationID is not strictly necessary here as the Service layer
			// will link it when doing the InsertLocation batch, but let's be safe.
			LangCode: place.DisplayName.LanguageCode,
			Name:     place.DisplayName.Text,
		})
	}

	return locations, translations, nil
}
