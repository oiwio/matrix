package handlers

import (
	"encoding/json"
	"matrix/auth"
	"matrix/modules/db"
	"matrix/modules/protocol"
	"matrix/modules/tools"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// PostFeed receive POST methods and store them in MongoDB
func PostFeed(w http.ResponseWriter, r *http.Request) {
	var (
		userId   bson.ObjectId
		err      error
		response *db.FeedResponse
		feed     *db.Feed
	)

	response = new(db.FeedResponse)
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

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	_, err = db.NewFeed(session, feed)
	if err != nil {
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

// GetFeedById need a feedId and return the detail of feed
func GetFeedById(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		feed     *db.Feed
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
		err      error
		feed     *db.Feed
		userId   bson.ObjectId
		feedId   bson.ObjectId
		response *db.FeedResponse
	)
	response = new(db.FeedResponse)

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
			err = db.DeleteFeed(session, feedId)
			if err != nil {
				HandleError(err)
				response.Success = false
				response.Error = protocol.ERROR_INTERNAL_ERROR
				JSONResponse(response, w)
				return
			}
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
	type RequestBody struct {
		UserId    string `json:"userId"`
		Timestamp int64  `json:"timestamp"`
	}

	var (
		rb       *RequestBody
		err      error
		userId   bson.ObjectId
		feeds    []*db.Feed
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

	err = json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INVALID_REQUEST
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
	userId = bson.ObjectIdHex(rb.UserId)

	feeds, err = db.GetFeedsByUserId(session, userId, rb.Timestamp)
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
