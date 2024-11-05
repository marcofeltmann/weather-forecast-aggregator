package api

import "log/slog"

type Config struct {
	WeatherApiKey string
	Logger        *slog.Logger
}
