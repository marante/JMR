package Spotify

import (
	"fmt"
	"net/url"
	"strconv"
)

type Followers struct {
	// The total number of followers.
	Count uint `json:"total"`
	// A link to the Web API endpoint providing full details of the followers,
	// or the empty string if this data is not available.
	Endpoint string `json:"href"`
}

// Image identifies an image associated with an item.
type Image struct {
	// The image height, in pixels.
	Height int `json:"height"`
	// The image width, in pixels.
	Width int `json:"width"`
	// The source URL of the image.
	URL string `json:"url"`
}

// PlaylistTracks contains details about the tracks in a playlist.
type PlaylistTracks struct {
	// A link to the Web API endpoint where full details of
	// the playlist's tracks can be retrieved.
	Endpoint string `json:"href"`
	// The total number of tracks in the playlist.
	Total uint `json:"total"`
}

type basePage struct {
	// A link to the Web API Endpoint returning the full
	// result of this request.
	Endpoint string `json:"href"`
	// The maximum number of items in the response, as set
	// in the query (or default value if unset).
	Limit int `json:"limit"`
	// The offset of the items returned, as set in the query
	// (or default value if unset).
	Offset int `json:"offset"`
	// The total number of items available to return.
	Total int `json:"total"`
	// The URL to the next page of items (if available).
	Next string `json:"next"`
	// The URL to the previous page of items (if available).
	Previous string `json:"previous"`
}

// User contains the basic, publicly available information about a Spotify user.
type User struct {
	// The name displayed on the user's profile.
	// Note: Spotify currently fails to populate
	// this field when querying for a playlist.
	DisplayName string `json:"display_name"`
	// Known public external URLs for the user.
	ExternalURLs map[string]string `json:"external_urls"`
	// Information about followers of the user.
	Followers Followers `json:"followers"`
	// A link to the Web API endpoint for this user.
	Endpoint string `json:"href"`
	// The Spotify user ID for the user.
	ID string `json:"id"`
	// The user's profile image.
	Images []Image `json:"images"`
	// The Spotify URI for the user.
	URI URI `json:"uri"`
}

// SimplePlaylistPage contains SimplePlaylists returned by the Web API.
type SimplePlaylistPage struct {
	basePage
	Playlists []SimplePlaylist `json:"items"`
}

// SimplePlaylist contains basic info about a Spotify playlist.
type SimplePlaylist struct {
	// Indicates whether the playlist owner allows others to modify the playlist.
	// Note: only non-collaborative playlists are currently returned by Spotify's Web API.
	Collaborative bool              `json:"collaborative"`
	ExternalURLs  map[string]string `json:"external_urls"`
	// A link to the Web API endpoint providing full details of the playlist.
	Endpoint string `json:"href"`
	ID       ID     `json:"id"`
	// The playlist image.  Note: this field is only  returned for modified,
	// verified playlists. Otherwise the slice is empty.  If returned, the source
	// URL for the image is temporary and will expire in less than a day.
	Images   []Image `json:"images"`
	Name     string  `json:"name"`
	Owner    User    `json:"owner"`
	IsPublic bool    `json:"public"`
	// The version identifier for the current playlist. Can be supplied in other
	// requests to target a specific playlist version.
	SnapshotID string `json:"snapshot_id"`
	// A collection to the Web API endpoint where full details of the playlist's
	// tracks can be retrieved, along with the total number of tracks in the playlist.
	Tracks PlaylistTracks `json:"tracks"`
	URI    URI            `json:"uri"`
}

// PlaylistTrackPage contains information about tracks in a playlist.
type PlaylistTrackPage struct {
	basePage
	Tracks []PlaylistTrack `json:"items"`
}

// SimpleAlbum contains basic data about an album.
type SimpleAlbum struct {
	// The name of the album.
	Name string `json:"name"`
	// A slice of SimpleArtists
	Artists []SimpleArtist `json:"artists"`
	// The type of the album: one of "album",
	// "single", or "compilation".
	AlbumType string `json:"album_type"`
	// The SpotifyID for the album.
	ID ID `json:"id"`
	// The SpotifyURI for the album.
	URI URI `json:"uri"`
	// The markets in which the album is available,
	// identified using ISO 3166-1 alpha-2 country
	// codes.  Note that al album is considered
	// available in a market when at least 1 of its
	// tracks is available in that market.
	AvailableMarkets []string `json:"available_markets"`
	// A link to the Web API enpoint providing full
	// details of the album.
	Endpoint string `json:"href"`
	// The cover art for the album in various sizes,
	// widest first.
	Images []Image `json:"images"`
	// Known external URLs for this album.
	ExternalURLs map[string]string `json:"external_urls"`
}

