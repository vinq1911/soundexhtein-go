package soundex

// levenshtein computes the Levenshtein distance between two byte slices.
// Uses single-row Wagner-Fischer DP. The row buffer must be at least len(b)+1.
// This is the shared core used by both Distance (Code) and FullDistance (FullCode).
func levenshtein(a, b []byte, row []int) int {
	aLen := len(a)
	bLen := len(b)

	if aLen == 0 {
		return bLen
	}
	if bLen == 0 {
		return aLen
	}

	// Ensure a is the shorter one for the single-row optimization.
	if aLen > bLen {
		a, b = b, a
		aLen, bLen = bLen, aLen
	}

	for j := 1; j <= bLen; j++ {
		row[j] = j
	}

	for i := 1; i <= aLen; i++ {
		prev := i - 1
		row[0] = i
		ai := a[i-1]

		for j := 1; j <= bLen; j++ {
			cost := 1
			if ai == b[j-1] {
				cost = 0
			}

			// min of insert, delete, substitute
			val := row[j] + 1 // delete
			if ins := row[j-1] + 1; ins < val {
				val = ins // insert
			}
			if sub := prev + cost; sub < val {
				val = sub // substitute
			}

			prev = row[j]
			row[j] = val
		}
	}

	return row[bLen]
}

// Distance computes the Levenshtein distance between two Codes.
// Optimized for short fixed-size codes (4-8 bytes). Zero allocations.
func Distance(a, b Code) int {
	var row [8]int
	return levenshtein(a.Bytes(), b.Bytes(), row[:])
}
