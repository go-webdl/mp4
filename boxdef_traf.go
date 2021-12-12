package mp4

import (
	"io"
)

// 8.8.6 Track Fragment Box

// Box Type: ‘traf’
// Container: Movie Fragment Box ('moof')
// Mandatory: No
// Quantity: Zero or more

// Within the movie fragment there is a set of track fragments, zero or more per
// track. The track fragments in turn contain zero or more track runs, each of
// which document a contiguous run of samples for that track. Within these
// structures, many fields are optional and can be defaulted.
//
// It is possible to add 'empty time' to a track using these structures, as well
// as adding samples. Empty inserts can be used in audio tracks doing silence
// suppression, for example.
type TrackFragmentBox struct {
	Header
	Container
}

var _ Box = (*TrackFragmentBox)(nil)

func init() {
	BoxRegistry[TrafBoxType] = func() Box { return &TrackFragmentBox{} }
}

func (b TrackFragmentBox) Mp4BoxType() BoxType {
	return TrafBoxType
}

func (b *TrackFragmentBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += b.Mp4BoxUpdateChildren()
	return b.Size
}

func (b *TrackFragmentBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.Mp4BoxReadChildren(r, b.Size-b.HeaderSize()); err != nil {
		return
	}
	return
}

func (b *TrackFragmentBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.Mp4BoxWriteChildren(w); err != nil {
		return
	}
	return
}
