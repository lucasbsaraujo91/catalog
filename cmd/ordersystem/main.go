package main

import (
	"catalog/configs"
	"catalog/internal/infra/database"
	repo "catalog/internal/infra/repository/combonamerepo"
	handler "catalog/internal/infra/web/handler/combonamehandler"
	"catalog/internal/infra/web/webserver"
	service "catalog/internal/usecase/combonameusecase"
	"log"
	"net/http"
)

func main() {
	// Carregar configs do ambiente
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Conectar ao banco
	db := database.NewPostgresConnection(cfg)
	defer db.Close()

	// Repository
	comboNameRepo := repo.NewPostgresRepository(db)

	// Service (Usecase)
	comboNameService := service.NewComboNameService(comboNameRepo)

	// Handler HTTP
	comboNameHandler := handler.NewWebComboNameHandler(comboNameService)

	// Criar servidor web
	ws := webserver.NewWebServer(cfg.WebServerPort)

	// Adicionar rota de teste
	ws.AddHandler("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}))

	ws.AddHandler("/combo-names", comboNameHandler.Routes())

	// Iniciar servidor
	ws.Start()
}
