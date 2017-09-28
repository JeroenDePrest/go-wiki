package router

import (
	"github.com/gorilla/mux"
	"gowiki/page"
	"gowiki/user"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func Create() *mux.Router {

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
		authMiddleware(page.PagesHandler),
	},
	Route{
		"CreatePage",
		"POST",
		"/api/pages",
		authMiddleware(page.CreatePageHandler),
	},
	Route{
		"GetOnePage",
		"GET",
		"/api/pages/{title}",
		authMiddleware(page.PageHandler),
	},
	Route{
		"RegisterUser",
		"POST",
		"/api/users",
		user.UserHandler,
	},
	Route{
		"AuthUser",
		"POST",
		"/api/auth",
		user.AuthHandler,
	},
}
