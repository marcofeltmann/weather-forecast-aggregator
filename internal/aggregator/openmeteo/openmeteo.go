package openmeteo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/marcofeltmann/weather-forecast-aggregator/internal/types"
	"golang.org/x/sync/errgroup"
)

// Caller shall implement the api.Aggregator interface to call the OpenMeteo API.
type Caller struct {
	clock  func() time.Time
	ctx    context.Context
	client *http.Client
}

// NewCaller creates a pre-configured OpenMeteo API caller.
func NewCaller() *Caller {
	return DebuggingCaller(
		context.Background(),
		&http.Client{},
		time.Now,
	)
}

// DebuggingCaller let define some specific types for the internal structure.
// This makes it useful for testing or debugging sessions.
func DebuggingCaller(ctx context.Context, client *http.Client, tf func() time.Time) *Caller {
	if client == nil {
		// It is said to be bad style panicking out of a package. I agree.
		// Since we're inside of the business layer and introducing an error for one
		// constructor that's supposed to be used only by developers I think it's fine
		panic(errors.New("http client is required for configured caller as we do http requests"))
	}
	return &Caller{
		clock:  tf,
		ctx:    ctx,
		client: client,
	}
}

// wrapper reflects the top level object of the API response.
// The hierarchy is used to unmarshal JSON in it for easier access.
type wrapper struct {
	Daily forecast `json:"daily"`
}

// forecast reflects the forecast data of the API response.
type forecast struct {
	Time    []string  `json:"time"`
	MaxTemp []float32 `json:"temperature_2m_max"`
}

const daysToFetch int = 5

// AggegrateWeather implements the api.Aggregator interface for the OpenMeteo API
func (c Caller) AggregateWeather(lat, lon float64) (types.FiveDayForecast, error) {
	results := make(chan types.Forecast)

	g := new(errgroup.Group)
	requestUrls := urlsToFetchIncluding(c.clock(), daysToFetch, lat, lon)

	// First error-prone go routine:
	// Request one URL after the other and put the result into the results channel
	g.Go(func() error {
		defer close(results)
		for _, u := range requestUrls {
			resp, err := c.client.Get(u.String())
			if err != nil {
				return fmt.Errorf("Get %s failed: %w", u.String(), err)
			}

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf(
					"GET %s unexpected status, want %d, got %d",
					u.String(), http.StatusOK, resp.StatusCode,
				)
			}

			bb, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf(
					"Reads response data from %s failed: %w",
					u.String(), err,
				)
			}

			var tmp wrapper
			if err = json.Unmarshal(bb, &tmp); err != nil {
				return fmt.Errorf(
					"Unmarshal response data from %s failed: %w",
					u.String(), err,
				)
			}

			d := tmp.Daily.Time[0]
			t := tmp.Daily.MaxTemp[0]
			res := types.Forecast{Date: d, MaxTemp: t}

			results <- res
		}
		return nil
	})

	// error groups work like wait groups: They wait until the provided function
	// returns. In addition they keep track of all errors that occured.
	go func() {
		if err := g.Wait(); err != nil {
			slog.Default().Info("Downloading data failed.", slog.Any("err", err))
		}
	}()

	var res types.FiveDayForecast
	// Second Go Routine converting the downloaded results into the exchange format
	// with 5 forecasts in one struct.
	g.Go(func() error {
		rr := make([]types.Forecast, 0, 5)

		for r := range results {
			rr = append(rr, r)
		}

		res.Day1 = types.Forecast{Date: rr[0].Date, MaxTemp: rr[0].MaxTemp}
		res.Day2 = types.Forecast{Date: rr[1].Date, MaxTemp: rr[1].MaxTemp}
		res.Day3 = types.Forecast{Date: rr[2].Date, MaxTemp: rr[2].MaxTemp}
		res.Day4 = types.Forecast{Date: rr[3].Date, MaxTemp: rr[3].MaxTemp}
		res.Day5 = types.Forecast{Date: rr[4].Date, MaxTemp: rr[4].MaxTemp}

		return nil
	})

	// Synchronously wait for all downloaded data to be converted.
	if err := g.Wait(); err != nil {
		slog.Default().Error("Converting failed.", slog.Any("err", err))
	}

	return res, nil
}

// urlsToFetchIncluding helps to generate the requested amount of API endpoint
// URLs with the provided start date `d`, counting one day up `amount` times.
func urlsToFetchIncluding(d time.Time, amount int, lat, lon float64) []url.URL {
	res := make([]url.URL, 0, amount)

	latest := d

	for i := 0; i < amount; i++ {
		date := latest.Format(time.DateOnly)

		u := url.URL{
			Scheme: "https",
			Host:   "api.open-meteo.com",
			Path:   "/v1/forecast",
			RawQuery: fmt.Sprintf(
				"latitude=%.6f&longitude=%.6f&start_date=%s&end_date=%s&daily=temperature_2m_max",
				lat, lon, date, date,
			),
		}
		res = append(res, u)

		latest = latest.AddDate(0, 0, 1)
	}

	return res
}
