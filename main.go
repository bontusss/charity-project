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

	// Projects page
	router.GET("/projects", func(c *gin.Context) {
		err := components.Projects().Render(c, c.Writer) // Use your actual projects template
		if err != nil {
			c.String(500, "Error rendering template: %v", err)
			return
		}
	})

	router.GET("/about", func(c *gin.Context) {
		err := components.About("About us").Render(c, c.Writer)
		if err != nil {
			c.String(500, "Error rendering template: %v", err)
			return
		}
	})

	router.GET("/blog", func(c *gin.Context) {
		err := components.Blog().Render(c, c.Writer)
		if err != nil {
			c.String(500, "Error rendering template: %v", err)
			return
		}
	})

	router.GET("/volunteer", func(c *gin.Context) {
		err := components.Volunteer().Render(c, c.Writer)
		if err != nil {
			c.String(500, "Error rendering template: %v", err)
			return
		}
	})

	router.GET("/donate", func(c *gin.Context) {
		err := components.Donate().Render(c, c.Writer)
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

	// Start the server on port 8080
	router.Run(":8080")
}
