package combonamehandler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"catalog/internal/usecase/combonameusecase"

	comboname "catalog/internal/entity/comboname"

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

	response := ComboNameResponse{
		ID:            result.ID,
		Name:          result.Name,
		ComboNameUuid: result.ComboNameUuid,
		Nickname:      result.Nickname,
		IsAvailable:   result.IsAvailable,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (h *WebComboNameHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Query Params
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	combos, total, err := h.ComboNameService.GetAll(r.Context(), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Monta resposta
	var response []ComboNameResponse
	for _, combo := range combos {
		response = append(response, ComboNameResponse{
			ID:            combo.ID,
			Name:          combo.Name,
			ComboNameUuid: combo.ComboNameUuid,
			Nickname:      combo.Nickname,
			IsAvailable:   combo.IsAvailable,
		})
	}

	paginated := PaginatedComboNameResponse{
		Data:       response,
		Page:       page,
		Limit:      limit,
		TotalItems: total,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(paginated)
}

func (h *WebComboNameHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID da rota
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Faz o decode do JSON de entrada
	var req UpdateComboNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Monta a entidade ComboName
	combo := &comboname.ComboName{
		ID:          id,
		Name:        req.Name,
		Nickname:    req.Nickname,
		IsAvailable: req.IsAvailable,
	}

	// Executa o update no usecase
	err = h.ComboNameService.Update(r.Context(), combo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retorna 204 No Content (atualização bem-sucedida)
	w.WriteHeader(http.StatusNoContent)
}

func (h *WebComboNameHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateComboNameRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	combo := &comboname.ComboName{
		Name:        req.Name,
		Nickname:    req.Nickname,
		IsAvailable: req.IsAvailable,
	}

	id, err := h.ComboNameService.Create(r.Context(), combo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(CreateComboNameResponse{
		ID: id,
	})
}

func (h *WebComboNameHandler) Disable(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.ComboNameService.Disable(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content → sucesso sem body
}
