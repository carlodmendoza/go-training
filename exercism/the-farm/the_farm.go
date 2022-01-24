package thefarm

import (
	"errors"
	"fmt"
)

// See types.go for the types defined for this exercise.

type SillyNephewError struct {
	cows int
}

// Error is a method of SillyNephewError that returns a message
// in case the number of cows is negative.
func (e *SillyNephewError) Error() string {
	return fmt.Sprintf("silly nephew, there cannot be %d cows", e.cows)
}

// DivideFood computes the fodder amount per cow for the given cows.
func DivideFood(weightFodder WeightFodder, cows int) (float64, error) {
	amt, err := weightFodder.FodderAmount()
	if err == ErrScaleMalfunction && amt > 0 {
		amt *= 2
		return amt / float64(cows), nil
	} else if amt < 0 {
		if err == ErrScaleMalfunction || err == nil {
			return 0, errors.New("negative fodder")
		}
		return 0, err
	} else if cows == 0 {
		return 0, errors.New("division by zero")
	} else if cows < 0 {
		return 0, &SillyNephewError{cows: cows}
	} else if err != nil {
		return 0, err
	} else {
		return amt / float64(cows), nil
	}
}
