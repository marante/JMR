package models

import "gopkg.in/mgo.v2/bson"

// PlayerLog represents log information that is saved in DocumentDB.
type PlayerLog struct {
	ID         bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID     string        `json:"userId" bson:"userId"`
	Activities struct {
		InVehicle struct {
			PlayedTracks  []string `json:"playedTracks,omitempty" bson:"playedTracks,omitempty"`
			ActivityCount int      `json:"activityCount,omitempty" bson:"activityCount,omitempty"`
			SkipCount     int      `json:"skipCount,omitempty" bson:"skipCount,omitempty"`
		} `json:"in_vehicle,omitempty" bson:"in_vehicle,omitempty"`
		InBicycle struct {
			PlayedTracks  []string `json:"playedTracks,omitempty" bson:"playedTracks,omitempty"`
			ActivityCount int      `json:"activityCount,omitempty" bson:"activityCount,omitempty"`
			SkipCount     int      `json:"skipCount,omitempty" bson:"skipCount,omitempty"`
		} `json:"on_bicycle,omitempty" bson:"bicycle,omitempty"`
		OnFoot struct {
			PlayedTracks  []string `json:"playedTracks,omitempty" bson:"playedTracks,omitempty"`
			ActivityCount int      `json:"activityCount,omitempty" bson:"activityCount,omitempty"`
			SkipCount     int      `json:"skipCount,omitempty" bson:"skipCount,omitempty"`
		} `json:"on_foot,omitempty" bson:"on_foot,omitempty"`
		Walking struct {
			PlayedTracks  []string `json:"playedTracks,omitempty" bson:"playedTracks,omitempty"`
			ActivityCount int      `json:"activityCount,omitempty" bson:"activityCount,omitempty"`
			SkipCount     int      `json:"skipCount,omitempty" bson:"skipCount,omitempty"`
		} `json:"walking,omitempty" bson:"walking,omitempty"`
		Running struct {
			PlayedTracks  []string `json:"playedTracks,omitempty" bson:"playedTracks,omitempty"`
			ActivityCount int      `json:"activityCount,omitempty" bson:"activityCount,omitempty"`
			SkipCount     int      `json:"skipCount,omitempty" bson:"skipCount,omitempty"`
		} `json:"running,omitempty" bson:"running,omitempty"`
		Still struct {
			PlayedTracks  []string `json:"playedTracks,omitempty" bson:"playedTracks,omitempty"`
			ActivityCount int      `json:"activityCount,omitempty" bson:"activityCount,omitempty"`
			SkipCount     int      `json:"skipCount,omitempty" bson:"skipCount,omitempty"`
		} `json:"still,omitempty" bson:"still,omitempty"`
		Unknown struct {
			PlayedTracks  []string `json:"playedTracks,omitempty" bson:"playedTracks,omitempty"`
			ActivityCount int      `json:"activityCount,omitempty" bson:"activityCount,omitempty"`
			SkipCount     int      `json:"skipCount,omitempty" bson:"skipCount,omitempty"`
		} `json:"unknown,omitempty" bson:"unknown,omitempty"`
	} `json:"activities,omitempty" bson:"activities,omitempty"`
	TotalSkipCount int    `json:"totalSkipCount,omitempty" bson:"totalSkipCount,omitempty"`
	TimeOfDay      string `json:"timeOfDay,omitempty" bson:"timeOfDay,omitempty"`
	ElapsedTime    int    `json:"elapsedTime,omitempty" bson:"elapsedTime,omitempty"`
}
