package soundex

// Metaphone computes a simplified Metaphone code (max 4 characters).
// Zero allocations.
func Metaphone(name []byte) Code {
	var c Code
	n := len(name)
	if n == 0 {
		return c
	}

	// Find first letter.
	start := 0
	for start < n && !isLetter(name[start]) {
		start++
	}
	if start == n {
		return c
	}

	// Work with uppercase copy in a fixed buffer.
	var buf [64]byte
	bLen := 0
	for i := start; i < n && bLen < 64; i++ {
		if isLetter(name[i]) {
			buf[bLen] = upper(name[i])
			bLen++
		}
	}

	// Drop initial silent letters.
	pos := 0
	if bLen >= 2 {
		pair := [2]byte{buf[0], buf[1]}
		switch pair {
		case [2]byte{'A', 'E'}, [2]byte{'G', 'N'}, [2]byte{'K', 'N'}, [2]byte{'P', 'N'}, [2]byte{'W', 'R'}:
			pos = 1
		}
	}

	codeLen := byte(0)
	var last byte

	emit := func(ch byte) {
		if codeLen < 4 && ch != last {
			codeLen++
			c[codeLen] = ch
			last = ch
		}
	}

	// If first char is a vowel, emit it.
	if isVowel(buf[pos]) {
		emit(buf[pos])
		pos++
	}

	for pos < bLen && codeLen < 4 {
		ch := buf[pos]
		prev := byte(0)
		if pos > 0 {
			prev = buf[pos-1]
		}
		next := byte(0)
		if pos+1 < bLen {
			next = buf[pos+1]
		}
		next2 := byte(0)
		if pos+2 < bLen {
			next2 = buf[pos+2]
		}

		// Skip duplicate adjacent letters (except C).
		if ch == prev && ch != 'C' {
			pos++
			continue
		}

		switch ch {
		case 'B':
			if prev != 'M' {
				emit('P')
			}
		case 'C':
			if next == 'I' || next == 'E' || next == 'Y' {
				if next == 'I' && next2 == 'A' {
					emit('X') // CIA -> X
				} else {
					emit('S')
				}
			} else {
				emit('K')
			}
		case 'D':
			if next == 'G' && (next2 == 'I' || next2 == 'E' || next2 == 'Y') {
				emit('J')
			} else {
				emit('T')
			}
		case 'F':
			emit('F')
		case 'G':
			if next == 'H' && pos+2 < bLen && !isVowel(next2) {
				pos++ // GH before non-vowel: silent
			} else if pos > 0 && next == 'H' && pos+2 >= bLen {
				// GH at end: silent
				pos++
			} else if prev == 'G' {
				// skip double G handled by previous
			} else {
				emit('K')
			}
		case 'H':
			if isVowel(next) && !isVowel(prev) {
				emit('H')
			}
		case 'J':
			emit('J')
		case 'K':
			if prev != 'C' {
				emit('K')
			}
		case 'L':
			emit('L')
		case 'M':
			emit('M')
		case 'N':
			emit('N')
		case 'P':
			if next == 'H' {
				emit('F')
				pos++
			} else {
				emit('P')
			}
		case 'Q':
			emit('K')
		case 'R':
			emit('R')
		case 'S':
			if next == 'C' && next2 == 'H' {
				emit('S')
				emit('K')
				pos += 2 // skip CH
			} else if next == 'H' || (next == 'I' && (next2 == 'O' || next2 == 'A')) {
				emit('X')
				if next != 'H' {
					pos++ // skip I in SIO/SIA
				}
			} else {
				emit('S')
			}
		case 'T':
			if next == 'H' {
				emit('0') // 0 represents 'th' in Metaphone
				pos++
			} else if next == 'I' && (next2 == 'O' || next2 == 'A') {
				emit('X')
			} else {
				emit('T')
			}
		case 'V':
			emit('F')
		case 'W', 'Y':
			if isVowel(next) {
				emit(ch)
			}
		case 'X':
			emit('K')
			emit('S')
		case 'Z':
			emit('S')
		default:
			// Vowels in non-initial position are dropped.
		}

		pos++
	}

	c[0] = codeLen
	return c
}

func isVowel(b byte) bool {
	return b == 'A' || b == 'E' || b == 'I' || b == 'O' || b == 'U'
}
