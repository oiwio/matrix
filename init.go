package main

import "matrix/config"

var configuration config.Config

func init() {
	configuration = config.New()
	// log.Infoln(configuration)
}
