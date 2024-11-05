package openmeteo_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/marcofeltmann/weather-forecast-aggregator/internal/aggregator/openmeteo"
	"github.com/marcofeltmann/weather-forecast-aggregator/internal/types"
)

func TestOpenMeteoAggregation_HistoricalData(t *testing.T) {
	lat, lon := 42.6493934, -8.8201753

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	// Technical Debt: This runs a test with the default HTTP client against the
	// real endpoint.
	// Using an net/httptest server for reproducable responses would be better.
	// Currently the historical data is still fetched, so it looks good enough.
	sut := openmeteo.ConfiguredCaller(ctx, &http.Client{}, func() time.Time {
		res, err := time.Parse(time.DateOnly, "2024-10-25")
		if err != nil {
			t.Fatalf("Cannot test hard-coded past value, got %+v", err)
		}
		return res
	})

	got, err := sut.AggregateWeather(lat, lon)
	if err != nil {
		t.Fatalf(
			"Error while aggregate from OpenMeteo for lat %.8f, lon %.8f, got: %+v",
			lat, lon, err,
		)
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
