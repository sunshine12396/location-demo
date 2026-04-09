package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/example/location-demo/internal/domain"
	"github.com/example/location-demo/pkg/rest/client"
)

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
	Types                 []string `json:"types"`
	FormattedAddress      string   `json:"formattedAddress"`
	ShortFormattedAddress string   `json:"shortFormattedAddress"`
	AddressComponents     []struct {
		LongText              string   `json:"longText"`
		ShortText             string   `json:"shortText"`
		Types                 []string `json:"types"`
		LanguageCode          string   `json:"languageCode"`
		FormattedAddress      string   `json:"formattedAddress"`
		ShortFormattedAddress string   `json:"shortFormattedAddress"`
	} `json:"addressComponents"`
}

// Search queries the Google Places API for locations.
func (c *GooglePlacesClient) Search(ctx context.Context, query string, lang string) ([]domain.Location, []domain.LocationTranslation, error) {
	if c.apiKey == "" {
		return []domain.Location{}, []domain.LocationTranslation{}, nil
	}

	reqBody := map[string]interface{}{
		"textQuery":    query,
		"languageCode": lang,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	payload := client.Payload{
		Body: bodyBytes,
	}

	res, err := c.searchTextClient.Send(ctx, payload)
	if err != nil {
		return nil, nil, fmt.Errorf("google_client: search: %w", err)
	}

	if res.Status != http.StatusOK {
		return nil, nil, fmt.Errorf("google_client: search: received non-200 status code: %d, body: %s", res.Status, string(res.Body))
	}

	var gResponse GooglePlacesResponse
	if err := json.Unmarshal(res.Body, &gResponse); err != nil {
		return nil, nil, fmt.Errorf("google_client: search: failed to decode response: %w", err)
	}

	var locations []domain.Location
	var translations []domain.LocationTranslation

	for _, place := range gResponse.Places {
		locType := mapGoogleTypes(place.Types)
		placeExternalType := ""
		if len(place.Types) > 0 {
			placeExternalType = place.Types[0]
		}

		// Parse hierarchy components (Store reversed for top-down iteration)
		var addressComponents []domain.AddressComponent
		for i := len(place.AddressComponents) - 1; i >= 0; i-- {
			comp := place.AddressComponents[i]
			compType := mapGoogleTypes(comp.Types)

			externalType := ""
			if len(comp.Types) > 0 {
				externalType = comp.Types[0]
			}
			// Only include important administrative layers to keep hierarchy clean
			if compType == domain.LocationTypeCountry || compType == domain.LocationTypeCity || compType == domain.LocationTypeDistrict || compType == domain.LocationTypeWard {
				addressComponents = append(addressComponents, domain.AddressComponent{
					LongName:     comp.LongText,
					ShortName:    comp.ShortText,
					ExternalType: externalType,
					Type:         compType,
				})
			}
		}

		locations = append(locations, domain.Location{
			ExternalID:        place.Id,
			Type:              locType,
			ExternalType:      placeExternalType,
			Lat:               place.Location.Latitude,
			Lng:               place.Location.Longitude,
			AddressComponents: addressComponents,
		})

		translations = append(translations, domain.LocationTranslation{
			LangCode:              place.DisplayName.LanguageCode,
			Name:                  place.DisplayName.Text,
			FormattedAddress:      place.FormattedAddress,
			ShortFormattedAddress: place.ShortFormattedAddress,
		})
	}

	return locations, translations, nil
}
