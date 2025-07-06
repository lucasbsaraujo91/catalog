package authhandler

import (
	"catalog/configs"
	"catalog/pkg/auth"
	"encoding/json"
	"net/http"
	"time"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	// Carrega as configs com usuário e senha do .env
	cfg, err := configs.LoadConfig()
	if err != nil {
		http.Error(w, "Erro ao carregar configurações", http.StatusInternalServerError)
		return
	}

	if creds.Username != cfg.AuthUsername || creds.Password != cfg.AuthPassword {
		http.Error(w, "Usuário ou senha inválidos", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(creds.Username, "user", time.Hour)
	if err != nil {
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{Token: token})
}
