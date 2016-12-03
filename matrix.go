package main

import (
	"github.com/gorilla/mux"
	"github.com/meatballhat/negroni-logrus"
	"github.com/urfave/negroni"
	"matrix/controllers"
)

func main() {
	r := mux.NewRouter()

	// Router binding
	r.HandleFunc("/", controllers.HomeHandler)
	r.HandleFunc("/feed",controllers.PostFeed).Methods("POST")
	r.HandleFunc("/music/{MusicId}",controllers.GetMusic).Methods("GET")
	
	n := negroni.Classic()
	n.Use(negronilogrus.NewMiddleware())
	n.UseHandler(r)

	n.Run(":1234")
}
