package main

import (
	"charity/components"
	db "charity/db/sqlc"
	"charity/handlers"
	"charity/internal/auth"
	"charity/internal/config"
	p "charity/internal/db"
	"charity/internal/services"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	c, err := config.LoadConfig(".")
	if err != nil {
		panic(fmt.Sprintf("Could not load config: %v", err))
	}

	dbConn, err := sql.Open(c.DBDriver, p.GetDBSource(c, c.DBName))
	if err != nil {
		panic(fmt.Sprintf("Could not load DB: %v", err))
	}

	m, err := migrate.New(
		"file://db/migrations",
		p.GetDBSource(c, c.DBName),
	)
	if err != nil {
		log.Fatalf("Unable to instantiate the database schema migrator - %v", err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Unable to migrate up to the latest database schema - %v", err)
		}
	}

	defer func(dbConn *sql.DB) {
		err := dbConn.Close()
		if err != nil {
			panic(fmt.Sprintf("Could not close DB: %v", err))
		}
	}(dbConn)

	queries := db.New(dbConn)
	// Create a new Gin router
	router := gin.Default()

	router.Static("/static", "./static")
	authService := auth.NewAuthService(queries, c)
	adminHandlers := handlers.NewAuthHandler(authService)

	// Initialize project service and handlers
	projectService := services.NewProjectService(queries)
	projectHandlers := handlers.NewProjectHandler(projectService)

	// Initialize blog service and handlers
	blogService := services.NewBlogService(queries)
	blogHandlers := handlers.NewBlogHandler(blogService)

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

	router.GET("/blog-detail", func(c *gin.Context) {
		err := components.BlogDetails().Render(c, c.Writer)
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

	router.GET("/contact", func(c *gin.Context) {
		err := components.Contact().Render(c, c.Writer)
		if err != nil {
			c.String(500, "Error rendering template: %v", err)
			return
		}
	})

	router.GET("/login", func(c *gin.Context) {
		err := components.Login("").Render(c, c.Writer)
		if err != nil {
			c.String(500, "Error rendering template: %v", err)
			return
		}
	})

	// Protected routes
	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware(authService))
	{
		protected.GET("/dashboard", func(c *gin.Context) {
			err := components.Dashboard().Render(c, c.Writer)
			if err != nil {
				c.String(500, "Error rendering templates: %v", err)
				return
			}
		})

		// Project management routes (admin only)
		protected.POST("/api/projects", projectHandlers.CreateProject)
		protected.PUT("/api/projects/:id", projectHandlers.UpdateProject)
		protected.DELETE("/api/projects/:id", projectHandlers.DeleteProject)
		protected.POST("/api/projects/:id/before", projectHandlers.CreateProjectBefore)
		protected.PUT("/api/projects/:id/before", projectHandlers.UpdateProjectBefore)
		protected.POST("/api/projects/:id/after", projectHandlers.CreateProjectAfter)
		protected.PUT("/api/projects/:id/after", projectHandlers.CreateProjectAfter)
		protected.POST("/api/projects/:id/images", projectHandlers.UploadProjectImage)
		protected.DELETE("/api/projects/images/:image_id", projectHandlers.DeleteProjectImage)
		protected.PUT("/api/projects/:id/status", projectHandlers.UpdateProjectStatus)

		// Blog management routes (admin only)
		protected.POST("/api/blog-posts", blogHandlers.CreateBlogPost)
		protected.PUT("/api/blog-posts/:id", blogHandlers.UpdateBlogPost)
		protected.DELETE("/api/blog-posts/:id", blogHandlers.DeleteBlogPost)
	}

	// Public project routes (no authentication required)
	router.GET("/api/projects", projectHandlers.ListProjects)
	router.GET("/api/projects/:id", projectHandlers.GetProject)
	router.GET("/api/projects/:id/before", projectHandlers.GetProjectBefore)
	router.GET("/api/projects/:id/after", projectHandlers.GetProjectAfter)
	router.GET("/api/projects/:id/images", projectHandlers.ListProjectImages)
	router.GET("/api/projects/:id/images/phase", projectHandlers.ListProjectImagesByPhase)

	// Public blog routes (no authentication required)
	router.GET("/api/blog-posts", blogHandlers.ListBlogPosts)
	router.GET("/api/blog-posts/:id", blogHandlers.GetBlogPost)

	router.POST("/api/login", adminHandlers.Login)

	router.GET("/admin/login", adminHandlers.ShowLogin)

	// Start the server on port 8080
	router.Run(":8081")
}
