package mp4

import (
	"encoding/binary"
	"io"
)

// 8.5.3 Degradation Priority Box

// Box Type: ‘stdp’
// Container: Sample Table Box (‘stbl’).
// Mandatory: No.
// Quantity: Zero or one.

// This box contains the degradation priority of each sample. The values are
// stored in the table, one for each sample. The size of the table, sample_count
// is taken from the sample_count in the Sample Size Box ('stsz').
// Specifications derived from this define the exact meaning and acceptable
// range of the priority field.
type DegradationPriorityBox struct {
	FullHeader
	NullContainer

	SamplePriority []uint16
}

var _ Box = (*DegradationPriorityBox)(nil)

func init() {
	BoxRegistry[StdpBoxType] = func() Box { return &DegradationPriorityBox{} }
}

func (b DegradationPriorityBox) Mp4BoxType() BoxType {
	return StdpBoxType
}

func (b *DegradationPriorityBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	// int i;
	// for (i=0; i < sample_count; i++) {
	// 	unsigned int(16) priority;
	// 	}
	b.Size += 2 * uint32(len(b.SamplePriority))
	return b.Size
}

func (b *DegradationPriorityBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	b.SamplePriority = make([]uint16, (b.Size-b.headerSize())/2)
	if err = binary.Read(r, binary.BigEndian, b.SamplePriority); err != nil {
		return
	}
	return
}

func (b *DegradationPriorityBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.SamplePriority); err != nil {
		return
	}
	return
}
