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

type Caller struct {
	client *http.Client
	apikey string
	clock  func() time.Time
}

func DefaultCaller(apikey string) (*Caller, error) {
	return DebuggingCaller(apikey, &http.Client{}, time.Now)
}

func DebuggingCaller(key string, c *http.Client, cf func() time.Time) (*Caller, error) {
	fmt.Println(key)
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

type result struct {
	Time    time.Time
	MaxTemp float32
}

type wrapper struct {
	Forecast collection `json:"forecast"`
}

type collection struct {
	ForecastDay []forecast `json:"forecastDay"`
}

type forecast struct {
	Date string `json:"date"`
	Day  day    `json:"day"`
}

type day struct {
	MaxTemp float32 `json:"maxtemp_c"`
}

func (c *Caller) AggregateWeather(lat, lon float64) (types.FiveDayForecast, error) {
	results := make(chan result)

	g := new(errgroup.Group)
	requestUrls := c.urlsToFetchIncluding(c.clock(), daysToFetch, lat, lon)

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

			d, err := time.Parse(time.DateOnly, tmp.Forecast.ForecastDay[0].Date)
			if err != nil {
				return fmt.Errorf("Generate time from %s failed: %w", tmp.Forecast.ForecastDay[0].Date, err)
			}
			t := tmp.Forecast.ForecastDay[0].Day.MaxTemp
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
