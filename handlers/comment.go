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

	"errors"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// PostComment receive POST methods and store them in MongoDB
func PostComment(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		FeedId      string `json:"feedId"`
		ReferenceId string `json:"referenceId,omitempty"`
		Conetent    string `json:"content,omitempty"`
	}
	var (
		rb          *RequestBody
		userId      bson.ObjectId
		referenceId bson.ObjectId
		err         error
		response    *db.CommentResponse
		comment     *db.Comment
		feedEvent   *event.FeedEvent
	)

	response = new(db.CommentResponse)
	feedEvent = new(event.FeedEvent)

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

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	//处理ObjectIdHex在接收错误id之后抛出的异常
	defer func() {
		if e := recover(); e != nil {
			HandleError(e.(error))
			response.Success = false
			response.Error = protocol.ERROR_INVALID_REQUEST
			json.NewEncoder(w).Encode(response)
		}
	}()

	comment = new(db.Comment)
	comment.CommentId = bson.NewObjectId()
	comment.FeedId = bson.ObjectIdHex(rb.FeedId)
	if rb.ReferenceId != "" {
		referenceId = bson.ObjectIdHex(rb.ReferenceId)
		comment.Reference, err = db.GetCommentUser(session, referenceId)
	}
	comment.Content = rb.Conetent
	comment.Author, err = db.GetCommentUser(session, userId)
	if comment.Author == nil {
		panic(errors.New("Can't find this user"))
	}
	comment.CreateAt = time.Now().Unix()

	feedEvent.EventId = event.EVENT_FEED_COMMENT_POST
	feedEvent.Comment = comment
	go producer.PublishJSONAsync("feed", feedEvent, nil)

	response.Success = true
	response.Comment = comment
	JSONResponse(response, w)
}

// GetCommentsByFeedId return the feed's comments
func GetCommentsByFeedId(w http.ResponseWriter, r *http.Request) {

	var (
		err       error
		feedId    bson.ObjectId
		timestamp int64
		comments  []*db.Comment
		response  *db.CommentResponse
	)

	response = new(db.CommentResponse)

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
			log.Errorln(e)
			response.Success = false
			response.Error = protocol.ERROR_INVALID_REQUEST
			json.NewEncoder(w).Encode(response)
		}
	}()

	feedId = bson.ObjectIdHex(r.FormValue("f"))
	timestamp, err = strconv.ParseInt(r.FormValue("t"), 10, 64)
	if err != nil {
		panic(err)
	}

	comments, err = db.GetCommentsByFeedId(session, feedId, timestamp)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INTERNAL_ERROR
		JSONResponse(response, w)
		return
	}
	response.Success = true
	response.Comments = comments
	JSONResponse(response, w)
}

// DelComment need a commentId and delete it
func DelComment(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		comment   *db.Comment
		feed      *db.Feed
		userId    bson.ObjectId
		commentId bson.ObjectId
		response  *db.Response
		feedEvent *event.FeedEvent
	)
	response = new(db.Response)
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

	commentId = bson.ObjectIdHex(mux.Vars(r)["CommentId"])
	comment, err = db.GetCommentById(session, commentId)
	feed, err = db.GetFeedById(session, comment.FeedId)
	if err == nil {
		if comment.Author.UserId == userId || feed.UserId == userId {
			feedEvent.EventId = event.EVENT_FEED_COMMENT_REMOVE
			feedEvent.CommentId = commentId
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
	JSONResponse(response, w)
}
