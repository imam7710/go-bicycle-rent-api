package service

import "testing"

func TestFrankfurterCurrencyService_Convert(t *testing.T) {
	service := NewFrankfurterCurrencyService()

	amount, err := service.Convert(
		50000,
		"IDR",
		"USD",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if amount <= 0 {
		t.Fatalf(
			"expected converted amount greater than 0, got %.2f",
			amount,
		)
	}

	t.Logf(
		"IDR 50,000 = USD %.2f",
		amount,
	)
}
