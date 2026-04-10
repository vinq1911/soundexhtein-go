package soundex

// Estonian computes a phonetic code optimized for Estonian words.
// Handles õ, ä, ö, ü, š, ž. Collapses three-way vowel/consonant length
// (single/double/overlong) to a single code. Variable-length up to 7 digits.
// Zero allocations.
//
// Phonetic codes:
//   '1'=B,P  '2'=D,T  '3'=G,K,C,Q  '4'=L  '5'=M,N  '6'=R  '7'=V,W
//   '8'=S,Š,Z,Ž  '9'=H  'A'=J  'C'=F
func Estonian(name []byte) Code {
	return encodeWithConfig(name, &estonianConfig)
}

var estonianConfig = langEncoder{
	firstLetter: runeToBaseLetter,
	isVowel:     estonianVowelCheck,
	digit:       estonianDigitFn,
	rules:       nil, // no special digraphs
}

func estonianDigitFn(_, r, _ rune) byte {
	switch r {
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
	case 'S', 'Š', 'Z', 'Ž':
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

func estonianVowelCheck(r rune) bool {
	switch r {
	case 'A', 'E', 'I', 'O', 'U', 'Õ', 'Ä', 'Ö', 'Ü':
		return true
	}
	return false
}

