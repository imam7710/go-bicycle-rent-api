package service

import "testing"

func TestOpenMeteoWeatherService_GetWeather(t *testing.T) {
	service := NewOpenMeteoWeatherService()

	weather, err := service.GetWeather("Jakarta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if weather.City != "Jakarta" {
		t.Fatalf(
			"expected Jakarta, got %s",
			weather.City,
		)
	}

	t.Logf(
		"Temperature: %.2f°C, Description: %s",
		weather.Temperature,
		weather.Description,
	)
}
