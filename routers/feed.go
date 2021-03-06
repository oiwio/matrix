package routers

import (
	"matrix/handlers"

	"github.com/gorilla/mux"
)

func SetFeedRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/feed", handlers.PostFeed).Methods("POST")
	router.HandleFunc("/feed/{FeedId}", handlers.GetFeedById).Methods("GET")
	router.HandleFunc("/feed/{FeedId}", handlers.DelFeed).Methods("DELETE")
	router.HandleFunc("/music", handlers.SearchMusic).Methods("GET")
	router.HandleFunc("/music/{MusicId}", handlers.GetMusic).Methods("GET")
	router.HandleFunc("/u/feed", handlers.GetFeedsByUserId).Methods("GET")
	router.HandleFunc("/feeds", handlers.GetNewestFeeds).Methods("GET")
	router.HandleFunc("/comment", handlers.PostComment).Methods("POST")
	router.HandleFunc("/comments", handlers.GetCommentsByFeedId).Methods("GET")
	router.HandleFunc("/comment/{CommentId}", handlers.DelComment).Methods("DELETE")

	router.HandleFunc("/feed/upload/{Type}", handlers.GetUploadToken).Methods("GET")

	return router
}
