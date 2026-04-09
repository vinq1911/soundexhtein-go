package soundex

// Latvian computes a phonetic code optimized for Latvian words.
// Handles č, š, ž, ģ, ķ, ļ, ņ (palatalized consonants) and long vowels (ā, ē, ī, ū).
// Variable-length up to 7 digits. Zero allocations.
//
// Phonetic codes:
//   '1'=B,P  '2'=D,T  '3'=G,Ģ,K,Ķ,C,Q  '4'=L,Ļ  '5'=M,N,Ņ  '6'=R  '7'=V,W
//   '8'=S,Š,Z,Ž  '9'=H  'A'=J  'B'=Č (affricate)  'C'=F
//   'D'=DZ,DŽ (affricates)
func Latvian(name []byte) Code {
	return encodeWithConfig(name, &latvianConfig)
}

var latvianConfig = langEncoder{
	firstLetter: latvianFirstLetter,
	isVowel:     latvianVowelCheck,
	digit:       latvianDigitFn,
	rules: []encoderRule{
		// DZ/DŽ affricates
		{func(r, n, _ rune) bool { return r == 'D' && (n == 'Z' || n == 'Ž') }, 'D', 2},
	},
}

func latvianDigitFn(_, r, _ rune) byte {
	switch r {
	case 'Č':
		return 'B' // Č is a distinct affricate
	case 'Š':
		return '8'
	case 'Ž':
		return '8'
	case 'Ģ', 'Ķ':
		return '3'
	case 'Ļ':
		return '4'
	case 'Ņ':
		return '5'
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

func latvianVowelCheck(r rune) bool {
	switch r {
	case 'A', 'Ā', 'E', 'Ē', 'I', 'Ī', 'O', 'U', 'Ū':
		return true
	}
	return false
}

func latvianFirstLetter(r rune) byte {
	switch r {
	case 'Ā':
		return 'A'
	case 'Ē':
		return 'E'
	case 'Ī':
		return 'I'
	case 'Ū':
		return 'U'
	case 'Č':
		return 'C'
	case 'Š':
		return 'S'
	case 'Ž':
		return 'Z'
	case 'Ģ':
		return 'G'
	case 'Ķ':
		return 'K'
	case 'Ļ':
		return 'L'
	case 'Ņ':
		return 'N'
	}
	if r >= 'A' && r <= 'Z' {
		return byte(r)
	}
	return '?'
}
