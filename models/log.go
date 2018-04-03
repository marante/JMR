package models

import "gopkg.in/mgo.v2/bson"

type ActivityData struct {
	PlayedTracks  []string `json:"playedTracks,omitempty" bson:"playedTracks,omitempty"`
	ActivityCount int      `json:"activityCount,omitempty" bson:"activityCount,omitempty"`
	SkipCount     int      `json:"skipCount,omitempty" bson:"skipCount,omitempty"`
}

type PlayerLog struct {
	ID         bson.ObjectId `json:"_id,omitempty"bson:"_id,omitempty"`
	UserID     string        `json:"userId" bson:"userId"`
	Activities struct {
		InVehicle struct {
			ActivityData `json:"activityData,omitempty" bson:"activityData,omitempty"`
		} `json:"in_vehicle,omitempty" bson:"in_vehicle,omitempty"`
		InBicycle struct {
			ActivityData `json:"activityData,omitempty" bson:"activityData,omitempty"`
		} `json:"on_bicycle,omitempty" bson:"bicycle,omitempty"`
		OnFoot struct {
			ActivityData `json:"activityData,omitempty" bson:"activityData,omitempty"`
		} `json:"on_foot,omitempty" bson:"on_foot,omitempty"`
		Walking struct {
			ActivityData `json:"activityData,omitempty" bson:"activityData,omitempty"`
		} `json:"walking,omitempty" bson:"walking,omitempty"`
		Running struct {
			ActivityData `json:"activityData,omitempty" bson:"activityData,omitempty"`
		} `json:"running,omitempty" bson:"running,omitempty"`
		Still struct {
			ActivityData `json:"activityData,omitempty" bson:"activityData,omitempty"`
		} `json:"still,omitempty" bson:"still,omitempty"`
		Unknown struct {
			ActivityData `json:"activityData,omitempty" bson:"activityData,omitempty"`
		} `json:"unknown,omitempty" bson:"unknown,omitempty"`
	} `json:"activities"`
	TotalSkipCount int    `json:"totalSkipCount,omitempty" bson:"totalSkipCount,omitempty"`
	TimeOfDay      string `json:"timeOfDay,omitempty" bson:"timeOfDay,omitempty"`
	ElapsedTime    int    `json:"elapsedTime,omitempty" bson:"elapsedTime,omitempty"`
	//Tracks         []string `json:"tracks,omitempty" bson:"tracks,omitempty"`
}
