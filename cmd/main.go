package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/grevtsevalex/system_monitoring/internal/collector"
	"github.com/grevtsevalex/system_monitoring/internal/logger"
	"github.com/grevtsevalex/system_monitoring/internal/scouts"
	"github.com/grevtsevalex/system_monitoring/internal/scouts/cpuScout"
	lascout "github.com/grevtsevalex/system_monitoring/internal/scouts/laScout"
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

	ctx, cancel := context.WithCancel(context.Background())
	logg := logger.New(config.Logger.Level, os.Stdout)

	sRunner := scouts.NewScoutsRunner(logg)
	storages := make([]scouts.ScoutStorage, 0)

	if config.Metrics.LoadAverage {
		laSt := lascout.NewLAStorage()
		sRunner.RegisterScout("loadAverage", lascout.NewLoadAverageScout(ctx, logg, laSt))
		storages = append(storages, laSt)
	}

	if config.Metrics.CPU {
		cpuSt := cpuScout.NewCPUStorage()
		sRunner.RegisterScout("CPU", cpuScout.NewCPUScout(ctx, logg, cpuSt))
		storages = append(storages, cpuSt)
	}

	err = sRunner.RunScouts()
	if err != nil {
		logg.Error("failed to run scouts: " + err.Error())
		os.Exit(1)
	}

	grpcServer := server.NewServer(server.Config{Port: config.Server.Port}, logg, collector.NewCollector(storages))

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

	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := grpcServer.Stop(ctx); err != nil {
		logg.Error("failed to stop grpc server: " + err.Error())
	}
}
