package utils

import (
	"github.com/marante/JMR/Spotify"
	"math"
)

var (
	tempDance, tempEnergy, tempSpeech, tempLoud, tempAcou, tempInstru, tempLive, tempVal, tempTempo [2]float64
)

func GetTrackAttributes(user *Spotify.UserInfo) (*Spotify.TrackAttributes, error) {
	attr, err := Spotify.GetAudioFeatures(user.Context.AnalyzeTracks, user.Token)
	if err != nil {
		return nil, err
	}
	// lowest values
	for _, val := range attr {
		tempDance[0] = math.Min(tempDance[0], val.Danceability)
		tempEnergy[0] = math.Min(tempEnergy[0], val.Energy)
		tempSpeech[0] = math.Min(tempSpeech[0], val.Speechiness)
		tempLoud[0] = math.Min(tempLoud[0], val.Loudness)
		tempAcou[0] = math.Min(tempAcou[0], val.Acousticness)
		tempInstru[0] = math.Min(tempInstru[0], val.Instrumentalness)
		tempLive[0] = math.Min(tempLive[0], val.Liveness)
		tempVal[0] = math.Min(tempVal[0], val.Valence)
		tempTempo[0] = math.Min(tempTempo[0], val.Tempo)
	}
	// highest values
	for _, val := range attr {
		tempDance[1] = math.Max(tempDance[1], val.Danceability)
		tempEnergy[1] = math.Max(tempEnergy[1], val.Energy)
		tempSpeech[1] = math.Max(tempSpeech[1], val.Speechiness)
		tempLoud[1] = math.Max(tempLoud[1], val.Loudness)
		tempAcou[1] = math.Max(tempAcou[1], val.Acousticness)
		tempInstru[1] = math.Max(tempInstru[1], val.Instrumentalness)
		tempLive[1] = math.Max(tempLive[1], val.Liveness)
		tempVal[1] = math.Max(tempVal[1], val.Valence)
		tempTempo[1] = math.Max(tempTempo[1], val.Tempo)
	}

	return Spotify.
		NewTrackAttributes().
		MinDanceability(tempDance[0]).
		MaxDanceability(tempDance[1]).
		MinEnergy(tempEnergy[0]).
		MaxEnergy(tempEnergy[1]).
		MinSpeechiness(tempSpeech[0]).
		MaxSpeechiness(tempSpeech[1]).
		MinLoudness(tempLoud[0]).
		MaxLoudness(tempLoud[1]).
		MinAcousticness(tempAcou[0]).
		MaxAcousticness(tempAcou[1]).
		MinInstrumentalness(tempInstru[0]).
		MaxInstrumentalness(tempInstru[1]).
		MinLiveness(tempLive[0]).
		MaxLiveness(tempLive[1]).
		MinValence(tempVal[0]).
		MaxValence(tempVal[1]).
		MinTempo(tempTempo[0]).
		MaxTempo(tempTempo[1]), nil
}
