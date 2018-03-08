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
	"strings"
)

var (
	store *sessions.CookieStore
)

// Below code simplifies and makes error handling for handlers more concrete.
type appError struct {
	Error   error
	Message string
	Code    int
}

type trackObject struct {
	URI  Spotify.URI `json:"uri"`
	Name string      `json:"name"`
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := ah(w, r); err != nil {
		http.Error(w, err.Message, err.Code)
	}
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
	router.Handle("/analysis", appHandler(TrackAnalysis)).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, router)))
}

func Index(w http.ResponseWriter, r *http.Request) *appError {
	fmt.Fprintf(w, "The server is currently up and active.")
	return nil
}

func RecentlyPlayed(w http.ResponseWriter, r *http.Request) *appError {
	decoder := json.NewDecoder(r.Body)
	var t Spotify.UserInfo
	if err := decoder.Decode(&t); err != nil {
		fmt.Println(err)
		return &appError{err, "Error trying to decode JSON body.", 400}
	}
	defer r.Body.Close()

	opts := &Spotify.RecentlyPlayedOptions{Limit: 50}
	tracks, err := Spotify.GetRecentlyPlayedTracksOpt(t.Token, opts)
	if err != nil {
		fmt.Println(err)
		return &appError{err, "Error trying to decode JSON body.", 400}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tracks)
	if err != nil {
		fmt.Println(err)
		return &appError{err, "Error encoding data to JSON", 400}
	}
	return nil
}

func Recommendations(w http.ResponseWriter, r *http.Request) *appError {
	var t Spotify.UserInfo
	var attr *Spotify.TrackAttributes
	var err error
	var seeds Spotify.Seeds
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		fmt.Println(err)
		return &appError{err, "Error trying to decode JSON body.", 400}
	}
	defer r.Body.Close()

	if len(t.Context.AnalyzeTracks) > 0 {
		attr, err = utils.GetTrackAttributes(&t)
		if err != nil {
			fmt.Println(err)
			return &appError{err, "Error trying to retrieve results from track analysis.", 400}
		}
		seeds = utils.Seed(nil, t.Context.AnalyzeTracks)
	} else {
		recentlyPlayed := &Spotify.RecentlyPlayedOptions{Limit: 50}
		tracks, err := Spotify.GetRecentlyPlayedTracksOpt(t.Token, recentlyPlayed)
		if err != nil {
			fmt.Println(err)
			return &appError{err, "Error trying to retrieve recently played tracks.", 400}
		}
		seeds = utils.Seed(tracks, nil)
	}

	options := &Spotify.Options{
		Country: strings.ToUpper(t.Context.Country),
	}

	// There might be occasions when this returns > 5 values, which is OK.
	recommendations, err := Spotify.GetRecommendations(seeds, attr, options, t.Token)
	if err != nil {
		fmt.Println(err)
		return &appError{err, "Error trying to retrieve recommendations.", 400}
	}

	var trackObjs []trackObject

	for _, val := range recommendations.Tracks {
		item := trackObject{URI: val.URI, Name: val.Name}
		for _, artist := range val.Artists {
			item.Name += " - " + artist.Name
		}
		trackObjs = append(trackObjs, item)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(trackObjs)
	if err != nil {
		fmt.Println(err)
		return &appError{err, "Error encoding data to JSON", 400}
	}
	return nil
}

func TrackAnalysis(w http.ResponseWriter, r *http.Request) *appError {
	decoder := json.NewDecoder(r.Body)
	var t Spotify.UserInfo
	if err := decoder.Decode(&t); err != nil {
		fmt.Println(err)
		return &appError{err, "Error trying to decode JSON body.", 400}
	}
	defer r.Body.Close()

	if t.Context.AnalyzeTracks != nil {
		attributes, err := Spotify.GetAudioFeatures(t.Context.AnalyzeTracks, t.Token)
		if err != nil {
			fmt.Println(err)
			return &appError{err, "There was an error trying to analyze tracks.", 400}
		}
		err = json.NewEncoder(w).Encode(attributes)
		if err != nil {
			fmt.Println(err)
			return &appError{err, "Error encoding data to JSON", 400}
		}
	}
	return nil
}
