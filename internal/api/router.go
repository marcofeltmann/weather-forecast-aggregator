package api

import (
	"expvar"
	"log/slog"
	"net/http"
)

type result struct {
	externalApiId string
	day1          string
	day2          string
	day3          string
	day4          string
	day5          string
}

const MissingParameterErrorDescription = `Missing request parameter(s).
Please provide valid 'lat' and 'lon' in the URL.`

func addRoutes(mux *http.ServeMux, logger *slog.Logger) {
	mux.Handle("GET /weather", http.HandlerFunc(notImplementedHandler))
	mux.Handle("HEAD /weather", http.HandlerFunc(notImplementedHandler))

	mux.Handle("GET /debug/vars", expvar.Handler())
}

func notImplementedHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
