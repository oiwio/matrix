package handlers

import (
	"matrix/auth"
	"matrix/producer"
	"net/http"
	"zion/db"
	"zion/event"
	"zion/protocol"

	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2/bson"
)

func FollowTag(w http.ResponseWriter, r *http.Request) {
	var (
		initiatorId   bson.ObjectId
		tagId         bson.ObjectId
		err           error
		response      *db.TagResponse
		tagEvent      *event.TagEvent
		isTagFollowed bool
	)
	response = new(db.TagResponse)
	tagEvent = new(event.TagEvent)

	initiatorId, err = auth.GetTokenFromRequest(r)
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
			JSONResponse(response, w)
		}
	}()

	tagId = bson.ObjectIdHex(mux.Vars(r)["TagId"])

	isTagFollowed = db.IsTagFollowed(session, initiatorId, tagId)

	if isTagFollowed {
		log.Warningln("已经加过关注")
		response.Success = true
		response.IsTagFollowed = isTagFollowed
		JSONResponse(response, w)
		return
	}
	tagEvent.EventId = event.EVENT_TAG_FOLLOW
	tagEvent.InitiatorId = initiatorId
	tagEvent.TagId = tagId
	go producer.PublishJSONAsync("tag", tagEvent, nil)

	//新的关系
	isTagFollowed = db.IsTagFollowed(session, initiatorId, tagId)
	response.IsTagFollowed = isTagFollowed
	response.Success = true
	JSONResponse(response, w)
}

func UnFollowTag(w http.ResponseWriter, r *http.Request) {
	var (
		initiatorId   bson.ObjectId
		tagId         bson.ObjectId
		err           error
		response      *db.TagResponse
		tagEvent      *event.TagEvent
		isTagFollowed bool
	)
	response = new(db.TagResponse)
	tagEvent = new(event.TagEvent)

	initiatorId, err = auth.GetTokenFromRequest(r)
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
			JSONResponse(response, w)
		}
	}()

	tagId = bson.ObjectIdHex(mux.Vars(r)["TagId"])

	isTagFollowed = db.IsTagFollowed(session, initiatorId, tagId)

	if !isTagFollowed {
		log.Warningln("已经取消关注")
		response.Success = true
		response.IsTagFollowed = isTagFollowed
		JSONResponse(response, w)
		return
	}
	tagEvent.EventId = event.EVENT_TAG_UNFOLLOW
	tagEvent.InitiatorId = initiatorId
	tagEvent.TagId = tagId
	go producer.PublishJSONAsync("tag", tagEvent, nil)

	//新的关系
	isTagFollowed = db.IsTagFollowed(session, initiatorId, tagId)

	response.IsTagFollowed = isTagFollowed
	response.Success = true
	JSONResponse(response, w)
}
