package Spotify

import (
	"net/url"
	"strconv"
	"strings"
)

type Seeds struct {
	Artists []ID
	Tracks  []ID
	Genre   string
}

// count returns the total number of seeds contained in s
func (s Seeds) count() int {
	if s.Genre != "" {
		return len(s.Artists) + len(s.Tracks) + 1
	}
	return len(s.Artists) + len(s.Tracks)
}

// Recommendations contains a list of recommended tracks based on seeds
type Recommendations struct {
	Seeds  []RecommendationSeed `json:"seeds"`
	Tracks []SimpleTrack        `json:"tracks"`
}

type RecommendationSeed struct {
	AfterFilteringSize int    `json:"afterFilteringSize"`
	AfterRelinkingSize int    `json:"afterRelinkingSize"`
	Endpoint           string `json:"href"`
	ID                 ID     `json:"id"`
	InitialPoolSize    int    `json:"initialPoolSize"`
	Type               string `json:"type"`
}

// Options contains optional parameters that can be provided
// to various API calls.  Only the non-nil fields are used
// in queries.
type Options struct {
	// Country is an ISO 3166-1 alpha-2 country code.  Provide
	// this parameter if you want the list of returned items to
	// be relevant to a particular country.  If omitted, the
	// results will be relevant to all countries.
	Country string
	// Limit is the maximum number of items to return.
	Limit int
	// Offset is the index of the first item to return.  Use it
	// with Limit to get the next set of items.
	Offset int
	// Timerange is the period of time from which to return results
	// in certain API calls. The three options are the following string
	// literals: "short", "medium", and "long"
	Timerange string
}

func toStringSlice(ids []ID) []string {
	result := make([]string, len(ids))
	for i, str := range ids {
		result[i] = str.String()
	}
	return result
}

// setSeedValues sets url values into v for each seed in seeds
func setSeedValues(seeds Seeds, v url.Values) {
	if len(seeds.Artists) != 0 {
		v.Set("seed_artists", strings.Join(toStringSlice(seeds.Artists), ","))
	}
	if len(seeds.Tracks) != 0 {
		v.Set("seed_tracks", strings.Join(toStringSlice(seeds.Tracks), ","))
	}
	if seeds.Genre != "" {
		v.Set("seed_genres", seeds.Genre)
	}
}

// setTrackAttributesValues sets track attributes values to the given url values
func setTrackAttributesValues(trackAttributes *TrackAttributes, values url.Values) {
	if trackAttributes == nil {
		return
	}
	for attr, val := range trackAttributes.intAttributes {
		values.Set(attr, strconv.Itoa(val))
	}
	for attr, val := range trackAttributes.floatAttributes {
		values.Set(attr, strconv.FormatFloat(val, 'f', -1, 64))
	}
}
