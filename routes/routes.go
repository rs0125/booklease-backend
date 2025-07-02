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
		BookRoutes.POST("/", middleware.RequireAuth(app), api.CreateBook)
		BookRoutes.POST("/:id/wishlist", api.AddToWishlist)
		BookRoutes.GET("/wishlist", api.Wishlist)
	}

	UserRoutes := r.Group("/user", middleware.RequireAuth(app))
	{
		UserRoutes.POST("/signup", api.CreateOrFetchUser)
		UserRoutes.POST("/phone", api.UpdatePhoneNumber)
		UserRoutes.GET("/", api.GetUserProfile)
	}

}
