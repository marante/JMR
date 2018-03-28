package models

import "gopkg.in/mgo.v2/bson"

type PlayerLog struct {
	ID            bson.ObjectId `json:"_id,omitempty"bson:"_id,omitempty"`
	UserID        string        `json:"userId" bson:"userId"`
	ActivityCount struct {
		InVehicle int `json:"in_vehicle,omitempty" bson:"in_vehicle,omitempty"`
		InBicycle int `json:"on_bicycle,omitempty" bson:"bicycle,omitempty"`
		OnFoot    int `json:"on_foot,omitempty" bson:"on_foot,omitempty"`
		Walking   int `json:"walking,omitempty" bson:"walking,omitempty"`
		Running   int `json:"running,omitempty" bson:"running,omitempty"`
		Still     int `json:"still,omitempty" bson:"still,omitempty"`
		Unknown   int `json:"unknown,omitempty" bson:"unknown,omitempty"`
	} `json:"activityCount"`
	SkipCount   int `json:"skipCount,omitempty" bson:"skipCount,omitempty"`
	ElapsedTime int `json:"elapsedTime,omitempty" bson:"elapsedTime,omitempty"`
}
