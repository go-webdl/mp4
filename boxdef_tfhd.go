package mp4

import (
	"encoding/binary"
	"io"
)

// 8.8.7 Track Fragment Header Box

// Box Type: ‘tfhd’
// Container: Track Fragment Box ('traf')
// Mandatory: Yes
// Quantity: Exactly one

// Each movie fragment can add zero or more fragments to each track; and a track
// fragment can add zero or more contiguous runs of samples. The track fragment
// header sets up information and defaults used for those runs of samples.
//
// The base‐data‐offset, if explicitly provided, is a data offset that is
// identical to a chunk offset in the Chunk Offset Box, i.e. applying to the
// complete file (e.g. starting with a file‐type box and movie box). In
// circumstances when the complete file does not exist or its size is unknown,
// it may be impossible to use an explicit base‐data‐offset; then, offsets need
// to be established relative to the movie fragment.
type TrackFragmentHeaderBox struct {
	FullHeader
	NullContainer

	TrackID uint32

	// the base offset to use when calculating data offsets
	BaseDataOffset uint64

	SampleDescrptionIndex uint32
	DefaultSampleDuration uint32
	DefaultSampleSize     uint32
	DefaultSampleFlags    uint32
}

const (
	// Indicates the presence of the base‐data‐offset field. This provides an
	// explicit anchor for the data offsets in each track run (see below). If
	// not provided and if the default‐base‐is‐moof flag is not set, the
	// base‐data‐offset for the first track in the movie fragment is the
	// position of the first byte of the enclosing Movie Fragment Box, and for
	// second and subsequent track fragments, the default is the end of the data
	// defined by the preceding track fragment. Fragments 'inheriting' their
	// offset in this way must all use the same datareference (i.e., the data
	// for these tracks must be in the same file)
	FLAG_TFHD_BASE_DATA_OFFSET uint32 = 0x01

	// Indicates the presence of this field, which over‐rides, in this fragment,
	// the default set up in the Track Extends Box.
	FLAG_TFHD_SAMPLE_DESCRIPTION_INDEX uint32 = 0x02
	FLAG_TFHD_DEFAULT_SAMPLE_DURATION  uint32 = 0x08
	FLAG_TFHD_DEFAULT_SAMPLE_SIZE      uint32 = 0x10
	FLAG_TFHD_DEFAULT_SAMPLE_FLAGS     uint32 = 0x20

	// This indicates that the duration provided in either
	// default‐sampleduration, or by the default‐duration in the Track Extends
	// Box, is empty, i.e. that there are no samples for this time interval. It
	// is an error to make a presentation that has both edit lists in the Movie
	// Box, and empty‐duration fragments.
	FLAG_TFHD_DURATION_IS_EMPTY uint32 = 0x010000

	// If base‐data‐offset‐present is 1, this flag is ignored. If
	// base‐dataoffset‐ present is zero, this indicates that the
	// base‐data‐offset for this track fragment is the position of the first
	// byte of the enclosing Movie Fragment Box. Support for the
	// default‐base‐ismoof flag is required under the ‘iso5’ brand, and it shall
	// not be used in brands or compatible brands earlier than iso5.
	FLAG_TFHD_DEFAULT_BASE_IS_MOOF uint32 = 0x020000
)

var _ Box = (*TrackFragmentHeaderBox)(nil)

func init() {
	BoxRegistry[TfhdBoxType] = func() Box { return &TrackFragmentHeaderBox{} }
}

func (b TrackFragmentHeaderBox) Mp4BoxType() BoxType {
	return TfhdBoxType
}

func (b *TrackFragmentHeaderBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 4 // unsigned int(32) track_ID;
	flags := b.Mp4BoxFlags()
	if flags&FLAG_TFHD_BASE_DATA_OFFSET > 0 {
		b.Size += 8 // unsigned int(64) base_data_offset;
	}
	if flags&FLAG_TFHD_SAMPLE_DESCRIPTION_INDEX > 0 {
		b.Size += 4 // unsigned int(32) sample_description_index;
	}
	if flags&FLAG_TFHD_DEFAULT_SAMPLE_DURATION > 0 {
		b.Size += 4 // unsigned int(32) default_sample_duration;
	}
	if flags&FLAG_TFHD_DEFAULT_SAMPLE_SIZE > 0 {
		b.Size += 4 // unsigned int(32) default_sample_size;
	}
	if flags&FLAG_TFHD_DEFAULT_SAMPLE_FLAGS > 0 {
		b.Size += 4 // unsigned int(32) default_sample_flags;
	}
	return b.Size
}

func (b *TrackFragmentHeaderBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.TrackID); err != nil {
		return
	}
	flags := b.Mp4BoxFlags()
	if flags&FLAG_TFHD_BASE_DATA_OFFSET > 0 {
		if err = binary.Read(r, binary.BigEndian, &b.BaseDataOffset); err != nil {
			return
		}
	}
	if flags&FLAG_TFHD_SAMPLE_DESCRIPTION_INDEX > 0 {
		if err = binary.Read(r, binary.BigEndian, &b.SampleDescrptionIndex); err != nil {
			return
		}
	}
	if flags&FLAG_TFHD_DEFAULT_SAMPLE_DURATION > 0 {
		if err = binary.Read(r, binary.BigEndian, &b.DefaultSampleDuration); err != nil {
			return
		}
	}
	if flags&FLAG_TFHD_DEFAULT_SAMPLE_SIZE > 0 {
		if err = binary.Read(r, binary.BigEndian, &b.DefaultSampleSize); err != nil {
			return
		}
	}
	if flags&FLAG_TFHD_DEFAULT_SAMPLE_FLAGS > 0 {
		if err = binary.Read(r, binary.BigEndian, &b.DefaultSampleFlags); err != nil {
			return
		}
	}
	return
}

func (b *TrackFragmentHeaderBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.TrackID); err != nil {
		return
	}
	flags := b.Mp4BoxFlags()
	if flags&FLAG_TFHD_BASE_DATA_OFFSET > 0 {
		if err = binary.Write(w, binary.BigEndian, b.BaseDataOffset); err != nil {
			return
		}
	}
	if flags&FLAG_TFHD_SAMPLE_DESCRIPTION_INDEX > 0 {
		if err = binary.Write(w, binary.BigEndian, b.SampleDescrptionIndex); err != nil {
			return
		}
	}
	if flags&FLAG_TFHD_DEFAULT_SAMPLE_DURATION > 0 {
		if err = binary.Write(w, binary.BigEndian, b.DefaultSampleDuration); err != nil {
			return
		}
	}
	if flags&FLAG_TFHD_DEFAULT_SAMPLE_SIZE > 0 {
		if err = binary.Write(w, binary.BigEndian, b.DefaultSampleSize); err != nil {
			return
		}
	}
	if flags&FLAG_TFHD_DEFAULT_SAMPLE_FLAGS > 0 {
		if err = binary.Write(w, binary.BigEndian, b.DefaultSampleFlags); err != nil {
			return
		}
	}
	return
}
