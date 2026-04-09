package soundex

import "testing"

func TestLatvian(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"Rīga", "R63"},        // R→6, G→3
		{"čau", "CB"},           // Č→B (affricate)
		{"šķēres", "S8368"},     // Š→8, Ķ→3, R→6, S→8
		{"dzīvot", "DD72"},      // DZ→D(affricate), V→7, T→2
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Latvian([]byte(tt.input))
			if tt.want == "" {
				if got.Len() != 0 {
					t.Errorf("Latvian(%q) = %q, want empty", tt.input, got.String())
				}
				return
			}
			if got.String() != tt.want {
				t.Errorf("Latvian(%q) = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}

func TestLatvianZeroAllocs(t *testing.T) {
	name := []byte("Rīga")
	allocs := testing.AllocsPerRun(100, func() {
		_ = Latvian(name)
	})
	if allocs != 0 {
		t.Errorf("Latvian: got %v allocs, want 0", allocs)
	}
}
