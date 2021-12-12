package mp4

import (
	"encoding/binary"
	"io"
)

// 12.2.2 Sound media header

// Box Types: ‘smhd’
// Container: Media Information Box (‘minf’)
// Mandatory: Yes
// Quantity: Exactly one specific media header shall be present

// Audio tracks use the SoundMediaHeaderbox in the media information box as
// defined in 8.4.5. The sound media header contains general presentation
// information, independent of the coding, for audio media. This header is used
// for all tracks containing audio.
type SoundMediaHeaderBox struct {
	FullHeader
	NullContainer

	// is a fixed‐point 8.8 number that places mono audio tracks in a stereo
	// space; 0 is centre (the normal value); full left is ‐1.0 and full right
	// is 1.0.
	Balance int16
}

var _ Box = (*SoundMediaHeaderBox)(nil)

func init() {
	BoxRegistry[SmhdBoxType] = func() Box { return &SoundMediaHeaderBox{} }
}

func (b SoundMediaHeaderBox) Mp4BoxType() BoxType {
	return SmhdBoxType
}

func (b *SoundMediaHeaderBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 2 // template int(16) balance = 0;
	b.Size += 2 // const unsigned int(16) reserved = 0;
	return b.Size
}

func (b *SoundMediaHeaderBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.Balance); err != nil {
		return
	}
	var reserved uint16
	if err = binary.Read(r, binary.BigEndian, &reserved); err != nil {
		return
	}
	return
}

func (b *SoundMediaHeaderBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.Balance); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint16(0)); err != nil {
		return
	}
	return
}
