package soundex

// encoderRule defines a digraph or special multi-character encoding rule.
// Checked before single-character encoding in the main loop.
type encoderRule struct {
	// match returns true if this rule applies at position pos in the buffer.
	// r is buf[pos], next is buf[pos+1], next2 is buf[pos+2].
	match func(r, next, next2 rune) bool
	code  byte // phonetic code to emit
	skip  int  // number of runes to consume (including the first)
}

// langEncoder defines a language's phonetic encoding configuration.
// Used by encodeWithConfig to avoid duplicating the main encoding loop
// across Swedish, Norwegian, Danish, Estonian, Latvian, and Lithuanian.
type langEncoder struct {
	// firstLetter maps the first rune to the code's leading ASCII byte.
	firstLetter func(rune) byte

	// isVowel reports whether a rune is a vowel in this language.
	isVowel func(rune) bool

	// digit maps a consonant rune to its phonetic code byte.
	// Returns 0 for vowels or unknown characters.
	// prev is the preceding rune (0 at start), next is the following (0 at end).
	digit func(prev, r, next rune) byte

	// rules are digraph/multi-character rules checked before single-char encoding.
	// Rules are evaluated in order; first match wins.
	rules []encoderRule
}

// amimica-ignore: buffer init similar to Finnish() but Finnish encodes vowels differently and cannot use the generic loop
// encodeWithConfig runs the generic encoding loop with language-specific config.
// Produces a Code [8]byte (up to 7 phonetic digits). Zero allocations.
//
// The encoding algorithm:
//  1. Fill a fixed-size rune buffer from UTF-8 input (uppercase, letters only).
//  2. Try each digraph rule at current position (longest match first).
//  3. If no rule matches, look up the single-character digit.
//  4. Skip double consonants (emit once).
//  5. Vowels reset the last-digit dedup tracker but emit no code.
//  6. The first consonant encountered also emits firstLetter as c[1].
func encodeWithConfig(name []byte, cfg *langEncoder) Code {
	var c Code
	var buf runeBuffer
	fillRuneBuffer(name, &buf)
	if buf.len == 0 {
		return c
	}

	codeLen := byte(0)
	var lastDigit byte
	firstLetter := cfg.firstLetter(buf.at(0))
	emittedFirst := false

	pos := 0
	for pos < buf.len && codeLen < 7 {
		r := buf.at(pos)
		next := buf.at(pos + 1)
		next2 := buf.at(pos + 2)

		// Check digraph rules first.
		matched := false
		for ri := range cfg.rules {
			rule := &cfg.rules[ri]
			if rule.match(r, next, next2) {
				if !emittedFirst {
					c[1] = firstLetter
					c[2] = rule.code
					codeLen = 2
					lastDigit = rule.code
					emittedFirst = true
				} else if rule.code != lastDigit {
					codeLen++
					c[codeLen] = rule.code
					lastDigit = rule.code
				}
				pos += rule.skip
				matched = true
				break
			}
		}
		if matched {
			continue
		}

		// Single-character encoding.
		prev := buf.at(pos - 1)
		digit := cfg.digit(prev, r, next)
		if digit == 0 {
			if cfg.isVowel(r) {
				lastDigit = 0 // vowels reset dedup
			}
			pos++
			continue
		}

		// Skip double consonants.
		advance := 1
		if buf.at(pos+1) == r && !cfg.isVowel(r) {
			advance = 2
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
		pos += advance
	}

	c[0] = codeLen
	return c
}

// scandFirstLetter is an alias for runeToBaseLetter, used by Scandinavian configs.
var scandFirstLetter = runeToBaseLetter
