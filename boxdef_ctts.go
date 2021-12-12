package mp4

import (
	"encoding/binary"
	"io"
)

// 8.6.1.3 Composition Time to Sample Box

// Box Type: ‘ctts’
// Container: Sample Table Box (‘stbl’)
// Mandatory: No
// Quantity: Zero or one

// This box provides the offset between decoding time and composition time. In
// version 0 of this box the decoding time must be less than the composition
// time, and the offsets are expressed as unsigned numbers such that CT(n) =
// DT(n) + CTTS(n) where CTTS(n) is the (uncompressed) table entry for sample n.
// In version 1 of this box, the composition timeline and the decoding timeline
// are still derived from each other, but the offsets are signed. It is
// recommended that for the computed composition timestamps, there is exactly
// one with the value 0 (zero).
//
// For either version of the box, each sample must have a unique composition
// timestamp value, that is, the timestamp for two samples shall never be the
// same.
//
// It may be true that there is no frame to compose at time 0; the handling of
// this is unspecified (systems might display the first frame for longer, or a
// suitable fill colour).
//
// When version 1 of this box is used, the CompositionToDecodeBox may also be
// present in the sample table to relate the composition and decoding timelines.
// When backwards‐compatibility or compatibility with an unknown set of readers
// is desired, version 0 of this box should be used when possible. In either
// version of this box, but particularly under version 0, if it is desired that
// the media start at track time 0, and the first media sample does not have a
// composition time of 0, an edit list may be used to ‘shift’ the media to time
// 0.
//
// The composition time to sample table is optional and must only be present if
// DT and CT differ for any samples.
//
// Hint tracks do not use this box.
type CompositionOffsetBox struct {
	FullHeader
	NullContainer
	Entries []CompositionOffsetEntry
}

var _ Box = (*CompositionOffsetBox)(nil)

func init() {
	BoxRegistry[CttsBoxType] = func() Box { return &CompositionOffsetBox{} }
}

type CompositionOffsetEntry struct {
	// is an integer that counts the number of consecutive samples that have the
	// given offset.
	SampleCount uint32

	// is an integer that gives the offset between CT and DT, such that CT(n) =
	// DT(n) + CTTS(n).
	SampleOffset int64
}

func (b CompositionOffsetBox) Mp4BoxType() BoxType {
	return CttsBoxType
}

func (b *CompositionOffsetBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 4 // unsigned int(32) entry_count;
	if b.Version == 0 || b.Version == 1 {
		// for (i=0; i < entry_count; i++) {
		//     unsigned int(32) sample_count;
		//     unsigned int(32) sample_offset;
		// }
		b.Size += 8 * uint32(len(b.Entries))
	}
	return b.Size
}

func (b *CompositionOffsetBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	var entryCount uint32
	if err = binary.Read(r, binary.BigEndian, &entryCount); err != nil {
		return
	}
	b.Entries = make([]CompositionOffsetEntry, entryCount)
	if b.Version == 0 {
		for i := uint32(0); i < entryCount; i++ {
			var tmp uint32
			if err = binary.Read(r, binary.BigEndian, &b.Entries[i].SampleCount); err != nil {
				return
			}
			if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
				return
			}
			b.Entries[i].SampleOffset = int64(tmp)
		}
	} else if b.Version == 1 {
		for i := uint32(0); i < entryCount; i++ {
			var tmp int32
			if err = binary.Read(r, binary.BigEndian, &b.Entries[i].SampleCount); err != nil {
				return
			}
			if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
				return
			}
			b.Entries[i].SampleOffset = int64(tmp)
		}
	}
	return
}

func (b *CompositionOffsetBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint32(len(b.Entries))); err != nil {
		return
	}
	if b.Version == 0 {
		for _, entry := range b.Entries {
			if err = binary.Write(w, binary.BigEndian, entry.SampleCount); err != nil {
				return
			}
			if err = binary.Write(w, binary.BigEndian, uint32(entry.SampleOffset)); err != nil {
				return
			}
		}
	} else if b.Version == 1 {
		for _, entry := range b.Entries {
			if err = binary.Write(w, binary.BigEndian, entry.SampleCount); err != nil {
				return
			}
			if err = binary.Write(w, binary.BigEndian, int32(entry.SampleOffset)); err != nil {
				return
			}
		}
	}
	return
}
