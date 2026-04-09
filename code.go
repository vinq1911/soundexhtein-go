package soundex

// Algorithm selects the phonetic encoding algorithm.
type Algorithm uint8

const (
	AlgoAmerican Algorithm = iota
	AlgoCologne
	AlgoDaitchMokotoff
	AlgoMetaphone
	AlgoFinnish
	AlgoSwedish
	AlgoNorwegian
	AlgoDanish
	AlgoEstonian
	AlgoLatvian
	AlgoLithuanian
)

// Code is a fixed-size phonetic code. code[0] is the length, code[1..len] is the data.
// Stack-allocated, zero heap overhead.
type Code [8]byte

// Len returns the number of meaningful bytes in the code.
func (c Code) Len() int { return int(c[0]) }

// Bytes returns the code data as a slice (no allocation — points into array).
func (c Code) Bytes() []byte { return c[1 : 1+c[0]] }

// String returns the code as a string.
func (c Code) String() string { return string(c.Bytes()) }

// Equal reports whether two codes have the same content.
func (c Code) Equal(other Code) bool {
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

// DMCodes holds up to 8 Daitch-Mokotoff branching codes. count is the number of valid codes.
type DMCodes struct {
	Codes [8]Code
	Count int
}

// Match represents a search result from PhoneticIndex.
type Match struct {
	Index    int
	Distance int
}

// Encode dispatches to the appropriate encoder.
func Encode(name []byte, algo Algorithm) Code {
	switch algo {
	case AlgoAmerican:
		return American(name)
	case AlgoCologne:
		return Cologne(name)
	case AlgoMetaphone:
		return Metaphone(name)
	case AlgoFinnish:
		return Finnish(name)
	case AlgoSwedish:
		return Swedish(name)
	case AlgoNorwegian:
		return Norwegian(name)
	case AlgoDanish:
		return Danish(name)
	case AlgoEstonian:
		return Estonian(name)
	case AlgoLatvian:
		return Latvian(name)
	case AlgoLithuanian:
		return Lithuanian(name)
	default:
		return American(name)
	}
}

// SoundexDistance encodes both names with the given algorithm, then returns their Levenshtein distance.
func SoundexDistance(a, b []byte, algo Algorithm) int {
	return Distance(Encode(a, algo), Encode(b, algo))
}

// upper returns the uppercase ASCII version of b, or b itself if not lowercase.
func upper(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - 32
	}
	return b
}

// isLetter reports whether b is an ASCII letter.
func isLetter(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}
