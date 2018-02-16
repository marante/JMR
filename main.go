package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/marante/JMR/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
	"os"
)

// CHANGE AT YOUR OWN WILL
//const redirectURI = "http://localhost:8080/callback"
const redirectURI = "https://thawing-tor-40623.herokuapp.com/callback"

var (
	auth = spotify.NewAuthenticator(
		redirectURI,
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
		spotify.ScopeUserModifyPlaybackState,
		spotify.ScopeUserReadRecentlyPlayed,
		spotify.ScopeUserReadPlaybackState)
	state = uuid.NewV4().String()
)

// AuthorizedClient is a client ready to be used for API calls.
type AuthorizedClient struct {
	client spotify.Client
}

// SpotifySongs struct containing information we want to send to endpoints.
type SpotifySongs struct {
	TrackURIs []spotify.URI
	Songs     []string
}

func (c *AuthorizedClient) index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to front page")
	tok, err := c.client.Token()
	if err != nil {
		fmt.Println(err)
	}
	user, err := c.client.CurrentUser()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintln(w, "You are logged in as:", user.ID)
	fmt.Fprintln(w, "You are logged in as:", tok.AccessToken)
}

func (c *AuthorizedClient) completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		fmt.Println(err)
	}
	// Checks if the state matches.
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		fmt.Println("State mismatch: %s != %s\n", st, state)
	}
	// use token to get authenticated client
	c.client = auth.NewClient(tok)
	http.Redirect(w, r, "/", 302)
}

func (c *AuthorizedClient) recommendations(w http.ResponseWriter, r *http.Request) {
	// used for seeds.
	var genres []string
	var artists []*spotify.FullArtist
	// Getting the 50 recently played tracks for a given user
	opts := &spotify.RecentlyPlayedOptions{Limit: 50}
	tracks, err := c.client.PlayerRecentlyPlayedOpt(opts)
	errCheck(err)

	// Looping over the result to extract artists
	for _, val := range tracks {
		for _, artist := range val.Track.Artists {
			item, err := c.client.GetArtist(artist.ID)
			if err != nil {
				log.Fatal(err)
			}
			artists = append(artists, item)
		}
	}

	for _, item := range artists {
		genres = append(genres, item.Genres...)
	}

	// Passing tracks and genres.
	seeds := utils.Seed(tracks, genres)

	attr := spotify.
		NewTrackAttributes().
		MinTempo(120).
		MinEnergy(0.7).
		MinValence(0.6)

	rec, err := c.client.GetRecommendations(seeds, attr, nil)
	errCheck(err)

	// REGION TESTING WILL REMOVE LATER (WHY U NO HEFF REGIONS)
	var trackURIs []spotify.URI
	var trackNames []string

	for _, item := range rec.Tracks {
		trackURIs = append(trackURIs, item.URI)
		trackNames = append(trackNames, item.Name)
	}

	seedInfo := SpotifySongs{
		TrackURIs: trackURIs,
		Songs:     trackNames,
	}

	playOpts := spotify.PlayOptions{
		URIs: trackURIs,
	}

	err = c.client.PlayOpt(&playOpts)
	errCheck(err)
	// REGION TESTING WILL REMOVE LATER (WHY U NO HEFF REGIONS)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(seedInfo)
	errCheck(err)
}

func authorize(w http.ResponseWriter, r *http.Request) {
	url := auth.AuthURL(state)
	http.Redirect(w, r, url, 302)
}

func main() {
	// returns a router object from the Gorilla/mux package.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	var router = mux.NewRouter()
	client := &AuthorizedClient{}
	router.HandleFunc("/", client.index).Methods("GET")
	router.HandleFunc("/auth", authorize).Methods("GET")
	router.HandleFunc("/callback", client.completeAuth).Methods("GET")
	router.HandleFunc("/recommendations", client.recommendations).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func errCheck(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
