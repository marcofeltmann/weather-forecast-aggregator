package api

import (
	"encoding/json"
	"errors"
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

const ParameterOutOfBoundsErrorDescription = `Parameters out of bounds.
lat must be within -90 and 90, lon must be within -180 and 180.`

// addRoutes manages the endpoint routing for the API server.
// Having all the information at one place might make navigating through the
// server easier, as all the info you get is "this endpoint didn't work!"
func (s Server) addRoutes() {
	mux := s.mux

	// According to https://go.dev/blog/routing-enhancements the GET method also
	// handles the HEAD method, so we don't need the HEAD routes.
	mux.Handle("GET /weather", s.meterMiddleware(s.dataAggregation))

	mux.Handle("GET /debug/vars", expvar.Handler())
}

// meterMiddleware returns an http.Handler that wraps an error-aware pseudo-handler.
// It's doing the RED metrics via expvar, but shouldn't interfere with the http
// requests and responses.
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

// dataAggregation verifies the request parameters, hooks up the aggregators and
// responses with the aggregated data, or with error status and messages.
func (s Server) dataAggregation(w http.ResponseWriter, r *http.Request) error {
	transfer := make(map[string]types.FiveDayForecast)

	params := r.URL.Query()
	pLat := params.Get("lat")
	pLon := params.Get("lon")

	var empty string
	if empty == pLat || empty == pLon {
		err := errors.New(MissingParameterErrorDescription)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return err
	}

	lat, err := strconv.ParseFloat(pLat, 64)
	if err != nil {
		extErr := fmt.Errorf("parse latitude parameter %#v into float: %w", lat, err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, extErr.Error())
		return extErr
	}
	lon, err := strconv.ParseFloat(pLon, 64)
	if err != nil {
		extErr := fmt.Errorf("parse longitude parameter %#v into float: %w", lon, err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, extErr.Error())
		return extErr
	}

	if !saneInputs(lat, lon) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, ParameterOutOfBoundsErrorDescription)
		return fmt.Errorf("bad request: lat or lon params out of bounds: %+#v", params)
	}

	aa, err := s.aggregators()
	if err != nil {
		extErr := fmt.Errorf("receiving aggregators: %w", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, extErr.Error())
		return extErr
	}

	for i, a := range aa {
		key := fmt.Sprintf("weatherAPI%d", i)
		part, err := a.AggregateWeather(lat, lon)
		if err != nil {
			extErr := fmt.Errorf("Request API %d failed: %+v", i, err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, extErr.Error())
			return extErr
		}

		transfer[key] = part
	}

	data, err := json.Marshal(transfer)
	if err != nil {
		extErr := fmt.Errorf("Marshalling response %+v failed: %+v", transfer, err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, extErr.Error())
		return extErr
	}

	// TIL: w.Header().Add(...) must be called right before w.WriteHeader()
	w.Header().Add("Content-Type", "application/json")
	// w.Write implicitely calls w.WriteHeader(http.StatusOK) before writing data
	_, err = w.Write(data)
	return err
}

func saneInputs(lat, lon float64) bool {
	switch {
	case lat < -90.000000:
		return false
	case lat > 90.000000:
		return false
	case lon < -180.000000:
		return false
	case lon > 180.000000:
		return false
	default:
		return true
	}
}
