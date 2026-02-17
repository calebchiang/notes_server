package controllers

import (
	"net/http"

	"github.com/calebchiang/notes_server/database"
	"github.com/calebchiang/notes_server/models"
	"github.com/gin-gonic/gin"
)

func CreateNote(c *gin.Context) {
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

	// Verify notebook belongs to user
	var notebook models.Notebook
	if err := database.DB.
		Where("id = ? AND user_id = ?", notebookID, userID.(uint)).
		First(&notebook).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Notebook not found",
		})
		return
	}

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
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

	note := models.Note{
		UserID:     userID.(uint),
		NotebookID: notebook.ID,
		Title:      input.Title,
		Content:    input.Content,
	}

	if err := database.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create note",
		})
		return
	}

	c.JSON(http.StatusCreated, note)
}

func GetNotes(c *gin.Context) {
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

	// Verify notebook belongs to user
	var notebook models.Notebook
	if err := database.DB.
		Where("id = ? AND user_id = ?", notebookID, userID.(uint)).
		First(&notebook).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Notebook not found",
		})
		return
	}

	var notes []models.Note

	if err := database.DB.
		Where("notebook_id = ?", notebook.ID).
		Order("created_at desc").
		Find(&notes).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch notes",
		})
		return
	}

	c.JSON(http.StatusOK, notes)
}

func DeleteNote(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	notebookID := c.Param("id")
	noteID := c.Param("note_id")

	if notebookID == "" || noteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Notebook ID and Note ID required",
		})
		return
	}

	// Verify notebook belongs to user
	var notebook models.Notebook
	if err := database.DB.
		Where("id = ? AND user_id = ?", notebookID, userID.(uint)).
		First(&notebook).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Notebook not found",
		})
		return
	}

	// Verify note belongs to this notebook AND user
	var note models.Note
	if err := database.DB.
		Where("id = ? AND notebook_id = ? AND user_id = ?", noteID, notebook.ID, userID.(uint)).
		First(&note).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Note not found",
		})
		return
	}

	if err := database.DB.Delete(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete note",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Note deleted successfully",
	})
}

func UpdateNote(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	noteID := c.Param("note_id")
	notebookID := c.Param("id")

	var note models.Note

	// Ensure note belongs to the logged in user and notebook
	if err := database.DB.
		Where("id = ? AND notebook_id = ? AND user_id = ?", noteID, notebookID, userID.(uint)).
		First(&note).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Note not found",
		})
		return
	}

	var input struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if input.Title != nil {
		note.Title = *input.Title
	}

	if input.Content != nil {
		note.Content = *input.Content
	}

	if err := database.DB.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update note",
		})
		return
	}

	c.JSON(http.StatusOK, note)
}
