package api

import "log/slog"

// Config for the API server.
// It differs from the application config as it requires a logger but not host.
type Config struct {
	WeatherApiKey string
	Logger        *slog.Logger
}
