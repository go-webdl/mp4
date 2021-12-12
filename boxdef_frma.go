package mp4

import (
	"encoding/binary"
	"io"
)

// 8.12.2 Original Format Box

// Box Types: ‘frma’
// Container: Protection Scheme Information Box (‘sinf’), Restricted Scheme
//            Information Box (‘rinf’), or Complete Track Information Box
//            (‘cinf’)
// Mandatory: Yes when used in a protected sample entry, in a restricted sample
//            entry, or in a sample entry for an incomplete track.
// Quantity:  Exactly one.

// The Original Format Box ‘frma’ contains the four‐character‐code of the
// original un‐transformed sample description:
type OriginalFormatBox struct {
	Header
	NullContainer

	// format of decrypted, encoded data (in case of protection) or
	// un-transformed sample entry (in case of restriction and complete track
	// information)
	DataFormat FourCC
}

var _ Box = (*OriginalFormatBox)(nil)

func init() {
	BoxRegistry[FrmaBoxType] = func() Box { return &OriginalFormatBox{} }
}

func (b OriginalFormatBox) Mp4BoxType() BoxType {
	return FrmaBoxType
}

func (b *OriginalFormatBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += 4 // unsigned int(32) data_format = codingname;
	return b.Size
}

func (b *OriginalFormatBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.DataFormat); err != nil {
		return
	}
	return
}

func (b *OriginalFormatBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.DataFormat); err != nil {
		return
	}
	return
}
