package google

import (
	"net/http"

	"github.com/example/location-demo/pkg/rest/client"
)

// GooglePlacesClient is a real implementation of domain.ExternalLocationProvider
// that calls the Google Places Text Search (New) API.
type GooglePlacesClient struct {
	apiKey string
	pool   *client.SharedPool

	autocompleteClient *client.HttpClient
	searchTextClient   *client.HttpClient
	placeDetailClient  *client.HttpClient
}

// NewGooglePlacesClient initializes the external provider.
func NewGooglePlacesClient(apiKey string) *GooglePlacesClient {
	pool := client.NewSharedPool()

	// Autocomplete Client
	autocompleteClient, _ := client.NewHttpClient(client.HttpClientConfig{
		URL:    "https://places.googleapis.com/v1/places:autocomplete",
		Method: http.MethodPost,
		Headers: map[string]string{
			"X-Goog-Api-Key": apiKey,
		},
	}, pool, client.WithServiceName("google-places-autocomplete"))

	// Search Text Client
	searchTextClient, _ := client.NewHttpClient(client.HttpClientConfig{
		URL:    "https://places.googleapis.com/v1/places:searchText",
		Method: http.MethodPost,
		Headers: map[string]string{
			"X-Goog-Api-Key":   apiKey,
			"X-Goog-FieldMask": "places.id,places.displayName,places.location,places.types,places.addressComponents,places.formattedAddress,places.shortFormattedAddress",
		},
	}, pool, client.WithServiceName("google-places-search"))

	// Place Detail Client
	placeDetailClient, _ := client.NewHttpClient(client.HttpClientConfig{
		URL:    "https://places.googleapis.com/v1/places/:placeId",
		Method: http.MethodGet,
		Headers: map[string]string{
			"X-Goog-Api-Key":   apiKey,
			"X-Goog-FieldMask": "id,types,location,addressComponents,displayName,formattedAddress,shortFormattedAddress",
		},
	}, pool, client.WithServiceName("google-places-detail"))

	return &GooglePlacesClient{
		apiKey:             apiKey,
		pool:               pool,
		autocompleteClient: autocompleteClient,
		searchTextClient:   searchTextClient,
		placeDetailClient:  placeDetailClient,
	}
}
