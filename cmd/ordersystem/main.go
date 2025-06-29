package main

import (
	"catalog/configs"
	"catalog/internal/event/event"
	"catalog/internal/event/handler"
	"catalog/internal/infra/cache/redis"
	"catalog/internal/infra/database"
	"catalog/internal/infra/kafkahelper"
	repo "catalog/internal/infra/repository/combonamerepo"
	"catalog/internal/infra/web/handler/cachehandler"
	combonamehandler "catalog/internal/infra/web/handler/combonamehandler"
	"catalog/internal/infra/web/webserver"
	"catalog/internal/usecase/combonameusecase"
	"catalog/pkg/events"
	"log"
	"net/http"
)

func main() {
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

	db := database.NewPostgresConnection(cfg)
	defer db.Close()

	kafkaBroker := cfg.KafkaBrokerAddress
	if kafkaBroker == "" {
		kafkaBroker = "kafka:9092"
	}
	kafkaWriter := kafkahelper.GetKafkaWriter(kafkaBroker, cfg.KafkaTopicComboName)
	defer kafkaWriter.Close()

	dispatcher := events.NewEventDispatcher()
	comboEvent := event.NewComboNameCreatedEvent()
	comboHandler := handler.NewComboNameCreatedHandler(kafkaWriter)

	_ = dispatcher.Register(comboEvent.GetName(), comboHandler)

	comboNameRepo := repo.NewPostgresRepository(db, redisClient)
	comboNameService := combonameusecase.NewComboNameService(comboNameRepo, comboEvent, dispatcher)

	comboNameHandler := combonamehandler.NewWebComboNameHandler(comboNameService)
	cacheHandler := cachehandler.NewWebCacheHandler(redisClient)

	ws := webserver.NewWebServer(cfg.WebServerPort)

	// Rotas
	ws.AddHandler("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}))
	ws.AddHandler("/combo-names", comboNameHandler.Routes())
	ws.AddHandler("/limpa-cache", cacheHandler.Routes())

	// Inicia servidor
	ws.Start()
}
