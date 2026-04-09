package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/example/location-demo/internal/domain"
	"github.com/example/location-demo/pkg/rest/client"
)

// ResolveNameToID translates a location name (e.g., "District 1") into a Google Place ID,
// using coordinates for bias and a specific type to ensure geographical and categorical accuracy.
func (c *GooglePlacesClient) ResolveNameToID(ctx context.Context, input domain.ResolveNameToIDInput) (*domain.ResolveNameToIDOutput, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("google_client: resolve: no API key configured")
	}

	reqBody := map[string]interface{}{
		"textQuery":    input.Name,
		"languageCode": input.Lang,
		"locationBias": map[string]interface{}{
			"circle": map[string]interface{}{
				"center": map[string]interface{}{
					"latitude":  input.Lat,
					"longitude": input.Lng,
				},
				"radius": 5000.0, // 5km radius for disambiguation
			},
		},
	}

	if input.IncludedType != "" {
		reqBody["includedType"] = input.IncludedType
	}

	bodyBytes, _ := json.Marshal(reqBody)

	payload := client.Payload{
		Body: bodyBytes,
	}

	res, err := c.searchTextClient.Send(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("google_client: resolve: %w", err)
	}

	if res.Status != http.StatusOK {
		return nil, fmt.Errorf("google_client: resolve: received non-200 status code: %d, body: %s", res.Status, string(res.Body))
	}

	var gResponse GooglePlacesResponse
	if err := json.Unmarshal(res.Body, &gResponse); err != nil {
		return nil, fmt.Errorf("google_client: resolve: failed to decode response: %w", err)
	}

	if len(gResponse.Places) == 0 {
		return nil, fmt.Errorf("google_client: resolve: no results found for name: %s", input.Name)
	}

	return &domain.ResolveNameToIDOutput{
		PlaceID:               gResponse.Places[0].Id,
		Lat:                   gResponse.Places[0].Location.Latitude,
		Lng:                   gResponse.Places[0].Location.Longitude,
		FormattedAddress:      gResponse.Places[0].FormattedAddress,
		ShortFormattedAddress: gResponse.Places[0].ShortFormattedAddress,
	}, nil
}
