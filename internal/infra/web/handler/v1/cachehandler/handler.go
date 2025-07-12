package cachehandler

import (
	"catalog/internal/infra/cache/redis"
	"encoding/json"
	"net/http"
)

type WebCacheHandler struct {
	cache *redis.RedisCache
}

func NewWebCacheHandler(cache *redis.RedisCache) *WebCacheHandler {
	return &WebCacheHandler{cache: cache}
}

// ClearCacheHandler limpa o cache com base no prefixo informado
// @Summary Limpa o cache com o prefixo especificado
// @Tags Cache
// @Security ApiKeyAuth
// @Produce json
// @Param prefix query string true "Prefixo das chaves do Redis"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /limpa-cache [delete]
func (h *WebCacheHandler) ClearCacheHandler(w http.ResponseWriter, r *http.Request) {
	prefix := r.URL.Query().Get("prefix")
	if prefix == "" {
		http.Error(w, "prefix is required", http.StatusBadRequest)
		return
	}

	if err := h.cache.DeleteByPrefix(r.Context(), prefix); err != nil {
		http.Error(w, "failed to clear cache: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "cache cleared",
		"prefix":  prefix,
	})
}
