package main

import (
	"net"

	"github.com/mahmud3253/Project/Post_Service/config"
	pb "github.com/mahmud3253/Project/Post_Service/genproto"
	"github.com/mahmud3253/Project/Post_Service/pkg/db"
	"github.com/mahmud3253/Project/Post_Service/pkg/logger"
	"github.com/mahmud3253/Project/Post_Service/service"
	grpcClient "github.com/mahmud3253/Project/Post_Service/service/grpc_client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "post-service")
	defer logger.Cleanup(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	grpcC, err := grpcClient.New(cfg)
	if err != nil {
		log.Error("error establishing grpc connection", logger.Error(err))
		return
	}

	PostService := service.NewUserService(connDB, log, grpcC)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterPostServiceServer(s, PostService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
