package db

import "gopkg.in/mgo.v2/bson"

type Tag struct {
	TagId         bson.ObjectId   `json:"tagId,omitempty" bson:"_id"`
	TagName       string          `json:"feedName,omitempty" bson:"feedName,omitempty"`
	Avatar        string          `json:"avatar,omitempty" bson:"avatar,omitempty"`
	FollowedCount int             `json:"followedCount,omitempty" bson:"followedCount,omitempty"`
	Followers     []bson.ObjectId `json:"followers,omitempty" bson:"followers,omitempty"`
	Description   string          `json:"description,omitempty" bson:"description,omitempty"`
}
