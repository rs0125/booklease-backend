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
		apiGroup.GET("/FAQ", api.GetFAQ)
	}

	BookRoutes := r.Group("/book")
	{
		BookRoutes.GET("/", api.GetBooks)
		BookRoutes.GET("/:id", api.GetBook)
		BookRoutes.DELETE("/:id", api.DeleteBook)
		BookRoutes.POST("/", api.CreateBook)
		BookRoutes.POST("/:id", api.addToWishlist)
	}

	UserRoutes := r.Group("/user")
	{
		UserRoutes.POST("/signup", api.CreateOrFetchUser)
	}
}
