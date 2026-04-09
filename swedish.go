package soundex

// Swedish computes a phonetic code optimized for Swedish words.
// Handles å, ä, ö and Swedish-specific consonant clusters (SJ, SKJ, KJ, TJ).
// Variable-length up to 7 digits. Zero allocations.
//
// Phonetic codes:
//   '1'=B,P  '2'=D,T  '3'=G,K,C(hard),Q  '4'=L  '5'=M,N  '6'=R  '7'=V,W
//   '8'=S,Z  '9'=H  'A'=J,GJ,DJ,LJ,HJ  'B'=NG,NK  'C'=F,PH
//   'D'=SJ,SKJ,STJ,SK(front),SCH,SH (sje-sound)
//   'E'=KJ,TJ,K(front) (tje-sound)
func Swedish(name []byte) Code {
	return encodeWithConfig(name, &swedishConfig)
}

var swedishConfig = langEncoder{
	firstLetter: scandFirstLetter,
	isVowel:     swedishVowelCheck,
	digit:       swedishDigit,
	rules: []encoderRule{
		// SJ-sound digraphs (check longest first)
		{func(r, n, n2 rune) bool { return r == 'S' && n == 'K' && n2 == 'J' }, 'D', 3},
		{func(r, n, n2 rune) bool { return r == 'S' && n == 'T' && n2 == 'J' }, 'D', 3},
		{func(r, n, n2 rune) bool { return r == 'S' && n == 'C' && n2 == 'H' }, 'D', 3},
		{func(r, n, _ rune) bool { return r == 'S' && n == 'K' && isFrontVowelSwe(n) }, 'D', 2},
		{func(r, n, _ rune) bool { return r == 'S' && n == 'J' }, 'D', 2},
		{func(r, n, _ rune) bool { return r == 'S' && n == 'H' }, 'D', 2},
		// TJ-sound digraphs
		{func(r, n, _ rune) bool { return r == 'K' && n == 'J' }, 'E', 2},
		{func(r, n, _ rune) bool { return r == 'T' && n == 'J' }, 'E', 2},
		// J-sound digraphs
		{func(r, n, _ rune) bool { return r == 'G' && n == 'J' }, 'A', 2},
		{func(r, n, _ rune) bool { return r == 'D' && n == 'J' }, 'A', 2},
		{func(r, n, _ rune) bool { return r == 'L' && n == 'J' }, 'A', 2},
		{func(r, n, _ rune) bool { return r == 'H' && n == 'J' }, 'A', 2},
		// NG/NK nasal
		{func(r, n, _ rune) bool { return r == 'N' && (n == 'G' || n == 'K') }, 'B', 2},
		// CK
		{func(r, n, _ rune) bool { return r == 'C' && n == 'K' }, '3', 2},
		// PH
		{func(r, n, _ rune) bool { return r == 'P' && n == 'H' }, 'C', 2},
	},
}

// swedishDigit returns the phonetic code for a single Swedish character.
// Context-sensitive: K before front vowel → tje-sound ('E'), C before front vowel → S.
func swedishDigit(_, r, next rune) byte {
	switch r {
	case 'K':
		if isFrontVowelSwe(next) {
			return 'E'
		}
		return '3'
	case 'C':
		if isFrontVowelSwe(next) {
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

func swedishVowelCheck(r rune) bool {
	switch r {
	case 'A', 'E', 'I', 'O', 'U', 'Y', 'Å', 'Ä', 'Ö':
		return true
	}
	return false
}

func isFrontVowelSwe(r rune) bool {
	return r == 'E' || r == 'I' || r == 'Y' || r == 'Ä' || r == 'Ö'
}
