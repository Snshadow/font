package sfnt

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"

	"github.com/go-sw/text-codec/apple"
	johab "github.com/go-sw/text-codec/korean"
)

// GetEncoding is a best-effort attempt to return the text encoding for a given
// platformID/encodingID/langID, which might result in broken text.
// Returns nil if the encoding is already UTF-8 compatible (e.g. ASCII) or unsupported.
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
	switch encodingID {
	case 0: // Mac Roman
		switch langID {
		case 15:
			return apple.Iceland
		case 17:
			return apple.Turkish
		case 18:
			return apple.Croatian
		case 24, 25, 26, 27, 28, 36, 38, 39, 40: // mac_latin2
			return apple.CentralEuropean
		case 37:
			return apple.Romanian
		default:
			return charmap.Macintosh
		}

	case 1:
		return apple.Japanese
	case 2:
		return apple.ChineseTraditional
	case 3:
		return apple.Korean
	case 6:
		return apple.Greek
	case 7:
		return charmap.MacintoshCyrillic
	case 25:
		return apple.ChineseSimplified
	case 29: // mac_latin2
		return apple.CentralEuropean
	case 35:
		return apple.Turkish
	case 37:
		return apple.Iceland
	}

	return nil
}

// getISOEncoding returns the encoding for ISO platform entries.
func getISOEncoding(encodingID PlatformEncodingID) encoding.Encoding {
	switch encodingID {
	case 0: // 7-bit ASCII
		return nil // ASCII is valid UTF-8
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
	case 6:
		return johab.Johab
	case 10: // Unicode full repertoire (UCS-4)
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	}

	return nil
}
