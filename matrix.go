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
	http.ListenAndServeTLS(configuration.Server.Host, configuration.Server.ServerCertificatePath, configuration.Server.ServerKeyPath, n)
}
