package main

import (
	"log"
	"time"

	"bookapi/routes"
	"bookapi/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	services.InitFirebase()
	services.InitDatabase() // ← Initialize DB here
	r := gin.Default()
	r.Use(cors.New(cors.Config{

		AllowOrigins:     []string{"http://127.0.0.1:5500"}, // or "*" for all
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	r.Static("/static", "./static")

	routes.RegisterAPIRoutes(r, services.App)

	log.Println("🚀 Server running at http://localhost:8080")
	r.Run(":8080")
}
