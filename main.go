package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/marcofeltmann/weather-forecast-aggregator/internal/api"
)

type config struct {
	Host          string `conf:"default::8080"`
	WeatherApiKey string `conf:"required"`
}

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
	var cfg config
	help, err := conf.Parse("", &cfg)
	if err != nil {
		switch err {
		case conf.ErrHelpWanted:
			fmt.Println(help)
			return

		default:
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	if err := run(logger, cfg); err != nil {
		logger.Error("Server failed while running.", slog.Any("result", err))
		os.Exit(2)
	}
}

func run(logger *slog.Logger, cfg config) error {
	srvConf := api.Config{
		Logger:        logger,
		WeatherApiKey: cfg.WeatherApiKey,
	}
	srv := api.NewServer(srvConf)
	if err := http.ListenAndServe(cfg.Host, srv.Handler()); err != nil {
		return fmt.Errorf("ListenAndServe on %s failed: %w", cfg.Host, err)
	}
	return nil
}
