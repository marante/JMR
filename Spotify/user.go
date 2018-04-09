package Spotify

// UserInfo provides information from the users mobilephone, which are needed for recommendation
type UserInfo struct {
	Token     string   `json:"token,omitempty"`
	UriTracks []string `json:"uris,omitempty"`
}
