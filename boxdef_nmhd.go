package mp4

import (
	"io"
)

// 8.4.5.2 Null Media Header Box

// Box Types: ‘nmhd’
// Container: Media Information Box (‘minf’)
// Mandatory: Yes
// Quantity: Exactly one specific media header shall be present

// Streams for which no specific media header is identified use a null Media
// Header Box, as defined here.
type NullMediaHeaderBox struct {
	FullHeader
	NullContainer
}

var _ Box = (*NullMediaHeaderBox)(nil)

func init() {
	BoxRegistry[NmhdBoxType] = func() Box { return &NullMediaHeaderBox{} }
}

func (b NullMediaHeaderBox) Mp4BoxType() BoxType {
	return NmhdBoxType
}

func (b *NullMediaHeaderBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	return b.Size
}

func (b *NullMediaHeaderBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	return
}

func (b *NullMediaHeaderBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	return
}
