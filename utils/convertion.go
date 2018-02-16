package utils

import (
	"fmt"
	"github.com/zmb3/spotify"
)

func ConvertToSpotifyID(i interface{}) []spotify.ID {
	var list []spotify.ID
	switch t := i.(type) {
	case []string:
		for _, a := range t {
			list = append(list, spotify.ID(a))
		}
		return list
	default:
		fmt.Println("Unknown type, whoops...")
	}
	return nil
}
