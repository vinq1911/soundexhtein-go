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
// amimica-ignore: buffer init similar to encodeWithConfig but Finnish has unique vowel encoding and cannot use the generic loop
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

// finnishDigitTable maps ASCII uppercase (0x00-0x7F) to phonetic codes.
// Vowels get lowercase codes, consonants get digit/letter codes, 0 = skip.
// Ä and Ö are outside ASCII and handled inline.
var finnishDigitTable = [128]byte{
	// 0x00-0x40: unused
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 0x00
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 0x10
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 0x20
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 0x30
	0,                                                  // 0x40 (@)
	'a',  // A = 0x41 — vowel
	'1',  // B = 0x42
	'3',  // C = 0x43 (→ K)
	'2',  // D = 0x44
	'e',  // E = 0x45 — vowel
	'C',  // F = 0x46 (loanword)
	'3',  // G = 0x47
	'9',  // H = 0x48
	'i',  // I = 0x49 — vowel
	'A',  // J = 0x4A
	'3',  // K = 0x4B
	'4',  // L = 0x4C
	'5',  // M = 0x4D
	'5',  // N = 0x4E
	'o',  // O = 0x4F — vowel
	'1',  // P = 0x50
	'3',  // Q = 0x51 (→ K)
	'6',  // R = 0x52
	'8',  // S = 0x53
	'2',  // T = 0x54
	'u',  // U = 0x55 — vowel
	'7',  // V = 0x56
	'7',  // W = 0x57
	'3',  // X = 0x58 (→ K)
	'y',  // Y = 0x59 — vowel
	'8',  // Z = 0x5A
	0, 0, 0, 0, 0, 0, // 0x5B-0x60
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 0x61-0x70
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 0x71-0x7F
}

// finnishVowelTable marks ASCII uppercase vowels. Ä and Ö handled inline.
var finnishVowelTable = [128]byte{
	// Only A(0x41), E(0x45), I(0x49), O(0x4F), U(0x55), Y(0x59) are set
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 0x00
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 0x10
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 0x20
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 0x30
	0,                                                  // 0x40
	1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1,      // A..O
	0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0,                   // P..Z
	0, 0, 0, 0, 0, 0, // 0x5B-0x60
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 0x61-0x70
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // 0x71-0x7F
}

// finnishVowel reports whether r is a Finnish vowel. Single indexed load for ASCII.
func finnishVowel(r rune) bool {
	if r < 128 {
		return finnishVowelTable[r] != 0
	}
	return r == 'Ä' || r == 'Ö'
}

// finnishDigit returns the phonetic code for r. Single indexed load for ASCII.
func finnishDigit(r rune) byte {
	if r < 128 {
		return finnishDigitTable[r]
	}
	if r == 'Ä' {
		return 'w'
	}
	if r == 'Ö' {
		return 'x'
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
