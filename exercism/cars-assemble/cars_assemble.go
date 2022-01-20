package cars

// CalculateWorkingCarsPerHour calculates how many working cars are
// produced by the assembly line every hour
func CalculateWorkingCarsPerHour(productionRate int, successRate float64) float64 {
	return (float64(productionRate) * successRate)/100
}

// CalculateWorkingCarsPerMinute calculates how many working cars are
// produced by the assembly line every minute
func CalculateWorkingCarsPerMinute(productionRate int, successRate float64) int {
    productionRate = productionRate/60
	return int((float64(productionRate) * successRate)/100)
}

// CalculateCost works out the cost of producing the given number of cars
func CalculateCost(carsCount int) uint {
    cost := ((carsCount / 10) * 95000) + ((carsCount % 10) * 10000)
	return uint(cost)
}
