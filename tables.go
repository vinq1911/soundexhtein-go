package soundex

// americanTable maps uppercase ASCII letters A-Z to Soundex digits.
// 0 means the letter is dropped (A, E, I, O, U, H, W, Y).
// Index: byte - 'A'
var americanTable = [26]byte{
	0, // A
	'1', // B
	'2', // C
	'3', // D
	0, // E
	'1', // F
	'2', // G
	0, // H
	0, // I
	'2', // J
	'2', // K
	'4', // L
	'5', // M
	'5', // N
	0, // O
	'1', // P
	'2', // Q
	'6', // R
	'2', // S
	'3', // T
	0, // U
	'1', // V
	0, // W
	'2', // X
	0, // Y
	'2', // Z
}

// cologneTable maps uppercase ASCII letters A-Z to Cologne phonetic digits.
// Uses 0xFF for context-dependent letters that need special handling.
// '0' = code 0, '1' = code 1, ... '8' = code 8
// 0xFF = context-dependent (D, T, X handled separately)
var cologneTable = [26]byte{
	'0', // A - vowel = 0
	'1', // B
	'8', // C - default 8, context-dependent (handled in code)
	0xFF, // D - context-dependent (D before S,C,Z = 8, else 2)
	'0', // E - vowel = 0
	'3', // F
	'4', // G
	'0', // H - ignored (code 0 but typically dropped)
	'0', // I - vowel = 0
	'0', // J - 0
	'4', // K
	'5', // L
	'6', // M
	'6', // N
	'0', // O - vowel = 0
	'1', // P - default 1, PH = 3 (handled in code)
	'4', // Q
	'7', // R
	'8', // S
	0xFF, // T - context-dependent (T before S,C,Z = 8, else 2)
	'0', // U - vowel = 0
	'3', // V
	'3', // W
	'8', // X - default 48, handled in code
	'0', // Y - vowel = 0
	'8', // Z
}

// metaphoneTable maps uppercase letters to their primary Metaphone code.
// 0 means dropped, 0xFF means context-dependent.
var metaphoneTable = [26]byte{
	0,    // A - vowel, only kept at start
	'P',  // B
	0xFF, // C - context-dependent
	0xFF, // D - context-dependent
	0,    // E - vowel
	'F',  // F
	0xFF, // G - context-dependent
	0xFF, // H - context-dependent
	0,    // I - vowel
	'J',  // J
	'K',  // K
	'L',  // L
	'M',  // M
	'N',  // N
	0,    // O - vowel
	'P',  // P, PH -> F
	'K',  // Q
	'R',  // R
	0xFF, // S - context-dependent
	0xFF, // T - context-dependent
	0,    // U - vowel
	'F',  // V
	0xFF, // W - context-dependent
	'K',  // X -> KS
	0,    // Y - treated as vowel in metaphone
	'S',  // Z
}
