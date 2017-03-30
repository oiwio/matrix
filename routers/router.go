package routers

import (
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router = SetAccountRoutes(router)
	router = SetUserRoutes(router)
	router = SetFeedRoutes(router)
	router = SetFriendRoutes(router)
	router = SetTagRoutes(router)
	return router
}
