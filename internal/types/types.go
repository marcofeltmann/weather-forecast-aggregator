package types

type Result struct {
	// No clue how to dynamically weatherAPI1/weatherAPI2 as key for the value
	WeatherAPI1 FiveDayForecast
	WeatherAPI2 FiveDayForecast
}

type FiveDayForecast struct {
	Day1 Forecast
	Day2 Forecast
	Day3 Forecast
	Day4 Forecast
	Day5 Forecast
}

type Forecast struct {
	Date    string
	MaxTemp float32
}
