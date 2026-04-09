package soundex

import "testing"

func TestCologne(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Mueller", "657"},
		{"Müller", "657"},       // non-ASCII stripped, treated as "ller"
		{"Wikipedia", "3412"},
		{"Breschnew", "17863"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Cologne([]byte(tt.input))
			if tt.want == "" {
				if got.Len() != 0 {
					t.Errorf("Cologne(%q) = %q, want empty", tt.input, got.String())
				}
				return
			}
			if got.String() != tt.want {
				t.Errorf("Cologne(%q) = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}

func TestCologneZeroAllocs(t *testing.T) {
	name := []byte("Mueller")
	allocs := testing.AllocsPerRun(100, func() {
		_ = Cologne(name)
	})
	if allocs != 0 {
		t.Errorf("Cologne: got %v allocs, want 0", allocs)
	}
}
