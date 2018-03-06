package utils

import (
	"github.com/marante/JMR/Spotify"
	"strings"
)

func AttributeSelector(user *Spotify.UserInfo) *Spotify.TrackAttributes {
	var bpm float64
	if user.Context.Bpm != 0 {
		bpm = float64(user.Context.Bpm)
	}
	return Spotify.NewTrackAttributes().MinTempo(bpm)
}

func OptionsSelector(user *Spotify.UserInfo) *Spotify.Options {
	var options Spotify.Options
	limit := 5
	upper := strings.ToUpper(user.Context.Country)
	if user.Context.Country != "" {
		options.Country = &upper
		options.Limit = &limit
	}
	return &options
}
