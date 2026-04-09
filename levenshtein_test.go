package soundex

import "testing"

func makeCode(s string) Code {
	var c Code
	c[0] = byte(len(s))
	copy(c[1:], s)
	return c
}

func TestDistance(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"R163", "R163", 0},
		{"R163", "R150", 2},
		{"S530", "S530", 0},
		{"A000", "B000", 1},
		{"A000", "A100", 1},
		{"", "", 0},
		{"A000", "", 4},
		{"", "B200", 4},
		{"ABC", "ABCD", 1},
		{"ABCD", "ABC", 1},
		{"kitten", "sittin", 2},
	}

	for _, tt := range tests {
		t.Run(tt.a+"_"+tt.b, func(t *testing.T) {
			got := Distance(makeCode(tt.a), makeCode(tt.b))
			if got != tt.want {
				t.Errorf("Distance(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDistanceZeroAllocs(t *testing.T) {
	a := makeCode("R163")
	b := makeCode("S530")
	allocs := testing.AllocsPerRun(100, func() {
		_ = Distance(a, b)
	})
	if allocs != 0 {
		t.Errorf("Distance: got %v allocs, want 0", allocs)
	}
}

func TestSoundexDistance(t *testing.T) {
	// Robert and Rupert should have the same Soundex → distance 0
	d := SoundexDistance([]byte("Robert"), []byte("Rupert"), AlgoAmerican)
	if d != 0 {
		t.Errorf("SoundexDistance(Robert, Rupert) = %d, want 0", d)
	}

	// Smith and Lee should differ
	d = SoundexDistance([]byte("Smith"), []byte("Lee"), AlgoAmerican)
	if d == 0 {
		t.Errorf("SoundexDistance(Smith, Lee) = 0, want non-zero")
	}
}
