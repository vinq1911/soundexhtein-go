package soundex

// Daitch-Mokotoff Soundex rules.
// Each rule maps a character sequence to up to 3 codes depending on context:
// [0] = at start of name, [1] = before a vowel, [2] = other.
// A value of 0xFF means "not applicable" (use next rule).

type dmRule struct {
	prefix  string
	codes   [3]byte // start, before-vowel, other
	advance int     // how many characters to consume
}

// dmRules is organized by first character for fast dispatch.
// This is a simplified Daitch-Mokotoff that covers the most common cases.
var dmRules = map[byte][]dmRule{
	'A': {
		{prefix: "AI", codes: [3]byte{0, 1, 0xFF}, advance: 2},
		{prefix: "AJ", codes: [3]byte{0, 1, 0xFF}, advance: 2},
		{prefix: "AY", codes: [3]byte{0, 1, 0xFF}, advance: 2},
		{prefix: "AU", codes: [3]byte{0, 7, 0xFF}, advance: 2},
		{prefix: "A", codes: [3]byte{0, 0xFF, 0xFF}, advance: 1},
	},
	'B': {{prefix: "B", codes: [3]byte{7, 7, 7}, advance: 1}},
	'C': {
		{prefix: "CHS", codes: [3]byte{5, 54, 54}, advance: 3},
		{prefix: "CH", codes: [3]byte{5, 5, 5}, advance: 2},
		{prefix: "CK", codes: [3]byte{5, 5, 5}, advance: 2},
		{prefix: "CZ", codes: [3]byte{4, 4, 4}, advance: 2},
		{prefix: "CS", codes: [3]byte{4, 4, 4}, advance: 2},
		{prefix: "C", codes: [3]byte{5, 5, 5}, advance: 1},
	},
	'D': {
		{prefix: "DRS", codes: [3]byte{4, 4, 4}, advance: 3},
		{prefix: "DRZ", codes: [3]byte{4, 4, 4}, advance: 3},
		{prefix: "DS", codes: [3]byte{4, 4, 4}, advance: 2},
		{prefix: "DZ", codes: [3]byte{4, 4, 4}, advance: 2},
		{prefix: "DT", codes: [3]byte{3, 3, 3}, advance: 2},
		{prefix: "D", codes: [3]byte{3, 3, 3}, advance: 1},
	},
	'E': {
		{prefix: "EI", codes: [3]byte{0, 1, 0xFF}, advance: 2},
		{prefix: "EJ", codes: [3]byte{0, 1, 0xFF}, advance: 2},
		{prefix: "EY", codes: [3]byte{0, 1, 0xFF}, advance: 2},
		{prefix: "EU", codes: [3]byte{1, 1, 0xFF}, advance: 2},
		{prefix: "E", codes: [3]byte{0, 0xFF, 0xFF}, advance: 1},
	},
	'F': {
		{prefix: "FB", codes: [3]byte{7, 7, 7}, advance: 2},
		{prefix: "F", codes: [3]byte{7, 7, 7}, advance: 1},
	},
	'G': {{prefix: "G", codes: [3]byte{5, 5, 5}, advance: 1}},
	'H': {{prefix: "H", codes: [3]byte{5, 5, 0xFF}, advance: 1}},
	'I': {
		{prefix: "IA", codes: [3]byte{1, 0xFF, 0xFF}, advance: 2},
		{prefix: "IE", codes: [3]byte{1, 0xFF, 0xFF}, advance: 2},
		{prefix: "IO", codes: [3]byte{1, 0xFF, 0xFF}, advance: 2},
		{prefix: "IU", codes: [3]byte{1, 0xFF, 0xFF}, advance: 2},
		{prefix: "I", codes: [3]byte{0, 0xFF, 0xFF}, advance: 1},
	},
	'J': {{prefix: "J", codes: [3]byte{1, 1, 1}, advance: 1}},
	'K': {
		{prefix: "KS", codes: [3]byte{5, 54, 54}, advance: 2},
		{prefix: "KH", codes: [3]byte{5, 5, 5}, advance: 2},
		{prefix: "K", codes: [3]byte{5, 5, 5}, advance: 1},
	},
	'L': {{prefix: "L", codes: [3]byte{8, 8, 8}, advance: 1}},
	'M': {
		{prefix: "MN", codes: [3]byte{66, 66, 66}, advance: 2},
		{prefix: "M", codes: [3]byte{6, 6, 6}, advance: 1},
	},
	'N': {
		{prefix: "NM", codes: [3]byte{66, 66, 66}, advance: 2},
		{prefix: "N", codes: [3]byte{6, 6, 6}, advance: 1},
	},
	'O': {
		{prefix: "OI", codes: [3]byte{0, 1, 0xFF}, advance: 2},
		{prefix: "OJ", codes: [3]byte{0, 1, 0xFF}, advance: 2},
		{prefix: "OY", codes: [3]byte{0, 1, 0xFF}, advance: 2},
		{prefix: "O", codes: [3]byte{0, 0xFF, 0xFF}, advance: 1},
	},
	'P': {
		{prefix: "PH", codes: [3]byte{7, 7, 7}, advance: 2},
		{prefix: "P", codes: [3]byte{7, 7, 7}, advance: 1},
	},
	'Q': {{prefix: "Q", codes: [3]byte{5, 5, 5}, advance: 1}},
	'R': {
		{prefix: "RS", codes: [3]byte{4, 4, 4}, advance: 2},
		{prefix: "RZ", codes: [3]byte{4, 4, 4}, advance: 2},
		{prefix: "R", codes: [3]byte{9, 9, 9}, advance: 1},
	},
	'S': {
		{prefix: "SCHTSCH", codes: [3]byte{2, 4, 4}, advance: 7},
		{prefix: "SHTCH", codes: [3]byte{2, 4, 4}, advance: 5},
		{prefix: "SHCH", codes: [3]byte{2, 4, 4}, advance: 4},
		{prefix: "SCH", codes: [3]byte{4, 4, 4}, advance: 3},
		{prefix: "SHT", codes: [3]byte{2, 43, 43}, advance: 3},
		{prefix: "SH", codes: [3]byte{4, 4, 4}, advance: 2},
		{prefix: "ST", codes: [3]byte{2, 43, 43}, advance: 2},
		{prefix: "S", codes: [3]byte{4, 4, 4}, advance: 1},
	},
	'T': {
		{prefix: "TSCH", codes: [3]byte{4, 4, 4}, advance: 4},
		{prefix: "TCH", codes: [3]byte{4, 4, 4}, advance: 3},
		{prefix: "TH", codes: [3]byte{3, 3, 3}, advance: 2},
		{prefix: "TS", codes: [3]byte{4, 4, 4}, advance: 2},
		{prefix: "TZ", codes: [3]byte{4, 4, 4}, advance: 2},
		{prefix: "T", codes: [3]byte{3, 3, 3}, advance: 1},
	},
	'U': {
		{prefix: "UI", codes: [3]byte{0, 1, 0xFF}, advance: 2},
		{prefix: "UJ", codes: [3]byte{0, 1, 0xFF}, advance: 2},
		{prefix: "UY", codes: [3]byte{0, 1, 0xFF}, advance: 2},
		{prefix: "U", codes: [3]byte{0, 0xFF, 0xFF}, advance: 1},
	},
	'V': {{prefix: "V", codes: [3]byte{7, 7, 7}, advance: 1}},
	'W': {{prefix: "W", codes: [3]byte{7, 7, 7}, advance: 1}},
	'X': {{prefix: "X", codes: [3]byte{5, 54, 54}, advance: 1}},
	'Y': {{prefix: "Y", codes: [3]byte{1, 0xFF, 0xFF}, advance: 1}},
	'Z': {
		{prefix: "ZH", codes: [3]byte{4, 4, 4}, advance: 2},
		{prefix: "ZS", codes: [3]byte{4, 4, 4}, advance: 2},
		{prefix: "Z", codes: [3]byte{4, 4, 4}, advance: 1},
	},
}

