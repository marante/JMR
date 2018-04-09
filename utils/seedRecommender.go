package utils

import (
	"github.com/marante/JMR/Spotify"
)

// Seed returns a spotify.Seed object containing exactly 5 or less items.
// These will be used for the seed
func Seed(tracks []Spotify.RecentlyPlayedItem, contextTrackIds []string) Spotify.Seeds {
	var songs []string
	var artists []string
	songs = append(songs, contextTrackIds...)
	for _, items := range tracks {
		songs = append(songs, items.Track.ID.String())
		for _, item := range items.Track.Artists {
			artists = append(artists, item.ID.String())
		}
	}

	// Make map to hold items and keep which seeds are artists and songs.
	m := make(map[string]map[string]int)
	m["Songs"] = MakeMap(songs)
	m["Artists"] = MakeMap(artists)

	// Reduce it down to 5 of each
	songOrder := MapReduce(m["Songs"])
	artistOrder := MapReduce(m["Artists"])

	// Compare entries, and see which has the highest hitrate.
	seeds := Comparator(songOrder, artistOrder)
	return seeds
}
