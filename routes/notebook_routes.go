package routes

import (
	"github.com/calebchiang/notes_server/controllers"
	"github.com/calebchiang/notes_server/middlewares"
	"github.com/gin-gonic/gin"
)

func NotebookRoutes(r *gin.Engine) {
	auth := r.Group("/notebooks")
	auth.Use(middlewares.RequireAuth())
	{
		auth.POST("", controllers.CreateNotebook)
		auth.GET("", controllers.GetNotebooks)
		auth.DELETE("/:id", controllers.DeleteNotebook)
	}
}
