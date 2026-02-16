package routes

import (
	"github.com/calebchiang/notes_server/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/users", controllers.CreateUser)
	r.POST("/login", controllers.LoginUser)
}
