package routes

import (
	"github.com/calebchiang/notes_server/controllers"
	"github.com/calebchiang/notes_server/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/users", controllers.CreateUser)
	r.POST("/login", controllers.LoginUser)

	auth := r.Group("/users")
	auth.Use(middlewares.RequireAuth())
	{
		auth.GET("/me", controllers.GetCurrentUser)
	}
}
