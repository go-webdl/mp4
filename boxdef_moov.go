package mp4

import (
	"io"
)

// 8.2.1 Movie Box

// Box Type: ‘moov’
// Container: File
// Mandatory: Yes
// Quantity: Exactly one

// The metadata for a presentation is stored in the single Movie Box which
// occurs at the top‐level of a file. Normally this box is close to the
// beginning or end of the file, though this is not required.
type MovieBox struct {
	Header
	Container
}

var _ Box = (*MovieBox)(nil)

func init() {
	BoxRegistry[MoovBoxType] = func() Box { return &MovieBox{} }
}

func (b MovieBox) Mp4BoxType() BoxType {
	return MoovBoxType
}

func (b *MovieBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += b.Mp4BoxUpdateChildren()
	return b.Size
}

func (b *MovieBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.Mp4BoxReadChildren(r, b.Size-b.HeaderSize()); err != nil {
		return
	}
	return
}

func (b *MovieBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.Mp4BoxWriteChildren(w); err != nil {
		return
	}
	return
}
