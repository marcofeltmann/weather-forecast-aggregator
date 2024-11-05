package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/marcofeltmann/weather-forecast-aggregator/internal/aggregator/openmeteo"
	"github.com/marcofeltmann/weather-forecast-aggregator/internal/aggregator/weatherapi"
)

/*
Inspiration for the server type is token from Mat Ryer's talk at GopherCon 2019:

"How I write HTTP Web Services after 8 years", watch it on YouTube:
https://www.youtube.com/watch?v=rWBSMsLG8po

But he's constantly updating it kinda each year like on the Grafana blog:
https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
*/
type Server struct {
	logger        *slog.Logger
	mux           *http.ServeMux
	weatherapikey string
}

/*
NewServer returns a configured API server.
*/
func NewServer(c Config) *Server {
	logger := c.Logger
	if logger == nil {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		logger = slog.Default()
		logger.Debug("NewServer called without logger, using slog.Default()")
	}

	/*
		As http.DefaultServeMux is a package-global reference type variable every
		third-party package might manipulate it. Even a different package in this
		codebase. Prometheus for instance registers it's /metrics handler there.

		So creating a new ServeMux type prevents unexpected side effects.
	*/
	mux := http.NewServeMux()
	s := Server{
		mux:           mux,
		logger:        logger,
		weatherapikey: c.WeatherApiKey,
	}

	s.addRoutes()
	return &s
}

func (s Server) Handler() http.Handler {
	return s.mux
}

var meteo Aggregator
var weather Aggregator

func (s Server) aggregators() ([]Aggregator, error) {
	if meteo == nil {
		meteo = openmeteo.DefaultCaller()
	}
	if weather == nil {
		var err error
		weather, err = weatherapi.DefaultCaller(s.weatherapikey)
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
