package main

import (
	"charity/components"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	router.Static("/static", "./static")

	// Define a simple GET endpoint
	router.GET("/", func(c *gin.Context) {
		err := components.Index("Chinedu Onyeizu Foundation").Render(c, c.Writer)
		if err != nil {
			c.String(500, "Error rendering template: %v", err)
			return
		}
	})
	router.GET("/volunteer", func(c *gin.Context) {
		err := components.VolunteerForm("Become a volunteer").Render(c, c.Writer)
		if err != nil {
			c.String(500, "Error rendering template: %v", err)
			return
		}
	})

	router.GET("/contact", func(c *gin.Context) {
		err := components.Contact().Render(c, c.Writer)
		if err != nil {
			c.String(500, "Error rendering template: %v", err)
			return
		}
	})

	router.GET("/coming-soon", func(c *gin.Context) {
		err := components.ComingSoon().Render(c, c.Writer)
		if err != nil {
			c.String(500, "Error rendering template: %v", err)
			return
		}
	})

	router.GET("/cause", func(c *gin.Context) {
		err := components.Causes().Render(c, c.Writer)
		if err != nil {
			c.String(500, "Error rendering template: %v", err)
			return
		}
	})

	// Start the server on port 8080
	router.Run(":8080")
}
