package event

import (
	"matrix/modules/db"

	"gopkg.in/mgo.v2/bson"
)

type (
	FeedEvent struct {
		EventId int           `bson:"eventId,omitempty"`
		Comment *db.Comment   `bson:"comment,omitempty"`
		Feed    *db.Feed      `bson:"feed,omitempty"`
		FeedId  bson.ObjectId `bson:"FeedId,omitempty"`
		UserId  bson.ObjectId `bson:"userId,omitempty"`
	}
)
