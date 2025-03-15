package main

import (
	"budget-tracker/service"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.TODO()
	logger := slog.Default()
	router := http.NewServeMux()
	err := godotenv.Load()
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "Failed to load .env file", slog.Any("error", err))
	}

	service.Start(router, logger)

	port := os.Getenv("PORT")

	logger.LogAttrs(ctx, slog.LevelInfo, "Starting server", slog.String("port", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "Failed to start server", slog.Any("error", err))
		return
	}
}
