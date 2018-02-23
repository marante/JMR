package utils

import (
	"fmt"
	"github.com/marante/JMR/Spotify"
	"sort"
)

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
		if ss[0][0].Value > ss[1][0].Value {
			seeds.Tracks = append(seeds.Tracks, Spotify.ID(ss[0][0].Key))
			ss[0] = ss[0][1:]
		} else {
			seeds.Artists = append(seeds.Artists, Spotify.ID(ss[1][0].Key))
			ss[1] = ss[1][1:]
		}
	}

	return seeds
}
