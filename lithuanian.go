package soundex

// Lithuanian computes a phonetic code optimized for Lithuanian words.
// Handles č, š, ž and nasal/long vowels (ą, ę, ė, į, ų, ū).
// Variable-length up to 7 digits. Zero allocations.
//
// Phonetic codes:
//   '1'=B,P  '2'=D,T  '3'=G,K,C,Q  '4'=L  '5'=M,N  '6'=R  '7'=V,W
//   '8'=S,Š,Z,Ž  '9'=H  'A'=J  'B'=Č (affricate)  'C'=F
//   'D'=DZ,DŽ (affricates)
func Lithuanian(name []byte) Code {
	return encodeWithConfig(name, &lithuanianConfig)
}

var lithuanianConfig = langEncoder{
	firstLetter: runeToBaseLetter,
	isVowel:     lithuanianVowelCheck,
	digit:       lithuanianDigitFn,
	rules: []encoderRule{
		// DZ/DŽ affricates
		{func(r, n, _ rune) bool { return r == 'D' && (n == 'Z' || n == 'Ž') }, 'D', 2},
	},
}

func lithuanianDigitFn(_, r, _ rune) byte {
	switch r {
	case 'Č':
		return 'B'
	case 'Š':
		return '8'
	case 'Ž':
		return '8'
	case 'B', 'P':
		return '1'
	case 'D', 'T':
		return '2'
	case 'G', 'K', 'C', 'Q', 'X':
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

func lithuanianVowelCheck(r rune) bool {
	switch r {
	case 'A', 'Ą', 'E', 'Ę', 'Ė', 'I', 'Į', 'O', 'U', 'Ų', 'Ū', 'Y':
		return true
	}
	return false
}

