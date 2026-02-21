package controllers

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/calebchiang/notes_server/services"
	"github.com/gin-gonic/gin"
)

func GenerateTranscript(c *gin.Context) {
	var input struct {
		URL string `json:"url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "YouTube URL required",
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
		"--write-auto-sub",
		"--sub-lang", "en",
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

	c.JSON(http.StatusOK, structuredTranscript)
}
