package gross

// Units stores the Gross Store unit measurements.
func Units() map[string]int {
	units := map[string]int {
        "quarter_of_a_dozen": 3,
        "half_of_a_dozen": 6,
        "dozen": 12,
        "small_gross": 120,
        "gross": 144,
        "great_gross": 1728,
    }
	return units
}

// NewBill creates a new bill.
func NewBill() map[string]int {
	bill := make(map[string]int)
    return bill
}

// AddItem adds an item to customer bill.
func AddItem(bill, units map[string]int, item, unit string) bool {
	if _, unitExists := units[unit]; unitExists {
        if _, itemExists := bill[item]; itemExists {
            bill[item] += units[unit]
        } else {
        	bill[item] = units[unit]
        }
    	return true
    } 	
    return false
}

// RemoveItem removes an item from customer bill.
func RemoveItem(bill, units map[string]int, item, unit string) bool {
	if _, itemExists := bill[item]; itemExists {
		if val2, unitExists := units[unit]; unitExists {
            if bill[item] - val2 > 0 {
                bill[item] -= val2
                return true
            } else if bill[item] - val2 == 0 {
            	delete(bill, item)
                return true
            } else {
            	return false
            }
        }
        return false
    }
	return false
}

// GetItem returns the quantity of an item that the customer has in his/her bill.
func GetItem(bill map[string]int, item string) (int, bool) {
	if val, itemExists := bill[item]; itemExists {
        return val, true
    }
	return 0, false
}
