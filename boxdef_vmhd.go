package mp4

import (
	"encoding/binary"
	"io"
)

// 12.1.2 Video media header

// Box Types: ‘vmhd’
// Container: Media Information Box (‘minf’)
// Mandatory: Yes
// Quantity: Exactly one

// Video tracks use the VideoMediaHeaderbox in the media information box as
// defined in 8.4.5. The video media header contains general presentation
// information, independent of the coding, for video media. Note that the flags
// field has the value 1.
type VideoMediaHeaderBox struct {
	FullHeader
	NullContainer

	// specifies a composition mode for this video track, from the following
	// enumerated set, which may be extended by derived specifications:
	//
	// * copy = 0 copy over the existing image
	GraphicsMode uint16

	// is a set of 3 colour values (red, green, blue) available for use by
	// graphics modes
	OpColor [3]uint16
}

var _ Box = (*VideoMediaHeaderBox)(nil)

func (b VideoMediaHeaderBox) Mp4BoxType() BoxType {
	return VmhdBoxType
}

func (b *VideoMediaHeaderBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 2         // template unsigned int(16) graphicsmode = 0;
	b.Size += 2 * 3     // template unsigned int(16)[3] opcolor = {0, 0, 0};
	b.Mp4BoxSetFlags(1) // Note that the flags field has the value 1.
	return b.Size
}

func (b *VideoMediaHeaderBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.GraphicsMode); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.OpColor); err != nil {
		return
	}
	return
}

func (b *VideoMediaHeaderBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.GraphicsMode); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.OpColor); err != nil {
		return
	}
	return
}
