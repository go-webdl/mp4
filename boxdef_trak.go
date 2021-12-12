package mp4

import (
	"io"
)

// 8.3.1 Track Box

// Box Type: ‘trak’
// Container: Movie Box (‘moov’)
// Mandatory: Yes
// Quantity: One or more

// This is a container box for a single track of a presentation. A presentation
// consists of one or more tracks. Each track is independent of the other tracks
// in the presentation and carries its own temporal and spatial information.
// Each track will contain its associated Media Box.
//
// Tracks are used for two purposes: (a) to contain media data (media tracks)
// and (b) to contain packetization information for streaming protocols (hint
// tracks).
//
// There shall be at least one media track within an ISO file, and all the media
// tracks that contributed to the hint tracks shall remain in the file, even if
// the media data within them is not referenced by the hint tracks; after
// deleting all hint tracks, the entire un‐hinted presentation shall remain.
type TrackBox struct {
	Header
	Container
}

var _ Box = (*TrackBox)(nil)

func init() {
	BoxRegistry[TrakBoxType] = func() Box { return &TrackBox{} }
}

func (b TrackBox) Mp4BoxType() BoxType {
	return TrakBoxType
}

func (b *TrackBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += b.Mp4BoxUpdateChildren()
	return b.Size
}

func (b *TrackBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.Mp4BoxReadChildren(r, b.Size-b.HeaderSize()); err != nil {
		return
	}
	return
}

func (b *TrackBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.Mp4BoxWriteChildren(w); err != nil {
		return
	}
	return
}
