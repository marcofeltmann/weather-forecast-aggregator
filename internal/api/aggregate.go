package api

import "github.com/marcofeltmann/weather-forecast-aggregator/internal/types"

/*
Aggregator interface implementation is the input port for weather data from
any resource like OpenMeteo, WeatherAPI, OpenWeatherMap, MeteoGalicia and others.
*/
type Aggregator interface {
	AggregateWeather(lat, lon float64) (types.FiveDayForecast, error)
}
