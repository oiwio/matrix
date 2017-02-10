package handlers

import (
	"encoding/json"
	"matrix/auth"
	"matrix/producer"
	"net/http"
	"strconv"
	"time"
	"zion/db"
	"zion/event"
	"zion/protocol"
	"zion/tools"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	qiniuconf "qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
)

// PostFeed receive POST methods and store them in MongoDB
func PostFeed(w http.ResponseWriter, r *http.Request) {
	var (
		userId    bson.ObjectId
		err       error
		response  *db.FeedResponse
		feed      *db.Feed
		feedEvent *event.FeedEvent
	)

	response = new(db.FeedResponse)
	feedEvent = new(event.FeedEvent)

	userId, err = auth.GetTokenFromRequest(r)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_NEED_SIGNIN
		JSONResponse(response, w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&feed)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INVALID_REQUEST
		json.NewEncoder(w).Encode(response)
		return
	}

	feed.FeedId = bson.NewObjectId()
	feed.UserId = userId
	feed.CreateAt = time.Now().Unix()
	feed.UpdateAt = feed.CreateAt

	feedEvent.EventId = event.EVENT_FEED_CREATE
	feedEvent.Feed = feed
	go producer.PublishJSONAsync("feed", feedEvent, nil)

	response.Success = true
	response.Feed = feed
	JSONResponse(response, w)
}

// GetMusic reveive GET methods and return the details of music
func GetMusic(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		music    *db.Music
		response *db.MusicResponse
	)
	response = new(db.MusicResponse)
	music, err = tools.GetNeteaseSongList(mux.Vars(r)["MusicId"])
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INVALID_REQUEST
		JSONResponse(response, w)
		return
	}
	response.Success = true
	response.Music = music
	log.Infoln(response)
	JSONResponse(response, w)
}

// SearchMusic reveive GET methods and return the details of music
func SearchMusic(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		musics   []*db.Music
		response *db.MusicResponse
	)
	response = new(db.MusicResponse)
	musics, err = tools.GetSearchList(r.FormValue("s"), r.FormValue("limit"), r.FormValue("offset"))
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INVALID_REQUEST
		JSONResponse(response, w)
		return
	}
	response.Success = true
	response.Musics = musics
	JSONResponse(response, w)
}

// GetFeedById need a feedId and return the detail of feed
func GetFeedById(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		feed     *db.Feed
		comments []*db.Comment
		feedId   bson.ObjectId
		response *db.FeedResponse
	)
	response = new(db.FeedResponse)

	_, err = auth.GetTokenFromRequest(r)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_NEED_SIGNIN
		JSONResponse(response, w)
		return
	}

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	//处理ObjectIdHex在接收错误id之后抛出的异常
	defer func() {
		if e := recover(); e != nil {
			response.Success = false
			response.Error = protocol.ERROR_INVALID_REQUEST
			json.NewEncoder(w).Encode(response)
		}
	}()

	feedId = bson.ObjectIdHex(mux.Vars(r)["FeedId"])

	feed, err = db.GetFeedById(session, feedId)
	comments, err = db.GetCommentsByFeedId(session, feedId, time.Now().Unix())

	feed.Comments = comments
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INTERNAL_ERROR
		JSONResponse(response, w)
		return
	}

	response.Success = true
	response.Feed = feed
	JSONResponse(response, w)
}

