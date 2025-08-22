package handlers

import (
	"charity/components"
	"charity/internal/auth"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(a *auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: a}
}

func (h *AuthHandler) ShowLogin(c *gin.Context) {
	err := components.Login("").Render(c, c.Writer)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering login page")
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	log.Printf("Login attempt for email: %s", email)

	token, admin, err := h.authService.LoginAdmin(c.Request.Context(), email, password)
	if err != nil {
		log.Printf("Login failed for email %s: %v", email, err)
		// Show error message on login page
		err = components.Login("Invalid email or password").Render(c, c.Writer)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error rendering login page")
		}
		return
	}

	log.Printf("Login successful for admin %s, setting cookie", admin.Email)

	// setting tokens in cookies with better settings
	c.SetCookie("auth_token", token, 3600*24, "/", "", false, true)

	// Log the cookie being set
	log.Printf("Cookie set: auth_token=%s", token[:20]+"...")

	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", "/dashboard")
		c.Status(http.StatusOK)
	} else {
		c.Redirect(http.StatusFound, "/dashboard")
	}
}
