package soundex

import "testing"

func TestDanish(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"København", "K315975"}, // K→3, B→1, N→5, H→9, V→7, N→5
		{"sjov", "SD7"},          // SJ → D, V → 7
		{"sang", "S8B"},          // S→8, NG→B
		{"mad", "M5E"},           // M→5, D after vowel → soft D (E)
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Danish([]byte(tt.input))
			if tt.want == "" {
				if got.Len() != 0 {
					t.Errorf("Danish(%q) = %q, want empty", tt.input, got.String())
				}
				return
			}
			if got.String() != tt.want {
				t.Errorf("Danish(%q) = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}

func TestDanishZeroAllocs(t *testing.T) {
	name := []byte("København")
	allocs := testing.AllocsPerRun(100, func() {
		_ = Danish(name)
	})
	if allocs != 0 {
		t.Errorf("Danish: got %v allocs, want 0", allocs)
	}
}
