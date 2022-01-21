package chessboard

// Declare a type named Rank which stores if a square is occupied by a piece - this will be a slice of bools
type Rank []bool

// Declare a type named Chessboard which contains a map of eight Ranks, accessed with keys from "A" to "H"
type Chessboard map[string]Rank

// CountInRank returns how many squares are occupied in the chessboard,
// within the given rank
func CountInRank(cb Chessboard, rank string) int {
    var count int
	if row, rowExists := cb[rank]; rowExists {
        for _, col := range row {
            if col {
                count++
            }
        }
    	return count
    }
	return 0
}

// CountInFile returns how many squares are occupied in the chessboard,
// within the given file
func CountInFile(cb Chessboard, file int) int {
	var count int
    if file >= 1 && file <= 8 {
        for _, row := range cb {
            if row[file-1] {
                count++
            }
    	}
    	return count
    }
    return 0
}

// CountAll should count how many squares are present in the chessboard
func CountAll(cb Chessboard) int {
    var sum int
	for _, row := range cb {
        sum += len(row)
    }
	return sum
}

// CountOccupied returns how many squares are occupied in the chessboard
func CountOccupied(cb Chessboard) int {
	var count int
	for _, row := range cb {
        for _, col := range row {
            if col {
                count++
            }
        }
    }
	return count
}
