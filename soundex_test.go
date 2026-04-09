package soundex

import "testing"

func TestAmerican(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Robert", "R163"},
		{"Rupert", "R163"},
		{"Rubin", "R150"},
		{"Ashcraft", "A261"},
		{"Ashcroft", "A261"},
		{"Tymczak", "T522"},
		{"Pfister", "P236"},
		{"Honeyman", "H555"},
		{"Smith", "S530"},
		{"Smythe", "S530"},
		{"Washington", "W252"},
		{"Lee", "L000"},
		{"Gutierrez", "G362"},
		{"Jackson", "J250"},
		{"", ""},
		{"A", "A000"},
		{"123", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := American([]byte(tt.input))
			if tt.want == "" {
				if got.Len() != 0 {
					t.Errorf("American(%q) = %q, want empty", tt.input, got.String())
				}
				return
			}
			if got.String() != tt.want {
				t.Errorf("American(%q) = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}

func TestAmericanZeroAllocs(t *testing.T) {
	name := []byte("Washington")
	allocs := testing.AllocsPerRun(100, func() {
		_ = American(name)
	})
	if allocs != 0 {
		t.Errorf("American: got %v allocs, want 0", allocs)
	}
}
