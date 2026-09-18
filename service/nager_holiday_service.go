package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type NagerHolidayService struct {
	baseURL    string
	httpClient *http.Client
	country    string
}

type nagerHolidayResponse struct {
	Date            string   `json:"date"`
	Name            string   `json:"name"`
	CountryCode     string   `json:"countryCode"`
	NationalHoliday bool     `json:"nationalHoliday"`
	HolidayTypes    []string `json:"holidayTypes"`
}

func NewNagerHolidayService() HolidayService {
	return &NagerHolidayService{
		baseURL: "https://date.nager.at/api/v4",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		country: "ID",
	}
}

func (s *NagerHolidayService) IsHoliday(date time.Time) (bool, error) {
	url := fmt.Sprintf(
		"%s/Holidays/%s/%d",
		s.baseURL,
		s.country,
		date.Year(),
	)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return false, fmt.Errorf(
			"failed to request holiday API: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf(
			"holiday API returned status %d",
			resp.StatusCode,
		)
	}

	var holidays []nagerHolidayResponse

	if err := json.NewDecoder(resp.Body).Decode(&holidays); err != nil {
		return false, fmt.Errorf(
			"failed to decode holiday API response: %w",
			err,
		)
	}

	requestedDate := date.Format("2006-01-02")

	for _, holiday := range holidays {
		if holiday.Date == requestedDate &&
			holiday.NationalHoliday {
			return true, nil
		}
	}

	return false, nil
}
