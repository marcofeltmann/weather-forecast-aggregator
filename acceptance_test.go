package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/marcofeltmann/weather-forecast-aggregator/internal/api"
	"github.com/marcofeltmann/weather-forecast-aggregator/internal/types"
)

func TestGetWeatherEndpoint_ReturnsResult(t *testing.T) {
	t.Log("This test runs against the real forecast endpoints.")
	t.Log("So it might break every single day.")
	// Technical Debt: This runs a test with the default HTTP client against the
	// real endpoint as I cannot inject some pre-configured http.Client.
	// So the response data will change at least daily, maybe even within the day
	// as forecasts get updated.
	// Using another net/httptest server for reproducable responses would be better.
	key, err := apiKey()
	if err != nil {
		t.Errorf("get API key: %+v", err)
		t.Fatal("Aborting")
	}
	sut := api.NewServer(api.Config{WeatherApiKey: key})

	srv := httptest.NewServer(sut.Handler())
	c := srv.Client()

	resp, err := c.Get(fmt.Sprintf("%s/weather?lat=42.6493934&lon=-8.8201753", srv.URL))
	if err != nil {
		t.Errorf("Request to internal test server without response, got %+v.", err)
		t.Fatal("This is bad. Really bad. Technically it should never happen.")
	}

	got := resp.StatusCode
	want := http.StatusOK

	if got != want {
		t.Errorf(
			"Weather endpoint with valid geocordinates must respond %s, got %s",
			http.StatusText(want),
			http.StatusText(got),
		)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unable to read response data, got %+v", err)
		t.Fatal("Can't verify response integrity, aborting!")
	}

	var result types.Result
	if err := json.Unmarshal(data, &result); err != nil {
		t.Errorf("Unable to unmarshal API response data, got %+v", err)
		t.Fatal("Can't verify response integrity, aborting!")
	}

	expected := types.Result{
		WeatherAPI1: types.FiveDayForecast{
			Day1: types.Forecast{
				Date:    "2024-11-05",
				MaxTemp: 21.9,
			},
			Day2: types.Forecast{
				Date:    "2024-11-06",
				MaxTemp: 21,
			},
			Day3: types.Forecast{
				Date:    "2024-11-07",
				MaxTemp: 22.6,
			},
			Day4: types.Forecast{
				Date:    "2024-11-08",
				MaxTemp: 18.3,
			},
			Day5: types.Forecast{
				Date:    "2024-11-09",
				MaxTemp: 18.2,
			},
		},
		WeatherAPI2: types.FiveDayForecast{
			Day1: types.Forecast{
				Date:    "2024-11-05",
				MaxTemp: 19.8,
			},
			Day2: types.Forecast{
				Date:    "2024-11-06",
				MaxTemp: 22.1,
			},
			Day3: types.Forecast{
				Date:    "2024-11-07",
				MaxTemp: 21,
			},
			Day4: types.Forecast{
				Date:    "2024-11-08",
				MaxTemp: 19,
			},
			Day5: types.Forecast{
				Date:    "2024-11-09",
				MaxTemp: 19,
			},
		},
	}

	if !cmp.Equal(expected, result) {
		fmt.Println(cmp.Diff(expected, result))
		t.Error("output mismatch, see diff for details")
	}
}

func apiKey() (string, error) {
	res := os.Getenv("WAPI_KEY")
	var empty string

	if res == empty {
		return empty, errors.New("envvar WAPI_KEY has empty value")
	}

	return res, nil
}
