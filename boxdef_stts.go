package mp4

import (
	"encoding/binary"
	"io"
)

// 8.6.1.2 Decoding Time to Sample Box

// Box Type: ‘stts’
// Container: Sample Table Box (‘stbl’)
// Mandatory: Yes
// Quantity: Exactly one

// This box contains a compact version of a table that allows indexing from
// decoding time to sample number. Other tables give sample sizes and pointers,
// from the sample number. Each entry in the table gives the number of
// consecutive samples with the same time delta, and the delta of those samples.
// By adding the deltas a complete time‐to‐sample map may be built.
//
// The Decoding Time to Sample Box contains decode time delta's: DT(n+1) = DT(n)
// + STTS(n) where STTS(n) is the (uncompressed) table entry for sample n.
//
// The sample entries are ordered by decoding time stamps; therefore the deltas
// are all non‐negative.
//
// The DT axis has a zero origin; DT(i) = SUM(for j=0 to i‐1 of delta(j)), and
// the sum of all deltas gives the length of the media in the track (not mapped
// to the overall timescale, and not considering any edit list).
//
// The Edit List Box provides the initial CT value if it is non‐empty
// (non‐zero).
type TimeToSampleBox struct {
	FullHeader
	NullContainer
	Entries []TimeToSampleEntry
}

var _ Box = (*TimeToSampleBox)(nil)

func init() {
	BoxRegistry[SttsBoxType] = func() Box { return &TimeToSampleBox{} }
}

type TimeToSampleEntry struct {
	// is an integer that counts the number of consecutive samples that have the
	// given duration.
	SampleCount uint32

	// is an integer that gives the delta of these samples in the time‐scale of
	// the media.
	SampleDelta uint32
}

func (b TimeToSampleBox) Mp4BoxType() BoxType {
	return SttsBoxType
}

func (b *TimeToSampleBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 4 // unsigned int(32) entry_count;
	// for (i=0; i < entry_count; i++) {
	//     unsigned int(32) sample_count;
	//     unsigned int(32) sample_delta;
	// }
	b.Size += 8 * uint32(len(b.Entries))
	return b.Size
}

func (b *TimeToSampleBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	var entryCount uint32
	if err = binary.Read(r, binary.BigEndian, &entryCount); err != nil {
		return
	}
	b.Entries = make([]TimeToSampleEntry, entryCount)
	if err = binary.Read(r, binary.BigEndian, b.Entries); err != nil {
		return
	}
	return
}

func (b *TimeToSampleBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint32(len(b.Entries))); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.Entries); err != nil {
		return
	}
	return
}
