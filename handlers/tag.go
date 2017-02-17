package handlers

import (
	"encoding/json"
	"matrix/auth"
	"matrix/producer"
	"net/http"
	"time"
	"zion/db"
	"zion/event"
	"zion/protocol"

	"gopkg.in/mgo.v2/bson"
)

// PostComment receive POST methods and store them in MongoDB
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
