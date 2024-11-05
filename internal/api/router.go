package api

import (
	"encoding/json"
	"expvar"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/marcofeltmann/weather-forecast-aggregator/internal/types"
)

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

	transfer := struct {
		WeatherAPI1 types.FiveDayForecast
		// Extend with new APIs here
	}{}

	aa := s.aggregators()
	for i, a := range aa {
		part, err := a.AggregateWeather(lat, lon)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Request API %d failed: %+v", i, err)
			return
		}
		switch i {
		case 0:
			transfer.WeatherAPI1 = part

		// Extend with new APIs here

		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Unhandled WeatherAPI ID %d", i)
			return
		}
	}

	data, err := json.Marshal(transfer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Marshalling response %+v failed: %+v", transfer, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
