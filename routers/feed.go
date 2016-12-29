package routers

import (
	"matrix/handlers"

	"github.com/gorilla/mux"
)

func SetFeedRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/feed", handlers.PostFeed).Methods("POST")
	router.HandleFunc("/feed/{FeedId}", handlers.GetFeedById).Methods("GET")
	return router
}
