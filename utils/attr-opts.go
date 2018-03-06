package utils

import (
	"github.com/marante/JMR/Spotify"
	"strings"
)

func AttributeSelector(user *Spotify.UserInfo) *Spotify.TrackAttributes {
	if user.Context.Bpm != 0 {
		return Spotify.NewTrackAttributes().MinTempo(float64(user.Context.Bpm))
	}
	return nil
}

func OptionsSelector(user *Spotify.UserInfo) *Spotify.Options {
	var options Spotify.Options
	limit := 5
	upper := strings.ToUpper(user.Context.Country)
	if user.Context.Country != "" {
		options.Country = &upper
	}
	options.Limit = &limit
	return &options
}
