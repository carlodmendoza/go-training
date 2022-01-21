package booking

import ("time" 
        "fmt")

// Schedule returns a time.Time from a string containing a date
func Schedule(date string) time.Time {
    t, _ := time.Parse("1/2/2006 15:04:05", date)
	return t
}

// HasPassed returns whether a date has passed
func HasPassed(date string) bool {
	t, _ := time.Parse("January 2, 2006 15:04:05", date)
	return t.Before(time.Now())
}

// IsAfternoonAppointment returns whether a time is in the afternoon
func IsAfternoonAppointment(date string) bool {
    t, _ := time.Parse("Monday, January 2, 2006 15:04:05", date)
	return t.Hour() >= 12 && t.Hour() < 18
}

// Description returns a formatted string of the appointment time
func Description(date string) string {
	year, month, day := Schedule(date).Date()
    hour, min, _ := Schedule(date).Clock()
    dayOfWeek := Schedule(date).Weekday()
    return fmt.Sprintf("You have an appointment on %s, %s %d, %d, at %d:%d.", dayOfWeek.String(), month.String(), day, year, hour, min)
}

// AnniversaryDate returns a Time with this year's anniversary
func AnniversaryDate() time.Time {
    now := time.Now()
	anniv := time.Date(now.Year(), time.September, 15, 0, 0, 0, 0, time.UTC)
	return anniv
}
