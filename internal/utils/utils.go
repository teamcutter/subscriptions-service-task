package utils

import (
	"time"
)

func ParseDateFromRequest(dateStr string) (string, error) {
	t, err := time.Parse("01-2006", dateStr)
	if err != nil {
		return "", err
	}

	return t.Format("2006-01-02"), nil
}

func ParseDateFromDB(dateStr string) (string, error) {
	// to cut of T*
	if len(dateStr) >= 10 {
		dateStr = dateStr[:10]
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", err
	}

	return t.Format("01-2006"), nil
}