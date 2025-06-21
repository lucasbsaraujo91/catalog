package combonamehandler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"catalog/internal/usecase/combonameusecase"

	"github.com/go-chi/chi/v5"
)

type WebComboNameHandler struct {
	ComboNameService *combonameusecase.ComboNameService
}

func NewWebComboNameHandler(service *combonameusecase.ComboNameService) *WebComboNameHandler {
	return &WebComboNameHandler{
		ComboNameService: service,
	}
}

func (h *WebComboNameHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	log.Printf("🟦 Handler → Received id param: %s", idStr)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing ID: %v | Raw value: %s", err, idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	result, err := h.ComboNameService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func (h *WebComboNameHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", h.GetByID)
	return r
}
