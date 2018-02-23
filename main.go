package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/marante/JMR/Spotify"
	"github.com/marante/JMR/utils"
	_ "github.com/zmb3/spotify"
	"log"
	"net/http"
	"os"
)

var (
	store *sessions.CookieStore
)

// UserInfo provides information from the users mobilephone, which are needed for recommendation
type UserInfo struct {
	Token      string `json:"token,omitempty"`
	DeviceName string `json:"deviceName,omitempty"`
	Context    struct {
		Time   string `json:"time,omitempty"`
		Loc    string `json:"loc,omitempty"`
		Motion string `json:"motion,omitempty"`
		Bpm    string `json:"bpm,omitempty"`
	} `json:"context,omitempty"`
}

// Below code simplifies and makes error handling for handlers more concrete.
type appError struct {
	Error   error
	Message string
	Code    int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := ah(w, r); err != nil {
		http.Error(w, err.Message, err.Code)
	}
}

type something struct {
	Token    string
	DeviceID string
	Context  *context
}

type context struct {
	Activity string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	var router = mux.NewRouter()
	router.Handle("/", appHandler(Index)).Methods("GET")
	router.Handle("/recently", appHandler(RecentlyPlayed)).Methods("POST")
	router.Handle("/recommendations", appHandler(Recommendations)).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, router)))
}

func Index(w http.ResponseWriter, r *http.Request) *appError {
	fmt.Fprintf(w, "The server is currently up and active.")
	return nil
}

func RecentlyPlayed(w http.ResponseWriter, r *http.Request) *appError {
	decoder := json.NewDecoder(r.Body)
	var t UserInfo
	if err := decoder.Decode(&t); err != nil {
		return &appError{err, "Error trying to decode JSON body.", 415}
	}
	defer r.Body.Close()
	// Getting the 50 recently played tracks for a given user
	opts := &Spotify.RecentlyPlayedOptions{Limit: 50}
	tracks, err := Spotify.GetRecentlyPlayedTracksOpt(t.Token, opts)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tracks)
	if err != nil {
		return &appError{err, "Error encoding data to JSON", 415}
	}
	return nil
}

func Recommendations(w http.ResponseWriter, r *http.Request) *appError {
	decoder := json.NewDecoder(r.Body)
	var t UserInfo
	if err := decoder.Decode(&t); err != nil {
		return &appError{err, "Error trying to decode JSON body.", 415}
	}
	defer r.Body.Close()

	// used for seeds.
	var artists []Spotify.SimpleArtist
	// Getting the 50 recently played tracks for a given user
	opts := &Spotify.RecentlyPlayedOptions{Limit: 50}
	tracks, err := Spotify.GetRecentlyPlayedTracksOpt(t.Token, opts)
	if err != nil {
		return &appError{err, "Error trying to retrieve recently played tracks.", 400}
	}

	// Looping over the result to extract artists
	for _, val := range tracks {
		for _, artist := range val.Track.Artists {
			artists = append(artists, artist)
		}
	}

	seeds := utils.Seed(tracks)
	attr := Spotify.
		NewTrackAttributes().
		MinTempo(120).
		MinEnergy(0.7).
		MinValence(0.6)

	recommendations, err := Spotify.GetRecommendations(seeds, attr, nil, t.Token)
	if err != nil {
		return &appError{err, "Error trying to retrieve recommendations.", 400}
	}

	var names []string

	for _, val := range recommendations.Tracks {
		names = append(names, val.Name)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(names)
	if err != nil {
		return &appError{err, "Error encoding data to JSON", 415}
	}
	return nil
}
