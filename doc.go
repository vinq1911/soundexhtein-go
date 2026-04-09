// Package soundex provides multi-language phonetic encoding and
// phonetic-distance computation with zero-allocation hot paths.
//
// # Encoding Tracks
//
// Three encoding tracks with different speed/accuracy trade-offs:
//
// Fast (Code [8]byte, up to 7 digits):
//
//	~9-40ns encode, ~15ns distance, 0 allocs.
//	Best for pre-filtering millions of candidates.
//	Functions: American, Cologne, Metaphone, Finnish, Swedish, Norwegian,
//	Danish, Estonian, Latvian, Lithuanian, DaitchMokotoff.
//
// Dual (head + tail, 4+4 digits each):
//
//	~200-350ns encode, ~30ns distance, 0 allocs.
//	Captures both the start and end of a word in two PackedCode uint32s.
//	Best for batch operations on long/compound words.
//	Functions: DualEncode, DualDistance, PackDual, PackedDualDistance.
//
// Full (FullCode [64]byte, no truncation):
//
//	~140-900ns encode, ~200-1500ns distance, 0 allocs.
//	Encodes the entire word with normalize→encode pipeline.
//	Produces both strict (geminates preserved) and relaxed (geminates collapsed) codes.
//	Best for accurate matching of agglutinative languages (Finnish, Estonian).
//	Functions: FullEncode, FullFinnish, FullDistance, FullSoundexDistance.
//
// # Batch & Index APIs
//
// For high-throughput scenarios:
//
//	Pack, PackDistance, HammingDistance — bit-packed uint32 codes, ~2ns distance.
//	BatchPack, BatchDistance — process slices with caller-provided output buffers.
//	NewIndex, Search, SearchInto — precomputed sorted index with early termination.
//
// # Supported Languages
//
// Western European:
//
//	American Soundex — standard 4-char code (1 letter + 3 digits)
//	Cologne Phonetic — German-optimized encoding
//	Metaphone — improved English phonetic encoding
//	Daitch-Mokotoff — multi-code encoding for Slavic/Yiddish names
//
// Finnish:
//
//	Finnish — fast 7-digit code with vowel harmony preservation
//	FullFinnish — full-length normalize→encode pipeline with strict/relaxed modes
//
// Scandinavian:
//
//	Swedish — SJ/SKJ sje-sound, KJ/TJ tje-sound, K-fronting
//	Norwegian — KJ/SJ digraphs, RS retroflex
//	Danish — soft D (intervocalic), SJ digraph
//
// Baltic:
//
//	Estonian — š/ž, õ/ä/ö/ü, triple consonant collapse
//	Latvian — č/š/ž, ģ/ķ/ļ/ņ palatalized, DZ/DŽ affricates
//	Lithuanian — č/š/ž, DŽ affricate, nasal vowels
//
// # Zero Allocations
//
// All encoding and distance functions use fixed-size stack-allocated arrays.
// No heap allocations on any hot path. Verified with testing.AllocsPerRun.
package soundex
