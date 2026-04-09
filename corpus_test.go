package soundex

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// loadCorpus reads a word list from testdata/{lang}.txt, one word per line.
func loadCorpus(t *testing.T, lang string) []string {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(thisFile), "testdata", lang+".txt")

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("loadCorpus(%q): %v", lang, err)
	}
	defer f.Close()

	var words []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			words = append(words, line)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("loadCorpus(%q): scan error: %v", lang, err)
	}
	if len(words) == 0 {
		t.Fatalf("loadCorpus(%q): no words found", lang)
	}
	return words
}

// testCorpusSanity runs sanity checks on a corpus for a given encoder.
func testCorpusSanity(t *testing.T, name string, corpus []string, encode func([]byte) Code) {
	t.Helper()

	codeDist := make(map[string]int)
	emptyCount := 0
	panicCount := 0

	for _, word := range corpus {
		func() {
			defer func() {
				if r := recover(); r != nil {
					panicCount++
					t.Errorf("%s(%q) panicked: %v", name, word, r)
				}
			}()
			code := encode([]byte(word))
			s := code.String()
			if code.Len() == 0 {
				emptyCount++
			} else {
				codeDist[s]++
			}
		}()
	}

	if panicCount > 0 {
		t.Errorf("%s: %d panics out of %d words", name, panicCount, len(corpus))
	}

	maxEmpty := len(corpus) / 20 // allow up to 5% empty
	if emptyCount > maxEmpty {
		t.Errorf("%s: %d empty codes out of %d words (max %d allowed)", name, emptyCount, len(corpus), maxEmpty)
	}

	uniqueCodes := len(codeDist)
	minUnique := len(corpus) / 10
	if uniqueCodes < minUnique && len(corpus) > 10 {
		t.Errorf("%s: only %d unique codes for %d words (want >= %d)", name, uniqueCodes, len(corpus), minUnique)
	}

	maxCluster := len(corpus) / 4
	for code, count := range codeDist {
		if count > maxCluster {
			t.Errorf("%s: code %q has %d/%d words (degenerate cluster)", name, code, count, len(corpus))
		}
	}

	t.Logf("%s: %d words, %d unique codes, %d empty, %d max cluster",
		name, len(corpus), uniqueCodes, emptyCount, maxInMap(codeDist))
}

func maxInMap(m map[string]int) int {
	max := 0
	for _, v := range m {
		if v > max {
			max = v
		}
	}
	return max
}

func TestCorpusFinnish(t *testing.T) {
	testCorpusSanity(t, "Finnish", loadCorpus(t, "finnish"), Finnish)
}

func TestCorpusSwedish(t *testing.T) {
	testCorpusSanity(t, "Swedish", loadCorpus(t, "swedish"), Swedish)
}

func TestCorpusNorwegian(t *testing.T) {
	testCorpusSanity(t, "Norwegian", loadCorpus(t, "norwegian"), Norwegian)
}

func TestCorpusDanish(t *testing.T) {
	testCorpusSanity(t, "Danish", loadCorpus(t, "danish"), Danish)
}

func TestCorpusEstonian(t *testing.T) {
	testCorpusSanity(t, "Estonian", loadCorpus(t, "estonian"), Estonian)
}

func TestCorpusLatvian(t *testing.T) {
	testCorpusSanity(t, "Latvian", loadCorpus(t, "latvian"), Latvian)
}

func TestCorpusLithuanian(t *testing.T) {
	testCorpusSanity(t, "Lithuanian", loadCorpus(t, "lithuanian"), Lithuanian)
}

func TestCorpusCrossLanguageSimilarWords(t *testing.T) {
	finnishPairs := [][2]string{
		{"talo", "tallo"},
		{"Mäkinen", "Mäkkinen"},
		{"koulu", "kouli"},
	}
	for _, p := range finnishPairs {
		a := Finnish([]byte(p[0]))
		b := Finnish([]byte(p[1]))
		d := Distance(a, b)
		if d > 1 {
			t.Errorf("Finnish distance(%q=%q, %q=%q) = %d, want <= 1",
				p[0], a.String(), p[1], b.String(), d)
		}
	}

	sjWords := []string{"sjuk", "sjö", "sjukhus", "sjöman"}
	for i := 0; i < len(sjWords); i++ {
		for j := i + 1; j < len(sjWords); j++ {
			a := Swedish([]byte(sjWords[i]))
			b := Swedish([]byte(sjWords[j]))
			d := Distance(a, b)
			t.Logf("Swedish(%q)=%q vs Swedish(%q)=%q dist=%d",
				sjWords[i], a.String(), sjWords[j], b.String(), d)
		}
	}
}

func TestCorpusSoundexLevenshteinIntegration(t *testing.T) {
	algos := []struct {
		name    string
		algo    Algorithm
		wordA   string
		wordB   string
		maxDist int
	}{
		{"Finnish", AlgoFinnish, "talo", "tallo", 0},
		{"Swedish", AlgoSwedish, "sjuk", "sjö", 3},
		{"Norwegian", AlgoNorwegian, "kjøre", "kjøpe", 3},
		{"Danish", AlgoDanish, "sang", "sang", 0},
		{"Estonian", AlgoEstonian, "Tallinn", "Tallinn", 0},
		{"Latvian", AlgoLatvian, "Rīga", "Rīga", 0},
		{"Lithuanian", AlgoLithuanian, "Vilnius", "Vilnius", 0},
	}

	for _, tt := range algos {
		t.Run(fmt.Sprintf("%s_%s_%s", tt.name, tt.wordA, tt.wordB), func(t *testing.T) {
			d := SoundexDistance([]byte(tt.wordA), []byte(tt.wordB), tt.algo)
			if d > tt.maxDist {
				t.Errorf("SoundexDistance(%q, %q, %s) = %d, want <= %d",
					tt.wordA, tt.wordB, tt.name, d, tt.maxDist)
			}
		})
	}
}
