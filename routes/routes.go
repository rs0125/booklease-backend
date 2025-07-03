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
		BookRoutes.DELETE("/:id", middleware.RequireAuth(app), api.DeleteBook)
		BookRoutes.POST("/", middleware.RequireAuth(app), api.CreateBook)
		BookRoutes.POST("/:id/wishlist", middleware.RequireAuth(app), api.AddToWishlist)
		BookRoutes.GET("/wishlist", middleware.RequireAuth(app), api.Wishlist)
	}

	RentalRoutes := r.Group("/rentals")
	{
		RentalRoutes.POST("/", middleware.RequireAuth(app), api.PostRental)
		RentalRoutes.GET("/", middleware.RequireAuth(app), api.GetRentals)
	}

	UserRoutes := r.Group("/user", middleware.RequireAuth(app))
	{
		UserRoutes.POST("/signup", api.CreateOrFetchUser)
		UserRoutes.POST("/phone", api.UpdatePhoneNumber)
		UserRoutes.GET("/", api.GetUserProfile)
	}

}
