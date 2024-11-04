package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marcofeltmann/weather-forecast-aggregator/internal/api"
)

func TestGetNonExistingEndpoint_ReturnsNotFoundStatus(t *testing.T) {
	sut := api.NewServer(nil)

	srv := httptest.NewServer(sut.Handler())
	c := srv.Client()

	resp, err := c.Get(fmt.Sprintf("%s/nonExistingEndpoint", srv.URL))
	if err != nil {
		t.Error("Request to internal test server without response, aborting.")
		t.Error(err.Error())
		t.Fatal("This is bad. Really bad. Technically it should never happen.")
	}

	got := resp.StatusCode
	want := http.StatusNotFound

	if got != want {
		t.Errorf("Unregistered endpoints must return %s, got %s",
			http.StatusText(want),
			http.StatusText(got),
		)
	}
}

func TestGetExpvarsEndpoint_ReturnsOKStatus(t *testing.T) {
	sut := api.NewServer(nil)

	srv := httptest.NewServer(sut.Handler())
	c := srv.Client()

	resp, err := c.Get(fmt.Sprintf("%s/debug/vars", srv.URL))
	if err != nil {
		t.Errorf("Request to internal test server without response, got %+v.", err)
		t.Fatal("This is bad. Really bad. Technically it should never happen.")
	}

	got := resp.StatusCode
	want := http.StatusOK

	if got != want {
		t.Errorf(
			"Debug endpoint for metrics must return %s, got %s",
			http.StatusText(want),
			http.StatusText(got),
		)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unable to read response data, got %+v", err)
		t.Fatalf("Can't check response integrity, aborting!")
	}
	if len(data) == 0 {
		t.Errorf(
			"Debug endpoint for metrics must return any bytes to be helpful, got %d",
			len(data),
		)
	}
}

func TestGetWeatherEndpointWithoutParameters_ReturnsBadRequestStatus(t *testing.T) {
	sut := api.NewServer(nil)

	srv := httptest.NewServer(sut.Handler())
	c := srv.Client()

	resp, err := c.Get(fmt.Sprintf("%s/weather", srv.URL))
	if err != nil {
		t.Errorf("Request to internal test server without response, got %+v.", err)
		t.Fatal("This is bad. Really bad. Technically it should never happen.")
	}

	got := resp.StatusCode
	want := http.StatusBadRequest

	if got != want {
		t.Errorf(
			"Weather endpoint with missing geocordinates must respond %s, got %s",
			http.StatusText(want),
			http.StatusText(got),
		)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unable to read response data, got %+v", err)
		t.Fatalf("Can't verify response integrity, aborting!")
	}
	if string(data) != api.MissingParameterErrorDescription {
		t.Errorf(
			"Weather endpoint must return reasonable error description, got %#v",
			string(data),
		)
	}
}

func TestGetWeatherEndpoint_ReturnsResult(t *testing.T) {
	sut := api.NewServer(nil)

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

	var result api.Result
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Errorf("Unable to unmarshal API response data, got %+v", err)
		t.Fatal("Can't verify response integrity, aborting!")
	}

	t.Logf("Result: %+v", result)

	t.Fatal("NIY: Output Verification")
}

// Coordinates Boundary tests:
// lat:  -90.0000000000 --  90.000000000
// lon: -180.0000000000 -- 180.000000000

// What if date in the past?!