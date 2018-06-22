package utils

import (
	"fmt"

	"github.com/marante/JMR/Spotify"
)

// SeedTrack returns a spotify.Seed object containing exactly 5 or less items.
func SeedTrack(token string, tracks []string) Spotify.Seeds {
	var songs []string
	var artists []string
	fullTracks, err := Spotify.GetTracks(token, tracks)
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range fullTracks {
		songs = append(songs, string(item.ID))
		for _, trackArtists := range item.Artists {
			artists = append(artists, trackArtists.Name)
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

// SeedGenre returns a spotify.Seeds object with a genre for recommendations.
func SeedGenre(genre string) Spotify.Seeds {
	seeds := Spotify.Seeds{
		Genre: genre,
	}
	return seeds
}
