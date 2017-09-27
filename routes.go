package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func CreateRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"GetAllPages",
		"GET",
		"/api/pages",
		authMiddleware(pagesHandler),
	},
	Route{
		"CreatePage",
		"POST",
		"/api/pages",
		authMiddleware(createPageHandler),
	},
	Route{
		"GetOnePage",
		"GET",
		"/api/pages/{title}",
		authMiddleware(pageHandler),
	},
	Route{
		"RegisterUser",
		"POST",
		"/api/users",
		userHandler,
	},
	Route{
		"AuthUser",
		"POST",
		"/api/auth",
		authHandler,
	},
}
