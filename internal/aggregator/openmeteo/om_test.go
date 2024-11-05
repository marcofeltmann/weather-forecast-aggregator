package openmeteo_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/marcofeltmann/weather-forecast-aggregator/internal/aggregator/openmeteo"
)

func TestOpenMeteoAggregation_HistoricalData(t *testing.T) {
	lat, lon := 42.6493934, -8.8201753

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

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

	t.Logf("Got: %#v", got)

	t.Fatal("NYI")
}
