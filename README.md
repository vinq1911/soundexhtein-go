# soundex-levenshtein

Multi-language phonetic encoding with Soundex-Levenshtein distance. Zero allocations. Insanely fast.

## Install

```bash
go get github.com/vinq1911/soundexhtein-go
```

Requires Go 1.26+.

## Quick Start

```go
import soundex "github.com/vinq1911/soundexhtein-go"

// Encode a word
code := soundex.Finnish([]byte("Helsinki"))
fmt.Println(code.String()) // "H9e48iB"

// Compare two words
d := soundex.SoundexDistance([]byte("Robert"), []byte("Rupert"), soundex.AlgoAmerican)
fmt.Println(d) // 0 (same Soundex code)

// Full Finnish encoding (no truncation, strict + relaxed)
pair := soundex.FullFinnish([]byte("käyttöjärjestelmä"))
fmt.Println(pair.Strict.String())  // "KWYTTXJWRJESTELM W"
fmt.Println(pair.Relaxed.String()) // "KWYTXJWRJESTELMW"
```

## Three Encoding Tracks

Choose the right track for your use case:

| Track | Speed | Accuracy | Best For |
|-------|-------|----------|----------|
| **Fast** | ~9-40ns encode, ~15ns distance | 87% | Pre-filtering millions of candidates |
| **Dual** | ~200-350ns encode, ~30ns distance | 91% | Batch operations with uint32 packing |
| **Full** | ~140-900ns encode, ~220ns distance | 100% | Accurate matching, agglutinative languages |

All tracks: **0 allocations, 0 B/op**.

### Fast Track

Fixed-size `Code [8]byte`. Up to 7 phonetic digits. Fits in `PackedCode uint32` for batch ops.

```go
a := soundex.Finnish([]byte("talo"))
b := soundex.Finnish([]byte("tallo"))
d := soundex.Distance(a, b) // 0 — double consonant collapsed
```

### Dual Track

Head (first 4 digits) + tail (last 4 digits). Captures both ends of long words.

```go
a := soundex.DualEncode([]byte("lentokone"), soundex.AlgoFinnish)
b := soundex.DualEncode([]byte("tietokone"), soundex.AlgoFinnish)
d := soundex.DualDistance(a, b) // small — share "kone" tail
```

### Full Track

Variable-length `FullCode [64]byte`. No truncation. Normalize-then-encode pipeline. Strict (preserves geminates) + relaxed (collapses geminates) modes.

```go
a := soundex.FullFinnish([]byte("matto"))
b := soundex.FullFinnish([]byte("mato"))

// Strict: matto != mato (geminate preserved)
soundex.FullDistance(a.Strict, b.Strict) // 1

// Relaxed: matto == mato (geminate collapsed)
soundex.FullDistance(a.Relaxed, b.Relaxed) // 0
```

The full track normalizes foreign spellings to Finnish phonetics:

```go
a := soundex.FullFinnish([]byte("Schwarzenegger"))
b := soundex.FullFinnish([]byte("Svartsenekker"))
soundex.FullDistance(a.Relaxed, b.Relaxed) // 0 — identical after normalization
```

## Supported Languages

### Western European

| Algorithm | Function | Code Length | Notes |
|-----------|----------|-------------|-------|
| American Soundex | `American()` | 4 fixed | Standard 1-letter + 3-digit |
| Cologne Phonetic | `Cologne()` | Variable | German-optimized |
| Metaphone | `Metaphone()` | Up to 4 | Improved English |
| Daitch-Mokotoff | `DaitchMokotoff()` | 6 fixed | Multi-code, Slavic/Yiddish |

### Finnish

| Function | Track | Notes |
|----------|-------|-------|
| `Finnish()` | Fast | 7-digit, vowel harmony preserved (a != ä) |
| `FullFinnish()` | Full | No truncation, strict+relaxed, foreign normalization |

Finnish vowel harmony is fully preserved — `ä`, `ö`, `y` are encoded distinctly from `a`, `o`, `u`.

### Scandinavian

