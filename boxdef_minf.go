package mp4

import (
	"io"
)

// 8.4.4 Media Information Box

// Box Type: ‘minf’
// Container: Media Box (‘mdia’)
// Mandatory: Yes
// Quantity: Exactly one

// This box contains all the objects that declare characteristic information of the media in the track.
type MediaInformationBox struct {
	Header
	Container
}

var _ Box = (*MediaInformationBox)(nil)

func init() {
	BoxRegistry[MinfBoxType] = func() Box { return &MediaInformationBox{} }
}

func (b MediaInformationBox) Mp4BoxType() BoxType {
	return MinfBoxType
}

func (b *MediaInformationBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += b.Mp4BoxUpdateChildren()
	return b.Size
}

func (b *MediaInformationBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.Mp4BoxReadChildren(r, b.Size-b.HeaderSize()); err != nil {
		return
	}
	return
}

func (b *MediaInformationBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.Mp4BoxWriteChildren(w); err != nil {
		return
	}
	return
}
