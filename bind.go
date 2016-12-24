package main

import (
	"net/http"

	"fmt"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// NewRouter bind routes
func NewRouter() *mux.Router {

	router := mux.NewRouter()
	authBase := mux.NewRouter()

	jm := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("kiwee19920306"), nil
		},
		UserProperty: "userId",
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	router.PathPrefix("/auth").Handler(negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(jm.HandlerWithNext),
		negroni.NewLogger(),
		negroni.Wrap(authBase),
	))

	auth := authBase.PathPrefix("/auth").Subrouter()

	for _, route := range commonRoutes {
		var handler http.Handler
		handler = route.HandlerFunc
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	for _, route := range authRoutes {
		var handler = route.HandlerFunc
		auth.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	fmt.Println(router)
	return router
}
