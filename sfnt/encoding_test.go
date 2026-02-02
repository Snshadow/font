package sfnt

import (
	"os"
	"strings"
	"testing"
)

// TestGetEncoding tests that GetEncoding returns the correct encoding
// for various platform/encoding/language combinations.
func TestGetEncoding(t *testing.T) {
	tests := []struct {
		name       string
		platformID PlatformID
		encodingID PlatformEncodingID
		langID     PlatformLanguageID
		wantNil    bool
	}{
		// Unicode platform
		{"Unicode/default", PlatformUnicode, 0, 0, false},
		{"Unicode/BMP", PlatformUnicode, 3, 0, false},

		// Mac platform - Roman with language variants
		{"Mac/Roman/English", PlatformMac, 0, 0, false},
		{"Mac/Roman/Icelandic", PlatformMac, 0, 15, false},
		{"Mac/Roman/Turkish", PlatformMac, 0, 17, false},
		{"Mac/Roman/Croatian", PlatformMac, 0, 18, false},
		{"Mac/Roman/Lithuanian", PlatformMac, 0, 24, false},
		{"Mac/Roman/Romanian", PlatformMac, 0, 37, false},
		{"Mac/Japanese", PlatformMac, 1, 11, false},
		{"Mac/TradChinese", PlatformMac, 2, 19, false},
		{"Mac/Korean", PlatformMac, 3, 23, false},
		{"Mac/Greek", PlatformMac, 6, 14, false},
		{"Mac/Cyrillic", PlatformMac, 7, 32, false},
		{"Mac/SimpChinese", PlatformMac, 25, 33, false},
		{"Mac/CentralEuropean", PlatformMac, 29, 0, false},
		{"Mac/Turkish", PlatformMac, 35, 0, false},
		{"Mac/Iceland", PlatformMac, 37, 0, false},

		// ISO platform
		{"ISO/ASCII", PlatformISO, 0, 0, true}, // ASCII is valid UTF-8, nil signals to use string() directly
		{"ISO/Unicode", PlatformISO, 1, 0, false},
		{"ISO/Latin1", PlatformISO, 2, 0, false},

		// Microsoft platform
		{"MS/Symbol", PlatformMicrosoft, 0, 0x0409, false},
		{"MS/Unicode", PlatformMicrosoft, 1, 0x0409, false},
		{"MS/ShiftJIS", PlatformMicrosoft, 2, 0x0411, false},
		{"MS/GBK", PlatformMicrosoft, 3, 0x0804, false},
		{"MS/Big5", PlatformMicrosoft, 4, 0x0404, false},
		{"MS/Wansung", PlatformMicrosoft, 5, 0x0412, false},
		{"MS/Johab", PlatformMicrosoft, 6, 0x0412, false},
		{"MS/UCS4", PlatformMicrosoft, 10, 0x0409, false},

		// Unknown/unsupported
		{"Unknown/Platform", PlatformID(99), 0, 0, true},
		{"Mac/Unknown", PlatformMac, 99, 0, true},
		{"MS/Unknown", PlatformMicrosoft, 99, 0, true},
		{"ISO/Unknown", PlatformISO, 99, 0, true},
		{"Custom/Platform", PlatformCustom, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := GetEncoding(tt.platformID, tt.encodingID, tt.langID)
			if tt.wantNil && enc != nil {
				t.Errorf("GetEncoding(%d, %d, %d) = %v, want nil",
					tt.platformID, tt.encodingID, tt.langID, enc)
			}
			if !tt.wantNil && enc == nil {
				t.Errorf("GetEncoding(%d, %d, %d) = nil, want non-nil",
					tt.platformID, tt.encodingID, tt.langID)
			}
		})
	}
}