// PlaylistTrack contains info about a track in a playlist.
type PlaylistTrack struct {
	// The date and time the track was added to the playlist.
	// You can use the TimestampLayout constant to convert
	// this field to a time.Time value.
	// Warning: very old playlists may not populate this value.
	AddedAt string `json:"added_at"`
	// The Spotify user who added the track to the playlist.
	// Warning: vary old playlists may not populate this value.
	AddedBy User `json:"added_by"`
	// Information about the track.
	Track FullTrack `json:"track"`
}

// FullTrack provides extra track data in addition to what is provided by SimpleTrack.
type FullTrack struct {
	SimpleTrack
	// The album on which the track appears. The album object includes a link in href to full information about the album.
	Album SimpleAlbum `json:"album"`
	// Known external IDs for the track.
	ExternalIDs map[string]string `json:"external_ids"`
	// Popularity of the track.  The value will be between 0 and 100,
	// with 100 being the most popular.  The popularity is calculated from
	// both total plays and most recent plays.
	Popularity int `json:"popularity"`
}

// FullPlaylist provides extra playlist data in addition to the data provided by SimplePlaylist.
type FullPlaylist struct {
	SimplePlaylist
	// The playlist description.  Only returned for modified, verified playlists.
	Description string `json:"description"`
	// Information about the followers of this playlist.
	Followers Followers         `json:"followers"`
	Tracks    PlaylistTrackPage `json:"tracks"`
}

// PlaylistOptions contains optional parameters that can be used when querying
// for featured playlists.  Only the non-nil fields are used in the request.
type PlaylistOptions struct {
	Options
	// The desired language, consisting of a lowercase IO 639
	// language code and an uppercase ISO 3166-1 alpha-2
	// country code, joined by an underscore.  Provide this
	// parameter if you want the results returned in a particular
	// language.  If not specified, the result will be returned
	// in the Spotify default language (American English).
	Locale *string
	// A timestamp in ISO 8601 format (yyyy-MM-ddTHH:mm:ss).
	// use this paramter to specify the user's local time to
	// get results tailored for that specific date and time
	// in the day.  If not provided, the response defaults to
	// the current UTC time.
	Timestamp *string
}

// FeaturedPlaylistsOpt gets a list of playlists featured by Spotify.
// It accepts a number of optional parameters via the opt argument.
func FeaturedPlaylistsOpt(token string, opt *PlaylistOptions) (message string, playlists *SimplePlaylistPage, e error) {
	spotifyURL := baseURL + "browse/featured-playlists"
	if opt != nil {
		v := url.Values{}
		if opt.Locale != nil {
			v.Set("locale", *opt.Locale)
		}
		if opt.Country != "" {
			v.Set("country", opt.Country)
		}
		if opt.Timestamp != nil {
			v.Set("timestamp", *opt.Timestamp)
		}
		if opt.Limit != 0 {
			v.Set("limit", strconv.FormatInt(int64(opt.Limit), 10))
		}
		if opt.Offset != 0 {
			v.Set("offset", strconv.FormatInt(int64(opt.Offset), 10))
		}
		if params := v.Encode(); params != "" {
			spotifyURL += "?" + params
		}
	}
	var result struct {
		Playlists SimplePlaylistPage `json:"playlists"`
		Message   string             `json:"message"`
	}

	err := Get(spotifyURL, token, &result)
	if err != nil {
		return "", nil, err
	}

	return result.Message, &result.Playlists, nil
}

func GetPlaylistTracksOpt(token string, userID string, playlistID ID,
	opt *Options, fields string) (*PlaylistTrackPage, error) {

	spotifyURL := fmt.Sprintf("%susers/%s/playlists/%s/tracks", baseURL, userID, playlistID)
	v := url.Values{}
	if fields != "" {
		v.Set("fields", fields)
	}
	if opt != nil {
		if opt.Limit != 0 {
			v.Set("limit", strconv.FormatInt(int64(opt.Limit), 10))
		}
		if opt.Offset != 0 {
			v.Set("offset", strconv.FormatInt(int64(opt.Offset), 10))
		}
	}
	if params := v.Encode(); params != "" {
		spotifyURL += "?" + params
	}

	var result PlaylistTrackPage

	err := Get(spotifyURL, token, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}