// DaitchMokotoff computes the Daitch-Mokotoff Soundex for a name.
// Returns up to 8 branching codes (6 digits each). Zero allocations.
func DaitchMokotoff(name []byte) DMCodes {
	var result DMCodes
	n := len(name)
	if n == 0 {
		return result
	}

	// Uppercase the input into a fixed buffer.
	var buf [64]byte
	bLen := 0
	for i := 0; i < n && bLen < 64; i++ {
		if isLetter(name[i]) {
			buf[bLen] = upper(name[i])
			bLen++
		}
	}
	if bLen == 0 {
		return result
	}

	// We build codes digit by digit. Each "branch" is a code being constructed.
	// Use fixed arrays to avoid allocations.
	type branch struct {
		digits [6]byte
		dLen   byte
		last   byte // last code digit for dedup
	}

	var branches [8]branch
	var newBranches [8]branch
	bCount := 1
	branches[0] = branch{}

	pos := 0
	isStart := true

	for pos < bLen {
		ch := buf[pos]
		rules, ok := dmRules[ch]
		if !ok {
			pos++
			isStart = false
			continue
		}

		// Find matching rule (longest prefix match — rules are sorted longest first).
		var matched *dmRule
		for ri := range rules {
			r := &rules[ri]
			if pos+len(r.prefix) <= bLen {
				match := true
				for k := 1; k < len(r.prefix); k++ {
					if buf[pos+k] != r.prefix[k] {
						match = false
						break
					}
				}
				if match {
					matched = r
					break
				}
			}
		}
		if matched == nil {
			pos++
			isStart = false
			continue
		}

		// Determine context: start, before-vowel, or other.
		nextPos := pos + matched.advance
		beforeVowel := nextPos < bLen && isDMVowel(buf[nextPos])

		var codeIdx int
		if isStart {
			codeIdx = 0
		} else if beforeVowel {
			codeIdx = 1
		} else {
			codeIdx = 2
		}

		digit := matched.codes[codeIdx]
		if digit == 0xFF {
			// Not applicable in this context — skip
			pos = nextPos
			isStart = false
			continue
		}

		// Apply digit to all branches.
		nbCount := 0
		for bi := 0; bi < bCount; bi++ {
			b := &branches[bi]
			if b.dLen >= 6 {
				if nbCount < 8 {
					newBranches[nbCount] = *b
					nbCount++
				}
				continue
			}

			// Multi-digit codes (like 54, 43, 66) emit two digits.
			if digit >= 10 {
				d1 := digit / 10
				d2 := digit % 10
				if d1+'0' != b.last && nbCount < 8 {
					nb := *b
					nb.digits[nb.dLen] = d1 + '0'
					nb.dLen++
					nb.last = d1 + '0'
					if nb.dLen < 6 && d2+'0' != nb.last {
						nb.digits[nb.dLen] = d2 + '0'
						nb.dLen++
						nb.last = d2 + '0'
					}
					newBranches[nbCount] = nb
					nbCount++
				}
			} else {
				dByte := digit + '0'
				if dByte != b.last && nbCount < 8 {
					nb := *b
					nb.digits[nb.dLen] = dByte
					nb.dLen++
					nb.last = dByte
					newBranches[nbCount] = nb
					nbCount++
				} else if nbCount < 8 {
					newBranches[nbCount] = *b
					nbCount++
				}
			}
		}

		if nbCount > 0 {
			copy(branches[:nbCount], newBranches[:nbCount])
			bCount = nbCount
		}

		pos = nextPos
		isStart = false
	}

	// Pad codes to 6 digits with zeros and convert to Code type.
	for i := 0; i < bCount && i < 8; i++ {
		b := &branches[i]
		for b.dLen < 6 {
			b.digits[b.dLen] = '0'
			b.dLen++
		}
		var code Code
		code[0] = 6
		copy(code[1:7], b.digits[:6])
		result.Codes[i] = code
	}
	result.Count = bCount
	if result.Count > 8 {
		result.Count = 8
	}

	return result
}

func isDMVowel(b byte) bool {
	return b == 'A' || b == 'E' || b == 'I' || b == 'O' || b == 'U'
}
