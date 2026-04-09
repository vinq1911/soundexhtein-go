package soundex

import "testing"

func TestNorwegian(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"kjøre", "KE6"},       // KJ → E, R → 6
		{"sjø", "SD"},          // SJ → D
		{"norsk", "N5F3"},      // N→5, RS→F(retroflex), K→3
		{"sang", "S8B"},        // S→8, NG→B
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Norwegian([]byte(tt.input))
			if tt.want == "" {
				if got.Len() != 0 {
					t.Errorf("Norwegian(%q) = %q, want empty", tt.input, got.String())
				}
				return
			}
			if got.String() != tt.want {
				t.Errorf("Norwegian(%q) = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}

func TestNorwegianZeroAllocs(t *testing.T) {
	name := []byte("kjøre")
	allocs := testing.AllocsPerRun(100, func() {
		_ = Norwegian(name)
	})
	if allocs != 0 {
		t.Errorf("Norwegian: got %v allocs, want 0", allocs)
	}
}
