package soundex

// FullFinnish encodes a word using full-length Finnish phonetic encoding.
// Returns both strict (preserves geminates) and relaxed (collapses geminates) codes.
// No truncation — the entire word is encoded. Zero allocations.
//
// Pipeline: normalize(foreign → Finnish) → encode(Finnish → phonetic codes)
//
// Consonant codes:
//   P (p), T (d,t), K (k), V (f,v), S (s), H (h), J (j),
//   L (l), R (r), M (m), N (n), NK (ng,nk)
//
// Vowel codes (all 8 distinct — vowel harmony preserved):
//   A (a), E (e), I (i), O (o), U (u), Y (y), Ä (ä), Ö (ö)
//
// Strict mode: geminates preserved (kk→KK, aa→AA), different-letter same-code collapsed (d+t→T)
// Relaxed mode: geminates collapsed (kk→K, aa→A), same-code always collapsed
func FullFinnish(word []byte) FullCodePair {
	var pair FullCodePair

	var nb normBuf
	normalizeFinnish(word, &nb)
	if nb.len == 0 {
		return pair
	}

	data := nb.data[:nb.len]
	sLen := byte(0) // strict length
	rLen := byte(0) // relaxed length
	var lastStrict byte
	var lastRelaxed byte

	i := 0
	for i < len(data) && sLen < 63 && rLen < 63 {
		ch := data[i]

		// Skip hyphens (compound separators) — reset adjacency tracking
		if ch == '-' {
			lastStrict = 0
			lastRelaxed = 0
			i++
			continue
		}

		// Check for NK/NG digraph
		if ch == 'n' && i+1 < len(data) && (data[i+1] == 'k' || data[i+1] == 'g') {
			code1 := byte('N')
			code2 := byte('K')

			// Strict: emit NK (always 2 code bytes)
			if sLen < 62 {
				sLen++
				pair.Strict[sLen] = code1
				sLen++
				pair.Strict[sLen] = code2
			}
			lastStrict = code2

			// Relaxed: emit NK unless last was already NK
			if lastRelaxed != code2 || (rLen >= 2 && pair.Relaxed[rLen-1] != code1) {
				if rLen < 62 {
					rLen++
					pair.Relaxed[rLen] = code1
					rLen++
					pair.Relaxed[rLen] = code2
				}
			}
			lastRelaxed = code2

			i += 2
			continue
		}

		code := fullFinnishCode(ch)
		if code == 0 {
			i++
			continue
		}

		isGeminate := i+1 < len(data) && data[i+1] == ch

		if isGeminate {
			// Strict: emit twice
			if sLen < 62 {
				sLen++
				pair.Strict[sLen] = code
				sLen++
				pair.Strict[sLen] = code
			}
			lastStrict = code

			// Relaxed: emit once (collapse geminate)
			if code != lastRelaxed {
				rLen++
				pair.Relaxed[rLen] = code
			}
			lastRelaxed = code

			i += 2 // skip both chars
		} else {
			// Strict: emit unless same code as last (different-letter same-code collapse)
			if code != lastStrict {
				sLen++
				pair.Strict[sLen] = code
			}
			lastStrict = code

			// Relaxed: emit unless same code as last
			if code != lastRelaxed {
				rLen++
				pair.Relaxed[rLen] = code
			}
			lastRelaxed = code

			i++
		}
	}

	pair.Strict[0] = sLen
	pair.Relaxed[0] = rLen
	return pair
}

// fullFinnishCode maps a normalized Finnish byte to its phonetic code.
func fullFinnishCode(ch byte) byte {
	switch ch {
	// Consonants
	case 'p':
		return 'P'
	case 't', 'd':
		return 'T'
	case 'k':
		return 'K'
	case 'v', 'f':
		return 'V'
	case 's':
		return 'S'
	case 'h':
		return 'H'
	case 'j':
		return 'J'
	case 'l':
		return 'L'
	case 'r':
		return 'R'
	case 'm':
		return 'M'
	case 'n':
		return 'N'
	// Vowels — each distinct
	case 'a':
		return 'A'
	case 'e':
		return 'E'
	case 'i':
		return 'I'
	case 'o':
		return 'O'
	case 'u':
		return 'U'
	case 'y':
		return 'Y'
	case 1: // ä (internal marker from normalize)
		return 'W' // Ä → W (distinct from A)
	case 2: // ö (internal marker from normalize)
		return 'X' // Ö → X (distinct from O)
	}
	return 0
}
