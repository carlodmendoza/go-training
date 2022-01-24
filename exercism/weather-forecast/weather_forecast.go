// Package weather outputs the weather condition of any given location.
package weather

// CurrentCondition stores the weather condition.
var CurrentCondition string

// CurrentLocation stores the given location.
var CurrentLocation string

// Forecast returns the weather condition of the chosen location.
func Forecast(city, condition string) string {
	CurrentLocation, CurrentCondition = city, condition
	return CurrentLocation + " - current weather condition: " + CurrentCondition
}
