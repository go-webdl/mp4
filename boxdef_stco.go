package mp4

import (
	"encoding/binary"
	"io"
)

// 8.7.5 Chunk Offset Box

// Box Type: ‘stco’
// Container: Sample Table Box (‘stbl’)
// Mandatory: Yes
// Quantity: Exactly one variant must be present

// The chunk offset table gives the index of each chunk into the containing
// file. There are two variants, permitting the use of 32‐bit or 64‐bit offsets.
// The latter is useful when managing very large presentations. At most one of
// these variants will occur in any single instance of a sample table.
//
// Offsets are file offsets, not the offset into any box within the file (e.g.
// Media Data Box). This permits referring to media data in files without any
// box structure. It does also mean that care must be taken when constructing a
// self‐contained ISO file with its metadata (Movie Box) at the front, as the
// size of the Movie Box will affect the chunk offsets to the media data.
type ChunkOffsetBox struct {
	FullHeader
	NullContainer
	Entries []ChunkOffsetEntry
}

var _ Box = (*ChunkOffsetBox)(nil)

func init() {
	BoxRegistry[StcoBoxType] = func() Box { return &ChunkOffsetBox{} }
}

type ChunkOffsetEntry struct {
	// is a 32 bit integer that gives the offset of the start of a chunk into
	// its containing media file.
	ChunkOffset uint32
}

func (b ChunkOffsetBox) Mp4BoxType() BoxType {
	return StcoBoxType
}

func (b *ChunkOffsetBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 4 // unsigned int(32) entry_count;
	// for (i=0; i < entry_count; i++) {
	//     unsigned int(32) chunk_offset;
	// }
	b.Size += 4 * uint32(len(b.Entries))
	return b.Size
}

func (b *ChunkOffsetBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	var entryCount uint32
	if err = binary.Read(r, binary.BigEndian, &entryCount); err != nil {
		return
	}
	b.Entries = make([]ChunkOffsetEntry, entryCount)
	if err = binary.Read(r, binary.BigEndian, b.Entries); err != nil {
		return
	}
	return
}

func (b *ChunkOffsetBox) Mp4BoxWrite(w io.Writer) (err error) {
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
