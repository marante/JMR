package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/marante/JMR/Spotify"
	. "github.com/marante/JMR/dao"
	. "github.com/marante/JMR/models"
	"github.com/marante/JMR/utils"
	_ "github.com/zmb3/spotify"
	_ "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
)

var (
	store *sessions.CookieStore
	dao   = LogsDAO{}
)

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

func init() {
	// Checks to use the correct env variable (Heroku vs Localhost).
	server := os.Getenv("SERVER")
	if server == "" {
		server = os.Getenv("SERVER_LOCAL")
	}
	database := os.Getenv("DATABASE")
	if server == "" {
		database = os.Getenv("DATABASE_LOCAL")
	}
	dao.Server = server
	dao.Database = database
	dao.Connect()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	var router = mux.NewRouter()
	router.Handle("/", appHandler(Index)).Methods("GET")
	router.Handle("/recently", appHandler(RecentlyPlayed)).Methods("POST")
	router.Handle("/random", appHandler(Random)).Methods("POST")
	router.Handle("/recommendations", appHandler(Recommendations)).Methods("POST")
	router.Handle("/analysis", appHandler(TrackAnalysis)).Methods("POST")

	// Routes for handling logging.
	router.Handle("/logs", appHandler(FindAllLogs)).Methods("GET")
	router.Handle("/logs/{id}", appHandler(FindLogById)).Methods("GET")
	router.Handle("/logs", appHandler(CreateLog)).Methods("POST")
	router.Handle("/logs", appHandler(UpdateLog)).Methods("PUT")
	router.Handle("/logs", appHandler(DeleteLog)).Methods("DELETE")
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

func Random(w http.ResponseWriter, r *http.Request) *appError {
	decoder := json.NewDecoder(r.Body)
	var t Spotify.UserInfo
	if err := decoder.Decode(&t); err != nil {
		fmt.Println(err)
		return &appError{err, "Error trying to decode JSON body.", 400}
	}
	defer r.Body.Close()

	opts := &Spotify.PlaylistOptions{
		Options: Spotify.Options{
			Limit: 10,
		},
	}

	_, lists, err := Spotify.FeaturedPlaylistsOpt(t.Token, opts)
	if err != nil {
		fmt.Println(err)
		return &appError{err, "Unable to retrieve featured playlists", 400}
	}

	tracks := utils.Randomizer(utils.MapReduceRandom(t.Token, lists.Playlists))

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

	if len(t.UriTracks) > 0 {
		attr, err = utils.GetTrackAttributes(&t)
		if err != nil {
			if err.Error() == "Only valid bearer authentication supported" || err.Error() == "The access token expired" {
				fmt.Println(err)
				return &appError{err, err.Error(), 401}
			}
			fmt.Println(err)
			return &appError{err, err.Error(), 400}
		}
		fmt.Println("With analysis of tracks")
		seeds = utils.Seed(nil, t.UriTracks)
	} else {
		recentlyPlayed := &Spotify.RecentlyPlayedOptions{Limit: 50}
		tracks, err := Spotify.GetRecentlyPlayedTracksOpt(t.Token, recentlyPlayed)
		if err != nil {
			if err.Error() == "Only valid bearer authentication supported" || err.Error() == "The access token expired" {
				fmt.Println(err)
				return &appError{err, err.Error(), 401}
			}
			fmt.Println(err)
			return &appError{err, err.Error(), 400}
		}
		seeds = utils.Seed(tracks, nil)
	}
	options := &Spotify.Options{
		Limit: 20,
	}

	fmt.Println(seeds)

	// There might be occasions when this returns > 5 values, which is OK.
	recommendations, err := Spotify.GetRecommendations(seeds, attr, options, t.Token)
	if err != nil {
		fmt.Println(err)
		return &appError{err, "Error trying to retrieve recommendations.", 400}
	}

	var trackObjs []utils.TrackObject
	for _, val := range recommendations.Tracks {
		item := utils.TrackObject{URI: val.URI, Name: val.Name}
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

	if t.UriTracks != nil {
		attributes, err := Spotify.GetAudioFeatures(t.UriTracks, t.Token)
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

func CreateLog(w http.ResponseWriter, r *http.Request) *appError {
	defer r.Body.Close()
	var log PlayerLog
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		fmt.Println(err)
		return &appError{err, "Error trying to decode JSON body.", http.StatusBadRequest}
	}
	log.ID = bson.NewObjectId()
	if log.UserID == "" {
		return &appError{nil, "Cannot add log without userid.", http.StatusBadRequest}
	}
	if err := dao.Insert(log); err != nil {
		fmt.Println(err)
		return &appError{err, "Error trying to insert logentry into DB", http.StatusBadRequest}
	}
	return nil
}

func FindLogById(w http.ResponseWriter, r *http.Request) *appError {
	params := mux.Vars(r)
	log, err := dao.FindById(params["id"])
	if err != nil {
		fmt.Println(err)
		return &appError{err, "Could not find document with specified id.", http.StatusNoContent}
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(log)
	return nil
}

func FindAllLogs(w http.ResponseWriter, r *http.Request) *appError {
	logs, err := dao.FindAll()
	if len(logs) == 0 {
		return &appError{nil, "No documents in DB", http.StatusNoContent}
	}
	if err != nil {
		fmt.Println(err)
		return &appError{err, "Error trying to fetch all records from repo", http.StatusBadRequest}
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(logs)
	err = json.NewEncoder(w).Encode(logs)
	return nil
}

func UpdateLog(w http.ResponseWriter, r *http.Request) *appError {
	defer r.Body.Close()
	var log PlayerLog
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		fmt.Println(err)
		return &appError{err, "Error when trying to decode body from request.", http.StatusBadRequest}
	}
	if err := dao.Update(log); err != nil {
		fmt.Println(err)
		return &appError{err, "Error when trying to update record associated with given id.", http.StatusBadRequest}
	}
	return nil
}

func DeleteLog(w http.ResponseWriter, r *http.Request) *appError {
	defer r.Body.Close()
	var log PlayerLog
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		fmt.Println(err)
		return &appError{err, "Error when trying to decode body from request.", http.StatusBadRequest}
	}
	if err := dao.Delete(log); err != nil {
		fmt.Println(err)
		return &appError{err, "Error when trying to update record associated with given id.", http.StatusBadRequest}
	}
	return nil
}
