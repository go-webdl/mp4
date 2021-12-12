package mp4

import (
	"io"
)

// 8.4.6 Extended language tag

// Box Type: ‘elng’
// Container: Media Box (‘mdia’)
// Mandatory: No
// Quantity: Zero or one

// The extended language tag box represents media language information, based on
// RFC 4646 (Best Common Practices – BCP – 47) industry standard. It is an
// optional peer of the media header box, and must occur after the media header
// box.
//
// The extended language tag can provide better language information than the
// language field in the Media Header, including information such as region,
// script, variation, and so on, as parts (or subtags). The extended language
// tag box is optional, and if it is absent the media language should be used.
// The extended language tag overrides the media language if they are not
// consistent.
//
// For best compatibility with earlier players, if an extended language tag is
// specified, the most compatible language code should be specified in the
// language field of the Media Header box (for example, "eng" if the extended
// language tag is "en‐UK"). If there is no reasonably compatible tag, the
// packed form of 'und' can be used.
type ExtendedLanguageBox struct {
	FullHeader
	NullContainer

	// is a NULL‐terminated C string containing an RFC 4646 (BCP 47) compliant
	// language tag string, such as "en‐US", "fr‐FR", or "zh‐CN".
	ExtendedLanguage NullTerminatedString
}

var _ Box = (*ExtendedLanguageBox)(nil)

func init() {
	BoxRegistry[ElngBoxType] = func() Box { return &ExtendedLanguageBox{} }
}

func (b ExtendedLanguageBox) Mp4BoxType() BoxType {
	return ElngBoxType
}

func (b *ExtendedLanguageBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += b.ExtendedLanguage.Size() // string name;
	return b.Size
}

func (b *ExtendedLanguageBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.ExtendedLanguage.ReadOfSize(r, b.Size-b.headerSize()); err != nil {
		return
	}
	return
}

func (b *ExtendedLanguageBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.ExtendedLanguage.Write(w); err != nil {
		return
	}
	return
}
