package main

import (
	"catalog/configs"
	"catalog/internal/infra/cache/redis"
	"catalog/internal/infra/database"
	repo "catalog/internal/infra/repository/combonamerepo"
	"catalog/internal/infra/web/handler/cachehandler"
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

	redisClient := redis.NewRedisClient(
		cfg.RedisHost,
		cfg.RedisPort,
		cfg.RedisPassword,
		cfg.RedisDB,
	)
	defer redisClient.Client.Close()

	// Conectar ao banco
	db := database.NewPostgresConnection(cfg)
	defer db.Close()

	// Repository
	comboNameRepo := repo.NewPostgresRepository(db, redisClient)

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

	cacheHandler := cachehandler.NewWebCacheHandler(redisClient)
	ws.AddHandler("/limpa-cache", cacheHandler.Routes())

	// Iniciar servidor
	ws.Start()
}
