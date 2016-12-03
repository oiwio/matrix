package db

import "gopkg.in/mgo.v2/bson"

type Comment struct {
	UserId          bson.ObjectId `json:"userId,omitempty" bson:"userId,omitempty"`
	NickName        string        `json:"nickname,omitempty" bson:"nickname,omitempty"`
	ReplyToUserId   bson.ObjectId `json:"replyToUserId,omitempty" bson:"replyToUserId,omitempty"`
	ReplyToUserName string        `json:"replyToUserName,omitempty" bson:"replyToUserName,omitempty"`
	Content         string        `json:"content,omitempty" bson:"content,omitempty"`
	CreateAt        int64         `json:"createAt" bson:"createAt,omitempty"`
}
