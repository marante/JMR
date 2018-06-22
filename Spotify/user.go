package Spotify

// UserInfo provides information from the users mobilephone, which are needed for recommendation
type UserInfo struct {
	Token     string   `json:"token,omitempty"`
	URITracks []string `json:"uris,omitempty"`
	Genre     string   `json:"genre,omitempty"`
}
