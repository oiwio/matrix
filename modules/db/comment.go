package db

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Comment struct {
		CommentId  bson.ObjectId `json:"commentId" bson:"_id"`
		FeedId     bson.ObjectId `json:"feedId,omitempty" bson:"feedId,omitempty"`
		Author     *CommentUser  `json:"author,omitempty" bson:"author,omitempty"`
		Reference  *CommentUser  `json:"reference,omitempty" bson:"reference,omitempty"`
		Content    string        `json:"content,omitempty" bson:"content,omitempty"`
		LikedCount int64         `json:"likedCount,omitempty" bson:"likedCount,omitempty"`
		CreateAt   int64         `json:"createAt,omitempty" bson:"createAt,omitempty"`
	}

	CommentUser struct {
		UserId   bson.ObjectId `json:"userId,omitempty" bson:"userId,omitempty"`
		Avatar   string        `json:"avatar,omitempty" bson:"avatar,omitempty"`
		NickName string        `json:"nickname,omitempty" bson:"nickname,omitempty"`
	}
)

func NewComment(s *mgo.Session, comment *Comment) (*Comment, error) {
	var (
		err error
	)

	comment.CommentId = bson.NewObjectId()
	collection := Collection(s, comment)
	err = collection.Insert(comment)
	if err != nil {
		return nil, err
	}
	return comment, err
}

func GetCommentsByFeedId(s *mgo.Session, feedId bson.ObjectId, timestamp int64) ([]*Comment, error) {
	var (
		err      error
		comment  *Comment
		comments []*Comment
	)
	comments = make([]*Comment, 15)
	fmt.Println(feedId)
	err = Collection(s, comment).Find(
		bson.M{
			"feedId": feedId,
			"createAt": bson.M{
				"$lt": timestamp,
			},
		}).Limit(15).Sort("-createAt").All(&comments)
	return comments, err
}

func GetCommentById(s *mgo.Session, commentId bson.ObjectId) (*Comment, error) {
	var (
		err error
	)
	comment := new(Comment)
	err = Collection(s, comment).FindId(commentId).One(comment)
	return comment, err
}

func DeleteComment(s *mgo.Session, commentId bson.ObjectId) error {
	var (
		err     error
		comment *Comment
	)
	err = Collection(s, comment).Remove(bson.M{"_id": commentId})
	return err
}

func DeleteCommentByFeedId(s *mgo.Session, feedId bson.ObjectId) error {
	var (
		err     error
		comment *Comment
	)
	_, err = Collection(s, comment).RemoveAll(bson.M{"feedId": feedId})
	return err
}

func GetCommentUser(s *mgo.Session, userId bson.ObjectId) (*CommentUser, error) {
	var (
		err error
	)
	user := new(User)
	commentUser := new(CommentUser)
	err = Collection(s, user).FindId(userId).One(user)
	if err != nil {
		return nil, err
	}
	commentUser.UserId = userId
	commentUser.Avatar = user.Avatar
	commentUser.NickName = user.Nickname
	return commentUser, err
}
