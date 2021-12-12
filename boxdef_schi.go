package mp4

import (
	"io"
)

// 8.12.6 Scheme Information Box

// Box Types: ‘schi’
// Container: Protection Scheme Information Box (‘sinf’), Restricted Scheme
//            Information Box (‘rinf’), or SRTP Process box (‘srpp‘)
// Mandatory: No
// Quantity : Zero or one

// The Scheme Information Box is a container Box that is only interpreted by the
// scheme being used. Any information the encryption or restriction system needs
// is stored here. The content of this box is a series of boxes whose type and
// format are defined by the scheme declared in the Scheme Type Box.
type SchemeInformationBox struct {
	Header
	Container
}

var _ Box = (*SchemeInformationBox)(nil)

func init() {
	BoxRegistry[SchiBoxType] = func() Box { return &SchemeInformationBox{} }
}

func (b SchemeInformationBox) Mp4BoxType() BoxType {
	return SchiBoxType
}

func (b *SchemeInformationBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += b.Mp4BoxUpdateChildren()
	return b.Size
}

func (b *SchemeInformationBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.Mp4BoxReadChildren(r, b.Size-b.HeaderSize()); err != nil {
		return
	}
	return
}

func (b *SchemeInformationBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.Mp4BoxWriteChildren(w); err != nil {
		return
	}
	return
}