| Function | Notes |
|----------|-------|
| `Swedish()` | SJ/SKJ/STJ sje-sound, KJ/TJ tje-sound, K before front vowel |
| `Norwegian()` | KJ/SJ digraphs, RS retroflex |
| `Danish()` | Soft D (after vowel), SJ digraph |

### Baltic

| Function | Notes |
|----------|-------|
| `Estonian()` | š/ž, õ as distinct vowel, triple consonant collapse |
| `Latvian()` | č/š/ž, ģ/ķ/ļ/ņ palatalized, DZ/DŽ affricates |
| `Lithuanian()` | č/š/ž, DŽ affricate, ą/ę/ė/į/ų nasal vowels |

## Batch & Index APIs

### Packed Codes (uint32)

For high-throughput batch comparison. Pack codes into `uint32` for cache-friendly loops.

```go
query := soundex.Pack(soundex.American([]byte("Robert")))
corpus := make([]soundex.PackedCode, 10000)
soundex.BatchPack(names, soundex.AlgoAmerican, corpus)

results := make([]int, 10000)
soundex.BatchDistance(query, corpus, results)
```

### Precomputed Index

Build once, search many times. Sorted codes with early termination.

```go
idx := soundex.NewIndex(names, soundex.AlgoFinnish)
matches := idx.Search([]byte("Helsinki"), 1) // maxDist=1

// Zero-alloc search with pre-sized buffer
results := make([]soundex.Match, 0, 100)
results = idx.SearchInto([]byte("Helsinki"), 1, results)
```

## CLI Tool

```bash
go build -o encode ./cmd/encode/

# Encode words (file or stdin)
echo "Helsinki Tampere Turku" | ./encode -lang finnish
echo "Stockholm Göteborg" | ./encode -lang swedish -dual
cat words.txt | ./encode -lang finnish -full

# Compare two words
./encode -lang finnish -full -compare "matto,mato"
./encode -lang swedish -compare "sjuk,sjö"
./encode -lang finnish -dual -compare "lentokone,tietokone"
```

Flags:
- `-lang` — algorithm (american, cologne, metaphone, finnish, swedish, norwegian, danish, estonian, latvian, lithuanian)
- `-dual` — head+tail dual encoding
- `-full` — full-length strict+relaxed encoding
- `-compare word1,word2` — compare two words and show distance

## Benchmarks (Apple M4)

```
BenchmarkAmerican           9.1 ns/op    0 allocs
BenchmarkFinnish           38.0 ns/op    0 allocs
BenchmarkSwedish           42.0 ns/op    0 allocs
BenchmarkCologne           18.0 ns/op    0 allocs
BenchmarkDistance           15.3 ns/op    0 allocs
BenchmarkPackDistance       32.0 ns/op    0 allocs
BenchmarkHammingDistance     1.9 ns/op    0 allocs
BenchmarkDualEncode        237  ns/op    0 allocs
BenchmarkDualDistance       30   ns/op    0 allocs
BenchmarkFullEncode        889  ns/op    0 allocs
BenchmarkFullDistance     1470  ns/op    0 allocs
BenchmarkIndexSearch1K    5100  ns/op    0 allocs
BenchmarkBatchDistance1K 34800  ns/op    0 allocs
```

Run benchmarks:
```bash
go test -bench=. -benchmem -count=3
```

## Architecture

```
Code [8]byte          — fast fixed-size code (up to 7 digits)
FullCode [64]byte     — full variable-length code (up to 63 digits)
PackedCode uint32     — bit-packed code for batch operations
DualCode              — head + tail pair for long words
FullCodePair          — strict + relaxed pair for accurate matching

Encode()              — fast track dispatcher
DualEncode()          — dual track
FullEncode()          — full track dispatcher
FullFinnish()         — full Finnish encode with normalization

Distance()            — Levenshtein on Code
FullDistance()         — Levenshtein on FullCode
PackDistance()         — bit-parallel Levenshtein on PackedCode
HammingDistance()      — byte-level Hamming on PackedCode

NewIndex()            — precomputed sorted index
BatchPack/Distance()  — bulk operations with caller-provided buffers
```

## License

MIT
