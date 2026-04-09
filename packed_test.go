package soundex

import "testing"

func TestPackUnpack(t *testing.T) {
	cases := []string{"R163", "S530", "A000", "L000"}
	for _, s := range cases {
		c := makeCode(s)
		p := Pack(c)
		u := Unpack(p)
		if u.String() != s {
			t.Errorf("Pack/Unpack(%q): got %q", s, u.String())
		}
	}
}

func TestPackDistance(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"R163", "R163", 0},
		{"R163", "R150", 2},
		{"S530", "S530", 0},
		{"A000", "B000", 1},
	}
	for _, tt := range tests {
		t.Run(tt.a+"_"+tt.b, func(t *testing.T) {
			pa := Pack(makeCode(tt.a))
			pb := Pack(makeCode(tt.b))
			got := PackDistance(pa, pb)
			if got != tt.want {
				t.Errorf("PackDistance(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestPackDistanceZeroAllocs(t *testing.T) {
	a := Pack(makeCode("R163"))
	b := Pack(makeCode("S530"))
	allocs := testing.AllocsPerRun(100, func() {
		_ = PackDistance(a, b)
	})
	if allocs != 0 {
		t.Errorf("PackDistance: got %v allocs, want 0", allocs)
	}
}

func TestHammingDistance(t *testing.T) {
	a := Pack(makeCode("R163"))
	b := Pack(makeCode("R163"))
	if d := HammingDistance(a, b); d != 0 {
		t.Errorf("HammingDistance same codes = %d, want 0", d)
	}

	c := Pack(makeCode("S530"))
	if d := HammingDistance(a, c); d == 0 {
		t.Error("HammingDistance different codes = 0, want > 0")
	}
}

func TestBatchPack(t *testing.T) {
	names := [][]byte{[]byte("Robert"), []byte("Smith"), []byte("Lee")}
	out := make([]PackedCode, len(names))
	BatchPack(names, AlgoAmerican, out)

	for i, name := range names {
		expected := Pack(American(name))
		if out[i] != expected {
			t.Errorf("BatchPack[%d] mismatch", i)
		}
	}
}

func TestBatchDistance(t *testing.T) {
	query := Pack(American([]byte("Robert")))
	corpus := []PackedCode{
		Pack(American([]byte("Robert"))),
		Pack(American([]byte("Rupert"))),
		Pack(American([]byte("Smith"))),
	}
	out := make([]int, len(corpus))
	BatchDistance(query, corpus, out)

	if out[0] != 0 {
		t.Errorf("BatchDistance[0] = %d, want 0 (Robert=Robert)", out[0])
	}
	if out[1] != 0 {
		t.Errorf("BatchDistance[1] = %d, want 0 (Robert=Rupert)", out[1])
	}
	if out[2] == 0 {
		t.Error("BatchDistance[2] = 0, want non-zero (Robert!=Smith)")
	}
}
