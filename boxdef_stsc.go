package mp4

import (
	"encoding/binary"
	"io"
)

// 8.7.4 Sample To Chunk Box

// Box Type: ‘stsc’
// Container: Sample Table Box (‘stbl’)
// Mandatory: Yes
// Quantity: Exactly one

// Samples within the media data are grouped into chunks. Chunks can be of
// different sizes, and the samples within a chunk can have different sizes.
// This table can be used to find the chunk that contains a sample, its
// position, and the associated sample description.
//
// The table is compactly coded. Each entry gives the index of the first chunk
// of a run of chunks with the same characteristics. By subtracting one entry
// here from the previous one, you can compute how many chunks are in this run.
// You can convert this to a sample count by multiplying by the appropriate
// samples‐per‐chunk.
type SampleToChunkBox struct {
	FullHeader
	NullContainer
	Entries []SampleToChunkEntry
}

var _ Box = (*SampleToChunkBox)(nil)

func init() {
	BoxRegistry[StscBoxType] = func() Box { return &SampleToChunkBox{} }
}

type SampleToChunkEntry struct {
	// is an integer that gives the index of the first chunk in this run of
	// chunks that share the same samples‐per‐chunk and
	// sample‐description‐index; the index of the first chunk in a track has the
	// value 1 (the first_chunk field in the first record of this box has the
	// value 1, identifying that the first sample maps to the first chunk).
	FirstChunk uint32

	// is an integer that gives the number of samples in each of these chunks
	SamplesPerChunk uint32

	// is an integer that gives the index of the sample entry that describes the
	// samples in this chunk. The index ranges from 1 to the number of sample
	// entries in the Sample Description Box
	SampleDescrptionIndex uint32
}

func (b SampleToChunkBox) Mp4BoxType() BoxType {
	return StscBoxType
}

func (b *SampleToChunkBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 4 // unsigned int(32) entry_count;
	// for (i=0; i < entry_count; i++) {
	//     unsigned int(32) first_chunk;
	//     unsigned int(32) samples_per_chunk;
	//     unsigned int(32) sample_description_index;
	// }
	b.Size += 12 * uint32(len(b.Entries))
	return b.Size
}

func (b *SampleToChunkBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	var entryCount uint32
	if err = binary.Read(r, binary.BigEndian, &entryCount); err != nil {
		return
	}
	b.Entries = make([]SampleToChunkEntry, entryCount)
	if err = binary.Read(r, binary.BigEndian, b.Entries); err != nil {
		return
	}
	return
}

func (b *SampleToChunkBox) Mp4BoxWrite(w io.Writer) (err error) {
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
