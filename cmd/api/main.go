package main

import (
	"auth-api/internal/server"
	"fmt"
	"log/slog"
)

func main() {
	server := server.NewServer()
	server.RegisterRoutes()

	slog.Info("Server is running on:", slog.Int("port", server.Port))

	_ = server.App.Listen(fmt.Sprintf(":%d", server.Port))
}