// TestNameEntryDecode tests decoding name entries from the comprehensive test font.
func TestNameEntryDecode(t *testing.T) {
	file, err := os.Open("testdata/EncodingTest.ttf")
	if err != nil {
		t.Fatalf("Failed to open test font: %v", err)
	}
	defer file.Close()

	font, err := Parse(file)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	nameTable, err := font.NameTable()
	if err != nil {
		t.Fatalf("Failed to get name table: %v", err)
	}

	tests := []struct {
		name       string
		platformID PlatformID
		encodingID PlatformEncodingID
		langID     PlatformLanguageID
		nameID     NameID
		contains   string
	}{
		// Unicode entries
		{"Unicode/FontFamily", PlatformUnicode, 3, 0, NameFontFamily, "EncodingTest"},
		{"Unicode/Copyright", PlatformUnicode, 3, 0, NameCopyrightNotice, "Copyright"},

		// Mac Roman (English)
		{"Mac/Roman/FontFamily", PlatformMac, 0, 0, NameFontFamily, "EncodingTest"},

		// Mac Japanese
		{"Mac/Japanese/FontFamily", PlatformMac, 1, 11, NameFontFamily, "エンコーディングテスト"},

		// Mac Traditional Chinese
		{"Mac/TradChinese/FontFamily", PlatformMac, 2, 19, NameFontFamily, "編碼測試"},

		// Mac Korean
		{"Mac/Korean/FontFamily", PlatformMac, 3, 23, NameFontFamily, "인코딩테스트"},

		// Mac Greek
		{"Mac/Greek/FontFamily", PlatformMac, 6, 14, NameFontFamily, "ΤεστΚωδικοποίησης"},

		// Mac Cyrillic (Russian)
		{"Mac/Cyrillic/FontFamily", PlatformMac, 7, 32, NameFontFamily, "ТестКодировки"},

		// Microsoft Unicode (English)
		{"MS/Unicode/FontFamily", PlatformMicrosoft, 1, 0x0409, NameFontFamily, "EncodingTest"},

		// Microsoft ShiftJIS (Japanese)
		{"MS/ShiftJIS/FontFamily", PlatformMicrosoft, 2, 0x0411, NameFontFamily, "エンコーディングテスト"},
	}

	entries := nameTable.List()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var found *NameEntry
			for _, entry := range entries {
				if entry.PlatformID == tt.platformID &&
					entry.EncodingID == tt.encodingID &&
					entry.LanguageID == tt.langID &&
					entry.NameID == tt.nameID {
					found = entry
					break
				}
			}

			if found == nil {
				t.Skipf("Entry not found in test font: platform=%d, encoding=%d, lang=%d, name=%d",
					tt.platformID, tt.encodingID, tt.langID, tt.nameID)
				return
			}

			decoded := found.String()
			if !strings.Contains(decoded, tt.contains) {
				t.Errorf("Decoded string %q does not contain expected %q", decoded, tt.contains)
			}
		})
	}
}

// TestSpecialCharacters tests that special characters from extended encodings decode correctly.
func TestSpecialCharacters(t *testing.T) {
	file, err := os.Open("testdata/EncodingTest.ttf")
	if err != nil {
		t.Fatalf("Failed to open test font: %v", err)
	}
	defer file.Close()

	font, err := Parse(file)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	nameTable, err := font.NameTable()
	if err != nil {
		t.Fatalf("Failed to get name table: %v", err)
	}

	// Look for entries that may contain special characters like:
	// Copyright (©), Trademark (™), Ellipsis (…), Em dash (—)
	entries := nameTable.List()

	specialChars := []struct {
		char        string
		description string
	}{
		{"©", "Copyright symbol"},
		{"™", "Trademark symbol"},
	}

	// Find copyright notice entries and check for special characters
	for _, entry := range entries {
		if entry.NameID == NameCopyrightNotice {
			decoded := entry.String()
			for _, sc := range specialChars {
				if strings.Contains(decoded, sc.char) {
					t.Logf("Found %s (%s) in platform=%s, encoding=%d: %s",
						sc.description, sc.char, entry.PlatformID, entry.EncodingID, decoded)
				}
			}
		}
	}

	// Verify at least some entries decode without errors (don't contain replacement character)
	for _, entry := range entries {
		decoded := entry.String()
		// Check for Unicode replacement character which indicates decoding failure
		if strings.Contains(decoded, "\ufffd") {
			t.Logf("Warning: Entry contains replacement character: platform=%s, encoding=%d, name=%s",
				entry.PlatformID, entry.EncodingID, entry.NameID)
		}
	}
}

// TestNilEncodingFallback tests the behavior when GetEncoding returns nil.
func TestNilEncodingFallback(t *testing.T) {
	// Create a NameEntry with unsupported platform
	entry := &NameEntry{
		PlatformID: PlatformID(99), // Unknown platform
		EncodingID: 0,
		LanguageID: 0,
		NameID:     NameFontFamily,
		Value:      []byte("Test Font"),
	}

	// Verify that String() returns raw bytes as string (fallback behavior)
	result := entry.String()
	if result != "Test Font" {
		t.Errorf("Fallback string = %q, want %q", result, "Test Font")
	}

	// Test with Custom platform (also unsupported)
	entry2 := &NameEntry{
		PlatformID: PlatformCustom,
		EncodingID: 0,
		LanguageID: 0,
		NameID:     NameFontFamily,
		Value:      []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}, // "Hello" in ASCII
	}

	result2 := entry2.String()
	if result2 != "Hello" {
		t.Errorf("Fallback string = %q, want %q", result2, "Hello")
	}
}

// TestMacRomanLanguageVariants tests that Mac Roman encoding handles language variants correctly.
func TestMacRomanLanguageVariants(t *testing.T) {
	tests := []struct {
		langID      PlatformLanguageID
		description string
	}{
		{0, "English"},
		{15, "Icelandic"},
		{17, "Turkish"},
		{18, "Croatian"},
		{24, "Lithuanian"},
		{25, "Polish"},
		{26, "Hungarian"},
		{27, "Estonian"},
		{28, "Latvian"},
		{36, "Slovenian"},
		{37, "Romanian"},
		{38, "Czech"},
		{39, "Slovak"},
		{40, "Slovenian (alt)"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			enc := GetEncoding(PlatformMac, 0, tt.langID)
			if enc == nil {
				t.Errorf("GetEncoding(Mac, 0, %d) returned nil for %s", tt.langID, tt.description)
			}
		})
	}
}

