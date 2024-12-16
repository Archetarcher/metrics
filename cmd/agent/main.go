package main

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/agent/client/grpc"
	"github.com/Archetarcher/metrics.git/internal/agent/client/rest"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/logger"
	"github.com/Archetarcher/metrics.git/internal/agent/services"
	"go.uber.org/zap"
	"log"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildData()
	conf := config.NewConfig()
	conf.ParseConfig()

	if err := logger.Initialize(conf.LogLevel); err != nil {
		log.Fatal("failed to init logger")
	}
	s := services.NewMetricsService(conf)

	if conf.EnableGRPC {
		err := grpc.Run(conf, s)
		if err != nil {
			logger.Log.Error("failed to start grpc client with error, finishing app", zap.Error(err))
			return
		}
	} else {
		err := rest.Run(conf, s)
		if err != nil {
			logger.Log.Error("failed to start rest client with error, finishing app", zap.Error(err))
			return
		}
	}

}

func printBuildData() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
