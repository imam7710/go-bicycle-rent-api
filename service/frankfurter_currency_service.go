package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FrankfurterCurrencyService struct {
	baseURL    string
	httpClient *http.Client
}

type frankfurterResponse struct {
	Date  string  `json:"date"`
	Base  string  `json:"base"`
	Quote string  `json:"quote"`
	Rate  float64 `json:"rate"`
}

func NewFrankfurterCurrencyService() CurrencyService {
	return &FrankfurterCurrencyService{
		baseURL: "https://api.frankfurter.dev/v2",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *FrankfurterCurrencyService) Convert(
	amount float64,
	from string,
	to string,
) (float64, error) {

	if amount < 0 {
		return 0, fmt.Errorf("amount cannot be negative")
	}

	url := fmt.Sprintf(
		"%s/rate/%s/%s",
		s.baseURL,
		from,
		to,
	)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to request currency API: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(
			"currency API returned status %d",
			resp.StatusCode,
		)
	}

	var data frankfurterResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf(
			"failed to decode currency API response: %w",
			err,
		)
	}

	return amount * data.Rate, nil
}
