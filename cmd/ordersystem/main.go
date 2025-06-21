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
	// ✅ Carregar configs do ambiente
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// // ✅ Conectar ao banco
	db := database.NewPostgresConnection(cfg)
	defer db.Close()

	// ✅ Repository
	comboRepo := repo.NewPostgresRepository(db)

	// ✅ Service (Usecase)
	comboService := service.NewComboNameService(comboRepo)

	// ✅ Handler HTTP
	comboHandler := handler.NewWebComboNameHandler(comboService)

	// ✅ Criar servidor web
	ws := webserver.NewWebServer(cfg.WebServerPort)

	// ✅ Adicionar rota de teste
	ws.AddHandler("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}))
	// ✅ Rota de GetByID
	//ws.AddHandler("/combo-names/{id}", comboHandler.GetByID)

	comboHandler = handler.NewWebComboNameHandler(comboService)

	ws.AddHandler("/combo-names", comboHandler.Routes())

	// ✅ Iniciar servidor
	ws.Start()
}
