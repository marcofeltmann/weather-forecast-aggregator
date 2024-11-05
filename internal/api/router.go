package api

import (
	"encoding/json"
	"expvar"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/marcofeltmann/weather-forecast-aggregator/internal/types"
)

var reqs *expvar.Int
var durMin *expvar.Int
var durMax *expvar.Int
var errs *expvar.Int

func init() {
	reqs = expvar.NewInt("requests_sum")
	errs = expvar.NewInt("errors_sum")
	durMin = expvar.NewInt("duration_min")
	durMax = expvar.NewInt("duration_max")
	durMin.Set(int64(time.Hour))
}

const MissingParameterErrorDescription = `Missing request parameter(s).
Please provide valid 'lat' and 'lon' in the URL.`

func (s Server) addRoutes() {
	mux := s.mux

	// According to https://go.dev/blog/routing-enhancements the GET method also
	// handles the HEAD method, so we don't need the HEAD routes.
	mux.Handle("GET /weather", s.meterMiddleware(s.notImplementedHandler))

	mux.Handle("GET /debug/vars", expvar.Handler())
}

func (s Server) meterMiddleware(inner func(http.ResponseWriter, *http.Request) error) http.Handler {
	var res http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Handling incoming request.", slog.String("uri", r.RequestURI))
		// Update received request call metric
		reqs.Add(1)

		// Remember the call timestamp to calculate the request duration
		start := time.Now()

		// Call the real work
		if err := inner(w, r); err != nil {
			// If an error occurs update the corresponding metric
			errs.Add(1)
			// And put it in a structured log.
			s.logger.Error("Inner handler returned error.", slog.Any("err", err))
		}

		// Calculate the duration of this request
		stop := time.Now()
		dur := int64(stop.Sub(start))

		// Update metrics if needed.
		if dur < durMin.Value() {
			durMin.Set(dur)
		}

		if dur > durMax.Value() {
			durMax.Set(dur)
		}

		s.logger.Info("Finished request handling.", slog.String("uri", r.RequestURI))
	}

	return res
}

func (s Server) notImplementedHandler(w http.ResponseWriter, r *http.Request) error {

	transfer := struct {
		WeatherAPI1 types.FiveDayForecast
		WeatherAPI2 types.FiveDayForecast
		//TODO: Extend with new APIs here
	}{}

	params := r.URL.Query()

	if !params.Has("lat") || !params.Has("lon") {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, MissingParameterErrorDescription)
		return fmt.Errorf("bad request: missing lat or lon params: %+#v", params)
	}

	lat, err := strconv.ParseFloat(params.Get("lat"), 64)
	if err != nil {
		extErr := fmt.Errorf("parse latitude parameter %#v into float: %w", lat, err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, extErr.Error())
		return extErr
	}
	lon, err := strconv.ParseFloat(params.Get("lon"), 64)
	if err != nil {
		extErr := fmt.Errorf("parse longitude parameter %#v into float: %w", lon, err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, extErr.Error())
		return extErr
	}

	aa, err := s.aggregators()
	if err != nil {
		extErr := fmt.Errorf("receiving aggregators: %w", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, extErr.Error())
		return extErr
	}

	for i, a := range aa {
		part, err := a.AggregateWeather(lat, lon)
		if err != nil {
			extErr := fmt.Errorf("Request API %d failed: %+v", i, err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, extErr.Error())
			return extErr
		}
		switch i {
		case 0:
			transfer.WeatherAPI1 = part
		case 1:
			transfer.WeatherAPI2 = part

			//TODO: Extend with new APIs here

		default:
			extErr := fmt.Errorf("Unhandled WeatherAPI ID %d", i)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, extErr.Error())
			return extErr
		}
	}

	data, err := json.Marshal(transfer)
	if err != nil {
		extErr := fmt.Errorf("Marshalling response %+v failed: %+v", transfer, err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, extErr.Error())
		return extErr
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	return err
}
