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

	"gopkg.in/mgo.v2/bson"
)

// PostTag receive POST methods and store them in MongoDB
func PostTag(w http.ResponseWriter, r *http.Request) {

	var (
		userId   bson.ObjectId
		err      error
		response *db.TagResponse
		tag      *db.Tag
		tagEvent *event.TagEvent
	)

	response = new(db.TagResponse)
	tagEvent = new(event.TagEvent)

	userId, err = auth.GetTokenFromRequest(r)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_NEED_SIGNIN
		JSONResponse(response, w)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&tag)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INVALID_REQUEST
		json.NewEncoder(w).Encode(response)
		return
	}

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	_, err = db.IsTagExist(session, tag.Name)
	if err == nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_TAG_AKREADY_CREATED
		json.NewEncoder(w).Encode(response)
		return
	}

	tag.TagId = bson.NewObjectId()
	tag.CreateUser = userId
	tag.CreateAt = time.Now().Unix()
	tag.UpdateAt = tag.CreateAt

	tagEvent.EventId = event.EVENT_TAG_CREATE
	tagEvent.Tag = tag

	go producer.PublishJSONAsync("tag", tagEvent, nil)

	response.Success = true
	response.Tag = tag
	JSONResponse(response, w)
}

// AddTags receive POST methods and store them in MongoDB
func AddTags(w http.ResponseWriter, r *http.Request) {

	type RequestBody struct {
		FeedId string   `json:"feedId"`
		TagIds []string `json:"tags"`
	}
	var (
		rb       *RequestBody
		userId   bson.ObjectId
		feedId   bson.ObjectId
		err      error
		response *db.Response
		tagEvent *event.TagEvent
	)

	response = new(db.Response)
	tagEvent = new(event.TagEvent)

	userId, err = auth.GetTokenFromRequest(r)
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
		json.NewEncoder(w).Encode(response)
		return
	}

	//处理ObjectIdHex在接收错误id之后抛出的异常
	defer func() {
		if e := recover(); e != nil {
			response.Success = false
			response.Error = protocol.ERROR_INVALID_REQUEST
			json.NewEncoder(w).Encode(response)
		}
	}()

	feedId = bson.ObjectIdHex(rb.FeedId)

	tagEvent.EventId = event.EVENT_ADD_TAGS
	tagEvent.TagIds = rb.TagIds
	tagEvent.AddUser = userId
	tagEvent.FeedId = feedId

	go producer.PublishJSONAsync("tag", tagEvent, nil)

	response.Success = true
	JSONResponse(response, w)
}

func GetFeedsByTagId(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		timestamp int64
		tagId     bson.ObjectId
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

	tagId = bson.ObjectIdHex(r.FormValue("id"))
	timestamp, err = strconv.ParseInt(r.FormValue("timestamp"), 10, 64)
	if err != nil {
		panic(err)
	}

	feeds, err = db.GetFeedsByTagId(session, tagId, timestamp)
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

func FuzzySearchByTagName(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		tagName  string
		tags     []*db.Tag
		response *db.TagResponse
	)

	response = new(db.TagResponse)

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

	tagName = r.FormValue("name")

	tags, err = db.FuzzySearchByTagName(session, tagName)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INTERNAL_ERROR
		JSONResponse(response, w)
		return
	}
	response.Success = true
	response.Tags = tags
	JSONResponse(response, w)
}
