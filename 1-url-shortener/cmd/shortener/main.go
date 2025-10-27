package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"valentino7504/1-url-shortener/internal/api"
	"valentino7504/1-url-shortener/internal/db"
	"valentino7504/1-url-shortener/internal/service"
)

func main() {
	path := *flag.String("db", "urls.db", "path to the db file")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Info("Starting a new server", slog.Any("port", 4000))

	sqliteDB, err := db.GetConnection(path)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	svc := service.NewShortenService(sqliteDB, logger)
	err = http.ListenAndServe(":4000", api.Routes(svc))
	if err != nil {
		logger.Error("App unable to start")
		os.Exit(1)
	}
}
