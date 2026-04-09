package soundex

import "testing"

func TestDualEncode(t *testing.T) {
	dc := DualEncode([]byte("Helsinki"), AlgoFinnish)
	if dc.Head.Len() == 0 || dc.Tail.Len() == 0 {
		t.Errorf("DualEncode(Helsinki): head=%q tail=%q, both should be non-empty",
			dc.Head.String(), dc.Tail.String())
	}
	if dc.Head.Len() > 4 {
		t.Errorf("DualEncode head length %d > 4", dc.Head.Len())
	}
	if dc.Tail.Len() > 4 {
		t.Errorf("DualEncode tail length %d > 4", dc.Tail.Len())
	}
	t.Logf("Helsinki: head=%q tail=%q", dc.Head.String(), dc.Tail.String())
}

func TestDualDistinction(t *testing.T) {
	// Words with same start but different end should differ in tail.
	a := DualEncode([]byte("käyttäjä"), AlgoFinnish)
	b := DualEncode([]byte("käytäntö"), AlgoFinnish)
	if a.Head.Equal(b.Head) && a.Tail.Equal(b.Tail) {
		t.Errorf("käyttäjä and käytäntö should differ in at least head or tail")
	}
	t.Logf("käyttäjä: head=%q tail=%q", a.Head.String(), a.Tail.String())
	t.Logf("käytäntö: head=%q tail=%q", b.Head.String(), b.Tail.String())
}

func TestDualDistanceSameWord(t *testing.T) {
	a := DualEncode([]byte("talo"), AlgoFinnish)
	b := DualEncode([]byte("tallo"), AlgoFinnish)
	d := DualDistance(a, b)
	if d != 0 {
		t.Errorf("DualDistance(talo, tallo) = %d, want 0 (double consonant)", d)
	}
}

func TestDualDistanceDifferent(t *testing.T) {
	a := DualEncode([]byte("älli"), AlgoFinnish)
	b := DualEncode([]byte("alli"), AlgoFinnish)
	d := DualDistance(a, b)
	if d == 0 {
		t.Errorf("DualDistance(älli, alli) = 0, want > 0 (ä≠a)")
	}
	t.Logf("älli: head=%q tail=%q", a.Head.String(), a.Tail.String())
	t.Logf("alli: head=%q tail=%q", b.Head.String(), b.Tail.String())
}

func TestDualLongWord(t *testing.T) {
	// Long compound word should still produce meaningful head + tail.
	word := []byte("lentokonesuihkuturbiinimoottoriapumekaanikkoaliupseerioppilas")
	dc := DualEncode(word, AlgoFinnish)
	if dc.Head.Len() == 0 || dc.Tail.Len() == 0 {
		t.Error("DualEncode on long word: empty head or tail")
	}
	t.Logf("long word: head=%q tail=%q", dc.Head.String(), dc.Tail.String())
}

func TestDualZeroAllocs(t *testing.T) {
	word := []byte("Helsinki")
	allocs := testing.AllocsPerRun(100, func() {
		_ = DualEncode(word, AlgoFinnish)
	})
	if allocs != 0 {
		t.Errorf("DualEncode: got %v allocs, want 0", allocs)
	}
}

func TestPackedDualDistance(t *testing.T) {
	a := PackDual(DualEncode([]byte("talo"), AlgoFinnish))
	b := PackDual(DualEncode([]byte("tallo"), AlgoFinnish))
	d := PackedDualDistance(a, b)
	if d != 0 {
		t.Errorf("PackedDualDistance(talo, tallo) = %d, want 0", d)
	}
}

func BenchmarkDualEncode(b *testing.B) {
	word := []byte("käyttöjärjestelmä")
	b.ReportAllocs()
	for b.Loop() {
		_ = DualEncode(word, AlgoFinnish)
	}
}

func BenchmarkDualDistance(b *testing.B) {
	a := DualEncode([]byte("käyttäjä"), AlgoFinnish)
	c := DualEncode([]byte("käytäntö"), AlgoFinnish)
	b.ReportAllocs()
	for b.Loop() {
		_ = DualDistance(a, c)
	}
}
