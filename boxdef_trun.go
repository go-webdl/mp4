package mp4

import (
	"encoding/binary"
	"io"
)

// 8.8.8 Track Fragment Run Box

// Box Type: ‘trun’
// Container: Track Fragment Box ('traf')
// Mandatory: No
// Quantity: Zero or more

// Within the Track Fragment Box, there are zero or more Track Run Boxes. If the
// duration‐is‐empty flag is set in the tf_flags, there are no track runs. A
// track run documents a contiguous set of samples for a track.
//
// The number of optional fields is determined from the number of bits set in
// the lower byte of the flags, and the size of a record from the bits set in
// the second byte of the flags. This procedure shall be followed, to allow for
// new fields to be defined.
//
// If the data‐offset is not present, then the data for this run starts
// immediately after the data of the previous run, or at the base‐data‐offset
// defined by the track fragment header if this is the first run in a track
// fragment, If the data‐offset is present, it is relative to the
// base‐data‐offset established in the track fragment header.
type TrackRunBox struct {
	FullHeader
	NullContainer

	// the number of samples being added in this run; also the number of rows in
	// the following table (the rows can be empty)
	SampleCount uint32

	// is added to the implicit or explicit data_offset established in the track
	// fragment header.
	DataOffset int32

	// provides a set of flags for the first sample only of this run.
	FirstSampleFlags uint32

	Samples []TrackRunSampleEntry
}

type TrackRunSampleEntry struct {
	SampleDuration              uint32
	SampleSize                  uint32
	SampleFlags                 uint32
	SampleCompositionTimeOffset int64
}

const (
	FLAG_TRUN_DATA_OFFSET uint32 = 0x01

	// this over‐rides the default flags for the first sample only. This makes
	// it possible to record a group of frames where the first is a key and the
	// rest are difference frames, without supplying explicit flags for every
	// sample. If this flag and field are used, sampleflags shall not be
	// present.
	FLAG_TRUN_FIRST_SAMPLE_FLAGS uint32 = 0x04

	// indicates that each sample has its own duration, otherwise the default is
	// used.
	FLAG_TRUN_SAMPLE_DURATION uint32 = 0x100

	// each sample has its own size, otherwise the default is used.
	FLAG_TRUN_SAMPLE_SIZE uint32 = 0x200

	// each sample has its own flags, otherwise the default is used.
	FLAG_TRUN_SAMPLE_FLAGS uint32 = 0x400

	// each sample has a composition time offset (e.g. as used for I/P/B video
	// in MPEG).
	FLAG_TRUN_SAMPLE_COMPOSITION_TIME_OFFSET uint32 = 0x800
)

var _ Box = (*TrackRunBox)(nil)

func init() {
	BoxRegistry[TrunBoxType] = func() Box { return &TrackRunBox{} }
}

func (b TrackRunBox) Mp4BoxType() BoxType {
	return TrunBoxType
}

func (b *TrackRunBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 4 // unsigned int(32) sample_count;
	flags := b.Mp4BoxFlags()
	if flags&FLAG_TRUN_DATA_OFFSET > 0 {
		b.Size += 4 // signed int(32) data_offset;
	}
	if flags&FLAG_TRUN_FIRST_SAMPLE_FLAGS > 0 {
		b.Size += 4 // unsigned int(32) first_sample_flags;
	}
	var entrySize uint32
	if flags&FLAG_TRUN_SAMPLE_DURATION > 0 {
		entrySize += 4 // unsigned int(32) sample_duration;
	}
	if flags&FLAG_TRUN_SAMPLE_SIZE > 0 {
		entrySize += 4 // unsigned int(32) sample_size;
	}
	if flags&FLAG_TRUN_SAMPLE_FLAGS > 0 {
		entrySize += 4 // unsigned int(32) sample_flags;
	}
	if flags&FLAG_TRUN_SAMPLE_COMPOSITION_TIME_OFFSET > 0 {
		entrySize += 4 // unsigned int(32) sample_composition_time_offset;
	}
	b.Size += entrySize * uint32(len(b.Samples))
	return b.Size
}

func (b *TrackRunBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.SampleCount); err != nil {
		return
	}
	flags := b.Mp4BoxFlags()
	if flags&FLAG_TRUN_DATA_OFFSET > 0 {
		if err = binary.Read(r, binary.BigEndian, &b.DataOffset); err != nil {
			return
		}
	}
	if flags&FLAG_TRUN_FIRST_SAMPLE_FLAGS > 0 {
		if err = binary.Read(r, binary.BigEndian, &b.FirstSampleFlags); err != nil {
			return
		}
	}
	b.Samples = make([]TrackRunSampleEntry, b.SampleCount)
	for i := uint32(0); i < b.SampleCount; i++ {
		if flags&FLAG_TRUN_SAMPLE_DURATION > 0 {
			if err = binary.Read(r, binary.BigEndian, &b.Samples[i].SampleDuration); err != nil {
				return
			}
		}
		if flags&FLAG_TRUN_SAMPLE_SIZE > 0 {
			if err = binary.Read(r, binary.BigEndian, &b.Samples[i].SampleSize); err != nil {
				return
			}
		}
		if flags&FLAG_TRUN_SAMPLE_FLAGS > 0 {
			if err = binary.Read(r, binary.BigEndian, &b.Samples[i].SampleFlags); err != nil {
				return
			}
		}
		if flags&FLAG_TRUN_SAMPLE_COMPOSITION_TIME_OFFSET > 0 {
			if b.Version == 0 {
				var tmp uint32
				if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
					return
				}
				b.Samples[i].SampleCompositionTimeOffset = int64(tmp)
			} else {
				var tmp int32
				if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
					return
				}
				b.Samples[i].SampleCompositionTimeOffset = int64(tmp)
			}
		}
	}
	return
}

func (b *TrackRunBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.SampleCount); err != nil {
		return
	}
	flags := b.Mp4BoxFlags()
	if flags&FLAG_TRUN_DATA_OFFSET > 0 {
		if err = binary.Write(w, binary.BigEndian, b.DataOffset); err != nil {
			return
		}
	}
	if flags&FLAG_TRUN_FIRST_SAMPLE_FLAGS > 0 {
		if err = binary.Write(w, binary.BigEndian, b.FirstSampleFlags); err != nil {
			return
		}
	}
	for _, sample := range b.Samples {
		if flags&FLAG_TRUN_SAMPLE_DURATION > 0 {
			if err = binary.Write(w, binary.BigEndian, sample.SampleDuration); err != nil {
				return
			}
		}
		if flags&FLAG_TRUN_SAMPLE_SIZE > 0 {
			if err = binary.Write(w, binary.BigEndian, sample.SampleSize); err != nil {
				return
			}
		}
		if flags&FLAG_TRUN_SAMPLE_FLAGS > 0 {
			if err = binary.Write(w, binary.BigEndian, sample.SampleFlags); err != nil {
				return
			}
		}
		if flags&FLAG_TRUN_SAMPLE_COMPOSITION_TIME_OFFSET > 0 {
			if b.Version == 1 {
				if err = binary.Write(w, binary.BigEndian, int32(sample.SampleCompositionTimeOffset)); err != nil {
					return
				}
			} else {
				if err = binary.Write(w, binary.BigEndian, uint32(sample.SampleCompositionTimeOffset)); err != nil {
					return
				}
			}
		}
	}
	return
}
