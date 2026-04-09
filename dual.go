package soundex

// DualCode holds head (start) and tail (end) phonetic codes for a word.
// Together they capture both ends of any word, regardless of length.
// Each code is 4 digits, packable into uint32 for fast batch operations.
type DualCode struct {
	Head Code
	Tail Code
}

// DualEncode computes head and tail phonetic codes for a word.
// Head encodes from the start, Tail from the end. Both are capped at 4 digits.
// Zero allocations.
func DualEncode(word []byte, algo Algorithm) DualCode {
	var dc DualCode

	var buf runeBuffer
	fillRuneBuffer(word, &buf)
	if buf.len == 0 {
		return dc
	}

	// Head: encode normally, cap at 4.
	full := Encode(word, algo)
	dc.Head = truncateCode(full, 4)

	// Tail: reverse the rune buffer, encode, cap at 4.
	var rev runeBuffer
	rev.len = buf.len
	for i := 0; i < buf.len; i++ {
		rev.runes[i] = buf.runes[buf.len-1-i]
	}

	// Encode the reversed buffer using the language-specific encoder.
	var tmpBuf [256]byte
	revBytes := runeBufferToBytes(&rev, tmpBuf[:])
	fullTail := Encode(revBytes, algo)
	dc.Tail = truncateCode(fullTail, 4)

	return dc
}

// truncateCode caps a Code to maxLen digits.
func truncateCode(c Code, maxLen byte) Code {
	if c[0] > maxLen {
		c[0] = maxLen
	}
	return c
}

// runeBufferToBytes writes a runeBuffer as UTF-8 into dst and returns the used portion.
func runeBufferToBytes(rb *runeBuffer, dst []byte) []byte {
	pos := 0
	for i := 0; i < rb.len && pos < len(dst)-3; i++ {
		r := rb.runes[i]
		if r < 0x80 {
			dst[pos] = byte(r)
			pos++
		} else if r < 0x800 {
			dst[pos] = byte(0xC0 | (r >> 6))
			dst[pos+1] = byte(0x80 | (r & 0x3F))
			pos += 2
		} else if r < 0x10000 {
			dst[pos] = byte(0xE0 | (r >> 12))
			dst[pos+1] = byte(0x80 | ((r >> 6) & 0x3F))
			dst[pos+2] = byte(0x80 | (r & 0x3F))
			pos += 3
		}
	}
	return dst[:pos]
}

// DualDistance computes the combined distance between two DualCodes.
// Returns headDist + tailDist, giving equal weight to start and end of word.
func DualDistance(a, b DualCode) int {
	return Distance(a.Head, b.Head) + Distance(a.Tail, b.Tail)
}

// DualSoundexDistance encodes both words with dual encoding and returns combined distance.
func DualSoundexDistance(a, b []byte, algo Algorithm) int {
	return DualDistance(DualEncode(a, algo), DualEncode(b, algo))
}

// PackedDualCode holds head and tail as packed uint32s for fast batch operations.
type PackedDualCode struct {
	Head PackedCode
	Tail PackedCode
}

// PackDual converts a DualCode to PackedDualCode.
func PackDual(dc DualCode) PackedDualCode {
	return PackedDualCode{
		Head: Pack(dc.Head),
		Tail: Pack(dc.Tail),
	}
}

// PackedDualDistance computes combined distance on packed codes.
func PackedDualDistance(a, b PackedDualCode) int {
	return PackDistance(a.Head, b.Head) + PackDistance(a.Tail, b.Tail)
}

// BatchDualPack encodes multiple words into PackedDualCodes.
func BatchDualPack(words [][]byte, algo Algorithm, out []PackedDualCode) {
	for i, w := range words {
		if i >= len(out) {
			break
		}
		out[i] = PackDual(DualEncode(w, algo))
	}
}

// BatchDualDistance computes dual distances from query to all corpus entries.
func BatchDualDistance(query PackedDualCode, corpus []PackedDualCode, out []int) {
	for i, c := range corpus {
		if i >= len(out) {
			break
		}
		out[i] = PackedDualDistance(query, c)
	}
}
