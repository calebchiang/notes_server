package controllers

import (
	"net/http"

	"github.com/calebchiang/notes_server/database"
	"github.com/calebchiang/notes_server/models"
	"github.com/gin-gonic/gin"
)

func CreateNotebook(c *gin.Context) {
	var input struct {
		Title    string `json:"title"`
		Color    string `json:"color"`
		Category string `json:"category"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if input.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title is required",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	notebook := models.Notebook{
		UserID:   userID.(uint),
		Title:    input.Title,
		Color:    input.Color,
		Category: input.Category,
	}

	if err := database.DB.Create(&notebook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create notebook",
		})
		return
	}

	c.JSON(http.StatusCreated, notebook)
}

func GetNotebooks(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	type NotebookWithCount struct {
		ID        uint   `json:"id"`
		UserID    uint   `json:"user_id"`
		Title     string `json:"title"`
		Color     string `json:"color"`
		Category  string `json:"category"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		NoteCount int    `json:"note_count"`
	}

	var notebooks []NotebookWithCount

	err := database.DB.
		Table("notebooks").
		Select(`
			notebooks.id,
			notebooks.user_id,
			notebooks.title,
			notebooks.color,
			notebooks.category,
			notebooks.created_at,
			notebooks.updated_at,
			COUNT(notes.id) as note_count
		`).
		Joins("LEFT JOIN notes ON notes.notebook_id = notebooks.id").
		Where("notebooks.user_id = ?", userID.(uint)).
		Group("notebooks.id").
		Order("notebooks.created_at desc").
		Scan(&notebooks).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch notebooks",
		})
		return
	}

	c.JSON(http.StatusOK, notebooks)
}

func DeleteNotebook(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	notebookID := c.Param("id")
	if notebookID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Notebook ID required",
		})
		return
	}

	var notebook models.Notebook

	if err := database.DB.
		Where("id = ? AND user_id = ?", notebookID, userID.(uint)).
		First(&notebook).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Notebook not found",
		})
		return
	}

	if err := database.DB.Delete(&notebook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete notebook",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notebook deleted successfully",
	})
}
