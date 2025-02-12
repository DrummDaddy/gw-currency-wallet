package grpc

import (
	"log"
	"net"

	"gw-currency-wallet/config"
	"gw-currency-wallet/internal/storages"

	walletpb "gw-currency-wallet/internal/grpc/proto"

	"google.golang.org/grpc"
)

type Server struct {
	cfg        *config.Config
	db         storages.Storage
	grpcServer *grpc.Server
}

// Создаем новый gRPC-сервер с зависимостями
func NewServer(cfg *config.Config, db storages.Storage) (*Server, error) {
	s := &Server{
		cfg: cfg,
		db:  db,
	}

	grpcServer := grpc.NewServer()
	walletpb.RegisterWalletServiceServer(grpcServer, nil)

	s.grpcServer = grpcServer
	return s, nil
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+s.cfg.GRPCPort)
	if err != nil {
		return err
	}

	log.Printf("gRPC сервер запущен на порту %s\n", s.cfg.GRPCPort)
	if err := s.grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() {
	log.Println("Остановка gRPC сервера...")
	s.grpcServer.GracefulStop()
}
