package soundex

import "testing"

func TestDaitchMokotoff(t *testing.T) {
	tests := []struct {
		input     string
		wantCount int
		wantFirst string
	}{
		{"Cohen", 1, "560000"},
		{"Schwartz", 1, "479400"},
		{"", 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := DaitchMokotoff([]byte(tt.input))
			if tt.wantCount == 0 {
				if got.Count != 0 {
					t.Errorf("DM(%q): count = %d, want 0", tt.input, got.Count)
				}
				return
			}
			if got.Count < 1 {
				t.Fatalf("DM(%q): count = 0, want >= 1", tt.input)
			}
			first := got.Codes[0].String()
			if first != tt.wantFirst {
				t.Errorf("DM(%q)[0] = %q, want %q", tt.input, first, tt.wantFirst)
			}
		})
	}
}

func TestDaitchMokotoffZeroAllocs(t *testing.T) {
	name := []byte("Schwartz")
	allocs := testing.AllocsPerRun(100, func() {
		_ = DaitchMokotoff(name)
	})
	if allocs != 0 {
		t.Errorf("DaitchMokotoff: got %v allocs, want 0", allocs)
	}
}
