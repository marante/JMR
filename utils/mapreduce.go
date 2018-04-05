package utils

import (
	"fmt"
	"github.com/marante/JMR/Spotify"
	"math/rand"
	"sort"
	"time"
)

// TrackObject is an aggregate of a spotify URI and name of track.
type TrackObject struct {
	URI  Spotify.URI `json:"uri"`
	Name string      `json:"name"`
}

// Pair represents a custom map
type Pair struct {
	Key   string
	Value int
}

func (p *Pair) String() string {
	return fmt.Sprintf("Key: %s | Val: %v", p.Key, p.Value)
}

// MakeMap initializes map with values from an array.
func MakeMap(arr []string) map[string]int {
	m := make(map[string]int)
	for _, val := range arr {
		m[val]++
	}
	return m
}

// MapReduce is a custom mapreduce function to determine what inputs are most common
// and sorts them according to the value. Only returning top 5 values
func MapReduce(wordFrequencies map[string]int) []Pair {
	var ss []Pair
	for k, v := range wordFrequencies {
		ss = append(ss, Pair{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	if len(ss) < 5 {
		return ss[:len(ss)]
	}

	return ss[:5]
}

// Comparator compares values in []Pair structs to determine which has the highest values.
// Song artist genre
func Comparator(ss ...[]Pair) Spotify.Seeds {
	seeds := Spotify.Seeds{}
	for i := 0; i < 5; i++ {
		if len(ss[0]) == 0 && len(ss[1]) == 0 {
			break
		} else if len(ss[0]) == 0 {
			seeds.Artists = append(seeds.Artists, Spotify.ID(ss[1][0].Key))
			ss[1] = ss[1][1:]
		} else if len(ss[1]) == 0 {
			seeds.Tracks = append(seeds.Tracks, Spotify.ID(ss[0][0].Key))
			ss[0] = ss[0][1:]
		} else {
			if ss[0][0].Value > ss[1][0].Value {
				seeds.Tracks = append(seeds.Tracks, Spotify.ID(ss[0][0].Key))
				ss[0] = ss[0][1:]
			} else {
				seeds.Artists = append(seeds.Artists, Spotify.ID(ss[1][0].Key))
				ss[1] = ss[1][1:]
			}
		}
	}

	return seeds
}

func MapReduceRandom(token string, playlists []Spotify.SimplePlaylist) []TrackObject {
	var playlistTracks []*Spotify.PlaylistTrackPage
	var tracks []TrackObject
	for _, v := range playlists {
		trackPage, err := Spotify.GetPlaylistTracksOpt(token, "spotify", v.ID, nil, "")
		if err != nil {
			fmt.Println(err)
		}
		playlistTracks = append(playlistTracks, trackPage)
	}
	var artistsNames string
	for _, val := range playlistTracks {
		for _, value := range val.Tracks {
			artistsNames = ""
			for _, artists := range value.Track.Artists {
				if artistsNames == "" {
					artistsNames = artists.Name
				} else {
					artistsNames += " - " + artists.Name
				}
			}
			trackObj := TrackObject{URI: value.Track.URI, Name: value.Track.Name + " - " + artistsNames}
			tracks = append(tracks, trackObj)
		}
	}
	return tracks
}

func Randomizer(tracks []TrackObject) []TrackObject {
	if len(tracks) == 0 {
		return nil
	}
	var uris []TrackObject
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	for len(uris) < 20 {
		uri := tracks[rand.Intn(len(tracks))]
		if contains(uris, uri) {
			continue
		}
		uris = append(uris, uri)
	}
	return uris
}

func contains(s []TrackObject, e TrackObject) bool {
	for _, a := range s {
		if a.URI == e.URI {
			return true
		}
	}
	return false
}
