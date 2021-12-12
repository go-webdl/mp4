package mp4

import (
	"encoding/binary"
	"io"
)

// 8.7.3.2 Sample Size Box

// Box Type: ‘stsz’
// Container: Sample Table Box (‘stbl’)
// Mandatory: Yes
// Quantity: Exactly one variant must be present

// This box contains the sample count and a table giving the size in bytes of
// each sample. This allows the media data itself to be unframed. The total
// number of samples in the media is always indicated in the sample count.
//
// There are two variants of the sample size box. The first variant has a fixed
// size 32‐bit field for representing the sample sizes; it permits defining a
// constant size for all samples in a track. The second variant permits smaller
// size fields, to save space when the sizes are varying but small. One of these
// boxes must be present; the first version is preferred for maximum
// compatibility.
type SampleSizeBox struct {
	FullHeader
	NullContainer

	// is integer specifying the default sample size. If all the samples are the
	// same size, this field contains that size value. If this field is set to
	// 0, then the samples have different sizes, and those sizes are stored in
	// the sample size table. If this field is not 0, it specifies the constant
	// sample size, and no array follows.
	SampleSize uint32
	Entries    []SampleSizeEntry
}

var _ Box = (*SampleSizeBox)(nil)

func init() {
	BoxRegistry[StszBoxType] = func() Box { return &SampleSizeBox{} }
}

type SampleSizeEntry struct {
	// is an integer specifying the size of a sample, indexed by its number.
	EntrySize uint32

	// is an integer that gives the number of samples in each of these chunks
	SamplesPerChunk uint32

	// is an integer that gives the index of the sample entry that describes the
	// samples in this chunk. The index ranges from 1 to the number of sample
	// entries in the Sample Description Box
	SampleDescrptionIndex uint32
}

func (b SampleSizeBox) Mp4BoxType() BoxType {
	return StszBoxType
}

func (b *SampleSizeBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 4 // unsigned int(32) sample_size;
	b.Size += 4 // unsigned int(32) sample_count;
	// if (sample_size==0) {
	//     for (i=1; i <= sample_count; i++) {
	//         unsigned int(32) entry_size;
	//     }
	// }
	if b.SampleSize == 0 {
		b.Size += 4 * uint32(len(b.Entries))
	}
	return b.Size
}

func (b *SampleSizeBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.SampleSize); err != nil {
		return
	}
	var sampleCount uint32
	if err = binary.Read(r, binary.BigEndian, &sampleCount); err != nil {
		return
	}
	if b.SampleSize == 0 {
		b.Entries = make([]SampleSizeEntry, sampleCount)
		if err = binary.Read(r, binary.BigEndian, b.Entries); err != nil {
			return
		}
	}
	return
}

func (b *SampleSizeBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.SampleSize); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint32(len(b.Entries))); err != nil {
		return
	}
	if b.SampleSize == 0 {
		if err = binary.Write(w, binary.BigEndian, b.Entries); err != nil {
			return
		}
	}
	return
}
