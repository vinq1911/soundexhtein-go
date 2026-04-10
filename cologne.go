package soundex

// amimica-ignore: init boilerplate shared with American/Metaphone is structurally similar but uses different rules
// Cologne computes the Cologne phonetic code for the given name.
// Optimized for German names. Zero allocations.
func Cologne(name []byte) Code {
	var c Code
	n := len(name)
	if n == 0 {
		return c
	}

	pos := 0
	for pos < n && !isLetter(name[pos]) {
		pos++
	}
	if pos == n {
		return c
	}

	codeLen := byte(0)
	var lastCode byte = 0xFF

	for pos < n && codeLen < 7 {
		b := name[pos]
		pos++
		if !isLetter(b) {
			continue
		}
		u := upper(b)
		idx := u - 'A'
		if idx >= 26 {
			continue
		}

		var digit byte

		switch u {
		case 'D', 'T':
			// Before C, S, Z → 8, else 2
			if pos < n {
				next := upper(name[pos])
				if next == 'C' || next == 'S' || next == 'Z' {
					digit = '8'
				} else {
					digit = '2'
				}
			} else {
				digit = '2'
			}
		case 'C':
			if codeLen == 0 {
				// Initial C: before A, H, K, L, O, Q, R, U, X → 4, else 8
				if pos < n {
					next := upper(name[pos])
					if next == 'A' || next == 'H' || next == 'K' || next == 'L' ||
						next == 'O' || next == 'Q' || next == 'R' || next == 'U' || next == 'X' {
						digit = '4'
					} else {
						digit = '8'
					}
				} else {
					digit = '8'
				}
			} else {
				// After S, Z → 8; before A, H, K, O, Q, U, X → 4; else 8
				var prev byte
				if pos >= 2 {
					prev = upper(name[pos-2])
				}
				if prev == 'S' || prev == 'Z' {
					digit = '8'
				} else if pos < n {
					next := upper(name[pos])
					if next == 'A' || next == 'H' || next == 'K' || next == 'O' ||
						next == 'Q' || next == 'U' || next == 'X' {
						digit = '4'
					} else {
						digit = '8'
					}
				} else {
					digit = '8'
				}
			}
		case 'P':
			if pos < n && upper(name[pos]) == 'H' {
				digit = '3'
				pos++ // skip H
			} else {
				digit = '1'
			}
		case 'X':
			// X = 48 in Cologne. Emit '4' then '8'.
			if lastCode != '4' {
				codeLen++
				c[codeLen] = '4'
				lastCode = '4'
			}
			digit = '8'
		case 'H':
			continue // H is always ignored
		default:
			digit = cologneTable[idx]
		}

		// Skip duplicate consecutive codes and drop '0' (vowels) except at start.
		if digit == '0' {
			lastCode = '0'
			continue
		}
		if digit != lastCode {
			codeLen++
			c[codeLen] = digit
			lastCode = digit
		}
	}

	c[0] = codeLen
	return c
}
