package soundex

// decodeRune decodes a single UTF-8 rune from b at position pos.
// Returns the rune and the number of bytes consumed. Zero allocations.
// This avoids importing unicode/utf8 for minimal dependency.
func decodeRune(b []byte, pos int) (rune, int) {
	if pos >= len(b) {
		return 0, 0
	}
	c := b[pos]
	if c < 0x80 {
		return rune(c), 1
	}
	if c < 0xC0 {
		return 0xFFFD, 1 // invalid continuation byte
	}
	if c < 0xE0 {
		if pos+1 >= len(b) {
			return 0xFFFD, 1
		}
		return rune(c&0x1F)<<6 | rune(b[pos+1]&0x3F), 2
	}
	if c < 0xF0 {
		if pos+2 >= len(b) {
			return 0xFFFD, 1
		}
		return rune(c&0x0F)<<12 | rune(b[pos+1]&0x3F)<<6 | rune(b[pos+2]&0x3F), 3
	}
	if pos+3 >= len(b) {
		return 0xFFFD, 1
	}
	return rune(c&0x07)<<18 | rune(b[pos+1]&0x3F)<<12 | rune(b[pos+2]&0x3F)<<6 | rune(b[pos+3]&0x3F), 4
}

// runeUpper returns the uppercase version of a rune (ASCII + common Nordic/Baltic).
func runeUpper(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 32
	}
	switch r {
	case 'ä':
		return 'Ä'
	case 'ö':
		return 'Ö'
	case 'ü':
		return 'Ü'
	case 'å':
		return 'Å'
	case 'æ':
		return 'Æ'
	case 'ø':
		return 'Ø'
	case 'õ':
		return 'Õ'
	case 'š':
		return 'Š'
	case 'ž':
		return 'Ž'
	case 'č':
		return 'Č'
	case 'ģ':
		return 'Ģ'
	case 'ķ':
		return 'Ķ'
	case 'ļ':
		return 'Ļ'
	case 'ņ':
		return 'Ņ'
	case 'ā':
		return 'Ā'
	case 'ē':
		return 'Ē'
	case 'ī':
		return 'Ī'
	case 'ū':
		return 'Ū'
	case 'ą':
		return 'Ą'
	case 'ę':
		return 'Ę'
	case 'ė':
		return 'Ė'
	case 'į':
		return 'Į'
	case 'ų':
		return 'Ų'
	}
	return r
}

// isRuneLetter reports whether r is a letter we handle (ASCII + Nordic/Baltic extended).
func isRuneLetter(r rune) bool {
	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
		return true
	}
	switch r {
	case 'Ä', 'ä', 'Ö', 'ö', 'Ü', 'ü', 'Å', 'å',
		'Æ', 'æ', 'Ø', 'ø', 'Õ', 'õ',
		'Š', 'š', 'Ž', 'ž', 'Č', 'č',
		'Ģ', 'ģ', 'Ķ', 'ķ', 'Ļ', 'ļ', 'Ņ', 'ņ',
		'Ā', 'ā', 'Ē', 'ē', 'Ī', 'ī', 'Ū', 'ū',
		'Ą', 'ą', 'Ę', 'ę', 'Ė', 'ė', 'Į', 'į', 'Ų', 'ų':
		return true
	}
	return false
}

// nordicVowel reports whether r is a vowel in Nordic/Baltic languages.
func nordicVowel(r rune) bool {
	switch r {
	case 'A', 'E', 'I', 'O', 'U', 'Y',
		'Ä', 'Ö', 'Ü', 'Å', 'Æ', 'Ø', 'Õ',
		'Ā', 'Ē', 'Ī', 'Ū',
		'Ą', 'Ę', 'Ė', 'Į', 'Ų':
		return true
	}
	return false
}

// runeBuffer is a fixed-size buffer for uppercased runes. Stack-allocated.
type runeBuffer struct {
	runes [64]rune
	len   int
}

// fillRuneBuffer decodes UTF-8 input into uppercase runes, filtering to letters only.
func fillRuneBuffer(name []byte, buf *runeBuffer) {
	buf.len = 0
	pos := 0
	for pos < len(name) && buf.len < 64 {
		r, size := decodeRune(name, pos)
		pos += size
		if isRuneLetter(r) {
			buf.runes[buf.len] = runeUpper(r)
			buf.len++
		}
	}
}

// peekRune returns the rune at index i in the buffer, or 0 if out of bounds.
func (rb *runeBuffer) at(i int) rune {
	if i >= 0 && i < rb.len {
		return rb.runes[i]
	}
	return 0
}
