package handlers

import (
	"encoding/json"
	"fmt"
	"matrix/config"
	"matrix/producer"
	"net/http"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
)

var (
	mgoSession    *mgo.Session
	RedisConn     redis.Conn
	configuration config.Config
)

func init() {
	var (
		err error
	)

	configuration = config.New()

	producer.Connect(configuration.NSQ.Host)

	mgoSession, err = mgo.Dial(configuration.MongoDB.Host)
	if err != nil {
		panic(err)
	}
	mgoSession.SetMode(mgo.Monotonic, true)

}

func JSONResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// HomeHandler is used to handle homepage request
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!\n"))
}

// HandleError is used to handle error in controllers
func HandleError(err error) {
	if err != nil {
		// notice that we're using 1, so it will actually log where
		// the error happened, 0 = this function, we don't want that.
		_, fn, line, _ := runtime.Caller(1)
		log.Errorln(fmt.Sprintf("[error] %s:%d %v", fn, line, err.Error()))
	}
}