// TestAllNameEntriesDecode ensures all entries in the test font can be decoded without panic.
func TestAllNameEntriesDecode(t *testing.T) {
	file, err := os.Open("testdata/EncodingTest.ttf")
	if err != nil {
		t.Fatalf("Failed to open test font: %v", err)
	}
	defer file.Close()

	font, err := Parse(file)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	nameTable, err := font.NameTable()
	if err != nil {
		t.Fatalf("Failed to get name table: %v", err)
	}

	entries := nameTable.List()
	t.Logf("Found %d name entries in test font", len(entries))

	for i, entry := range entries {
		// This should not panic
		decoded := entry.String()
		t.Logf("[%d] Platform=%s Encoding=%d Lang=%d Name=%s: %q",
			i, entry.PlatformID, entry.EncodingID, entry.LanguageID, entry.NameID, decoded)
	}
}

// BenchmarkGetEncoding benchmarks the GetEncoding function for different platforms.
func BenchmarkGetEncoding(b *testing.B) {
	benchmarks := []struct {
		name       string
		platformID PlatformID
		encodingID PlatformEncodingID
		langID     PlatformLanguageID
	}{
		{"Unicode", PlatformUnicode, 3, 0},
		{"Mac/Roman", PlatformMac, 0, 0},
		{"Mac/Japanese", PlatformMac, 1, 11},
		{"Mac/Cyrillic", PlatformMac, 7, 32},
		{"ISO/Latin1", PlatformISO, 2, 0},
		{"MS/Unicode", PlatformMicrosoft, 1, 0x0409},
		{"MS/ShiftJIS", PlatformMicrosoft, 2, 0x0411},
		{"Unknown", PlatformID(99), 0, 0},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				GetEncoding(bm.platformID, bm.encodingID, bm.langID)
			}
		})
	}
}

// BenchmarkNameEntryString benchmarks decoding name entries to strings.
func BenchmarkNameEntryString(b *testing.B) {
	file, err := os.Open("testdata/EncodingTest.ttf")
	if err != nil {
		b.Fatalf("Failed to open test font: %v", err)
	}
	defer file.Close()

	font, err := Parse(file)
	if err != nil {
		b.Fatalf("Failed to parse font: %v", err)
	}

	nameTable, err := font.NameTable()
	if err != nil {
		b.Fatalf("Failed to get name table: %v", err)
	}

	entries := nameTable.List()
	if len(entries) == 0 {
		b.Fatal("No name entries found")
	}

	// Group entries by platform for targeted benchmarks
	var unicodeEntry, macRomanEntry, macJapaneseEntry, msUnicodeEntry, msShiftJISEntry *NameEntry
	for _, entry := range entries {
		if entry.NameID == NameFontFamily {
			switch {
			case entry.PlatformID == PlatformUnicode && unicodeEntry == nil:
				unicodeEntry = entry
			case entry.PlatformID == PlatformMac && entry.EncodingID == 0 && entry.LanguageID == 0 && macRomanEntry == nil:
				macRomanEntry = entry
			case entry.PlatformID == PlatformMac && entry.EncodingID == 1 && macJapaneseEntry == nil:
				macJapaneseEntry = entry
			case entry.PlatformID == PlatformMicrosoft && entry.EncodingID == 1 && msUnicodeEntry == nil:
				msUnicodeEntry = entry
			case entry.PlatformID == PlatformMicrosoft && entry.EncodingID == 2 && msShiftJISEntry == nil:
				msShiftJISEntry = entry
			}
		}
	}

	benchmarks := []struct {
		name  string
		entry *NameEntry
	}{
		{"Unicode", unicodeEntry},
		{"Mac/Roman", macRomanEntry},
		{"Mac/Japanese", macJapaneseEntry},
		{"MS/Unicode", msUnicodeEntry},
		{"MS/ShiftJIS", msShiftJISEntry},
	}

	for _, bm := range benchmarks {
		if bm.entry == nil {
			continue
		}
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = bm.entry.String()
			}
		})
	}
}

// BenchmarkAllNameEntries benchmarks decoding all name entries in a font.
func BenchmarkAllNameEntries(b *testing.B) {
	file, err := os.Open("testdata/EncodingTest.ttf")
	if err != nil {
		b.Fatalf("Failed to open test font: %v", err)
	}
	defer file.Close()

	font, err := Parse(file)
	if err != nil {
		b.Fatalf("Failed to parse font: %v", err)
	}

	nameTable, err := font.NameTable()
	if err != nil {
		b.Fatalf("Failed to get name table: %v", err)
	}

	entries := nameTable.List()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, entry := range entries {
			_ = entry.String()
		}
	}
}
