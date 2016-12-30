package routers

import (
	"matrix/handlers"

	"github.com/gorilla/mux"
)

func SetFeedRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/feed", handlers.PostFeed).Methods("POST")
	router.HandleFunc("/feed/{FeedId}", handlers.GetFeedById).Methods("GET")
	router.HandleFunc("/feed/{FeedId}", handlers.DelFeed).Methods("DELETE")
	router.HandleFunc("/u/feeds", handlers.GetFeedsByUserId).Methods("POST")
	return router
}
