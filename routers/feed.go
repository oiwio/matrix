package routers

import (
	"matrix/auth"
	"matrix/handlers"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func SetFeedRoutes(router *mux.Router) *mux.Router {
	router.Handle("/feed",
		negroni.New(
			negroni.HandlerFunc(auth.RequireTokenAuthentication),
			negroni.HandlerFunc(handlers.PostFeed),
		)).Methods("POST")
	router.Handle("/feed/{FeedId}",
		negroni.New(
			negroni.HandlerFunc(auth.RequireTokenAuthentication),
			negroni.HandlerFunc(handlers.GetFeedById),
		)).Methods("GET")
	return router
}
