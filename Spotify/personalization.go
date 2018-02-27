package Spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	// Base URL for API calls
	baseURL = "https://api.spotify.com/v1/"
	// defaultRetryDurationS helps us fix an apparent server bug whereby we will
	// be told to retry but not be given a wait-interval.
	defaultRetryDuration = time.Second * 5
	// rateLimitExceededStatusCode is the code that the server returns when our
	// request frequency is too high.
	rateLimitExceededStatusCode = 429
	// Maximum number of seeds
	MaxNumberOfSeeds = 5
)

type PlaybackContext struct {
	// ExternalURLs of the context, or null if not available.
	ExternalURLs map[string]string `json:"external_urls"`
	// Endpoint of the context, or null if not available.
	Endpoint string `json:"href"`
	// Type of the item's context. Can be one of album, artist or playlist.
	Type string `json:"type"`
	// URI is the Spotify URI for the context.
	URI URI `json:"uri"`
}

type RecentlyPlayedItem struct {
	// Track is the track information
	Track SimpleTrack `json:"track"`

	// PlayedAt is the time that this song was played
	PlayedAt time.Time `json:"played_at"`

	// PlaybackContext is the current playback context
	PlaybackContext PlaybackContext `json:"context"`
}

type RecentlyPlayedResult struct {
	Items []RecentlyPlayedItem `json:"items"`
}

type RecentlyPlayedOptions struct {
	// Limit is the maximum number of items to return. Must be no greater than
	// fifty.
	Limit int

	// AfterEpochMs is a Unix epoch in milliseconds that describes a time after
	// which to return songs.
	AfterEpochMs int64

	// BeforeEpochMs is a Unix epoch in milliseconds that describes a time
	// before which to return songs.
	BeforeEpochMs int64
}

type seeds struct {
	//...
}

// Error represents an error returned by the Spotify Web API.
type Error struct {
	// A short description of the error.
	Message string `json:"message"`
	// The HTTP status code.
	Status int `json:"status"`
}

func (e Error) Error() string {
	return e.Message
}

// decodeError decodes an Error from an io.Reader.
func decodeError(resp *http.Response) error {
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(responseBody) == 0 {
		return fmt.Errorf("spotify: HTTP %d: %s (body empty)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	buf := bytes.NewBuffer(responseBody)

	var e struct {
		E Error `json:"error"`
	}
	err = json.NewDecoder(buf).Decode(&e)
	if err != nil {
		return fmt.Errorf("spotify: couldn't decode error: (%d) [%s]", len(responseBody), responseBody)
	}

	if e.E.Message == "" {
		// Some errors will result in there being a useful status-code but an
		// empty message, which will confuse the user (who only has access to
		// the message and not the code). An example of this is when we send
		// some of the arguments directly in the HTTP query and the URL ends-up
		// being too long.

		e.E.Message = fmt.Sprintf("spotify: unexpected HTTP %d: %s (empty error)",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return e.E
}

// shouldRetry determines whether the status code indicates that the
// previous operation should be retried at a later time
func shouldRetry(status int) bool {
	return status == http.StatusAccepted || status == http.StatusTooManyRequests
}

func retryDuration(resp *http.Response) time.Duration {
	raw := resp.Header.Get("Retry-After")
	if raw == "" {
		return defaultRetryDuration
	}
	seconds, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return defaultRetryDuration
	}
	return time.Duration(seconds) * time.Second
}

func get(url string, token string, result interface{}) error {
	for {
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating a new request", err)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == rateLimitExceededStatusCode {
			time.Sleep(retryDuration(resp))
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return decodeError(resp)
		}

		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return err
		}

		break
	}
	return nil
}

// GetRecentlyPlayedTracks gets the latest 20 played tracks from users Spotify history
func GetRecentlyPlayedTracks(token string) ([]RecentlyPlayedItem, error) {
	return GetRecentlyPlayedTracksOpt(token, nil)
}

// GetRecentlyPlayedTracksOpt does the same thing GetRecentlyPlayedTracks, but with options
// (If they are provided)
func GetRecentlyPlayedTracksOpt(token string, opt *RecentlyPlayedOptions) ([]RecentlyPlayedItem, error) {
	spotifyURL := baseURL + "me/player/recently-played"
	if opt != nil {
		v := url.Values{}
		if opt.Limit != 0 {
			v.Set("limit", strconv.FormatInt(int64(opt.Limit), 10))
		}
		if opt.BeforeEpochMs != 0 {
			v.Set("before", strconv.FormatInt(int64(opt.BeforeEpochMs), 10))
		}
		if opt.AfterEpochMs != 0 {
			v.Set("after", strconv.FormatInt(int64(opt.AfterEpochMs), 10))
		}
		if params := v.Encode(); params != "" {
			spotifyURL += "?" + params
		}
	}

	result := RecentlyPlayedResult{}
	err := get(spotifyURL, token, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

func GetRecommendations(seeds Seeds, trackAttributes *TrackAttributes, opt *Options, token string) (*Recommendations, error) {
	v := url.Values{}

	if seeds.count() == 0 {
		return nil, fmt.Errorf("spotify: at least one seed is required")
	}
	if seeds.count() > MaxNumberOfSeeds {
		return nil, fmt.Errorf("spotify: exceeded maximum of %d seeds", MaxNumberOfSeeds)
	}

	setSeedValues(seeds, v)
	setTrackAttributesValues(trackAttributes, v)

	if opt != nil {
		if opt.Limit != nil {
			v.Set("limit", strconv.Itoa(*opt.Limit))
		}
		if opt.Country != nil {
			v.Set("market", *opt.Country)
		}
	}

	spotifyURL := baseURL + "recommendations?" + v.Encode()

	fmt.Println(spotifyURL)

	var recommendations Recommendations
	err := get(spotifyURL, token, &recommendations)
	if err != nil {
		return nil, err
	}

	return &recommendations, err
}
