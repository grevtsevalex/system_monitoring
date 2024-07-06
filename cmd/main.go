package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/grevtsevalex/system_monitoring/internal/logger"
	"github.com/grevtsevalex/system_monitoring/internal/server"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "env.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := NewConfig(configFile)
	if err != nil {
		err = fmt.Errorf("config initialization: %w", err)
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	logg := logger.New(config.Logger.Level, os.Stdout)

	grpcServer := server.NewServer(server.Config{Port: config.Server.Port}, logg)

	logg.Info("monitoring is running...")

	go func() {
		if err := grpcServer.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := grpcServer.Stop(ctx); err != nil {
		logg.Error("failed to stop grpc server: " + err.Error())
	}
}
