package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/example/location-demo/internal/domain"
	"github.com/example/location-demo/pkg/rest/client"
)

// FetchByExternalID retrieves a highly detailed single location's data from the external provider.
func (c *GooglePlacesClient) FetchByExternalID(ctx context.Context, externalID string, lang string) (*domain.Location, *domain.LocationTranslation, []domain.LocationAlias, error) {
	if c.apiKey == "" {
		return nil, nil, nil, fmt.Errorf("google_client: no API key")
	}

	payload := client.Payload{
		PathVars: map[string]string{
			"placeId": externalID,
		},
		Header: map[string]string{
			"Accept-Language": lang,
		},
	}

	res, err := c.placeDetailClient.Send(ctx, payload)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("google_client: fetch_by_external_id: %w", err)
	}

	if res.Status != http.StatusOK {
		return nil, nil, nil, fmt.Errorf("google_client: fetch_by_external_id: received non-200 status code: %d, body: %s", res.Status, string(res.Body))
	}

	var place GooglePlace
	if err := json.Unmarshal(res.Body, &place); err != nil {
		return nil, nil, nil, fmt.Errorf("google_client: fetch_by_external_id: failed to decode response: %w", err)
	}

	// Parse hierarchy components (Store reversed for top-down iteration)
	var addressComponents []domain.AddressComponent
	for i := len(place.AddressComponents) - 1; i >= 0; i-- {
		comp := place.AddressComponents[i]

		if len(comp.Types) == 0 || !IsAllowSearchTextType(GoogleType(comp.Types[0])) {
			continue
		}

		compType := mapGoogleTypes(comp.Types)
		if compType == domain.LocationTypeCountry ||
			compType == domain.LocationTypeCity ||
			compType == domain.LocationTypeDistrict ||
			compType == domain.LocationTypeWard {

			addressComponents = append(addressComponents, domain.AddressComponent{
				LongName:     comp.LongText,
				ShortName:    comp.ShortText,
				ExternalType: comp.Types[0],
				Type:         compType,
				LanguageCode: comp.LanguageCode,
			})
		}
	}

	locType := mapGoogleTypes(place.Types)

	loc := &domain.Location{
		ExternalID:        place.Id,
		Type:              locType,
		ExternalType:      place.Types[0],
		Lat:               place.Location.Latitude,
		Lng:               place.Location.Longitude,
		Provider:          domain.ProviderGoogle,
		AddressComponents: addressComponents,
	}

	translation := &domain.LocationTranslation{
		LangCode:              lang,
		Name:                  place.DisplayName.Text,
		FormattedAddress:      place.FormattedAddress,
		ShortFormattedAddress: place.ShortFormattedAddress,
	}

	var aliases []domain.LocationAlias
	if place.FormattedAddress != "" {
		aliases = append(aliases, domain.LocationAlias{
			Alias: place.FormattedAddress,
		})
	}

	return loc, translation, aliases, nil
}
