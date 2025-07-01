package routes

import (
	"bookapi/api"
	"bookapi/middleware"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(r *gin.Engine, app *firebase.App) {
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/hello", middleware.RequireAuth(app), api.HelloHandler)
	}

	BookRoutes := r.Group("/book")
	{
		BookRoutes.GET("/", api.GetBooks)
		BookRoutes.GET("/:id", api.GetBook)
		BookRoutes.DELETE("/:id", api.DeleteBook)
		BookRoutes.POST("/", api.CreateBook)
	}
}
