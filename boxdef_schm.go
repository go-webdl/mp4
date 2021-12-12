package mp4

import (
	"encoding/binary"
	"io"
)

// 8.12.5 Scheme Type Box

// Box Types: ‘schm’
// Container: Protection Scheme Information Box (‘sinf’), Restricted Scheme
//            Information Box (‘rinf’), or SRTP Process box (‘srpp‘)
// Mandatory: No
// Quantity : Zero or one in ‘sinf’, depending on the protection structure;
//            Exactly one in ‘rinf’ and ‘srpp’

// The Scheme Type Box (‘schm’) identifies the protection or restriction scheme.
type SchemeTypeBox struct {
	FullHeader
	NullContainer

	// is the code defining the protection or restriction scheme.
	SchemeType FourCC

	// is the version of the scheme (used to create the content)
	SchemeVersion uint32

	// allows for the option of directing the user to a web‐page if they do not
	// have the scheme installed on their system. It is an absolute URI formed
	// as a null‐terminated string in UTF‐8 characters.
	SchemeURI NullTerminatedString
}

const (
	FLAG_SCHM_SCHEME_URI uint32 = 0x000001
)

var _ Box = (*SchemeTypeBox)(nil)

func init() {
	BoxRegistry[SchmBoxType] = func() Box { return &SchemeTypeBox{} }
}

func (b SchemeTypeBox) Mp4BoxType() BoxType {
	return SchmBoxType
}

func (b *SchemeTypeBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 4 // unsigned int(32) scheme_type;
	b.Size += 4 // unsigned int(32) scheme_version;
	if b.Mp4BoxFlags()&FLAG_SCHM_SCHEME_URI > 0 {
		b.Size += b.SchemeURI.Size() // unsigned int(8) scheme_uri[];

	}
	return b.Size
}

func (b *SchemeTypeBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.SchemeType); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.SchemeVersion); err != nil {
		return
	}
	if b.Mp4BoxFlags()&FLAG_SCHM_SCHEME_URI > 0 {
		if err = b.SchemeURI.ReadOfSize(r, b.Size-b.headerSize()-8); err != nil {
			return
		}
	}
	return
}

func (b *SchemeTypeBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.SchemeType); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.SchemeVersion); err != nil {
		return
	}
	if b.Mp4BoxFlags()&FLAG_SCHM_SCHEME_URI > 0 {
		if err = b.SchemeURI.Write(w); err != nil {
			return
		}
	}
	return
}
