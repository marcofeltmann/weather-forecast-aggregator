package weatherapi_test

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	openweathermap "github.com/marcofeltmann/weather-forecast-aggregator/internal/aggregator/weatherapi"
	"github.com/marcofeltmann/weather-forecast-aggregator/internal/types"
)

func TestHappyPath(t *testing.T) {
	lat, lon := 42.6493934, -8.8201753

	key, err := apiKey()
	if err != nil {
		t.Fatalf("OpenWeatherMap API Key retrieval: %+v", err)
	}

	// Technical debt: This uses the default http.Client and calls the API endpoint.
	// As nobody knows how long the historical data is stored there might be false
	// negatives in the future.
	sut := openweathermap.DebuggingCaller(key, &http.Client{}, time.Now)

	got, err := sut.AggregateWeather(lat, lon)
	if err != nil {
		t.Errorf("aggregate: %+v", err)
		t.Fatal("Cannot verify result, aborting.")
	}

	want := types.FiveDayForecast{
		Day1: types.Forecast{Date: "2024-10-25", MaxTemp: 14.6},
		Day2: types.Forecast{Date: "2024-10-26", MaxTemp: 14.9},
		Day3: types.Forecast{Date: "2024-10-27", MaxTemp: 18.2},
		Day4: types.Forecast{Date: "2024-10-28", MaxTemp: 21.2},
		Day5: types.Forecast{Date: "2024-10-29", MaxTemp: 22.3},
	}

	if !cmp.Equal(want, got) {
		fmt.Println(cmp.Diff(want, got))
		t.Error("output mismatch, see diff")
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
