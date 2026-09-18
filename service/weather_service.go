package service

type WeatherService interface {
	GetWeather(city string) (WeatherInfo, error)
}

type WeatherInfo struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
}
