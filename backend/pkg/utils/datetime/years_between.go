package datetime

import "time"

func YearBetween(from, to time.Time) int {
	years := to.Year() - from.Year()

	if to.Month() < from.Month() || (to.Month() == from.Month() && to.Day() < from.Day()) {
		years--
	}

	return years
}
