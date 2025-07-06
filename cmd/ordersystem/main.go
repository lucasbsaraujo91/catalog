package main

import (
	"catalog/configs"
	"catalog/internal/event/event"
	"catalog/internal/event/handler"
	"catalog/internal/infra/cache/redis"
	"catalog/internal/infra/database"
	"catalog/internal/infra/kafkahelper"
	repo "catalog/internal/infra/repository/combonamerepo"
	authhandler "catalog/internal/infra/web/handler/authhandler"
	"catalog/internal/infra/web/handler/cachehandler"
	combonamehandler "catalog/internal/infra/web/handler/combonamehandler"
	webmiddleware "catalog/internal/infra/web/middleware"
	"catalog/internal/infra/web/webserver"
	"catalog/internal/usecase/combonameusecase"
	"catalog/pkg/events"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Carrega configurações
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Redis
	redisClient := redis.NewRedisClient(
		cfg.RedisHost,
		cfg.RedisPort,
		cfg.RedisPassword,
		cfg.RedisDB,
	)
	defer redisClient.Client.Close()

	// PostgreSQL
	db := database.NewPostgresConnection(cfg)
	defer db.Close()

	// Kafka
	kafkaBroker := cfg.KafkaBrokerAddress
	if kafkaBroker == "" {
		kafkaBroker = "kafka:9092"
	}
	kafkaWriter := kafkahelper.GetKafkaWriter(kafkaBroker, cfg.KafkaTopicComboName)
	defer kafkaWriter.Close()

	// Dispatcher e eventos
	dispatcher := events.NewEventDispatcher()
	comboEvent := event.NewComboNameCreatedEvent()
	comboHandler := handler.NewComboNameCreatedHandler(kafkaWriter)
	_ = dispatcher.Register(comboEvent.GetName(), comboHandler)

	// Repositório e serviço
	comboNameRepo := repo.NewPostgresRepository(db, redisClient)
	comboNameService := combonameusecase.NewComboNameService(comboNameRepo, comboEvent, dispatcher)

	// Handlers HTTP
	comboNameHandler := combonamehandler.NewWebComboNameHandler(comboNameService)
	cacheHandler := cachehandler.NewWebCacheHandler(redisClient)

	// Servidor web
	ws := webserver.NewWebServer(cfg.WebServerPort)

	// Rota pública
	ws.AddHandler("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}))

	// Rota pública de login
	ws.AddHandler("/login", http.HandlerFunc(authhandler.LoginHandler))

	// Rotas protegidas com JWT
	protected := chi.NewRouter()
	protected.Use(webmiddleware.JWTAuthMiddleware)
	protected.Mount("/combo-names", comboNameHandler.Routes())
	protected.Mount("/limpa-cache", cacheHandler.Routes())
	ws.AddHandler("/", protected)

	// Start server
	ws.Start()
}
