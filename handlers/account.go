package handlers

import (
	"encoding/json"
	"matrix/modules/db"
	"modules/protocol"
	"net/http"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {

	type UserRegister struct {
		Nickname    string `json:"nickname"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		Gender      int    `json:"gender"`
		DeviceToken string `json:"deviceToken"`
		Account     string `json:"account"`
		Avatar      string `json:"avatar"`
	}

	var (
		err      error
		register *UserRegister
		response *db.UserResponse
		user     *db.User
		t        string
	)

	response = new(db.UserResponse)
	err = json.NewDecoder(r.Body).Decode(&register)
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INVALID_REQUEST
		JSONResponse(response, w)
		return
	}
	log.Infoln("来用户了！！")
	user = new(db.User)
	user.Nickname = register.Nickname
	user.Username = register.Username
	user.Gender = register.Gender
	user.Avatar = register.Avatar
	user.DeviceToken = register.DeviceToken

	//Create session for every request
	session := mgoSession.Copy()
	defer session.Close()

	if isValidPhone(register.Account) {
		//手机注册
		_, err = db.GetUserByPhone(session, register.Account)
		if err != nil {
			user.Phone = register.Account
			user, err = db.NewUser(session, user, register.Password)
			if err != nil {
				HandleError(err)
				response.Success = false
				response.Error = protocol.ERROR_CANNOT_REGISTRY
				JSONResponse(response, w)
				return
			}
			response.Success = true
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["userId"] = user.UserId
			claims["exp"] = time.Now().Add(time.Hour * 2160)
			t, err = token.SignedString([]byte("kiwee19920306"))
			if err != nil {
				HandleError(err)
				response.Success = false
				response.Error = protocol.ERROR_INTERNAL_ERROR
				JSONResponse(response, w)
				return
			}
			response.Token = t
			JSONResponse(response, w)

		} else {
			response.Success = false
			response.Error = protocol.ERROR_PHONE_ALREADY_REGISTRIED
			JSONResponse(response, w)
			return
		}
	} else if isValidEmail(register.Account) {
		//邮箱注册
		_, err = db.GetUserByEmail(session, register.Account)
		if err != nil {
			user.Email = register.Account
			user, err = db.NewUser(session, user, register.Password)
			if err != nil {
				HandleError(err)
				response.Success = false
				response.Error = protocol.ERROR_CANNOT_REGISTRY
				JSONResponse(response, w)
				return
			}
			response.Success = true
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["userId"] = user.UserId
			claims["exp"] = time.Now().Add(time.Hour * 2160)
			t, err = token.SignedString([]byte("kiwee19920306"))
			if err != nil {
				HandleError(err)
				response.Success = false
				response.Error = protocol.ERROR_INTERNAL_ERROR
				JSONResponse(response, w)
				return
			}
			response.Token = t
			JSONResponse(response, w)

		} else {
			response.Success = false
			response.Error = protocol.ERROR_EMAIL_ALREADY_REGISTRIED
			JSONResponse(response, w)
			return
		}
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	type SignIn struct {
		Account     string `json:"account"`
		Password    string `json:"password"`
		DeviceToken string `json:"deviceToken"`
	}
	var (
		err      error
		response *db.UserResponse
		user     *db.User
		signIn   *SignIn
		t        string
	)
	response = new(db.UserResponse)
	err = json.NewDecoder(r.Body).Decode(&signIn)
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

	if isValidPhone(signIn.Account) {
		user, err = db.GetUserByPhone(session, signIn.Account)
		HandleError(err)
	} else if isValidEmail(signIn.Account) {
		user, err = db.GetUserByEmail(session, signIn.Account)
		HandleError(err)
	} else {
		response.Success = false
		response.Error = protocol.ERROR_INVALID_ACCOUNT
		JSONResponse(response, w)
		return
	}
	if user == nil {
		response.Success = false
		response.Error = protocol.ERROR_INVALID_USER
		JSONResponse(response, w)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(signIn.Password))
	if err != nil {
		response.Success = false
		response.Error = protocol.ERROR_PASSWORD_NOT_MATCH
		JSONResponse(response, w)
		return
	}
	response.Success = true
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.UserId
	claims["exp"] = time.Now().Add(time.Hour * 2160)
	t, err = token.SignedString([]byte("kiwee19920306"))
	if err != nil {
		HandleError(err)
		response.Success = false
		response.Error = protocol.ERROR_INTERNAL_ERROR
		JSONResponse(response, w)
		return
	}
	response.Token = t
	JSONResponse(response, w)
}

func isValidPhone(phone string) bool {
	m, _ := regexp.MatchString("\\+(9[976]\\d|8[987530]\\d|6[987]\\d|5[90]\\d|42\\d|3[875]\\d|2[98654321]\\d|9[8543210]|8[6421]|6[6543210]|5[87654321]|4[987654310]|3[9643210]|2[70]|7|1)\\d{1,14}$", phone)
	return m
}

func isValidEmail(email string) bool {
	m, _ := regexp.MatchString("^([a-zA-Z0-9_\\-\\.])+@([a-zA-Z0-9_-])+\\.([a-zA-Z0-9_-])+", email)
	return m
}
