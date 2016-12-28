package main

import (
	"matrix/routers"
	"net/http"

	"github.com/urfave/negroni"
)

func main() {
	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)
	http.ListenAndServeTLS(":1234", "./secure/server.crt", "./secure/server.key", n)
}
