package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Feed struct {
		FeedId   bson.ObjectId `json:"feedId,omitempty" bson:"_id"`
		UserId   bson.ObjectId `json:"userId,omitempty" bson:"userId,omitempty"`
		FeedType int           `json:"feedType,omitempty" bson:"feedType,omitempty"`

		Nickname string `json:"nickname,omitempty" bson:"nickname,omitempty"`
		Avatar   string `json:"avatar,omitempty" bson:"avatar,omitempty"`

		Coordinate []float64  `json:"coordinate,omitempty" bson:"coordinate,omitempty"` //坐标
		Tags       []*FeedTag `json:"tags,omitempty" bson:"tags,omitempty"`
		Text       string     `json:"text,omitempty" bson:"text,omitempty"`
		Photo      *Photo     `json:"photo,omitempty" bson:"photo,omitempty"`
		Video      *Video     `json:"video,omitempty" bson:"video,omitempty"`
		Audio      *Audio     `json:"audio,omitempty" bson:"audio,omitempty"`
		Music      *Music     `json:"music,omitempty" bson:"music,omitempty"`

		LikedCount    int64           `json:"likedCount,omitempty" bson:"likedCount,omitempty"`
		LikedUsers    []bson.ObjectId `json:"likeUsers,omitempty" bson:"likeUsers,omitempty"`
		CommentsCount int64           `json:"commentsCount,omitempty" bson:"commentsCount,omitempty"`
		Comments      []*Comment      `json:"comments,omitempty" bson:"comments,omitempty"`
		ShareCount    int64           `json:"shareCount,omitempty" bson:"shareCount,omitempty"`
		ViewsCount    int64           `json:"viewsCount,omitempty" bson:"viewsCount,omitempty"`

		CreateAt int64 `json:"createAt,omitempty" bson:"createAt,omitempty"`
		UpdateAt int64 `json:"updateAt,omitempty" bson:"updateAt,omitempty"`
	}

	FeedTag struct {
		TagId   bson.ObjectId `json:"tagId,omitempty" bson:"tagId"`
		TagName string        `json:"feedName,omitempty" bson:"feedName,omitempty"`
	}

	Audio struct {
		Url      string `json:"url,omitempty" structs:"url,omitempty"`
		CoverUrl string `json:"coverUrl,omitempty" structs:"coverUrl,omitempty"`
		Source   int    `json:"source,omitempty" structs:"source,omitempty"` // 来源 0:本地音乐，1:app音乐，2:用户录音
		Name     string `json:"name,omitempty" structs:"name,omitempty"`     // 音乐名称
		Volume   int    `json:"volume,omitempty" structs:"volume,omitempty"` // 音量
		Duration int    `json:"duration,omitempty" structs:"duration,omitempty"`
	}

	Music struct {
		Id       int    `json:"id,omitempty"`
		Duration int    `json:"duration,omitempty"`
		Name     string `json:"name,omitempty"`
		Url      string `json:"url,omitempty"`
		Artists  string `json:"artists,omitempty"`
		Source   string `json:"source,omitempty"`
		CoverUrl string `json:"coverUrl,omitempty"`
	}

	Photo struct {
		Url    string  `json:"url,omitempty" structs:"url,omitempty"`
		Width  float64 `json:"width,omitempty" structs:"width,omitempty"`
		Height float64 `json:"height,omitempty" structs:"height,omitempty"`
	}

	Video struct {
		VideoUrl string  `json:"videoUrl,omitempty" structs:"videoUrl,omitempty"`
		CoverUrl string  `json:"coverUrl,omitempty" structs:"coverUrl,omitempty"`
		Duration int     `json:"duration,omitempty" structs:"duration,omitempty"`
		Width    float64 `json:"width,omitempty" structs:"width,omitempty"`
		Height   float64 `json:"height,omitempty" structs:"height,omitempty"`
	}
)

func NewFeed(s *mgo.Session, feed *Feed) (*Feed, error) {
	var (
		err error
	)

	feed.FeedId = bson.NewObjectId()
	collection := Collection(s, feed)
	err = collection.Insert(feed)
	if err != nil {
		return nil, err
	}
	return feed, err
}

//Получить содержимое на основе ID
func GetFeedById(s *mgo.Session, feedId bson.ObjectId) (*Feed, error) {
	var (
		err error
	)
	//Определить содержание
	feed := new(Feed)
	//получать данные из базы данных
	err = Collection(s, feed).FindId(feedId).One(feed)
	//отправляется клиенту
	return feed, err
}

func GetNewestFeeds(s *mgo.Session, page int) ([]*Feed, error) {
	var (
		err   error
		feed  *Feed
		feeds []*Feed
	)
	feeds = make([]*Feed, 15)
	err = Collection(s, feed).Find(nil).Limit(15).Skip(page * 15).All(&feeds)
	return feeds, err
}

func DeleteFeed(s *mgo.Session, feedId bson.ObjectId) error {
	var (
		err  error
		feed *Feed
	)
	err = Collection(s, feed).Remove(bson.M{"_id": feedId})
	return err
}
