package handlers

import (
	"charity/internal/services"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	blogService *services.BlogService
}

func NewBlogHandler(blogService *services.BlogService) *BlogHandler {
	return &BlogHandler{
		blogService: blogService,
	}
}

// handleImageUpload processes image upload and returns the image link
func (h *BlogHandler) handleImageUpload(c *gin.Context) string {
	if file, err := c.FormFile("image"); err == nil && file != nil {
		// Save the file with URL-safe filename
		ext := filepath.Ext(file.Filename)
		timestamp := time.Now().Unix()
		filename := fmt.Sprintf("blog_%d_%d%s", file.Size, timestamp, ext)
		if err := c.SaveUploadedFile(file, "./static/images/"+filename); err != nil {
			return ""
		}
		return "/static/images/" + filename
	}
	return ""
}

func (h *BlogHandler) CreateBlogPost(c *gin.Context) {
	// Get form data
	title := c.PostForm("title")
	body := c.PostForm("body")

	if title == "" || body == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and body are required"})
		return
	}

	// Handle image upload if provided
	imageLink := h.handleImageUpload(c)

	post, err := h.blogService.CreateBlogPost(c.Request.Context(), title, body, imageLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"blog_post": post})
}

func (h *BlogHandler) GetBlogPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog post ID"})
		return
	}

	post, err := h.blogService.GetBlogPost(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blog_post": post})
}

func (h *BlogHandler) ListBlogPosts(c *gin.Context) {
	posts, err := h.blogService.ListBlogPosts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blog_posts": posts})
}

func (h *BlogHandler) UpdateBlogPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog post ID"})
		return
	}

	// Get form data
	title := c.PostForm("title")
	body := c.PostForm("body")

	if title == "" || body == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and body are required"})
		return
	}

	// Handle image upload if provided
	imageLink := h.handleImageUpload(c)

	post, err := h.blogService.UpdateBlogPost(c.Request.Context(), int32(id), title, body, imageLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blog_post": post})
}

func (h *BlogHandler) DeleteBlogPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog post ID"})
		return
	}

	if err := h.blogService.DeleteBlogPost(c.Request.Context(), int32(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog post deleted successfully"})
}
