package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/marcofeltmann/weather-forecast-aggregator/internal/api"
)

/*
Inspired by Mat Ryer's "How I Write HTTP Web Services After n Years" series

Defer he heavy lifting into a function that returns an error if something breaks.
You can easily log that error at one place and return it in the setup for main.
This avoids a lot of the annoying 'if err != nil { println(err); return }' we've
all seen out there
*/
func main() {
	logger := slog.Default()
	host := ":8080"
	logger.Info("Starting server", slog.String("address", host))
	if err := run(logger, host); err != nil {
		logger.Error("Server failed while running.", slog.Any("result", err))
		os.Exit(2)
	}
}

func run(logger *slog.Logger, host string) error {
	srv := api.NewServer(logger)
	if err := http.ListenAndServe(host, srv.Handler()); err != nil {
		return fmt.Errorf("ListenAndServe on %s failed: %w", host, err)
	}
	return nil
}
