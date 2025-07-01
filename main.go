package main

import (
	"log"

	"bookapi/routes"
	"bookapi/services"

	"github.com/gin-gonic/gin"
)

func main() {
	services.InitFirebase()
	services.InitDatabase() // ← Initialize DB here
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	r.Static("/static", "./static")

	routes.RegisterAPIRoutes(r, services.App)

	log.Println("🚀 Server running at http://localhost:8080")
	r.Run(":8080")
}
