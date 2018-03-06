package Spotify

import (
	"fmt"
	"strings"
)

// BaseAddress for Spotify endpoints.
const baseAddress = "https://api.spotify.com/v1/"

// URI identifies an artist, album, track, or category.  For example,
// spotify:track:6rqhFgbbKwnb9MLmUQDhG6
type URI string

// ID is a base-62 identifier for an artist, track, album, etc.
// It can be found at the end of a spotify.URI.
type ID string

func (id *ID) String() string {
	return string(*id)
}

// SimpleTrack contains basic info about a track.
type SimpleTrack struct {
	Artists []SimpleArtist `json:"artists"`
	// A list of the countries in which the track can be played,
	// identified by their ISO 3166-1 alpha-2 codes.
	AvailableMarkets []string `json:"available_markets"`
	// The disc number (usually 1 unless the album consists of more than one disc).
	DiscNumber int `json:"disc_number"`
	// The length of the track, in milliseconds.
	Duration int `json:"duration_ms"`
	// Whether or not the track has explicit lyrics.
	// true => yes, it does; false => no, it does not.
	Explicit bool `json:"explicit"`
	// External URLs for this track.
	ExternalURLs map[string]string `json:"external_urls"`
	// A link to the Web API endpoint providing full details for this track.
	Endpoint string `json:"href"`
	ID       ID     `json:"id"`
	Name     string `json:"name"`
	// A URL to a 30 second preview (MP3) of the track.
	PreviewURL string `json:"preview_url"`
	// The number of the track.  If an album has several
	// discs, the track number is the number on the specified
	// DiscNumber.
	TrackNumber int `json:"track_number"`
	URI         URI `json:"uri"`
}

// Trackattributes needed to hold return value from Track Analysis endpoint on Spotify.
type AnalysisAttributes struct {
	Danceability     float64 `json:"danceability"`
	Energy           float64 `json:"energy"`
	Key              int     `json:"key"`
	Loudness         float64 `json:"loudness"`
	Mode             int     `json:"mode"`
	Speechiness      float64 `json:"speechiness"`
	Acousticness     float64 `json:"acousticness"`
	Instrumentalness float64 `json:"instrumentalness"`
	Liveness         float64 `json:"liveness"`
	Valence          float64 `json:"valence"`
	Tempo            float64 `json:"tempo"`
	Type             string  `json:"type"`
	ID               string  `json:"id"`
	URI              string  `json:"uri"`
	TrackHref        string  `json:"track_href"`
	AnalysisURL      string  `json:"analysis_url"`
	DurationMs       int     `json:"duration_ms"`
	TimeSignature    int     `json:"time_signature"`
}

type Mode int

// Mode of music, 0 or 1, equivalent to "Dur/moll"
const (
	Minor Mode = iota
	Major
)

type Key int

const (
	C Key = iota
	CSharp
	D
	DSharp
	E
	F
	FSharp
	G
	GSharp
	A
	ASharp
	B
	DFlat = CSharp
	EFlat = DSharp
	GFlat = FSharp
	AFlat = GSharp
	BFlat = ASharp
)

func (st SimpleTrack) String() string {
	return fmt.Sprintf("TRACK<[%s] [%s]>", st.ID, st.Name)
}

func GetAudioFeatures(tracks []string, token string) ([]*AnalysisAttributes, error) {
	url := fmt.Sprintf("%saudio-features?ids=%s", baseAddress, strings.Join(tracks, ","))
	temp := struct {
		F []*AnalysisAttributes `json:"audio_features"`
	}{}

	err := Get(url, token, &temp)
	if err != nil {
		return nil, err
	}

	return temp.F, nil
}
