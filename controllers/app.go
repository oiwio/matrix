package controllers

import(
    "net/http"
    "gopkg.in/mgo.v2"
    log "github.com/Sirupsen/logrus"
    "runtime"
    "fmt"
	"matrix/modules/db"
	"github.com/garyburd/redigo/redis"
)

var (
	mgoSession *mgo.Session
	ActiveUser *db.User
	RedisConn redis.Conn
)

func init() {
    var(
        err error
    )
	mgoSession, err = mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	mgoSession.SetMode(mgo.Monotonic, true)

	RedisConn,err=redis.Dial("tcp","127.0.0.1:6379")
	if err!=nil{
		log.Errorln("Connect to redis error",err)
        return
	}
}

// HomeHandler is used to handle homepage request
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!\n"))
}

// HandleError is used to handle error in controllers
func HandleError(err error) {
	if err != nil{
		// notice that we're using 1, so it will actually log where
		// the error happened, 0 = this function, we don't want that.
		_, fn, line, _ := runtime.Caller(1)
        log.Errorln(fmt.Sprintf("[error] %s:%d %v", fn, line, err.Error()))
	}
}

// func checkActiveUser (r *http.Request){
// 	var (
// 		err   error
// 		token string
// 		user  *db.User
// 	)
// 	token=r.Header.Get()
// }