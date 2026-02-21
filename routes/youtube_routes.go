package routes

import (
	"github.com/calebchiang/notes_server/controllers"
	"github.com/calebchiang/notes_server/middlewares"
	"github.com/gin-gonic/gin"
)

func YouTubeRoutes(r *gin.Engine) {
	auth := r.Group("/youtube")
	auth.Use(middlewares.RequireAuth())
	{
		auth.POST("/transcripts", controllers.GenerateTranscript)
	}
}
