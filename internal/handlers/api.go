package handlers

import (
	LecturerHandler "lumen/go-siakad/internal/handlers/lecturers"
	StudentHandler "lumen/go-siakad/internal/handlers/students"

	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
)

func Handler(r *chi.Mux) {
	r.Use(chiMiddleware.StripSlashes)
	r.Use(gziphandler.GzipHandler)
	r.Route("/student", func(router chi.Router) {
		router.Get("/list", StudentHandler.LoadList)
		router.Get("/detail", StudentHandler.LoadByID)
	})
	r.Route("/teacher", func(router chi.Router) {
		router.Get("/list", LecturerHandler.GetLecturers)
	})
}
