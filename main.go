package main

import (
	"os"

	"github.com/calebchiang/notes_server/database"
	"github.com/calebchiang/notes_server/models"
	"github.com/calebchiang/notes_server/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	database.Connect()
	database.DB.AutoMigrate(
		&models.User{},
		&models.Notebook{},
		&models.Note{},
		&models.Transcript{},
	)

	r := gin.Default()
	routes.UserRoutes(r)
	routes.NotebookRoutes(r)
	routes.NoteRoutes(r)
	routes.YouTubeRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
