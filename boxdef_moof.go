package mp4

import (
	"io"
)

// 8.8.4 Movie Fragment Box

// Box Type: ‘moof’
// Container: File
// Mandatory: No
// Quantity: Zero or more

// The movie fragments extend the presentation in time. They provide the
// information that would previously have been in the Movie Box. The actual
// samples are in Media Data Boxes, as usual, if they are in the same file. The
// data reference index is in the sample description, so it is possible to build
// incremental presentations where the media data is in files other than the
// file containing the Movie Box.
//
// The Movie Fragment Box is a top‐level box, (i.e. a peer to the Movie Box and
// Media Data boxes). It contains a Movie Fragment Header Box, and then one or
// more Track Fragment Boxes.
//
// > NOTE There is no requirement that any particular movie fragment extend all
// tracks present in the movie header, and there is no restriction on the
// location of the media data referred to by the movie fragments. However,
// derived specifications may make such restrictions.
type MovieFragmentBox struct {
	Header
	Container
}

var _ Box = (*MovieFragmentBox)(nil)

func init() {
	BoxRegistry[MoofBoxType] = func() Box { return &MovieFragmentBox{} }
}

func (b MovieFragmentBox) Mp4BoxType() BoxType {
	return MoofBoxType
}

func (b *MovieFragmentBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += b.Mp4BoxUpdateChildren()
	return b.Size
}

func (b *MovieFragmentBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.Mp4BoxReadChildren(r, b.Size-b.HeaderSize()); err != nil {
		return
	}
	return
}

func (b *MovieFragmentBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.Mp4BoxWriteChildren(w); err != nil {
		return
	}
	return
}
