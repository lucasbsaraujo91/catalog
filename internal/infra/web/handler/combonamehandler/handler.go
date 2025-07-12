package combonamehandler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	comboname "catalog/internal/entity/comboname"
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

// GetByID retorna um combo pelo ID
// @Summary Busca um ComboName pelo ID
// @Tags ComboNames
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "ID do combo"
// @Success 200 {object} ComboNameResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Router /combo-names/{id} [get]
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

// GetAll retorna todos os combos
// @Summary Lista todos os ComboNames com paginação
// @Tags ComboNames
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "Número da página" default(1)
// @Param limit query int false "Limite de itens por página" default(10)
// @Success 200 {object} PaginatedComboNameResponse
// @Failure 401 {string} string "Unauthorized"
// @Router /combo-names [get]
func (h *WebComboNameHandler) GetAll(w http.ResponseWriter, r *http.Request) {
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

// Update atualiza um ComboName existente
// @Summary Atualiza um ComboName existente
// @Tags ComboNames
// @Security ApiKeyAuth
// @Accept json
// @Param id path int true "ID do combo"
// @Param request body UpdateComboNameRequest true "Dados do combo"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /combo-names/{id} [put]
func (h *WebComboNameHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateComboNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	combo := &comboname.ComboName{
		ID:          id,
		Name:        req.Name,
		Nickname:    req.Nickname,
		IsAvailable: req.IsAvailable,
	}

	err = h.ComboNameService.Update(r.Context(), combo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Create cria um novo ComboName
// @Summary Cria um novo ComboName
// @Tags ComboNames
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body CreateComboNameRequest true "Dados do combo"
// @Success 201 {object} CreateComboNameResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /combo-names [post]
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

// Disable desativa um ComboName pelo ID
// @Summary Desativa (soft delete) um ComboName
// @Tags ComboNames
// @Security ApiKeyAuth
// @Param id path int true "ID do combo"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /combo-names/{id} [delete]
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

	w.WriteHeader(http.StatusNoContent)
}
