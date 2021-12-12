package mp4

import (
	"io"
)

// 8.8.1 Movie Extends Box

// Box Type: ‘mvex’
// Container: Movie Box (‘moov’)
// Mandatory: No
// Quantity: Zero or one

// This box warns readers that there might be Movie Fragment Boxes in this file.
// To know of all samples in the tracks, these Movie Fragment Boxes must be
// found and scanned in order, and their information logically added to that
// found in the Movie Box.
//
// There is a narrative introduction to Movie Fragments in Annex A.
type MovieExtendsBox struct {
	Header
	Container
}

var _ Box = (*MovieExtendsBox)(nil)

func init() {
	BoxRegistry[MvexBoxType] = func() Box { return &MovieExtendsBox{} }
}

func (b MovieExtendsBox) Mp4BoxType() BoxType {
	return MvexBoxType
}

func (b *MovieExtendsBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += b.Mp4BoxUpdateChildren()
	return b.Size
}

func (b *MovieExtendsBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.Mp4BoxReadChildren(r, b.Size-b.HeaderSize()); err != nil {
		return
	}
	return
}

func (b *MovieExtendsBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.Mp4BoxWriteChildren(w); err != nil {
		return
	}
	return
}
