package main

import (
	"fmt"
	"net/http"

	"lumen/go-siakad/internal/handlers"
	"lumen/go-siakad/internal/middleware"

	"github.com/NYTimes/gziphandler"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetReportCaller(true)

	engine := gin.New()
	engine.RedirectTrailingSlash = false
	engine.RedirectFixedPath = false

	handlers.Handler(engine)

	var handler http.Handler = gziphandler.GzipHandler(middleware.StripSlashes(engine))

	fmt.Println("Starting Backend API service...")

	err := http.ListenAndServe("localhost:8000", handler)
	if err != nil {
		log.Error(err)
	}
}
