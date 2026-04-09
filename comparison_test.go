package soundex

import (
	"fmt"
	"testing"
)

// Speed/accuracy comparison between the three encoding tracks:
// 1. Fast (Code [8]byte, 4-7 digits) — ultra-fast, lossy
// 2. Dual (head+tail, 4+4 digits) — fast, captures both ends
// 3. Full (FullCode [64]byte, no truncation) — accurate, variable-length

// --- Accuracy comparison ---

type testPair struct {
	a, b     string
	relation string // "same", "similar", "different"
}

var finnishAccuracyPairs = []testPair{
	// Geminates (same sound, different spelling weight)
	{"talo", "tallo", "same"},
	{"mato", "matto", "same"},
	{"tuli", "tuuli", "same"},
	{"kuka", "kukka", "same"},

	// Vowel harmony pairs (must distinguish)
	{"älli", "alli", "different"},
	{"mäki", "maki", "different"},
	{"sää", "saa", "different"},
	{"käsi", "kasi", "different"},
	{"pöytä", "poyta", "different"},

	// Shared morpheme compounds (should be similar)
	{"lentokone", "tietokone", "similar"},
	{"lentokone", "lentokenttä", "similar"},
	{"ruokakauppa", "vaatekauppa", "similar"},
	{"makuuhuone", "kylpyhuone", "similar"},
	{"ilmavoimat", "merivoimat", "similar"},
	{"perustuslaki", "rikoslaki", "similar"},

	// Totally different words
	{"lentokone", "hammaslääkäri", "different"},
	{"käyttöjärjestelmä", "hautausmaa", "different"},
	{"sähköposti", "pelastuslaitos", "different"},
	{"Helsinki", "Rovaniemi", "different"},
	{"tietokone", "sanomalehti", "different"},

	// Foreign name handling (normalization makes these identical or near-identical)
	{"Schmidt", "Smitti", "similar"},
	{"Schwarzenegger", "Svartsenekker", "same"},
	{"Müller", "Myller", "same"},
}

func TestAccuracyComparison(t *testing.T) {
	type result struct {
		mode     string
		correct  int
		total    int
		details  []string
	}

	modes := []struct {
		name     string
		distance func(a, b []byte) int
	}{
		{"fast", func(a, b []byte) int {
			return Distance(Finnish(a), Finnish(b))
		}},
		{"dual", func(a, b []byte) int {
			return DualDistance(DualEncode(a, AlgoFinnish), DualEncode(b, AlgoFinnish))
		}},
		{"full-strict", func(a, b []byte) int {
			pa := FullFinnish(a)
			pb := FullFinnish(b)
			return FullDistance(pa.Strict, pb.Strict)
		}},
		{"full-relaxed", func(a, b []byte) int {
			pa := FullFinnish(a)
			pb := FullFinnish(b)
			return FullDistance(pa.Relaxed, pb.Relaxed)
		}},
	}

	for _, m := range modes {
		correct := 0
		total := len(finnishAccuracyPairs)
		var failures []string

		for _, p := range finnishAccuracyPairs {
			d := m.distance([]byte(p.a), []byte(p.b))
			ok := false

			switch p.relation {
			case "same":
				ok = d == 0
			case "similar":
				ok = d > 0 && d <= 8
			case "different":
				ok = d >= 1
			}

			if ok {
				correct++
			} else {
				failures = append(failures, fmt.Sprintf(
					"  %s(%q,%q)=%d want=%s", m.name, p.a, p.b, d, p.relation))
			}
		}

		pct := 100.0 * float64(correct) / float64(total)
		t.Logf("%-14s accuracy: %d/%d (%.0f%%)", m.name, correct, total, pct)
		for _, f := range failures {
			t.Log(f)
		}
	}
}

// --- Speed comparison ---

var benchWords = [][]byte{
	[]byte("Helsinki"),
	[]byte("käyttöjärjestelmä"),
	[]byte("lentokonesuihkuturbiinimoottoriapumekaanikkoaliupseerioppilas"),
	[]byte("Mäkinen"),
	[]byte("sanomalehti"),
}

func BenchmarkFastEncode(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		for _, w := range benchWords {
			_ = Finnish(w)
		}
	}
}

func BenchmarkDualEncodeTrack(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		for _, w := range benchWords {
			_ = DualEncode(w, AlgoFinnish)
		}
	}
}

func BenchmarkFullEncode(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		for _, w := range benchWords {
			_ = FullFinnish(w)
		}
	}
}

func BenchmarkFastDistance(b *testing.B) {
	codes := make([]Code, len(benchWords))
	for i, w := range benchWords {
		codes[i] = Finnish(w)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		for i := 0; i < len(codes)-1; i++ {
			_ = Distance(codes[i], codes[i+1])
		}
	}
}

func BenchmarkDualDistanceTrack(b *testing.B) {
	codes := make([]DualCode, len(benchWords))
	for i, w := range benchWords {
		codes[i] = DualEncode(w, AlgoFinnish)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		for i := 0; i < len(codes)-1; i++ {
			_ = DualDistance(codes[i], codes[i+1])
		}
	}
}

func BenchmarkFullDistance(b *testing.B) {
	codes := make([]FullCodePair, len(benchWords))
	for i, w := range benchWords {
		codes[i] = FullFinnish(w)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		for i := 0; i < len(codes)-1; i++ {
			_ = FullDistance(codes[i].Relaxed, codes[i+1].Relaxed)
		}
	}
}

// --- Combined encode+distance benchmark ---

func BenchmarkFastEndToEnd(b *testing.B) {
	a := []byte("lentokone")
	c := []byte("tietokone")
	b.ReportAllocs()
	for b.Loop() {
		_ = Distance(Finnish(a), Finnish(c))
	}
}

func BenchmarkDualEndToEnd(b *testing.B) {
	a := []byte("lentokone")
	c := []byte("tietokone")
	b.ReportAllocs()
	for b.Loop() {
		_ = DualDistance(DualEncode(a, AlgoFinnish), DualEncode(c, AlgoFinnish))
	}
}

func BenchmarkFullEndToEnd(b *testing.B) {
	a := []byte("lentokone")
	c := []byte("tietokone")
	b.ReportAllocs()
	for b.Loop() {
		pa := FullFinnish(a)
		pb := FullFinnish(c)
		_ = FullDistance(pa.Relaxed, pb.Relaxed)
	}
}
