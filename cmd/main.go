package main

import (
	"net"
	"post-service/config"
	pbp "post-service/genproto/post"
	"post-service/pkg/db"
	"post-service/pkg/logger"
	"post-service/service"
	"post-service/service/grpcClient"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "template-service")
	defer logger.Cleanup(log)

	log.Info("main:sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, err := db.ConnectToDB(cfg)

	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	grpcClient, err := grpcClient.New(cfg)

	if err != nil {
		log.Fatal("grpc connection error", logger.Error(err))
	}

	postService := service.NewPostService(connDB, log, grpcClient)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pbp.RegisterPostServiceServer(s, postService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

}

