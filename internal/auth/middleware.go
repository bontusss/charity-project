package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie
		token, err := c.Cookie("auth_token")
		if err != nil {
			log.Printf("No auth_token cookie found: %v", err)
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		log.Printf("Found auth_token cookie: %s", token[:20]+"...")

		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			// Clear invalid cookie
			c.SetCookie("auth_token", "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		log.Printf("Token validated for user: %v", claims["email"])

		// Add user info to context
		c.Set("user_id", claims["sub"])
		c.Set("user_email", claims["email"])
		c.Next()
	}
}
