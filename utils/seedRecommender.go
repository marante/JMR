package utils

import (
	"fmt"
	"github.com/marante/JMR/Spotify"
)

var (
	songs   []string
	artists []string
)

// Seed returns a spotify.Seed object containing exactly 5 or less items.
// These will be used for the seed
func Seed(tracks []Spotify.RecentlyPlayedItem) Spotify.Seeds {
	// artists & songs
	for _, items := range tracks {
		songs = append(songs, items.Track.ID.String())
		for _, item := range items.Track.Artists {
			artists = append(artists, item.ID.String())
		}
	}

	m := make(map[string]map[string]int)

	m["Songs"] = MakeMap(songs)
	m["Artists"] = MakeMap(artists)

	songOrder := MapReduce(m["Songs"])
	artistOrder := MapReduce(m["Artists"])

	seeds := Comparator(songOrder, artistOrder)

	fmt.Println("# of artistseeds: ", len(seeds.Artists))
	fmt.Println("# of trackseeds: ", len(seeds.Tracks))
	return seeds
}
