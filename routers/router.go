package routers

import (
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router = SetAccountRoutes(router)
	router = SetFeedRoutes(router)
	return router
}
