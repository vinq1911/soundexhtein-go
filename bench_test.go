package soundex

import "testing"

// --- Encoding Benchmarks ---

func BenchmarkAmerican(b *testing.B) {
	name := []byte("Washington")
	b.ReportAllocs()
	for b.Loop() {
		_ = American(name)
	}
}

func BenchmarkCologne(b *testing.B) {
	name := []byte("Mueller")
	b.ReportAllocs()
	for b.Loop() {
		_ = Cologne(name)
	}
}

func BenchmarkMetaphone(b *testing.B) {
	name := []byte("Washington")
	b.ReportAllocs()
	for b.Loop() {
		_ = Metaphone(name)
	}
}

func BenchmarkDaitchMokotoff(b *testing.B) {
	name := []byte("Schwartz")
	b.ReportAllocs()
	for b.Loop() {
		_ = DaitchMokotoff(name)
	}
}

// --- Distance Benchmarks ---

func BenchmarkDistance(b *testing.B) {
	a := makeCode("R163")
	c := makeCode("S530")
	b.ReportAllocs()
	for b.Loop() {
		_ = Distance(a, c)
	}
}

func BenchmarkSoundexDistance(b *testing.B) {
	a := []byte("Robert")
	c := []byte("Smith")
	b.ReportAllocs()
	for b.Loop() {
		_ = SoundexDistance(a, c, AlgoAmerican)
	}
}

// --- Packed Benchmarks ---

// --- Nordic/Baltic Encoding Benchmarks ---

func BenchmarkFinnish(b *testing.B) {
	name := []byte("Heikkilä")
	b.ReportAllocs()
	for b.Loop() {
		_ = Finnish(name)
	}
}

func BenchmarkSwedish(b *testing.B) {
	name := []byte("Stockholm")
	b.ReportAllocs()
	for b.Loop() {
		_ = Swedish(name)
	}
}

func BenchmarkNorwegian(b *testing.B) {
	name := []byte("kjøre")
	b.ReportAllocs()
	for b.Loop() {
		_ = Norwegian(name)
	}
}

func BenchmarkDanish(b *testing.B) {
	name := []byte("København")
	b.ReportAllocs()
	for b.Loop() {
		_ = Danish(name)
	}
}

func BenchmarkEstonian(b *testing.B) {
	name := []byte("Tallinn")
	b.ReportAllocs()
	for b.Loop() {
		_ = Estonian(name)
	}
}

func BenchmarkLatvian(b *testing.B) {
	name := []byte("šķēres")
	b.ReportAllocs()
	for b.Loop() {
		_ = Latvian(name)
	}
}

func BenchmarkLithuanian(b *testing.B) {
	name := []byte("Vilnius")
	b.ReportAllocs()
	for b.Loop() {
		_ = Lithuanian(name)
	}
}

func BenchmarkPack(b *testing.B) {
	code := American([]byte("Washington"))
	b.ReportAllocs()
	for b.Loop() {
		_ = Pack(code)
	}
}

func BenchmarkPackDistance(b *testing.B) {
	a := Pack(American([]byte("Robert")))
	c := Pack(American([]byte("Smith")))
	b.ReportAllocs()
	for b.Loop() {
		_ = PackDistance(a, c)
	}
}

func BenchmarkHammingDistance(b *testing.B) {
	a := Pack(American([]byte("Robert")))
	c := Pack(American([]byte("Smith")))
	b.ReportAllocs()
	for b.Loop() {
		_ = HammingDistance(a, c)
	}
}

func BenchmarkBatchDistance1000(b *testing.B) {
	query := Pack(American([]byte("Robert")))
	corpus := make([]PackedCode, 1000)
	for i := range corpus {
		corpus[i] = PackedCode(uint32(i) * 12345)
	}
	out := make([]int, 1000)
	b.ReportAllocs()
	for b.Loop() {
		BatchDistance(query, corpus, out)
	}
}

// --- Index Benchmarks ---

func BenchmarkNewIndex1000(b *testing.B) {
	names := make([][]byte, 1000)
	for i := range names {
		names[i] = []byte("Washington")
	}
	b.ReportAllocs()
	for b.Loop() {
		_ = NewIndex(names, AlgoAmerican)
	}
}

func BenchmarkIndexSearch1000(b *testing.B) {
	names := generateNames(1000)
	idx := NewIndex(names, AlgoAmerican)
	results := make([]Match, 0, 100)
	query := []byte("Robert")
	b.ReportAllocs()
	for b.Loop() {
		results = idx.SearchInto(query, 1, results)
	}
}

func BenchmarkIndexSearch10000(b *testing.B) {
	names := generateNames(10000)
	idx := NewIndex(names, AlgoAmerican)
	results := make([]Match, 0, 100)
	query := []byte("Robert")
	b.ReportAllocs()
	for b.Loop() {
		results = idx.SearchInto(query, 1, results)
	}
}

func generateNames(n int) [][]byte {
	pool := []string{
		"Robert", "Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia",
		"Miller", "Davis", "Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez",
		"Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin",
		"Lee", "Perez", "Thompson", "White", "Harris", "Sanchez", "Clark",
		"Ramirez", "Lewis", "Robinson", "Walker", "Young", "Allen", "King",
		"Wright", "Scott", "Torres", "Nguyen", "Hill", "Flores", "Green",
		"Adams", "Nelson", "Baker", "Hall", "Rivera", "Campbell", "Mitchell",
	}
	names := make([][]byte, n)
	for i := range names {
		names[i] = []byte(pool[i%len(pool)])
	}
	return names
}
