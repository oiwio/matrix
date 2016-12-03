package controllers

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
func PostFeed(w http.ResponseWriter, r *http.Request) {
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
		response.Error = protocol.ERROR_INVALID_REQUEST
		json.NewEncoder(w).Encode(response)
		return
	}
	response.Success = true
	response.Feed = feed
	json.NewEncoder(w).Encode(response)
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
		json.NewEncoder(w).Encode(response)
		return
	}
	response.Success = true
	response.Music = music
	log.Infoln(response)
	json.NewEncoder(w).Encode(response)
}

