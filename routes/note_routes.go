package routes

import (
	"github.com/calebchiang/notes_server/controllers"
	"github.com/calebchiang/notes_server/middlewares"
	"github.com/gin-gonic/gin"
)

func NoteRoutes(r *gin.Engine) {
	auth := r.Group("/notebooks/:id/notes")
	auth.Use(middlewares.RequireAuth())
	{
		auth.POST("", controllers.CreateNote)
		auth.GET("", controllers.GetNotes)
		auth.DELETE("/:note_id", controllers.DeleteNote)
		auth.PATCH("/:note_id", controllers.UpdateNote)
	}
}
