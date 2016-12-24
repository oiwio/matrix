package main

import "net/http"

func main() {

	router := NewRouter()

	http.ListenAndServeTLS(":1234", "./secure/server.crt", "./secure/server.key", router)
}
