package soundex

// Distance computes the Levenshtein distance between two Codes.
// Optimized for short fixed-size codes (4-8 bytes). Zero allocations.
func Distance(a, b Code) int {
	aLen := int(a[0])
	bLen := int(b[0])

	if aLen == 0 {
		return bLen
	}
	if bLen == 0 {
		return aLen
	}

	// Ensure a is the shorter one for the single-row optimization.
	aData := a[1 : 1+aLen]
	bData := b[1 : 1+bLen]
	if aLen > bLen {
		aData, bData = bData, aData
		aLen, bLen = bLen, aLen
	}

	// Single-row DP. Max code length is 7 (8-byte array minus length byte).
	// Use a fixed [8]int array to avoid any heap allocation.
	var row [8]int
	for j := 1; j <= bLen; j++ {
		row[j] = j
	}

	for i := 1; i <= aLen; i++ {
		prev := i - 1
		row[0] = i
		ai := aData[i-1]

		for j := 1; j <= bLen; j++ {
			cost := 1
			if ai == bData[j-1] {
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
