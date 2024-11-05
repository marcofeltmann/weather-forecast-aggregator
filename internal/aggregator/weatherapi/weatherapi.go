package weatherapi

import (
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

var ErrNoApiKeyProvided = errors.New("API key missing")

// Caller shall implement the api.Aggregator interface to make WeatherAPI calls.
type Caller struct {
	client *http.Client
	apikey string
	clock  func() time.Time
}

// NewCaller creates a pre-configured caller that uses the provided API key.
func NewCaller(apikey string) (*Caller, error) {
	return DebuggingCaller(apikey, &http.Client{}, time.Now)
}

// DebuggingCaller lets inject non-default implementation for testing and
// debugging sessions.
func DebuggingCaller(key string, c *http.Client, cf func() time.Time) (*Caller, error) {
	var empty string
	if empty == key {
		return nil, ErrNoApiKeyProvided
	}
	return &Caller{
		apikey: key,
		client: c,
		clock:  cf,
	}, nil
}

const daysToFetch = 5

// wrapper is the upper data structure of WeatherAPI result.
// The whole structure is here to unmarshal the received JSON into.
type wrapper struct {
	Forecast collection `json:"forecast"`
}

// collection holds an array of forecasts in the WeatherAPI result.
type collection struct {
	ForecastDay []forecast `json:"forecastDay"`
}

// forecast contains the date and the grouped forecasts of WeatherAPI result.
type forecast struct {
	Date string `json:"date"`
	Day  day    `json:"day"`
}

// day finally contains all the forecast information for the given day.
type day struct {
	MaxTemp float32 `json:"maxtemp_c"`
}

// AggregateWeather implements the api.Aggregator interface on Caller
func (c *Caller) AggregateWeather(lat, lon float64) (types.FiveDayForecast, error) {
	results := make(chan types.Forecast)

	g := new(errgroup.Group)
	requestUrls := c.urlsToFetchIncluding(c.clock(), daysToFetch, lat, lon)

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
				return fmt.Errorf("GET %s unexpected status, want %d, got %d", u.String(), http.StatusOK, resp.StatusCode)
			}

			bb, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("Reads response data from %s failed: %w", u.String(), err)
			}

			var tmp wrapper
			if err = json.Unmarshal(bb, &tmp); err != nil {
				return fmt.Errorf("Unmarshal response data from %s failed: %w", u.String(), err)
			}

			d := tmp.Forecast.ForecastDay[0].Date
			t := tmp.Forecast.ForecastDay[0].Day.MaxTemp
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
func (c *Caller) urlsToFetchIncluding(d time.Time, amount int, lat, lon float64) []url.URL {
	res := make([]url.URL, 0, amount)

	latest := d

	for i := 0; i < amount; i++ {
		date := latest.Format(time.DateOnly)

		u := url.URL{
			Scheme: "https",
			Host:   "api.weatherapi.com",
			Path:   "/v1/forecast.json",
			RawQuery: fmt.Sprintf(
				"key=%s&q=%f,%f&date=%s&day=maxtemp_c",
				c.apikey, lat, lon, date,
			),
		}
		res = append(res, u)

		latest = latest.AddDate(0, 0, 1)
	}

	return res
}
