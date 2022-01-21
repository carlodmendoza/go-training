package elon

import "fmt"

// TODO: define the 'Drive()' method
func (c *Car) Drive() {
    if (c.battery >= c.batteryDrain) {
        c.battery -= c.batteryDrain
        c.distance += c.speed
    }
}

// TODO: define the 'DisplayDistance() string' method
func (c Car) DisplayDistance() string {
    return fmt.Sprintf("Driven %d meters", c.distance)
}

// TODO: define the 'DisplayBattery() string' method
func (c Car) DisplayBattery() string {
    return fmt.Sprintf("Battery at %d%%", c.battery)
}

// TODO: define the 'CanFinish(trackDistance int) bool' method
func (c Car) CanFinish(trackDistance int) bool {
    var runs, batteryLeft int
    runs = trackDistance / c.speed
    batteryLeft = c.battery - (runs*c.batteryDrain)
    if batteryLeft >= 0 {
        return true
    } 
	return false
}
