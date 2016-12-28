package handlers

import (
	"encoding/json"
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
func PostFeed(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var (
		err      error
		response *db.FeedResponse
		feed     *db.Feed
	)

	response = new(db.FeedResponse)
	err = json.NewDecoder(r.Body).Decode(&feed)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INVALID_REQUEST
		json.NewEncoder(w).Encode(response)
		return
	}
	feed.FeedId = bson.NewObjectId()
	// feed.UserId =
	feed.CreateAt = time.Now().Unix()
	feed.UpdateAt = feed.CreateAt

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	_, err = db.NewFeed(session, feed)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_FEED_CANTCREATE
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
func GetFeedById(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var (
		err      error
		feed     *db.Feed
		feedId   bson.ObjectId
		response *db.FeedResponse
	)
	response = new(db.FeedResponse)

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	//处理ObjectIdHex在接收错误id之后抛出的异常
	defer func() {
		if e := recover(); e != nil {
			response.Success = false
			response.Error = protocol.ERROR_FEED_CANTGET
			json.NewEncoder(w).Encode(response)
		}
	}()

	feedId = bson.ObjectIdHex(mux.Vars(r)["FeedId"])

	feed, err = db.GetFeedById(session, feedId)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_FEED_CANTGET
		JSONResponse(response, w)
		return
	}
	response.Success = true
	response.Feed = feed
	JSONResponse(response, w)
}
