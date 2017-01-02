package routers

import (
	"matrix/handlers"

	"github.com/gorilla/mux"
)

func SetAccountRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/register", handlers.Register).Methods("POST")
	router.HandleFunc("/signin", handlers.SignIn).Methods("POST")
	router.HandleFunc("/token", handlers.RefreshToken).Methods("GET")
	router.HandleFunc("/exist", handlers.IsExist).Methods("GET")
	return router
}
