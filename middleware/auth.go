package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
)

func RequireAuth(app *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}
		idToken := strings.TrimPrefix(authHeader, "Bearer ")

		authClient, err := app.Auth(context.Background())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to init Firebase"})
			return
		}

		token, err := authClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// email, ok := token.Claims["email"].(string)
		// if !ok || !strings.HasSuffix(email, "@vitstudent.ac.in") {
		// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized domain"})
		// 	return
		// }

		email, ok := token.Claims["email"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Email not found" + err.Error()})
			return
		}

		uid, ok := token.Claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "UID not found in token"})
			return
		}

		name, ok := token.Claims["name"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Name not found in token"})
			return
		}

		log.Println("✅ Authenticated:", email)
		c.Set("email", email)
		c.Set("uid", uid)
		c.Set("name", name)
		c.Next()
	}
}
