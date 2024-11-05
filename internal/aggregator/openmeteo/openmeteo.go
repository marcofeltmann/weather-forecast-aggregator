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

type Caller struct {
	clock  func() time.Time
	ctx    context.Context
	client *http.Client
}

func DefaultCaller() Caller {
	return ConfiguredCaller(
		context.Background(),
		&http.Client{},
		time.Now,
	)
}

func ConfiguredCaller(ctx context.Context, client *http.Client, tf func() time.Time) Caller {
	if client == nil {
		panic(errors.New("http client is required for configured caller as we do http requests"))
	}
	return Caller{
		clock:  tf,
		ctx:    ctx,
		client: client,
	}
}

/*
The result for a single day request looks like that:

	{
	  "latitude": 42.5625,
	  "longitude": -8.8125,
	  "generationtime_ms": 0.033020973205566406,
	  "utc_offset_seconds": 0,
	  "timezone": "GMT",
	  "timezone_abbreviation": "GMT",
	  "elevation": 2,
	  "daily_units": {
	    "time": "iso8601",
	    "temperature_2m_max": "°C"
	  },
	  "daily": {
	    "time": [
	      "2024-11-09"
	    ],
	    "temperature_2m_max": [
	      20.6
	    ]
	  }
	}

Guess I'm only interested in the time and temperature array, but I don't need
two arrays for this.
*/
type result struct {
	Time    time.Time
	MaxTemp float32
}
type wrapper struct {
	Daily forecast `json:"daily"`
}
type forecast struct {
	Time    []string  `json:"time"`
	MaxTemp []float32 `json:"temperature_2m_max"`
}

const daysToFetch int = 5

/*
AggegrateWeather aggregates weather data for five days from whatever the internal clock
decides is today.
*/
func (c Caller) AggregateWeather(lat, lon float64) (types.FiveDayForecast, error) {
	results := make(chan result)

	g := new(errgroup.Group)
	requestUrls := urlsToFetchIncluding(c.clock(), daysToFetch, lat, lon)

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

			d, err := time.Parse(time.DateOnly, tmp.Daily.Time[0])
			if err != nil {
				return fmt.Errorf("Generate time from %s failed: %w", tmp.Daily.Time[0], err)
			}
			t := tmp.Daily.MaxTemp[0]
			res := result{Time: d, MaxTemp: t}

			results <- res
		}
		return nil
	})

	go func() {
		if err := g.Wait(); err != nil {
			slog.Default().Info("Downloading data failed.", slog.Any("err", err))
		}
	}()

	var res types.FiveDayForecast
	g.Go(func() error {
		rr := make([]result, 0, 5)

		// I guess this is only working because of luck, as response of day 5 might
		// be downloaded way before response of day 1 due to internet latency
		for r := range results {
			rr = append(rr, r)
		}

		res.Day1 = types.Forecast{Date: rr[0].Time.Format(time.DateOnly), MaxTemp: rr[0].MaxTemp}
		res.Day2 = types.Forecast{Date: rr[1].Time.Format(time.DateOnly), MaxTemp: rr[1].MaxTemp}
		res.Day3 = types.Forecast{Date: rr[2].Time.Format(time.DateOnly), MaxTemp: rr[2].MaxTemp}
		res.Day4 = types.Forecast{Date: rr[3].Time.Format(time.DateOnly), MaxTemp: rr[3].MaxTemp}
		res.Day5 = types.Forecast{Date: rr[4].Time.Format(time.DateOnly), MaxTemp: rr[4].MaxTemp}

		return nil
	})

	// Wait for all downloaded data to be converted.
	if err := g.Wait(); err != nil {
		slog.Default().Error("Converting failed.", slog.Any("err", err))
	}

	return res, nil
}

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
