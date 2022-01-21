package lasagna

// TODO: define the 'PreparationTime()' function
func PreparationTime(layers []string, time int) int {
    if time == 0 {
        time = 2
    }
	return len(layers) * time
}

// TODO: define the 'Quantities()' function
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

// TODO: define the 'AddSecretIngredient()' function
func AddSecretIngredient(friendsList, myList []string) {
    myList[len(myList)-1] = friendsList[len(friendsList)-1]
}

// TODO: define the 'ScaleRecipe()' function
func ScaleRecipe(recipe []float64, portions int) []float64 {
    newRecipe := make([]float64, len(recipe))
    for i, qty := range recipe {
        newRecipe[i] = qty * float64(portions)/2
    }
	return newRecipe
}
