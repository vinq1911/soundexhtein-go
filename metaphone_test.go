package soundex

import "testing"

func TestMetaphone(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Smith", "SM0"}, // TH -> 0
		{"Schmidt", "SKMT"},
		{"Phone", "FN"},
		{"", ""},
		{"A", "A"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Metaphone([]byte(tt.input))
			if tt.want == "" {
				if got.Len() != 0 {
					t.Errorf("Metaphone(%q) = %q, want empty", tt.input, got.String())
				}
				return
			}
			if got.String() != tt.want {
				t.Errorf("Metaphone(%q) = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}

func TestMetaphoneZeroAllocs(t *testing.T) {
	name := []byte("Washington")
	allocs := testing.AllocsPerRun(100, func() {
		_ = Metaphone(name)
	})
	if allocs != 0 {
		t.Errorf("Metaphone: got %v allocs, want 0", allocs)
	}
}
