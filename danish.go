package soundex

// Danish computes a phonetic code optimized for Danish words.
// Handles æ, ø, å and Danish-specific phonetics (soft D, SJ).
// Variable-length up to 7 digits. Zero allocations.
//
// Phonetic codes:
//   '1'=B,P  '2'=D(hard),T  '3'=G,K,C(hard),Q  '4'=L  '5'=M,N  '6'=R  '7'=V,W
//   '8'=S,Z  '9'=H  'A'=J  'B'=NG,NK  'C'=F,PH
//   'D'=SJ (sj-sound)
//   'E'=soft D (D between vowels or after vowel before consonant/end)
func Danish(name []byte) Code {
	return encodeWithConfig(name, &danishConfig)
}

var danishConfig = langEncoder{
	firstLetter: scandFirstLetter,
	isVowel:     danishVowelCheck,
	digit:       danishDigit,
	rules: []encoderRule{
		// SJ digraph
		{func(r, n, _ rune) bool { return r == 'S' && n == 'J' }, 'D', 2},
		// NG/NK nasal
		{func(r, n, _ rune) bool { return r == 'N' && (n == 'G' || n == 'K') }, 'B', 2},
		// CK, PH
		{func(r, n, _ rune) bool { return r == 'C' && n == 'K' }, '3', 2},
		{func(r, n, _ rune) bool { return r == 'P' && n == 'H' }, 'C', 2},
	},
}

// danishDigit handles the Danish soft D rule: D after a vowel and not before
// another vowel is a "soft D" (approximant [ð]), encoded as 'E' to distinguish
// it from the hard D/T coded as '2'.
func danishDigit(prev, r, next rune) byte {
	switch r {
	case 'D':
		// Soft D: after a vowel, and not before another vowel
		if danishVowelCheck(prev) && !danishVowelCheck(next) {
			return 'E' // soft D
		}
		return '2'
	case 'C':
		if isFrontVowelDan(next) {
			return '8'
		}
		return '3'
	case 'B', 'P':
		return '1'
	case 'T':
		return '2'
	case 'G', 'K', 'Q', 'X':
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
	}
	return 0
}

func danishVowelCheck(r rune) bool {
	switch r {
	case 'A', 'E', 'I', 'O', 'U', 'Y', 'Æ', 'Ø', 'Å':
		return true
	}
	return false
}

func isFrontVowelDan(r rune) bool {
	return r == 'E' || r == 'I' || r == 'Y' || r == 'Æ' || r == 'Ø'
}
