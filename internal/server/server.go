package server

import (
	"context"
	"fmt"
	"net"

	serverpb "github.com/grevtsevalex/system_monitoring/internal/server/pb"
	"google.golang.org/grpc"
)

// Logger тип логгера.
type Logger interface {
	Log(msg string)
	Error(msg string)
}

// Server модель сервера.
type Server struct {
	config Config
	logger Logger
}

// Config конфиг сервера.
type Config struct {
	Port int
}

// NewServer конструктор сервера.
func NewServer(config Config, logger Logger) *Server {
	return &Server{config: config, logger: logger}
}

// Start старт сервера.
func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		return fmt.Errorf("create grpc listener: %w", err)
	}

	server := grpc.NewServer()
	serverpb.RegisterMonitoringServiceServer(server, NewService())
	s.logger.Log(fmt.Sprintf("starting grpc server on %s", lsn.Addr().String()))

	if err := server.Serve(lsn); err != nil {
		return fmt.Errorf("serve grpc connections: %w", err)
	}

	return nil
}

// Stop остановка сервера.
func (s *Server) Stop(ctx context.Context) error {
	<-ctx.Done()
	s.logger.Log("Stopping grpc server...")
	return nil
}
