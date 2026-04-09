package soundex

import "testing"

func TestSwedish(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"sjuk", "SD3"},         // SJ → D, K → 3
		{"skjorta", "SD62"},     // SKJ → D, R → 6, T → 2
		{"kärlek", "KE643"},     // K before ä → E (tje), R → 6, L → 4, K → 3
		{"tjuv", "TE7"},         // TJ → E, V → 7
		{"kung", "K3B"},         // K → 3, NG → B
		{"Stockholm", "S823945"},// S→8, T→2, K→3, H→9, L→4, M→5
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Swedish([]byte(tt.input))
			if tt.want == "" {
				if got.Len() != 0 {
					t.Errorf("Swedish(%q) = %q, want empty", tt.input, got.String())
				}
				return
			}
			if got.String() != tt.want {
				t.Errorf("Swedish(%q) = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}

func TestSwedishZeroAllocs(t *testing.T) {
	name := []byte("Stockholm")
	allocs := testing.AllocsPerRun(100, func() {
		_ = Swedish(name)
	})
	if allocs != 0 {
		t.Errorf("Swedish: got %v allocs, want 0", allocs)
	}
}
