package combonamehandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *WebComboNameHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.GetAll)
	r.Get("/{id}", h.GetByID)
	return r
}
