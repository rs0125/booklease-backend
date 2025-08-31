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
		BookRoutes.GET("/mybooks", middleware.RequireAuth(app), api.MyBooks)
		BookRoutes.DELETE("/:id", middleware.RequireAuth(app), api.DeleteBook)
		BookRoutes.POST("/", middleware.RequireAuth(app), api.CreateBook)
		BookRoutes.POST("/:id/wishlist", middleware.RequireAuth(app), api.AddToWishlist)
		BookRoutes.GET("/wishlist", middleware.RequireAuth(app), api.Wishlist)
	}

	RentalRoutes := r.Group("/rentals")
	{
		RentalRoutes.POST("/", middleware.RequireAuth(app), api.PostRental)
		RentalRoutes.GET("/", middleware.RequireAuth(app), api.GetRentals)

		RentalRoutes.GET("/lent", middleware.RequireAuth(app), api.LentMaterials)         // books listed by user that someone has rented - lent
		RentalRoutes.GET("/borrowed", middleware.RequireAuth(app), api.BorrowedMaterials) //books that user has borrowed
		RentalRoutes.POST("/:id/decision", middleware.RequireAuth(app), api.DecideRental)
		RentalRoutes.DELETE("/delete/:id", middleware.RequireAuth(app), api.DeleteRental)

		RentalRoutes.PATCH("/:id/return", middleware.RequireAuth(app), api.ReturnRental)
	}

	UserRoutes := r.Group("/user", middleware.RequireAuth(app))
	{
		UserRoutes.POST("/signup", api.CreateOrFetchUser)
		UserRoutes.POST("/phone", api.UpdatePhoneNumber)
		UserRoutes.GET("/", api.GetUserProfile)

	}

	NotiRoutes := r.Group("/notifications")
	{
		NotiRoutes.GET("/", middleware.RequireAuth(app), api.Notifications)
		NotiRoutes.DELETE("/:id", middleware.RequireAuth(app), api.DeleteNotification)
		NotiRoutes.PATCH("/:id/seen", middleware.RequireAuth(app), api.MarkNotificationSeen)
		NotiRoutes.DELETE("/", middleware.RequireAuth(app), api.DeleteAllNotifications)
	}

}
