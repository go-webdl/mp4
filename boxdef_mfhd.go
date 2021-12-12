package mp4

import (
	"encoding/binary"
	"io"
)

// 8.8.5 Movie Fragment Header Box

// Box Type: ‘mfhd’
// Container: Movie Fragment Box ('moof')
// Mandatory: Yes
// Quantity: Exactly one

// The movie fragment header contains a sequence number, as a safety check. The
// sequence number usually starts at 1 and increases for each movie fragment in
// the file, in the order in which they occur. This allows readers to verify
// integrity of the sequence in environments where undesired re‐ordering might
// occur.
type MovieFragmentHeaderBox struct {
	FullHeader
	NullContainer

	// a number associated with this fragment
	SequenceNumber uint32
}

var _ Box = (*MovieFragmentHeaderBox)(nil)

func init() {
	BoxRegistry[MfhdBoxType] = func() Box { return &MovieFragmentHeaderBox{} }
}

func (b MovieFragmentHeaderBox) Mp4BoxType() BoxType {
	return MfhdBoxType
}

func (b *MovieFragmentHeaderBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 4 // unsigned int(32) sequence_number;
	return b.Size
}

func (b *MovieFragmentHeaderBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.SequenceNumber); err != nil {
		return
	}
	return
}

func (b *MovieFragmentHeaderBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.SequenceNumber); err != nil {
		return
	}
	return
}
