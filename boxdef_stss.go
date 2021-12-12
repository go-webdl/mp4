package mp4

import (
	"encoding/binary"
	"io"
)

// 8.6.2 Sync Sample Box

// Box Type: ‘stss’
// Container: Sample Table Box (‘stbl’)
// Mandatory: No
// Quantity: Zero or one

// This box provides a compact marking of the sync samples within the stream.
// The table is arranged in strictly increasing order of sample number.
//
// If the sync sample box is not present, every sample is a sync sample.
type SyncSampleBox struct {
	FullHeader
	NullContainer
	SampleNumbers []uint32 // gives the numbers of the samples that are sync samples in the stream.
}

var _ Box = (*SyncSampleBox)(nil)

func init() {
	BoxRegistry[StssBoxType] = func() Box { return &SyncSampleBox{} }
}

func (b SyncSampleBox) Mp4BoxType() BoxType {
	return StssBoxType
}

func (b *SyncSampleBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 4 // unsigned int(32) entry_count;
	// for (i=0; i < entry_count; i++) {
	//     unsigned int(32) sample_number;
	// }
	b.Size += 4 * uint32(len(b.SampleNumbers))
	return b.Size
}

func (b *SyncSampleBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	var entryCount uint32
	if err = binary.Read(r, binary.BigEndian, &entryCount); err != nil {
		return
	}
	b.SampleNumbers = make([]uint32, (b.Size-b.headerSize()-4)/4)
	if err = binary.Read(r, binary.BigEndian, b.SampleNumbers); err != nil {
		return
	}
	return
}

func (b *SyncSampleBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint32(len(b.SampleNumbers))); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.SampleNumbers); err != nil {
		return
	}
	return
}
