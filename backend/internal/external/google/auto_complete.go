package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/example/location-demo/internal/domain"
	"github.com/example/location-demo/pkg/rest/client"
)

type GoogleAutocompleteResponse struct {
	Suggestions []struct {
		PlacePrediction struct {
			Text struct {
				Text string `json:"text"`
			} `json:"text"`
			PlaceId string   `json:"placeId"`
			Types   []string `json:"types"`
		} `json:"placePrediction"`
	} `json:"suggestions"`
}

// Autocomplete provides suggestions as the user types using Google Places Autocomplete API.
func (c *GooglePlacesClient) Autocomplete(ctx context.Context, query string, lang string) ([]domain.SearchResult, error) {
	if c.apiKey == "" {
		return []domain.SearchResult{}, nil
	}

	reqBody := map[string]interface{}{
		"input":        query,
		"languageCode": lang,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	payload := client.Payload{
		Body: bodyBytes,
	}

	res, err := c.autocompleteClient.Send(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("google_client: autocomplete: %w", err)
	}

	if res.Status != http.StatusOK {
		return nil, fmt.Errorf("google_client: autocomplete: received status code: %d, body: %s", res.Status, string(res.Body))
	}

	var gResponse GoogleAutocompleteResponse
	if err := json.Unmarshal(res.Body, &gResponse); err != nil {
		return nil, fmt.Errorf("google_client: autocomplete: failed to decode response: %w", err)
	}

	var results []domain.SearchResult
	for _, sugg := range gResponse.Suggestions {
		pred := sugg.PlacePrediction
		results = append(results, domain.SearchResult{
			ExternalID: pred.PlaceId,
			Name:       pred.Text.Text,
			Type:       mapGoogleTypes(pred.Types),
		})
	}

	return results, nil
}
