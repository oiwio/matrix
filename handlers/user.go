package handlers

import (
	"encoding/json"
	"matrix/auth"
	"net/http"
	"zion/db"
	"zion/protocol"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		user     *db.User
		userId   bson.ObjectId
		response *db.UserResponse
	)
	response = new(db.UserResponse)

	if mux.Vars(r)["UserId"] != "" {
		_, err = auth.GetTokenFromRequest(r)
		if err != nil {
			HandleError(err)
			response.Success = false
			response.Error = protocol.ERROR_NEED_SIGNIN
			JSONResponse(response, w)
			return
		}

		//处理ObjectIdHex在接收错误id之后抛出的异常
		defer func() {
			if e := recover(); e != nil {
				response.Success = false
				response.Error = protocol.ERROR_INVALID_USER
				json.NewEncoder(w).Encode(response)
			}
		}()
		userId = bson.ObjectIdHex(mux.Vars(r)["UserId"])

	} else {
		userId, err = auth.GetTokenFromRequest(r)
		if err != nil {
			HandleError(err)
			response.Success = false
			response.Error = protocol.ERROR_NEED_SIGNIN
			JSONResponse(response, w)
			return
		}
	}

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()
	user, err = db.GetUserById(session, userId)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INVALID_USER
		JSONResponse(response, w)
		return
	}
	response.Success = true
	response.User = &db.User{
		Avatar:    user.Avatar,
		Nickname:  user.Nickname,
		Username:  user.Username,
		Gender:    user.Gender,
		Signature: user.Signature,
		Follower:  user.Follower,
		Following: user.Following,
	}
	JSONResponse(response, w)
}
