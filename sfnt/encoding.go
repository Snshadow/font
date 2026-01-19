package sfnt

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
)

// GetEncoding is a best-effort attempt to return the text encoding for a given
// platformID/encodingID/langID, which might result in broken text.
// Returns nil if the encoding is unknown or unsupported.
func GetEncoding(platformID PlatformID, encodingID PlatformEncodingID, langID PlatformLanguageID) encoding.Encoding {
	switch platformID {
	case PlatformUnicode:
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	case PlatformMac:
		return getMacEncoding(encodingID, langID)
	case PlatformISO:
		return getISOEncoding(encodingID)
	case PlatformMicrosoft:
		return getMicrosoftEncoding(encodingID)
	}

	return nil
}

// getMacEncoding returns the encoding for Mac platform entries.
func getMacEncoding(encodingID PlatformEncodingID, langID PlatformLanguageID) encoding.Encoding {
	// TODO: implement custom charmap for mac encodings
	switch encodingID {
	case 0: // Mac Roman
		switch langID {
		case 15: // mac_iceland
			return charmap.ISO8859_1
		case 17: // mac_turkish
			return charmap.ISO8859_9
		case 18: // mac_croatian
			return charmap.ISO8859_2
		case 24, 25, 26, 27, 28, 36, 38, 39, 40: // mac_latin2
			return charmap.ISO8859_2
		case 37: // mac_romanian
			return charmap.ISO8859_16
		default:
			return charmap.Macintosh
		}

	case 1: // x_mac_japanese_ttx
		return japanese.ShiftJIS
	case 2: // x_mac_trad_chinese_ttx
		return traditionalchinese.Big5
	case 3: // x_mac_korean_ttx
		return korean.EUCKR
	case 6: // mac_greek
		return charmap.ISO8859_7
	case 7: // mac_cyrillic
		return charmap.MacintoshCyrillic
	case 25: // x_mac_simp_chinese_ttx
		return simplifiedchinese.HZGB2312
	case 29: // mac_latin2
		return charmap.ISO8859_2
	case 35: // mac_turkish
		return charmap.ISO8859_9
	case 37: // mac_iceland
		return charmap.ISO8859_1
	}

	return nil
}

// getISOEncoding returns the encoding for ISO platform entries.
func getISOEncoding(encodingID PlatformEncodingID) encoding.Encoding {
	switch encodingID {
	case 0: // 7-bit ASCII
		// ASCII is a subset of ISO-8859-1
		return charmap.ISO8859_1
	case 1: // ISO 10646 (Unicode)
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	case 2: // ISO 8859-1 (Latin 1)
		return charmap.ISO8859_1
	}

	return nil
}

// getMicrosoftEncoding returns the encoding for Microsoft platform entries.
func getMicrosoftEncoding(encodingID PlatformEncodingID) encoding.Encoding {
	switch encodingID {
	case 0: // Symbol
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	case 1: // Unicode BMP (UCS-2)
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	case 2:
		return japanese.ShiftJIS
	case 3:
		return simplifiedchinese.GBK
	case 4:
		return traditionalchinese.Big5
	case 5:
		return korean.EUCKR
	case 6: // Johab
		// TODO: implement custom Johab charmap
		return nil
	case 10: // Unicode full repertoire (UCS-4)
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	}

	return nil
}
