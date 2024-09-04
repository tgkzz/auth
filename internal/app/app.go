package app

import (
	"context"
	grpcApp "github.com/tgkzz/auth/internal/app/grpc"
	"github.com/tgkzz/auth/internal/service/auth"
	"github.com/tgkzz/auth/internal/storage/postgresql"

	"log/slog"
)

type App struct {
	GRPCServer *grpcApp.App
}

func New(ctx context.Context, log *slog.Logger, grpcPort int, storagePath string) *App {
	pgStorage, err := postgresql.New(ctx, storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, pgStorage)

	grpcapp := grpcApp.New(log, authService, grpcPort)

	return &App{GRPCServer: grpcapp}
}
