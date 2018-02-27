package Spotify

// SimpleArtist returns a simple artist object explained in the Spotify web api.
type SimpleArtist struct {
	Name         string            `json:"name"`
	ID           ID                `json:"id"`
	URI          URI               `json:"uri"`
	Endpoint     string            `json:"href"`
	ExternalURLs map[string]string `json:"external_urls"`
}
