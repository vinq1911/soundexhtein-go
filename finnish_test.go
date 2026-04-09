package soundex

import "testing"

func TestFinnish(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"talo", "T2a4o"},
		{"tallo", "T2a4o"},        // double L → single
		{"Mäkinen", "M5w3i5e"},    // ä→w distinct from a
		{"Mäkkinen", "M5w3i5e"},   // double K → single
		{"Heikkilä", "H9ei3i4"},   // 7 chars max, trailing ä truncated
		{"Helsinki", "H9e48iB"},   // NK → B
		{"koulu", "K3ou4u"},
		{"kaupunki", "K3au1uB"},   // NK → B
		{"ääni", "Ww5i"},          // Ä→W first letter, ä→w vowel
		{"öljy", "Xx4Ay"},         // Ö→X first letter, ö→x vowel
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Finnish([]byte(tt.input))
			if tt.want == "" {
				if got.Len() != 0 {
					t.Errorf("Finnish(%q) = %q, want empty", tt.input, got.String())
				}
				return
			}
			if got.String() != tt.want {
				t.Errorf("Finnish(%q) = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}

func TestFinnishVowelDistinction(t *testing.T) {
	// Finnish ä/a, ö/o are distinct phonemes — codes must differ.
	distinctPairs := [][2]string{
		{"älli", "alli"},
		{"mäki", "maki"},
		{"sää", "saa"},
		{"pöytä", "poyta"},
		{"käsi", "kasi"},
	}
	for _, p := range distinctPairs {
		a := Finnish([]byte(p[0]))
		b := Finnish([]byte(p[1]))
		if a.Equal(b) {
			t.Errorf("Finnish(%q)=%q == Finnish(%q)=%q, want DIFFERENT (ä≠a, ö≠o)",
				p[0], a.String(), p[1], b.String())
		}
	}
}

func TestFinnishSimilarWords(t *testing.T) {
	// Words that sound the same should get the same code.
	samePairs := [][2]string{
		{"talo", "tallo"},       // single vs double L
		{"Mäkinen", "Mäkkinen"}, // single vs double K
		{"tuli", "tuuli"},       // short vs long vowel
	}
	for _, p := range samePairs {
		a := Finnish([]byte(p[0]))
		b := Finnish([]byte(p[1]))
		if !a.Equal(b) {
			t.Errorf("Finnish(%q)=%q != Finnish(%q)=%q, expected SAME",
				p[0], a.String(), p[1], b.String())
		}
	}
}

func TestFinnishZeroAllocs(t *testing.T) {
	name := []byte("Heikkilä")
	allocs := testing.AllocsPerRun(100, func() {
		_ = Finnish(name)
	})
	if allocs != 0 {
		t.Errorf("Finnish: got %v allocs, want 0", allocs)
	}
}
