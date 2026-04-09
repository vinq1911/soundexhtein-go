package soundex

import "sort"

// PhoneticIndex is a precomputed index for fast corpus search.
// Build once with NewIndex, then search many times.
type PhoneticIndex struct {
	codes   []PackedCode
	indices []int32 // original indices (int32 to save memory)
	algo    Algorithm
}

// NewIndex builds a PhoneticIndex from a corpus of names.
// Names are encoded with the given algorithm and sorted by code for early termination.
func NewIndex(names [][]byte, algo Algorithm) *PhoneticIndex {
	n := len(names)
	idx := &PhoneticIndex{
		codes:   make([]PackedCode, n),
		indices: make([]int32, n),
		algo:    algo,
	}

	// Encode all names.
	for i, name := range names {
		idx.codes[i] = Pack(Encode(name, algo))
		idx.indices[i] = int32(i)
	}

	// Sort by code value for locality and early termination.
	sort.Sort(idx)

	return idx
}

// Len, Less, Swap for sort.Interface.
func (idx *PhoneticIndex) Len() int { return len(idx.codes) }
func (idx *PhoneticIndex) Less(i, j int) bool {
	return uint32(idx.codes[i]) < uint32(idx.codes[j])
}
func (idx *PhoneticIndex) Swap(i, j int) {
	idx.codes[i], idx.codes[j] = idx.codes[j], idx.codes[i]
	idx.indices[i], idx.indices[j] = idx.indices[j], idx.indices[i]
}

// Search finds all corpus entries within maxDist Levenshtein distance of the query.
// Returns matches with original indices. Uses early termination on sorted codes.
// For zero allocations, use SearchInto with a pre-allocated result slice.
func (idx *PhoneticIndex) Search(query []byte, maxDist int) []Match {
	results := make([]Match, 0, 16)
	return idx.SearchInto(query, maxDist, results)
}

// SearchInto is like Search but appends to a caller-provided slice. Zero allocations
// if the slice has sufficient capacity.
func (idx *PhoneticIndex) SearchInto(query []byte, maxDist int, results []Match) []Match {
	qCode := Pack(Encode(query, idx.algo))
	results = results[:0]

	// Early termination: if the first byte of the code differs by more than maxDist,
	// the Levenshtein distance must be >= that difference.
	qFirst := byte(uint32(qCode) >> 24)

	for i, code := range idx.codes {
		// Quick reject: check first byte difference.
		cFirst := byte(uint32(code) >> 24)
		diff := int(cFirst) - int(qFirst)
		if diff < 0 {
			diff = -diff
		}
		if diff > maxDist {
			continue
		}

		// Full distance check.
		d := PackDistance(qCode, code)
		if d <= maxDist {
			results = append(results, Match{
				Index:    int(idx.indices[i]),
				Distance: d,
			})
		}
	}

	return results
}

// SearchPacked searches with a pre-encoded query. Zero allocations with pre-sized results.
func (idx *PhoneticIndex) SearchPacked(qCode PackedCode, maxDist int, results []Match) []Match {
	results = results[:0]
	qFirst := byte(uint32(qCode) >> 24)

	for i, code := range idx.codes {
		cFirst := byte(uint32(code) >> 24)
		diff := int(cFirst) - int(qFirst)
		if diff < 0 {
			diff = -diff
		}
		if diff > maxDist {
			continue
		}

		d := PackDistance(qCode, code)
		if d <= maxDist {
			results = append(results, Match{
				Index:    int(idx.indices[i]),
				Distance: d,
			})
		}
	}

	return results
}

// Size returns the number of entries in the index.
func (idx *PhoneticIndex) Size() int { return len(idx.codes) }
