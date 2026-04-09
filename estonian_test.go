package soundex

import "testing"

func TestEstonian(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"Tallinn", "T245"},     // T→2, L→4 (double L→single), N→5 (double N→single)
		{"šokolaad", "S8342"},   // Š→8, K→3, L→4, D→2 (double A→single vowel)
		{"õlu", "O4"},           // Õ→vowel, L→4
		{"küla", "K34"},         // K→3, L→4
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Estonian([]byte(tt.input))
			if tt.want == "" {
				if got.Len() != 0 {
					t.Errorf("Estonian(%q) = %q, want empty", tt.input, got.String())
				}
				return
			}
			if got.String() != tt.want {
				t.Errorf("Estonian(%q) = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}

func TestEstonianZeroAllocs(t *testing.T) {
	name := []byte("Tallinn")
	allocs := testing.AllocsPerRun(100, func() {
		_ = Estonian(name)
	})
	if allocs != 0 {
		t.Errorf("Estonian: got %v allocs, want 0", allocs)
	}
}
