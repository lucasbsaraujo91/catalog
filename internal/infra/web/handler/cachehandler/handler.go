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
