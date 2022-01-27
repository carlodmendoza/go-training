package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	fmt.Println("Server running in port 8080")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

// handler prints to a writer the required output
// or error depending on the input
func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/solve" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "Failed to parse form data: %s\n", err)
		}
		if val, ok := r.Form["coef"]; ok {
			m := [3][3]float64{}
			var c1, c2, c3 float64
			if n, _ := fmt.Sscanf(val[0], "%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,",
				&m[0][0], &m[0][1], &m[0][2], &c1,
				&m[1][0], &m[1][1], &m[1][2], &c2,
				&m[2][0], &m[2][1], &m[2][2], &c3,
			); n == 12 {
				fmt.Fprintln(w, computeThreeUnknowns(m, c1, c2, c3))
			} else {
				fmt.Fprintln(w, "Invalid format or incorrect number of coefficients")
			}
		} else {
			fmt.Fprintln(w, "'coef' parameter not found")
		}
	} else {
		fmt.Fprintln(w, "Invalid URL")
	}
}

// getDeterminantOf2x2 returns the determinant of a 2x2 matrix
func getDeterminantOf2x2(a, b, c, d float64) float64 {
	return a*d - b*c
}

// getDeterminantOf3x3 returns the determinant of a 3x3 matrix
func getDeterminantOf3x3(m [3][3]float64) float64 {
	return m[0][0]*getDeterminantOf2x2(m[1][1], m[2][1], m[1][2], m[2][2]) -
		m[1][0]*getDeterminantOf2x2(m[0][1], m[2][1], m[0][2], m[2][2]) +
		m[2][0]*getDeterminantOf2x2(m[0][1], m[1][1], m[0][2], m[1][2])
}

// replaceColIn3x3 replaces the values of a column of 3x3 matrix
// given the column index, three numbers, and the matrix itself.
// It returns a copy of the new matrix
func replaceColIn3x3(m [3][3]float64, colIndex int, c1, c2, c3 float64) [3][3]float64 {
	var newM [3][3]float64
	for i := range newM {
		for j := range newM[i] {
			switch {
			case i == 0 && j == colIndex:
				newM[i][j] = c1
			case i == 1 && j == colIndex:
				newM[i][j] = c2
			case i == 2 && j == colIndex:
				newM[i][j] = c3
			default:
				newM[i][j] = m[i][j]
			}
		}
	}
	return newM
}

// fmtFloat returns the formatted string of a
// given floating-point number
func fmtFloat(num float64) string {
	if num == float64(int(num)) {
		return strconv.Itoa(int(num))
	}
	return strconv.FormatFloat(num, 'f', 2, 64)
}

// computeThreeUnknowns solves a system of three equations
// with three variables using Cramer's Rule. It returns an output
// string which has the system of equations and its solution.
func computeThreeUnknowns(m [3][3]float64, c1, c2, c3 float64) string {
	var D, Dx, Dy, Dz, x, y, z float64
	var output string
	DxM := replaceColIn3x3(m, 0, c1, c2, c3)
	DyM := replaceColIn3x3(m, 1, c1, c2, c3)
	DzM := replaceColIn3x3(m, 2, c1, c2, c3)

	D = getDeterminantOf3x3(m)
	Dx = getDeterminantOf3x3(DxM)
	Dy = getDeterminantOf3x3(DyM)
	Dz = getDeterminantOf3x3(DzM)

	output += fmt.Sprintf("system:\n"+
		"%sx + %sy + %sz = %s\n"+
		"%sx + %sy + %sz = %s\n"+
		"%sx + %sy + %sz = %s\n\n",
		fmtFloat(m[0][0]), fmtFloat(m[0][1]), fmtFloat(m[0][2]), fmtFloat(c1),
		fmtFloat(m[1][0]), fmtFloat(m[1][1]), fmtFloat(m[1][2]), fmtFloat(c2),
		fmtFloat(m[2][0]), fmtFloat(m[2][1]), fmtFloat(m[2][2]), fmtFloat(c3),
	)
	if D == 0 {
		if Dx == 0 && Dy == 0 && Dz == 0 {
			output += "dependent - with multiple solutions"
		} else {
			output += "inconsistent - no solution"
		}
	} else {
		x = Dx / D
		y = Dy / D
		z = Dz / D
		output += fmt.Sprintf("solution:\n"+
			"x = %s, y = %s, z = %s",
			fmtFloat(x), fmtFloat(y), fmtFloat(z),
		)
	}
	return output
}
