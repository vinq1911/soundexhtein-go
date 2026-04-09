package soundex

import "math/bits"

// PackedCode stores a Soundex code in a single uint32.
// Bytes are packed big-endian: byte[0] in bits 31-24, byte[1] in 23-16, etc.
// This allows lexicographic comparison with simple integer comparison.
type PackedCode uint32

// PackedCodeLong stores a longer code (up to 8 chars) in a uint64.
type PackedCodeLong uint64

// Pack converts a Code to a PackedCode (first 4 bytes).
func Pack(c Code) PackedCode {
	n := c.Len()
	var p uint32
	for i := 0; i < 4; i++ {
		if i < n {
			p |= uint32(c[1+i]) << (24 - uint(i)*8)
		}
	}
	return PackedCode(p)
}

// PackLong converts a Code to a PackedCodeLong (up to 8 bytes).
func PackLong(c Code) PackedCodeLong {
	n := c.Len()
	var p uint64
	for i := 0; i < 7 && i < n; i++ {
		p |= uint64(c[1+i]) << (56 - uint(i)*8)
	}
	return PackedCodeLong(p)
}

// Unpack converts a PackedCode back to a Code.
func Unpack(p PackedCode) Code {
	var c Code
	v := uint32(p)
	n := byte(0)
	for i := 0; i < 4; i++ {
		b := byte(v >> (24 - uint(i)*8))
		if b != 0 {
			n++
			c[n] = b
		} else {
			break
		}
	}
	c[0] = n
	return c
}

// PackDistance computes the Levenshtein distance between two PackedCodes
// using a bit-parallel approach. For 4-byte codes this is extremely fast.
func PackDistance(a, b PackedCode) int {
	if a == b {
		return 0
	}

	// XOR to find differing bytes, then count.
	// For fixed-width 4-byte codes, we use a direct byte-by-byte approach
	// that the compiler can optimize heavily.
	av := uint32(a)
	bv := uint32(b)

	// Extract individual bytes and compare.
	aLen := packedLen(av)
	bLen := packedLen(bv)

	if aLen == 0 {
		return bLen
	}
	if bLen == 0 {
		return aLen
	}

	// For 4-byte fixed codes, unrolled Levenshtein is faster than bit-parallel.
	// 4x4 = 16 cells, fully unrollable.
	var ab [4]byte
	var bb [4]byte
	ab[0] = byte(av >> 24)
	ab[1] = byte(av >> 16)
	ab[2] = byte(av >> 8)
	ab[3] = byte(av)
	bb[0] = byte(bv >> 24)
	bb[1] = byte(bv >> 16)
	bb[2] = byte(bv >> 8)
	bb[3] = byte(bv)

	// Use bit-parallel Myers' algorithm for the general case.
	// For alphabet size <=64, one machine word suffices.
	return myersBitParallel(ab[:aLen], bb[:bLen])
}

// packedLen returns the number of non-zero bytes in a packed uint32.
func packedLen(v uint32) int {
	n := 0
	if v>>24 != 0 {
		n++
	}
	if (v>>16)&0xFF != 0 {
		n++
	}
	if (v>>8)&0xFF != 0 {
		n++
	}
	if v&0xFF != 0 {
		n++
	}
	return n
}

// myersBitParallel computes Levenshtein distance using Myers' bit-parallel algorithm.
// Optimal when len(b) <= 64 (fits in one machine word).
func myersBitParallel(a, b []byte) int {
	m := len(a)
	n := len(b)
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}

	// Build pattern bitmasks. For Soundex codes, alphabet is tiny (digits + letters).
	// We use a [256]uint64 but only touch used entries — stays in L1.
	var peq [256]uint64
	for i := 0; i < n; i++ {
		peq[b[i]] |= 1 << uint(i)
	}

	// Myers' algorithm
	var pv, mv uint64
	pv = ^uint64(0) // all 1s
	mv = 0
	score := n

	for i := 0; i < m; i++ {
		eq := peq[a[i]]
		xv := eq | mv
		xh := ((eq & pv) + pv) ^ pv | eq | mv
		ph := mv | ^(xh | pv)
		mh := pv & xh

		// Check last bit for score update
		if ph&(1<<uint(n-1)) != 0 {
			score++
		}
		if mh&(1<<uint(n-1)) != 0 {
			score--
		}

		// Shift
		ph = (ph << 1) | 1
		mh = mh << 1
		pv = mh | ^(xv | ph)
		mv = ph & xv
	}

	return score
}

// BatchPack encodes multiple names and packs them into PackedCodes.
// Caller provides the output slice. Zero allocations.
func BatchPack(names [][]byte, algo Algorithm, out []PackedCode) {
	for i, name := range names {
		if i >= len(out) {
			break
		}
		out[i] = Pack(Encode(name, algo))
	}
}

// BatchDistance computes distances from a query to all corpus entries.
// Caller provides the output slice. Zero allocations.
func BatchDistance(query PackedCode, corpus []PackedCode, out []int) {
	_ = out[len(corpus)-1] // bounds check elimination hint
	for i, c := range corpus {
		out[i] = PackDistance(query, c)
	}
}

// HammingDistance computes the Hamming distance between two PackedCodes.
// Faster than Levenshtein when codes are the same length (which they are for American Soundex).
func HammingDistance(a, b PackedCode) int {
	// XOR gives bits that differ, then count bytes that are non-zero.
	xor := uint32(a) ^ uint32(b)
	d := 0
	if xor>>24 != 0 {
		d++
	}
	if (xor>>16)&0xFF != 0 {
		d++
	}
	if (xor>>8)&0xFF != 0 {
		d++
	}
	if xor&0xFF != 0 {
		d++
	}
	return d
}

// PopcountDiff returns the popcount of XOR — a rough similarity metric.
// Useful for fast pre-filtering before exact distance computation.
func PopcountDiff(a, b PackedCode) int {
	return bits.OnesCount32(uint32(a) ^ uint32(b))
}
