package routers

import (
	"matrix/handlers"

	"github.com/gorilla/mux"
)

func SetUserRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/u/profile/{UserId}", handlers.GetProfile).Methods("GET")
	router.HandleFunc("/u/me", handlers.GetProfile).Methods("GET")
	router.HandleFunc("/u/profile", handlers.UpdateProfile).Methods("PUT")
	return router
}
