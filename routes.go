package main

import (
	"matrix/handlers"
	"net/http"
)

// Route struct
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes route array
type Routes []Route

var (
	authRoutes = Routes{
		Route{
			"PostFeed",
			"POST",
			"/feed",
			handlers.PostFeed,
		},
		Route{
			"GetMusic",
			"GET",
			"/music/{MusicId}",
			handlers.GetMusic,
		},
	}

	commonRoutes = Routes{
		Route{
			"Home",
			"GET",
			"/",
			handlers.HomeHandler,
		},
		Route{
			"SignIn",
			"POST",
			"/u/signin",
			handlers.SignIn,
		},
		Route{
			"Register",
			"POST",
			"/u/register",
			handlers.Register,
		},
	}
)
