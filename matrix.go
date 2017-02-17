package main

import (
	"log"
	"matrix/routers"
	"net/http"
	_ "net/http/pprof"

	"github.com/urfave/negroni"
)

func main() {
	//这里实现了远程获取pprof数据的接口
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)
	http.ListenAndServeTLS(configuration.Server.Host, configuration.Server.ServerCertificatePath, configuration.Server.ServerKeyPath, n)
}
