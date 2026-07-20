package handlers

import (
	"net/http"

	"lumen/go-siakad/internal/auth"
	LecturerHandler "lumen/go-siakad/internal/handlers/lecturers"
	StudentHandler "lumen/go-siakad/internal/handlers/students"
	"lumen/go-siakad/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Handler(engine *gin.Engine) {
	engine.Use(gin.Recovery())
	engine.Use(middleware.NewCORSMiddleware())
	engine.Use(middleware.AccessLog())

	engine.GET("/student/list", StudentHandler.LoadList)
	engine.GET("/student/detail", StudentHandler.LoadByID)
	engine.GET("/teacher/list", LecturerHandler.GetLecturers)

	// Temporary manual-verification route for the new auth middleware
	// (quickstart.md §2) — not a real feature endpoint.
	debug := engine.Group("/internal/_debug")
	debug.Use(middleware.RequireRole(auth.RoleAdministrator))
	debug.GET("/whoami", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
}
