package controllers

import (
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/calebchiang/notes_server/database"
	"github.com/calebchiang/notes_server/models"
	"github.com/calebchiang/notes_server/services"
	"github.com/gin-gonic/gin"
)

func GenerateTranscript(c *gin.Context) {
	var input struct {
		NoteID uint   `json:"note_id"`
		URL    string `json:"url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.URL == "" || input.NoteID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "note_id and YouTube URL required",
		})
		return
	}

	videoID := services.ExtractVideoID(input.URL)
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid YouTube URL",
		})
		return
	}

	outputTemplate := videoID + ".%(ext)s"

	cmd := exec.Command(
		"yt-dlp",
		"-o", outputTemplate,
		"--skip-download",
		"--write-subs",
		"--write-auto-subs",
		"--sub-lang", "en.*",
		"--sub-format", "json3",
		"--extractor-args", "youtube:player_client=android",
		input.URL,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to download captions",
		})
		return
	}

	rawFile := videoID + ".en.json3"

	structuredTranscript, err := services.ExtractStructuredTranscript(rawFile, videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse transcript",
		})
		return
	}

	os.Remove(rawFile)

	var builder strings.Builder
	for _, seg := range structuredTranscript.Segments {
		builder.WriteString(seg.Text)
		builder.WriteString(" ")
	}
	fullText := strings.TrimSpace(builder.String())

	transcript := models.Transcript{
		NoteID:   input.NoteID,
		FullText: fullText,
		Source:   "youtube",
		SourceID: videoID,
	}

	if err := database.DB.
		Where("note_id = ?", input.NoteID).
		Assign(transcript).
		FirstOrCreate(&transcript).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save transcript",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transcript successfully saved",
	})
}
