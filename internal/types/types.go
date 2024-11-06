package types

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
