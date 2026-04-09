package soundex

import "testing"

func TestLithuanian(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"Vilnius", "V7458"},      // V→7, L→4, N→5, S→8
		{"čia", "CB"},             // Č→B (affricate)
		{"šalis", "S848"},         // Š→8, L→4, S→8
		{"džiaugsmas", "DD3858"},  // DŽ→D, G→3, S→8, M→5, S→8
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Lithuanian([]byte(tt.input))
			if tt.want == "" {
				if got.Len() != 0 {
					t.Errorf("Lithuanian(%q) = %q, want empty", tt.input, got.String())
				}
				return
			}
			if got.String() != tt.want {
				t.Errorf("Lithuanian(%q) = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}

func TestLithuanianZeroAllocs(t *testing.T) {
	name := []byte("Vilnius")
	allocs := testing.AllocsPerRun(100, func() {
		_ = Lithuanian(name)
	})
	if allocs != 0 {
		t.Errorf("Lithuanian: got %v allocs, want 0", allocs)
	}
}
