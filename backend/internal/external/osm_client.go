package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/example/location-demo/internal/domain"
)

// OSMClient is an implementation of domain.ExternalLocationProvider
// that calls the OpenStreetMap (OSM) Nominatim API.
type OSMClient struct {
	httpClient *http.Client
}

// NewOSMClient initializes the open-source external provider.
func NewOSMClient() *OSMClient {
	return &OSMClient{
		httpClient: &http.Client{
			// Nominatim can sometimes be slow, so give it a slightly longer timeout
			Timeout: 15 * time.Second,
		},
	}
}

// NominatimResponse Structs
// https://nominatim.org/release-docs/develop/api/Search/
type NominatimResponse []NominatimPlace

type NominatimPlace struct {
	PlaceID     int64  `json:"place_id"`
	OSMID       int64  `json:"osm_id"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
	Class       string `json:"class"` // e.g., "place", "boundary", "highway"
	Type        string `json:"type"`  // e.g., "city", "country", "administrative"
}

// Search queries the OSM Nominatim API for locations.
func (c *OSMClient) Search(ctx context.Context, query string, lang string) ([]domain.Location, []domain.LocationTranslation, error) {
	// Nominatim endpoint for search
	// NOTE: Nominatim REQUIRES a User-Agent header, otherwise it may block the request.
	reqURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&accept-language=%s",
		url.QueryEscape(query), url.QueryEscape(lang))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("osm_client: failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "DemoApp/1.0")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("osm_client: http request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("osm_client: received non-200 status code: %d", res.StatusCode)
	}

	var nResponse NominatimResponse
	if err := json.NewDecoder(res.Body).Decode(&nResponse); err != nil {
		return nil, nil, fmt.Errorf("osm_client: failed to decode response: %w", err)
	}

	var locations []domain.Location
	var translations []domain.LocationTranslation

	// Map OSM results to our domain models
	for _, place := range nResponse {

		var lat, lng float64
		fmt.Sscanf(place.Lat, "%f", &lat)
		fmt.Sscanf(place.Lon, "%f", &lng)

		// Determine location type safely
		locType := domain.LocationTypeCity

		if place.Type == "country" || place.Type == "state" {
			locType = domain.LocationTypeCountry
		} else if place.Type == "city" || place.Type == "town" || place.Type == "municipality" {
			locType = domain.LocationTypeCity
		} else if place.Type == "district" || place.Type == "suburb" || place.Type == "borough" {
			locType = domain.LocationTypeDistrict
		} else if strings.Contains(place.Class, "tourism") || strings.Contains(place.Class, "historic") {
			locType = domain.LocationTypeLandmark
		} else if strings.Contains(place.Class, "amenity") || strings.Contains(place.Class, "shop") {
			locType = domain.LocationTypeVenue
		} else if place.Class == "boundary" && place.Type == "administrative" {
			// If it's an administrative boundary but not clearly a country, default to city to avoid everything being a country
			locType = domain.LocationTypeCity
		}

		// OSM's display name is often very long
		parts := strings.Split(place.DisplayName, ",")
		primaryName := strings.TrimSpace(parts[0])
		
		// Optional: Append Country/State to disambiguate identical names
		if len(parts) > 1 {
			primaryName = fmt.Sprintf("%s, %s", primaryName, strings.TrimSpace(parts[len(parts)-1]))
		}

		// Prevent adding multiple identical display names to avoid spamming the UI
		isDuplicate := false
		for _, existing := range translations {
			if existing.Name == primaryName {
				isDuplicate = true
				break
			}
		}
		if isDuplicate {
			continue
		}

		locations = append(locations, domain.Location{
			ExternalID: fmt.Sprintf("osm_%d", place.OSMID),
			Type:       locType,
			Lat:        lat,
			Lng:        lng,
		})

		translations = append(translations, domain.LocationTranslation{
			LangCode: lang,
			Name:     primaryName,
		})
	}

	return locations, translations, nil
}
