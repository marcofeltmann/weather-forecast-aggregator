package types

// Result is the data model for the API server response.
// I'm trying to flatlining it a little.
type Result struct {
	// No clue how to dynamically set weatherAPI1/weatherAPI2 keys for JSON
	WeatherAPI1 FiveDayForecast
	WeatherAPI2 FiveDayForecast
	//TODO: Add new APIs here
}

// FiveDayForecast holds the forecasts for five days.
// Having it identified with the name makes it a little easier to handle instead
// of using arrays.
type FiveDayForecast struct {
	Day1 Forecast
	Day2 Forecast
	Day3 Forecast
	Day4 Forecast
	Day5 Forecast
}

// Forecast holds the corresponding date and maximum temperature in ºC
// There is more to weather than that, but this is good enough for a quick start.
type Forecast struct {
	Date    string
	MaxTemp float32
}
