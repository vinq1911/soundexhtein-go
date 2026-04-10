package soundex

// FullCode is a variable-length phonetic code for agglutinative languages.
// No truncation — encodes the entire word. Stack-allocated, zero heap.
// fullcode[0] = length, fullcode[1..len] = data. Max 63 code bytes.
type FullCode [64]byte

// Len returns the number of meaningful bytes in the code.
func (c FullCode) Len() int { return int(c[0]) }

// Bytes returns the code data as a slice.
func (c FullCode) Bytes() []byte { return c[1 : 1+c[0]] }

// String returns the code as a string.
func (c FullCode) String() string { return string(c.Bytes()) }

// Equal reports whether two FullCodes have the same content.
func (c FullCode) Equal(other FullCode) bool {
	n := c[0]
	if n != other[0] {
		return false
	}
	for i := byte(1); i <= n; i++ {
		if c[i] != other[i] {
			return false
		}
	}
	return true
}

// FullCodePair holds strict and relaxed encodings of the same word.
// Strict preserves geminates (kk→KK) for exact matching.
// Relaxed collapses geminates (kk→K) for typo-tolerant fuzzy matching.
type FullCodePair struct {
	Strict  FullCode
	Relaxed FullCode
}

// FullDistance computes Levenshtein distance between two FullCodes.
// Uses the shared levenshtein core with a [64]int row buffer. Zero allocations.
func FullDistance(a, b FullCode) int {
	var row [64]int
	return levenshtein(a.Bytes(), b.Bytes(), row[:])
}

// FullSoundexDistance encodes both words fully, returns Levenshtein distance.
// Uses relaxed mode by default (better for fuzzy matching).
func FullSoundexDistance(a, b []byte, algo Algorithm) int {
	pa := FullEncode(a, algo)
	pb := FullEncode(b, algo)
	return FullDistance(pa.Relaxed, pb.Relaxed)
}

// FullEncode dispatches to the appropriate full-length encoder.
func FullEncode(word []byte, algo Algorithm) FullCodePair {
	switch algo {
	case AlgoFinnish:
		return FullFinnish(word)
	default:
		// For non-Finnish languages, wrap the short encoder in a FullCodePair.
		c := Encode(word, algo)
		var fc FullCode
		fc[0] = c[0]
		copy(fc[1:], c[1:1+c[0]])
		return FullCodePair{Strict: fc, Relaxed: fc}
	}
}
