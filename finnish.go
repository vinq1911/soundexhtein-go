package soundex

// Finnish computes a phonetic code optimized for Finnish words.
// Handles ä, ö, y as front vowels distinct from a, o, u.
// Variable-length up to 7 digits. Zero allocations.
//
// Finnish phonetic groupings:
//
// Consonants:
//
//	'1' = B, P
//	'2' = D, T
//	'3' = G, K
//	'4' = L
//	'5' = M, N
//	'6' = R
//	'7' = V, W
//	'8' = S, Z, TS
//	'9' = H
//	'A' = J
//	'B' = NG, NK (Finnish velar nasal)
//	'C' = F (loanword)
//
// Vowels (encoded, not dropped — Finnish vowel identity is phonemic):
//
//	'a' = A (back)
//	'e' = E (neutral)
//	'i' = I (neutral)
//	'o' = O (back)
//	'u' = U (back)
//	'y' = Y (front, pairs with U)
//	'w' = Ä (front, pairs with A)
//	'x' = Ö (front, pairs with O)
//
// Adjacent identical vowels (long vowels like aa, ää) are collapsed to one code.
func Finnish(name []byte) Code {
	var c Code
	var buf runeBuffer
	fillRuneBuffer(name, &buf)
	if buf.len == 0 {
		return c
	}

	codeLen := byte(0)
	var lastDigit byte
	firstLetter := finnishFirstLetter(buf.at(0))
	emittedFirst := false

	pos := 0
	for pos < buf.len && codeLen < 7 {
		r := buf.at(pos)
		next := buf.at(pos + 1)

		// Digraph: NG, NK → 'B' (velar nasal)
		if r == 'N' && (next == 'G' || next == 'K') {
			digit := byte('B')
			if !emittedFirst {
				c[1] = firstLetter
				c[2] = digit
				codeLen = 2
				lastDigit = digit
				emittedFirst = true
			} else if digit != lastDigit {
				codeLen++
				c[codeLen] = digit
				lastDigit = digit
			}
			pos += 2
			continue
		}

		// Digraph: TS → '8'
		if r == 'T' && next == 'S' {
			digit := byte('8')
			if !emittedFirst {
				c[1] = firstLetter
				c[2] = digit
				codeLen = 2
				lastDigit = digit
				emittedFirst = true
			} else if digit != lastDigit {
				codeLen++
				c[codeLen] = digit
				lastDigit = digit
			}
			pos += 2
			continue
		}

		// Double consonants → single code (skip the duplicate)
		if r == next && !finnishVowel(r) {
			pos++
			continue
		}

		digit := finnishDigit(r)

		// Adjacent identical vowels (long vowels) → single code
		if digit != 0 && finnishVowel(r) && r == next {
			pos++
		}

		if digit == 0 {
			// Unknown character — skip
			pos++
			continue
		}

		if !emittedFirst {
			c[1] = firstLetter
			c[2] = digit
			codeLen = 2
			lastDigit = digit
			emittedFirst = true
		} else if digit != lastDigit {
			codeLen++
			c[codeLen] = digit
			lastDigit = digit
		}
		pos++
	}

	c[0] = codeLen
	return c
}

func finnishVowel(r rune) bool {
	switch r {
	case 'A', 'E', 'I', 'O', 'U', 'Y', 'Ä', 'Ö':
		return true
	}
	return false
}

func finnishDigit(r rune) byte {
	// Consonants
	switch r {
	case 'B', 'P':
		return '1'
	case 'D', 'T':
		return '2'
	case 'G', 'K':
		return '3'
	case 'L':
		return '4'
	case 'M', 'N':
		return '5'
	case 'R':
		return '6'
	case 'V', 'W':
		return '7'
	case 'S', 'Z':
		return '8'
	case 'H':
		return '9'
	case 'J':
		return 'A'
	case 'F':
		return 'C'
	case 'C', 'Q', 'X':
		return '3'
	}
	// Vowels — encoded with distinct codes
	switch r {
	case 'A':
		return 'a'
	case 'E':
		return 'e'
	case 'I':
		return 'i'
	case 'O':
		return 'o'
	case 'U':
		return 'u'
	case 'Y':
		return 'y'
	case 'Ä':
		return 'w' // front vowel, distinct from A
	case 'Ö':
		return 'x' // front vowel, distinct from O
	}
	return 0
}

// finnishFirstLetter returns the leading character for the code.
// Ä and Ö get distinct representations to avoid collisions with A and O.
func finnishFirstLetter(r rune) byte {
	switch r {
	case 'Ä':
		return 'W' // Ä → W (unique, not used in native Finnish)
	case 'Ö':
		return 'X' // Ö → X (unique, not used in native Finnish)
	}
	if r >= 'A' && r <= 'Z' {
		return byte(r)
	}
	return '?'
}
