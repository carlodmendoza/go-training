package elon

import "fmt"

// Drive updates the number of meters driven based on
// the car's speed, and reduces the battery
// according to the battery drainage
func (c *Car) Drive() {
	if c.battery >= c.batteryDrain {
		c.battery -= c.batteryDrain
		c.distance += c.speed
	}
}

// DisplayDistance returns the distance as displayed on
// the LED display as a string
func (c Car) DisplayDistance() string {
	return fmt.Sprintf("Driven %d meters", c.distance)
}

// DisplayBattery returns the battery percentage as displayed
// on the LED display as a string
func (c Car) DisplayBattery() string {
	return fmt.Sprintf("Battery at %d%%", c.battery)
}

// CanFinish takes trackDistance as its parameter and returns
// true if the car can finish the race; otherwise, returns false
func (c Car) CanFinish(trackDistance int) bool {
	var runs, batteryLeft int
	runs = trackDistance / c.speed
	batteryLeft = c.battery - (runs * c.batteryDrain)
	return batteryLeft >= 0
}
