package handlers

import (
	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
)

func Handler(r *chi.Mux) {
	r.Use(chiMiddleware.StripSlashes)
	r.Use(gziphandler.GzipHandler)
	r.Route("/student", func(router chi.Router) {
		router.Get("/list", LoadList)
		router.Get("/get", LoadByID)
	})
}
