package mp4

import (
	"encoding/binary"
	"io"
)

// 8.4.3 Handler Reference Box

// Box Type: ‘hdlr’
// Container: Media Box (‘mdia’) or Meta Box (‘meta’)
// Mandatory: Yes
// Quantity: Exactly one

// This box within a Media Box declares media type of the track, and thus the
// process by which the mediadata in the track is presented. For example, a
// format for which the decoder delivers video would be stored in a video track,
// identified by being handled by a video handler. The documentation of the
// storage of a media format identifies the media type which that format uses.
//
// This box when present within a Meta Box, declares the structure or format of
// the 'meta' box contents.
//
// There is a general handler for metadata streams of any type; the specific
// format is identified by the sample entry, as for video or audio, for example.
type HandlerBox struct {
	FullHeader
	NullContainer

	// when present in a media box, contains a value as defined in clause 12, or
	// a value from a derived specification, or registration.
	//
	// when present in a meta box, contains an appropriate value to indicate the
	// format of the meta box contents. The value ‘null’ can be used in the
	// primary meta box to indicate that it is merely being used to hold
	// resources.
	HandlerType FourCC

	// is a null‐terminated string in UTF‐8 characters which gives a
	// human‐readable name for the track type (for debugging and inspection
	// purposes).
	Name NullTerminatedString
}

var _ Box = (*HandlerBox)(nil)

func init() {
	BoxRegistry[HdlrBoxType] = func() Box { return &HandlerBox{} }
}

func (b HandlerBox) Mp4BoxType() BoxType {
	return HdlrBoxType
}

func (b *HandlerBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 4             // unsigned int(32) pre_defined = 0;
	b.Size += 4             // unsigned int(32) handler_type;
	b.Size += 4 * 3         // const unsigned int(32)[3] reserved = 0;
	b.Size += b.Name.Size() // string name;
	return b.Size
}

func (b *HandlerBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	var tmp [20]byte
	if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
		return
	}
	copy(b.HandlerType[:], tmp[4:8])
	if err = b.Name.ReadOfSize(r, b.Size-b.headerSize()-20); err != nil {
		return
	}
	return
}

func (b *HandlerBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	var tmp [20]byte
	copy(tmp[4:8], b.HandlerType[:])
	if err = binary.Write(w, binary.BigEndian, tmp); err != nil {
		return
	}
	if err = b.Name.Write(w); err != nil {
		return
	}
	return
}
