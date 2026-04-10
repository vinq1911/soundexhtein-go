package soundex

// amimica-ignore: init boilerplate shared with Cologne/Metaphone is structurally similar but uses different tables and code lengths
// American computes the American Soundex code for the given name.
// Returns a 4-character code (1 letter + 3 digits). Zero allocations.
func American(name []byte) Code {
	var c Code
	n := len(name)
	if n == 0 {
		return c
	}

	// Find first letter.
	pos := 0
	for pos < n && !isLetter(name[pos]) {
		pos++
	}
	if pos == n {
		return c
	}

	// First character is the uppercase letter.
	first := upper(name[pos])
	c[1] = first
	codeLen := byte(1)
	pos++

	// Get the Soundex digit for the first letter (used to suppress duplicates).
	lastDigit := americanTable[first-'A']

	// Process remaining characters.
	for pos < n && codeLen < 4 {
		b := name[pos]
		pos++
		if !isLetter(b) {
			continue
		}
		u := upper(b)
		digit := americanTable[u-'A']
		if digit != 0 && digit != lastDigit {
			codeLen++
			c[codeLen] = digit
			lastDigit = digit
		} else if digit == 0 {
			// Vowels and H/W reset the last digit for adjacent consonant separation.
			// In standard American Soundex, H and W do NOT separate identical codes.
			// Only vowels (A, E, I, O, U) separate.
			if u != 'H' && u != 'W' {
				lastDigit = 0
			}
		}
	}

	// Pad with zeros.
	for codeLen < 4 {
		codeLen++
		c[codeLen] = '0'
	}

	c[0] = 4
	return c
}
