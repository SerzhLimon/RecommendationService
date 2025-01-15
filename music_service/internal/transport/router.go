package transport

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc gin.HandlerFunc
}

func NewRouter(handleFunctions ApiHandleFunctions) *gin.Engine {
	return NewRouterWithGinEngine(gin.Default(), handleFunctions)
}

func NewRouterWithGinEngine(router *gin.Engine, handleFunctions ApiHandleFunctions) *gin.Engine {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8070"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	for _, route := range getRoutes(handleFunctions) {
		if route.HandlerFunc == nil {
			route.HandlerFunc = DefaultHandleFunc
		}
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			router.PATCH(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		}
	}

	return router
}

func DefaultHandleFunc(c *gin.Context) {
	c.String(http.StatusNotImplemented, "501 not implemented")
}

type ApiHandleFunctions struct {
	Server Server
}

func getRoutes(handleFunctions ApiHandleFunctions) []Route {
	return []Route{
		{
			"listensong",
			http.MethodGet,
			"/song",
			handleFunctions.Server.ListenSong,
		},
		{
			"likesong",
			http.MethodPost,
			"/song/like",
			handleFunctions.Server.LikeSong,
		},
		{
			"createsong",
			http.MethodPost,
			"/song/create",
			handleFunctions.Server.CreateSong,
		},
		{
			"updatesong",
			http.MethodPatch,
			"/song/update",
			handleFunctions.Server.UpdateSong,
		},
		{
			"deletesong",
			http.MethodDelete,
			"/song/delete",
			handleFunctions.Server.DeleteSong,
		},
	}
}
