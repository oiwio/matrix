package routers

import (
	"matrix/handlers"

	"github.com/gorilla/mux"
)

func SetFriendRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/follow/{ResponderId}", handlers.FollowUser).Methods("PUT")
	router.HandleFunc("/unfollow/{ResponderId}", handlers.UnfollowUser).Methods("PUT")
	return router
}