// DelFeed need a feedId and delete it
func DelFeed(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		feed      *db.Feed
		feedEvent *event.FeedEvent
		userId    bson.ObjectId
		feedId    bson.ObjectId
		response  *db.FeedResponse
	)
	response = new(db.FeedResponse)
	feedEvent = new(event.FeedEvent)

	userId, err = auth.GetTokenFromRequest(r)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_NEED_SIGNIN
		JSONResponse(response, w)
		return
	}

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	//处理ObjectIdHex在接收错误id之后抛出的异常
	defer func() {
		if e := recover(); e != nil {
			response.Success = false
			response.Error = protocol.ERROR_INVALID_REQUEST
			json.NewEncoder(w).Encode(response)
		}
	}()

	feedId = bson.ObjectIdHex(mux.Vars(r)["FeedId"])

	feed, err = db.GetFeedById(session, feedId)
	if err == nil {
		if feed.UserId == userId {
			feedEvent.EventId = event.EVENT_FEED_REMOVE
			feedEvent.FeedId = feed.FeedId
			// push message to nsq
			go producer.PublishJSONAsync("feed", feedEvent, nil)
		} else {
			response.Success = false
			response.Error = protocol.ERROR_AUTH
			JSONResponse(response, w)
			return
		}
	} else {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INVALID_REQUEST
		JSONResponse(response, w)
		return
	}

	response.Success = true
	response.Feed = feed
	JSONResponse(response, w)
}

// GetFeedsByUserId return the personal feeds
func GetFeedsByUserId(w http.ResponseWriter, r *http.Request) {

	var (
		err       error
		timestamp int64
		userId    bson.ObjectId
		feeds     []*db.Feed
		response  *db.FeedResponse
	)

	response = new(db.FeedResponse)

	_, err = auth.GetTokenFromRequest(r)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_NEED_SIGNIN
		JSONResponse(response, w)
		return
	}

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	//处理ObjectIdHex在接收错误id之后抛出的异常
	defer func() {
		if e := recover(); e != nil {
			response.Success = false
			response.Error = protocol.ERROR_INVALID_REQUEST
			json.NewEncoder(w).Encode(response)
		}
	}()

	userId = bson.ObjectIdHex(r.FormValue("u"))
	timestamp, err = strconv.ParseInt(r.FormValue("t"), 10, 64)
	if err != nil {
		panic(err)
	}

	feeds, err = db.GetFeedsByUserId(session, userId, timestamp)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INTERNAL_ERROR
		JSONResponse(response, w)
		return
	}
	response.Success = true
	response.Feeds = feeds
	JSONResponse(response, w)
}

// GetNewestFeeds return the neweast feeds
func GetNewestFeeds(w http.ResponseWriter, r *http.Request) {

	var (
		timestamp int64
		err       error
		feeds     []*db.Feed
		response  *db.FeedResponse
	)

	response = new(db.FeedResponse)

	_, err = auth.GetTokenFromRequest(r)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_NEED_SIGNIN
		JSONResponse(response, w)
		return
	}

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	//处理接收错误的时间戳抛出的异常
	defer func() {
		if e := recover(); e != nil {
			response.Success = false
			response.Error = protocol.ERROR_INVALID_REQUEST
			json.NewEncoder(w).Encode(response)
		}
	}()
	timestamp, err = strconv.ParseInt(r.FormValue("t"), 10, 64)
	if err != nil {
		panic(err)
	}

	feeds, err = db.GetNewestFeeds(session, timestamp)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INTERNAL_ERROR
		JSONResponse(response, w)
		return
	}
	response.Success = true
	response.Feeds = feeds
	JSONResponse(response, w)
}

// GetFeedById need a feedId and return the detail of feed
func GetUploadToken(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		response *db.FeedResponse
		bucket   string
	)
	response = new(db.FeedResponse)

	_, err = auth.GetTokenFromRequest(r)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_NEED_SIGNIN
		JSONResponse(response, w)
		return
	}

	qiniuconf.ACCESS_KEY = configuration.Qiniu.AccessKey
	qiniuconf.SECRET_KEY = configuration.Qiniu.SecretKey

	switch mux.Vars(r)["Type"] {
	case "image":
		bucket = "aladdin-image"
	case "video":
		bucket = "aladdin-video"
	default:
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INVALID_REQUEST
		JSONResponse(response, w)
		return
	}

	c := kodo.New(0, nil)

	// 设置上传的策略
	policy := &kodo.PutPolicy{
		Scope: bucket,
		//设置Token过期时间
		Expires: 3600,
	}
	// 生成一个上传token
	token := c.MakeUptoken(policy)

	response.Success = true
	response.Token = token
	JSONResponse(response, w)
}
