package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rafawastaken/tick-storm/backend/pkg/logger"
)

func (a *App) routes() http.Handler {
	const apiPrefix = "/api/v1"
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(logger.Middleware(a.log))
	r.Use(middleware.Recoverer)

	r.Get(apiPrefix+"/health", a.handleHealth)
	r.Get(apiPrefix+"/ready", a.handleReady)

	// r.Mount(apiPrefix+"/crypto", cryptoHandler.Routes())
	return r
}
