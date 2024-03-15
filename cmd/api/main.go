package main

import (
	"auth-api/internal/server"
	"fmt"
	"log/slog"
)

func main() {
	server := server.NewServer()

	slog.Info("Server is running on:", slog.String("port", server.Addr))
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
