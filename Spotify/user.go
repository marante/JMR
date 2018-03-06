package Spotify

// UserInfo provides information from the users mobilephone, which are needed for recommendation
type UserInfo struct {
	Token   string `json:"token,omitempty"`
	Context struct {
		ContextTracks []string `json:"contextTracks,omitempty"`
		AnalyzeTracks []string `json:"analyzeTracks,omitempty"`
		Country       string   `json:"country,omitempty"`
		Bpm           int      `json:"bpm,omitempty"`
	} `json:"context,omitempty"`
}
