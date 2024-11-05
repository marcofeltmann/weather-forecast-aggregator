package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/marcofeltmann/weather-forecast-aggregator/internal/aggregator/openmeteo"
	"github.com/marcofeltmann/weather-forecast-aggregator/internal/aggregator/weatherapi"
)

// Server holds all the information that the API server needs.
//
// Inspiration for the server type is token from Mat Ryer's talk at GopherCon 2019:
// "How I write HTTP Web Services after 8 years", watch it on YouTube:
// https://www.youtube.com/watch?v=rWBSMsLG8po
//
// Zet he's constantly updating it kinda each year like on the Grafana blog:
// https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
type Server struct {
	logger        *slog.Logger
	mux           *http.ServeMux
	weatherapikey string
}

// NewServer returns an API server set up according to the configuration.
func NewServer(c Config) *Server {
	logger := c.Logger
	if logger == nil {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		logger = slog.Default()
		logger.Debug("NewServer called without logger, using slog.Default()")
	}

	// As http.DefaultServeMux is a package-global reference type variable every
	// third-party package might manipulate it. Even a different package in this
	// codebase. Prometheus for instance registers it's /metrics handler there.
	// So creating a new ServeMux type prevents unexpected side effects.
	mux := http.NewServeMux()
	s := Server{
		mux:           mux,
		logger:        logger,
		weatherapikey: c.WeatherApiKey,
	}

	s.addRoutes()
	return &s
}

// Handler makes the routed handler accessible from the outside.
// So we can easily hook it into our http.Server or httptest.Server.
func (s Server) Handler() http.Handler {
	return s.mux
}

// TODO: Extend new API aggregators here
var meteo Aggregator
var weather Aggregator

// aggregators lazy-loads the registered aggregators for later use.
// The lazy-loading approach is used to reduce startup time while increasing
// duration of the first request. Might be helpful inside of Kubernetes.
func (s Server) aggregators() ([]Aggregator, error) {
	if meteo == nil {
		meteo = openmeteo.NewCaller()
	}
	if weather == nil {
		var err error
		weather, err = weatherapi.NewCaller(s.weatherapikey)
		if err != nil {
			return nil, fmt.Errorf("initialize weatherapi caller: %w", err)
		}
	}

	res := []Aggregator{
		meteo,
		weather,
	}
	return res, nil
}
