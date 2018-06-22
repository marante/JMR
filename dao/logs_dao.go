package dao

import (
	"log"

	. "github.com/marante/JMR/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// LogsDAO represents a struct capable of connecting to the DB and correct server
type LogsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	collection = "logs"
)

// Connect attempts to connect to the DB with the given credentials from the .toml config file.
func (l *LogsDAO) Connect() {
	session, err := mgo.Dial(l.Server)
	if err != nil {
		log.Fatal(err)
	}
	session.SetSafe(&mgo.Safe{})
	db = session.DB(l.Database)
}

// FindAll finds all retrieves all documents from the DB.
func (l *LogsDAO) FindAll() ([]PlayerLog, error) {
	var logs []PlayerLog
	err := db.C(collection).Find(bson.M{}).All(&logs)
	return logs, err
}

// FindByID finds document based on ID
func (l *LogsDAO) FindByID(id string) (PlayerLog, error) {
	var log PlayerLog
	err := db.C(collection).Find(bson.M{"userId": id}).One(&log)
	return log, err
}

// Insert inserts a document into the DB
func (l *LogsDAO) Insert(log PlayerLog) error {
	err := db.C(collection).Insert(&log)
	return err
}

// Delete delete a specifc logentry from the DB
func (l *LogsDAO) Delete(log PlayerLog) error {
	err := db.C(collection).Remove(bson.M{"userId": log.UserID})
	return err
}

// Update updates/replaces a logentry in the DB
func (l *LogsDAO) Update(log PlayerLog) error {
	err := db.C(collection).Update(bson.M{"userId": log.UserID}, &log)
	return err
}
