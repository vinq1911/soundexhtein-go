package soundex

// runeToBaseLetter maps a rune to its ASCII base letter for code output.
// Shared by all language encoders for the first-letter position.
// Returns '?' for unrecognized runes.
func runeToBaseLetter(r rune) byte {
	if r >= 'A' && r <= 'Z' {
		return byte(r)
	}
	switch r {
	// Scandinavian
	case 'Ä', 'Æ':
		return 'A'
	case 'Ö', 'Ø':
		return 'O'
	case 'Å':
		return 'A'
	case 'Ü':
		return 'U'
	case 'Õ':
		return 'O'
	// Baltic diacriticals → base consonant
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
	// Baltic/Lithuanian long vowels → base
	case 'Ā':
		return 'A'
	case 'Ē':
		return 'E'
	case 'Ī':
		return 'I'
	case 'Ū':
		return 'U'
	case 'Ą':
		return 'A'
	case 'Ę':
		return 'E'
	case 'Ė':
		return 'E'
	case 'Į':
		return 'I'
	case 'Ų':
		return 'U'
	}
	return '?'
}
