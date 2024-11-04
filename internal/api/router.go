package api

import (
	"expvar"
	"log/slog"
	"net/http"
	"strconv"
)

type Result struct {
	ExternalApiId string
	Day1          string
	Day2          string
	Day3          string
	Day4          string
	Day5          string
}

const MissingParameterErrorDescription = `Missing request parameter(s).
Please provide valid 'lat' and 'lon' in the URL.`

func (s Server) addRoutes() {
	mux := s.mux

	// According to https://go.dev/blog/routing-enhancements the GET method also
	// matches the HEAD method, so we don't need the HEAD routes.
	mux.Handle("GET /weather", http.HandlerFunc(s.notImplementedHandler))

	mux.Handle("GET /debug/vars", expvar.Handler())
}

func (s Server) notImplementedHandler(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()

	s.logger.Debug("Got Weather Request", slog.Any("params", params))

	if !params.Has("lat") || !params.Has("lon") {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(MissingParameterErrorDescription))
		if err != nil {
			s.logger.Error(
				"Write response failed",
				slog.String("status", http.StatusText(http.StatusBadRequest)),
				slog.Any("error", err),
			)
		}
		return
	}

	lat, err := strconv.ParseFloat(params.Get("lat"), 64)
	if err != nil {
		s.logger.Error(
			"Unable to parse latitude parameter into float",
			slog.String("lat", params.Get("lat")),
			slog.Any("error", err),
		)
		return
	}
	lon, err := strconv.ParseFloat(params.Get("lon"), 64)
	if err != nil {
		s.logger.Error(
			"Unable to parse longitude parameter into float",
			slog.String("lat", params.Get("lon")),
			slog.Any("error", err),
		)
		return
	}

	s.logger.Debug(
		"Parsed coordinates.",
		slog.Float64("latitude", lat),
		slog.Float64("longitude", lon),
	)

	w.WriteHeader(http.StatusNotImplemented)
	_, err = w.Write([]byte("{}"))
	if err != nil {
		s.logger.Error(
			"Write response failed",
			slog.String("status", http.StatusText(http.StatusBadRequest)),
			slog.Any("error", err),
		)
	}
}
