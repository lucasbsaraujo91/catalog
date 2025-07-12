package cachehandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *WebCacheHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/", h.ClearCacheHandler) // POST /limpa-cache?prefix=combo:
	return r
}
