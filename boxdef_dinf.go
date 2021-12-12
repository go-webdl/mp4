package mp4

import (
	"io"
)

// 8.7.1 Data Information Box

// Box Type: ‘dinf’
// Container: Media Information Box (‘minf’) or Meta Box (‘meta’)
// Mandatory: Yes (required within ‘minf’ box) and No (optional within ‘meta’ box)
// Quantity: Exactly one

// The data information box contains objects that declare the location of the
// media information in a track.
type DataInformationBox struct {
	Header
	Container
}

var _ Box = (*DataInformationBox)(nil)

func init() {
	BoxRegistry[DinfBoxType] = func() Box { return &DataInformationBox{} }
}

func (b DataInformationBox) Mp4BoxType() BoxType {
	return DinfBoxType
}

func (b *DataInformationBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += b.Mp4BoxUpdateChildren()
	return b.Size
}

func (b *DataInformationBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.Mp4BoxReadChildren(r, b.Size-b.HeaderSize()); err != nil {
		return
	}
	return
}

func (b *DataInformationBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.Mp4BoxWriteChildren(w); err != nil {
		return
	}
	return
}
