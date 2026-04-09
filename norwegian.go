package soundex

// Norwegian computes a phonetic code optimized for Norwegian words.
// Handles æ, ø, å and Norwegian consonant clusters (KJ, SJ, RS retroflex).
// Variable-length up to 7 digits. Zero allocations.
//
// Phonetic codes:
//   '1'=B,P  '2'=D,T  '3'=G,K,C(hard),Q  '4'=L  '5'=M,N  '6'=R  '7'=V,W
//   '8'=S,Z  '9'=H  'A'=J,GJ,HJ  'B'=NG,NK  'C'=F,PH
//   'D'=SJ,SKJ,SH,SK(front) (sj-sound)
//   'E'=KJ,TJ (kj-sound)
//   'F'=RS (retroflex)
func Norwegian(name []byte) Code {
	return encodeWithConfig(name, &norwegianConfig)
}

var norwegianConfig = langEncoder{
	firstLetter: scandFirstLetter,
	isVowel:     norwegianVowelCheck,
	digit:       norwegianDigit,
	rules: []encoderRule{
		// RS retroflex (must check before single R)
		{func(r, n, _ rune) bool { return r == 'R' && n == 'S' }, 'F', 2},
		// SJ-sound digraphs
		{func(r, n, n2 rune) bool { return r == 'S' && n == 'K' && n2 == 'J' }, 'D', 3},
		{func(r, n, n2 rune) bool {
			return r == 'S' && n == 'K' && isFrontVowelNor(n2)
		}, 'D', 2},
		{func(r, n, _ rune) bool { return r == 'S' && n == 'J' }, 'D', 2},
		{func(r, n, _ rune) bool { return r == 'S' && n == 'H' }, 'D', 2},
		// KJ-sound
		{func(r, n, _ rune) bool { return r == 'K' && n == 'J' }, 'E', 2},
		{func(r, n, _ rune) bool { return r == 'T' && n == 'J' }, 'E', 2},
		// J-sound
		{func(r, n, _ rune) bool { return r == 'G' && n == 'J' }, 'A', 2},
		{func(r, n, _ rune) bool { return r == 'H' && n == 'J' }, 'A', 2},
		// NG/NK nasal
		{func(r, n, _ rune) bool { return r == 'N' && (n == 'G' || n == 'K') }, 'B', 2},
		// CK, PH
		{func(r, n, _ rune) bool { return r == 'C' && n == 'K' }, '3', 2},
		{func(r, n, _ rune) bool { return r == 'P' && n == 'H' }, 'C', 2},
	},
}

func norwegianDigit(_, r, next rune) byte {
	switch r {
	case 'K':
		if isFrontVowelNor(next) {
			return 'E'
		}
		return '3'
	case 'C':
		if isFrontVowelNor(next) {
			return '8'
		}
		return '3'
	case 'B', 'P':
		return '1'
	case 'D', 'T':
		return '2'
	case 'G', 'Q', 'X':
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

func norwegianVowelCheck(r rune) bool {
	switch r {
	case 'A', 'E', 'I', 'O', 'U', 'Y', 'Æ', 'Ø', 'Å':
		return true
	}
	return false
}

func isFrontVowelNor(r rune) bool {
	return r == 'E' || r == 'I' || r == 'Y' || r == 'Æ' || r == 'Ø'
}
