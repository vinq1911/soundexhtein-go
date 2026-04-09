package soundex

import "testing"

func TestIndex(t *testing.T) {
	names := [][]byte{
		[]byte("Robert"),
		[]byte("Rupert"),
		[]byte("Smith"),
		[]byte("Smythe"),
		[]byte("Lee"),
		[]byte("Washington"),
		[]byte("Jackson"),
	}

	idx := NewIndex(names, AlgoAmerican)
	if idx.Size() != len(names) {
		t.Fatalf("Size = %d, want %d", idx.Size(), len(names))
	}

	// Search for "Robert" — should find Robert and Rupert (same Soundex).
	matches := idx.Search([]byte("Robert"), 0)
	found := map[int]bool{}
	for _, m := range matches {
		found[m.Index] = true
	}
	if !found[0] || !found[1] {
		t.Errorf("Search(Robert, 0): expected indices 0,1 (Robert,Rupert), got %v", matches)
	}

	// Search with distance 1 should find more.
	matches = idx.Search([]byte("Robert"), 1)
	if len(matches) < 2 {
		t.Errorf("Search(Robert, 1): expected >= 2 matches, got %d", len(matches))
	}
}

func TestIndexSearchInto(t *testing.T) {
	names := [][]byte{
		[]byte("Robert"),
		[]byte("Rupert"),
		[]byte("Smith"),
	}

	idx := NewIndex(names, AlgoAmerican)
	results := make([]Match, 0, 10)

	results = idx.SearchInto([]byte("Robert"), 0, results)
	if len(results) < 1 {
		t.Error("SearchInto: expected at least 1 match")
	}
}

func TestIndexSearchZeroAllocsPreSized(t *testing.T) {
	names := make([][]byte, 100)
	for i := range names {
		names[i] = []byte("Name")
	}
	idx := NewIndex(names, AlgoAmerican)
	results := make([]Match, 0, 200)
	query := []byte("Name")

	allocs := testing.AllocsPerRun(10, func() {
		results = idx.SearchInto(query, 1, results)
	})
	if allocs != 0 {
		t.Errorf("SearchInto: got %v allocs, want 0", allocs)
	}
}
