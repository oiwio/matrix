package routers

import (
	"matrix/handlers"

	"github.com/gorilla/mux"
)

func SetTagRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/tag", handlers.PostTag).Methods("POST")
	router.HandleFunc("/feed/tags", handlers.AddTags).Methods("PUT")
	router.HandleFunc("/tag/feeds", handlers.GetFeedsByTagId).Methods("GET")
	router.HandleFunc("/tag/follow/{TagId}", handlers.FollowTag).Methods("PUT")
	router.HandleFunc("/tag/unfollow/{TagId}", handlers.UnFollowTag).Methods("PUT")
	router.HandleFunc("/s/tag", handlers.FuzzySearchByTagName).Methods("GET")
	return router
}
