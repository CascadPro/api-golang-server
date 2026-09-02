package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CascadePro/api-golang-server/internal/app"
	core_config "github.com/CascadePro/api-golang-server/internal/core/config"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	"go.uber.org/zap"

	_ "github.com/CascadePro/api-golang-server/docs"
)

// @title 				Cascade Pro App API
// @version 			1.0.0
// @description 	Cascade App Pro API
// @host 					127.0.0.1:8000
// @BasePath 			/api/v1
func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init app logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	App, err := app.New(ctx, cfg, logger)
	if err != nil {
		logger.Fatal("failed to create new App", zap.Error(err))
	}

	if err := App.Run(ctx); err != nil {
		logger.Fatal("failed to run App", zap.Error(err))
	}
}
