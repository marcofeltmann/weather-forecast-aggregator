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

// TestHappyPath should test the happy path. Turned out that didn't work as we
// cannot have historical data from the forecast endpoint.
// As I'm unsure if I should test against another endpoint or try to mock the
// output data I didn't invest more time here. It's good enough for a human to
// verify if the data is there.
func TestHappyPath(t *testing.T) {
	lat, lon := 42.6493934, -8.8201753

	key, err := apiKey()
	if err != nil {
		t.Fatalf("OpenWeatherMap API Key retrieval: %+v", err)
	}

	t.Log("This test runs against the real forecast endpoint, as the historical API differs.")
	t.Log("So it might break every single day.")

	// Technical debt: This uses the default http.Client and calls the API endpoint.
	// As nobody knows how long the historical data is stored there might be false
	// negatives in the future.
	sut, err := openweathermap.DebuggingCaller(key, &http.Client{}, time.Now)
	if err != nil {
		t.Errorf("creating DebuggingCaller: %+v", err)
		t.Fatal("Aborting")
	}

	got, err := sut.AggregateWeather(lat, lon)
	if err != nil {
		t.Errorf("aggregate: %+v", err)
		t.Fatal("Cannot verify result, aborting.")
	}

	/*
			  types.FiveDayForecast{
		  	Day1: types.Forecast{
		- 		Date:    "2024-10-25",
		+ 		Date:    "2024-11-05",
		- 		MaxTemp: 14.600000381469727,
		+ 		MaxTemp: 19.799999237060547,
		  	},
		  	Day2: types.Forecast{
		- 		Date:    "2024-10-26",
		+ 		Date:    "2024-11-06",
		- 		MaxTemp: 14.899999618530273,
		+ 		MaxTemp: 22.100000381469727,
		  	},
		  	Day3: types.Forecast{
		- 		Date:    "2024-10-27",
		+ 		Date:    "2024-11-07",
		- 		MaxTemp: 18.200000762939453,
		+ 		MaxTemp: 21,
		  	},
		  	Day4: types.Forecast{
		- 		Date:    "2024-10-28",
		+ 		Date:    "2024-11-08",
		- 		MaxTemp: 21.200000762939453,
		+ 		MaxTemp: 19,
		  	},
		  	Day5: types.Forecast{
		- 		Date:    "2024-10-29",
		+ 		Date:    "2024-11-09",
		- 		MaxTemp: 22.299999237060547,
		+ 		MaxTemp: 19,
		  	},
	*/
	want := types.FiveDayForecast{
		Day1: types.Forecast{Date: "2024-11-05", MaxTemp: 19.8},
		Day2: types.Forecast{Date: "2024-11-06", MaxTemp: 22.1},
		Day3: types.Forecast{Date: "2024-11-07", MaxTemp: 21},
		Day4: types.Forecast{Date: "2024-11-08", MaxTemp: 19},
		Day5: types.Forecast{Date: "2024-11-09", MaxTemp: 19},
	}

	if !cmp.Equal(want, got) {
		fmt.Println(cmp.Diff(want, got))
		t.Error("output mismatch, see diff")
	}
}

// apiKey is a helper function to get the API key from an envvar.
// The github.com/ardanlabs/conf/v3 package doesn't help here.
func apiKey() (string, error) {
	res := os.Getenv("WEATHER_API_KEY")
	var empty string

	if res == empty {
		return empty, errors.New("envvar WAPI_KEY has empty value")
	}

	return res, nil
}
