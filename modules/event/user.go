package event

import (
	"matrix/modules/db"
	"time"
)

type (
	UserRegistry struct {
		DeviceToken string
		Phone       string
		Code        string
		Email       string
		URL         string
		Timestamp   time.Time `json:"_"`
	}

	VerifyCode struct {
		Code      string
		Source    string
		User      db.User
		Timestamp time.Time
	}

	UserSignIn struct {
		Device    string
		Token     string
		User      db.User
		Timestamp time.Time
	}

	SMSService struct {
		Error   int    `json:"error"`
		Message string `json:"msg"`
	}

	UserEvent struct {
		EventId      int
		UserRegistry *UserRegistry
		VerifyCode   *VerifyCode
		User         *db.User
		UserSignIn   *UserSignIn
		Location     string
		Settings     *db.UserSettings
	}
)
