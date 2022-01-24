package lasagna

// PreparationTime accepts a slice of layers as a []string and
// the average preparation time per layer in minutes as an int.
// It returns the estimate for the total preparation time based
// on the number of layers as an int.
func PreparationTime(layers []string, time int) int {
	if time == 0 {
		time = 2
	}
	return len(layers) * time
}

// Quantities takes a slice of layers as parameter as a []string.
// It returns the quantity of noodles and sauce needed to make your meal.
func Quantities(layers []string) (int, float64) {
	var noodles int
	var sauce float64
	for _, layer := range layers {
		if layer == "noodles" {
			noodles++
		} else if layer == "sauce" {
			sauce++
		}
	}
	noodles *= 50
	sauce *= 0.2
	return noodles, sauce
}

// AddSecretIngredient accepts two slices of ingredients of type []string
// as parameters. The first parameter is the list your friend sent you,
// the second is the ingredient list of your own recipe. This replaces the
// last item in your list with the last item from your friends list.
func AddSecretIngredient(friendsList, myList []string) {
	myList[len(myList)-1] = friendsList[len(friendsList)-1]
}

// ScaleRecipe returns a slice of float64 of the amounts needed for
// the desired number of portions.
func ScaleRecipe(recipe []float64, portions int) []float64 {
	newRecipe := make([]float64, len(recipe))
	for i, qty := range recipe {
		newRecipe[i] = qty * float64(portions) / 2
	}
	return newRecipe
}
