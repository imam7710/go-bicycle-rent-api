package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type OpenMeteoWeatherService struct {
	baseURL    string
	httpClient *http.Client
}

type openMeteoResponse struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
		WeatherCode int     `json:"weather_code"`
	} `json:"current"`
}

func NewOpenMeteoWeatherService() WeatherService {
	return &OpenMeteoWeatherService{
		baseURL: "https://api.open-meteo.com/v1/forecast",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *OpenMeteoWeatherService) GetWeather(
	city string,
) (WeatherInfo, error) {

	latitude := "-6.2088"
	longitude := "106.8456"

	url := fmt.Sprintf(
		"%s?latitude=%s&longitude=%s&current=temperature_2m,weather_code",
		s.baseURL,
		latitude,
		longitude,
	)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return WeatherInfo{}, fmt.Errorf(
			"failed to request weather API: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WeatherInfo{}, fmt.Errorf(
			"weather API returned status %d",
			resp.StatusCode,
		)
	}

	var data openMeteoResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return WeatherInfo{}, fmt.Errorf(
			"failed to decode weather API response: %w",
			err,
		)
	}

	return WeatherInfo{
		City:        city,
		Temperature: data.Current.Temperature,
		Description: weatherDescription(data.Current.WeatherCode),
	}, nil
}

func weatherDescription(code int) string {
	switch {
	case code == 0:
		return "Clear sky"
	case code >= 1 && code <= 3:
		return "Partly cloudy"
	case code >= 45 && code <= 48:
		return "Fog"
	case code >= 51 && code <= 57:
		return "Drizzle"
	case code >= 61 && code <= 67:
		return "Rain"
	case code >= 71 && code <= 77:
		return "Snow"
	case code >= 80 && code <= 82:
		return "Rain showers"
	case code >= 95:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}
