package soundex

// normBuf is the output buffer for normalization. Stack-allocated.
type normBuf struct {
	data [128]byte
	len  int
}

func (nb *normBuf) emit(b byte) {
	if nb.len < 128 {
		nb.data[nb.len] = b
		nb.len++
	}
}

func (nb *normBuf) emitTwo(a, b byte) {
	if nb.len < 127 {
		nb.data[nb.len] = a
		nb.data[nb.len+1] = b
		nb.len += 2
	}
}

// normalizeFinnish normalizes a UTF-8 word for Finnish phonetic encoding.
// Foreign letters → Finnish equivalents. Output is lowercased Finnish alphabet.
// Zero allocations.
func normalizeFinnish(word []byte, out *normBuf) {
	out.len = 0
	n := len(word)
	pos := 0

	for pos < n {
		r, size := decodeRune(word, pos)
		pos += size

		// Lowercase
		if r >= 'A' && r <= 'Z' {
			r = r + 32
		}

		// Multi-char lookahead for digraphs (on lowercased input)
		if pos < n {
			next, nextSize := decodeRune(word, pos)
			if next >= 'A' && next <= 'Z' {
				next = next + 32
			}

			switch {
			case r == 's' && next == 'c' && pos+nextSize < n:
				// Check for "sch"
				nn, _ := decodeRune(word, pos+nextSize)
				if nn >= 'A' && nn <= 'Z' {
					nn = nn + 32
				}
				if nn == 'h' {
					out.emit('s')
					pos += nextSize
					// skip 'h'
					_, hs := decodeRune(word, pos)
					pos += hs
					continue
				}
			case r == 'c' && next == 'k':
				out.emitTwo('k', 'k')
				pos += nextSize
				continue
			case r == 's' && next == 'h':
				out.emit('s')
				pos += nextSize
				continue
			case r == 'c' && next == 'h':
				out.emit('k')
				pos += nextSize
				continue
			case r == 'p' && next == 'h':
				out.emit('f')
				pos += nextSize
				continue
			case r == 't' && next == 'h':
				out.emit('t')
				pos += nextSize
				continue
			}
		}

		switch {
		// Native Finnish vowels — keep as-is
		case r == 'a', r == 'e', r == 'i', r == 'o', r == 'u', r == 'y':
			out.emit(byte(r))
		case r == 0xE4: // ä — stored as 1 (internal marker, not ASCII)
			out.emit(1)
		case r == 0xF6: // ö — stored as 2 (internal marker)
			out.emit(2)

		// Native Finnish consonants — keep
		case r == 'd', r == 'f', r == 'h', r == 'j', r == 'k', r == 'l',
			r == 'm', r == 'n', r == 'p', r == 'r', r == 's', r == 't', r == 'v':
			out.emit(byte(r))

		// Foreign consonant substitutions
		case r == 'b':
			out.emit('p')
		case r == 'g':
			out.emit('k')
		case r == 'q':
			out.emit('k')
		case r == 'w':
			out.emit('v')
		case r == 'x':
			out.emitTwo('k', 's')
		case r == 'z':
			out.emitTwo('t', 's')
		case r == 'c':
			// Context: before e/i/y → s, else → k
			if pos < n {
				nx, _ := decodeRune(word, pos)
				if nx >= 'A' && nx <= 'Z' {
					nx = nx + 32
				}
				if nx == 'e' || nx == 'i' || nx == 'y' {
					out.emit('s')
				} else {
					out.emit('k')
				}
			} else {
				out.emit('k')
			}

		// Scandinavian
		case r == 0xE5: // å → o
			out.emit('o')
		case r == 0xF8: // ø → ö
			out.emit(2)
		case r == 0xE6: // æ → ä
			out.emit(1)

		// German
		case r == 0xFC: // ü → y
			out.emit('y')
		case r == 0xDF: // ß → ss
			out.emitTwo('s', 's')

		// Accented vowels → base
		case r == 0xE0 || r == 0xE1 || r == 0xE2 || r == 0xE3: // à á â ã → a
			out.emit('a')
		case r == 0xE8 || r == 0xE9 || r == 0xEA || r == 0xEB: // è é ê ë → e
			out.emit('e')
		case r == 0xEC || r == 0xED || r == 0xEE || r == 0xEF: // ì í î ï → i
			out.emit('i')
		case r == 0xF2 || r == 0xF3 || r == 0xF4 || r == 0xF5: // ò ó ô õ → o
			out.emit('o')
		case r == 0xF9 || r == 0xFA || r == 0xFB: // ù ú û → u
			out.emit('u')
		case r == 0xFD: // ý → y
			out.emit('y')

		// Diacritical consonants
		case r == 0xE7: // ç → s
			out.emit('s')
		case r == 0xF1: // ñ → n
			out.emit('n')
		case r == 0xF0: // ð → t
			out.emit('t')
		case r == 0xFE: // þ → t
			out.emit('t')

		// Central/Eastern European (multi-byte runes)
		case r == 0x161: // š → s
			out.emit('s')
		case r == 0x17E: // ž → ts
			out.emitTwo('t', 's')
		case r == 0x10D: // č → s
			out.emit('s')
		case r == 0x159: // ř → r
			out.emit('r')
		case r == 0x142: // ł → l
			out.emit('l')
		case r == 0x144: // ń → n
			out.emit('n')
		case r == 0x15B: // ś → s
			out.emit('s')
		case r == 0x17A: // ź → ts
			out.emitTwo('t', 's')
		case r == 0x17C: // ż → ts
			out.emitTwo('t', 's')
		case r == 0x105: // ą → a
			out.emit('a')
		case r == 0x119: // ę → e
			out.emit('e')
		case r == 0x107: // ć → s
			out.emit('s')
		case r == 0x111: // đ → t
			out.emit('t')
		case r == 0x103: // ă → a
			out.emit('a')
		case r == 0x21B: // ț → ts
			out.emitTwo('t', 's')
		case r == 0x219: // ș → s
			out.emit('s')
		case r == 0x11F: // ğ → k
			out.emit('k')
		case r == 0x131: // ı → i
			out.emit('i')
		case r == 0x15F: // ş → s
			out.emit('s')

		// Baltic long vowels → base
		case r == 0x101: // ā → a
			out.emit('a')
		case r == 0x113: // ē → e
			out.emit('e')
		case r == 0x12B: // ī → i
			out.emit('i')
		case r == 0x16B: // ū → u
			out.emit('u')

		// Baltic palatalized → base
		case r == 0x123: // ģ → k
			out.emit('k')
		case r == 0x137: // ķ → k
			out.emit('k')
		case r == 0x13C: // ļ → l
			out.emit('l')
		case r == 0x146: // ņ → n
			out.emit('n')

		// Lithuanian
		case r == 0x117: // ė → e
			out.emit('e')
		case r == 0x12F: // į → i
			out.emit('i')
		case r == 0x173: // ų → u
			out.emit('u')

		// Hyphen — keep (compound word separator)
		case r == '-':
			out.emit('-')

		// Everything else — drop
		default:
			continue
		}
	}
}
