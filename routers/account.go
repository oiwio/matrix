package routers

import (
	"matrix/auth"
	"matrix/handlers"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetAccountRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/register", handlers.Register).Methods("POST")
	router.HandleFunc("/signin", handlers.SignIn).Methods("POST")
	router.Handle("/token/{UserId}",
		negroni.New(
			negroni.HandlerFunc(auth.RequireTokenAuthentication),
			negroni.HandlerFunc(handlers.RefreshToken),
		)).Methods("GET")
	return router
}
