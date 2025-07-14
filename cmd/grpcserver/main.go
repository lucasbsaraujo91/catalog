package main

import (
	"catalog/configs"
	"catalog/internal/infra/cache/redis"
	"catalog/internal/infra/database"
	pb "catalog/internal/infra/grpc/pb"
	"catalog/internal/infra/grpc/service"
	repo "catalog/internal/infra/repository/combonamerepo"
	"catalog/internal/usecase/combonameusecase"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	redisClient := redis.NewRedisClient(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, cfg.RedisDB)
	defer redisClient.Client.Close()

	db := database.NewPostgresConnection(cfg)
	defer db.Close()

	comboRepo := repo.NewPostgresRepository(db, redisClient)
	comboService := combonameusecase.NewComboNameService(comboRepo, nil, nil)

	grpcServer := grpc.NewServer()

	// Habilita o Server Reflection:
	reflection.Register(grpcServer)

	pb.RegisterComboNameServiceServer(grpcServer, service.NewComboNameGrpcService(comboService))

	lis, err := net.Listen("tcp", ":"+cfg.GRPCServerPort)
	if err != nil {
		log.Fatalf("Erro ao iniciar o listener gRPC: %v", err)
	}

	log.Printf("🚀 Servidor gRPC rodando na porta %s", cfg.GRPCServerPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Erro ao iniciar servidor gRPC: %v", err)
	}
}
