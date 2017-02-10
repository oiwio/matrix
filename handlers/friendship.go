package handlers

import (
	"matrix/auth"
	"matrix/producer"
	"net/http"
	"zion/db"
	"zion/event"
	"zion/protocol"

	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func FollowUser(w http.ResponseWriter, r *http.Request) {
	var (
		initiatorId bson.ObjectId
		responderId bson.ObjectId
		err         error
		response    *db.FriendshipResponse
		friendEvent *event.FriendEvent
		relation    int
	)
	response = new(db.FriendshipResponse)
	friendEvent = new(event.FriendEvent)

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

	responderId = bson.ObjectIdHex(mux.Vars(r)["ResponderId"])

	relation = db.GetRelation(session, initiatorId, responderId)

	if db.IsFriendshipExist(session, initiatorId, responderId) {
		log.Warningln("已经加过关注")
		response.Success = true
		response.Relation = relation
		JSONResponse(response, w)
		return
	}
	friendEvent.EventId = event.EVENT_FRIEND_FOLLOW
	friendEvent.InitiatorId = initiatorId
	friendEvent.ResponderId = responderId
	go producer.PublishJSONAsync("friend", friendEvent, nil)

	//新的关系
	relation = db.GetRelation(session, initiatorId, responderId)
	//延迟50毫秒返回
	time.After(time.Millisecond * 50)
	response.Relation = relation
	response.Success = true
	JSONResponse(response, w)
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	var (
		initiatorId bson.ObjectId
		responderId bson.ObjectId
		err         error
		response    *db.FriendshipResponse
		friendEvent *event.FriendEvent
		relation    int
	)
	response = new(db.FriendshipResponse)
	friendEvent = new(event.FriendEvent)

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

	responderId = bson.ObjectIdHex(mux.Vars(r)["ResponderId"])

	relation = db.GetRelation(session, initiatorId, responderId)

	if !db.IsFriendshipExist(session, initiatorId, responderId) {
		log.Warningln("已经取消关注")
		response.Success = true
		response.Relation = relation
		JSONResponse(response, w)
		return
	}
	friendEvent.EventId = event.EVENT_FRIEND_UNFOLLOW
	friendEvent.InitiatorId = initiatorId
	friendEvent.ResponderId = responderId
	go producer.PublishJSONAsync("friend", friendEvent, nil)

	//新的关系
	relation = db.GetRelation(session, initiatorId, responderId)
	//延迟50毫秒返回
	time.After(time.Millisecond * 50)
	response.Relation = relation
	response.Success = true
	JSONResponse(response, w)
}
