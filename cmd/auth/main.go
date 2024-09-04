package main

import (
	"context"
	"github.com/tgkzz/auth/internal/app"
	"github.com/tgkzz/auth/internal/config"
	"github.com/tgkzz/auth/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	ctx := context.Background()

	application := app.New(ctx, log, cfg.GRPC.Port, cfg.StoragePath)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	log.Info("Gracefully stopped")

}
