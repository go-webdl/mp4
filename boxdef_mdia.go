package mp4

import (
	"io"
)

// 8.4.1 Media Box

// Box Type: ‘mdia’
// Container: Track Box (‘trak’)
// Mandatory: Yes
// Quantity: Exactly one

// The media declaration container contains all the objects that declare
// information about the media data within a track.
type MediaBox struct {
	Header
	Container
}

var _ Box = (*MediaBox)(nil)

func init() {
	BoxRegistry[MdiaBoxType] = func() Box { return &MediaBox{} }
}

func (b MediaBox) Mp4BoxType() BoxType {
	return MdiaBoxType
}

func (b *MediaBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += b.Mp4BoxUpdateChildren()
	return b.Size
}

func (b *MediaBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.Mp4BoxReadChildren(r, b.Size-b.HeaderSize()); err != nil {
		return
	}
	return
}

func (b *MediaBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.Mp4BoxWriteChildren(w); err != nil {
		return
	}
	return
}
