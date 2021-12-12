package mp4

import (
	"encoding/binary"
	"io"
)

// 8.8.3 Track Extends Box

// Box Type: ‘trex’
// Container: Movie Extends Box (‘mvex’)
// Mandatory: Yes
// Quantity: Exactly one for each track in the Movie Box

// This sets up default values used by the movie fragments. By setting defaults
// in this way, space and complexity can be saved in each Track Fragment Box.
//
// The sample flags field in sample fragments (default_sample_flags here and in
// a Track Fragment Header Box, and sample_flags and first_sample_flags in a
// Track Fragment Run Box) is coded as a 32‐bit value. It has the following
// structure:
//
//     bit(4) reserved=0;
//     unsigned int(2) is_leading;
//     unsigned int(2) sample_depends_on;
//     unsigned int(2) sample_is_depended_on;
//     unsigned int(2) sample_has_redundancy;
//     bit(3) sample_padding_value;
//     bit(1) sample_is_non_sync_sample;
//     unsigned int(16) sample_degradation_priority;
//
// The is_leading, sample_depends_on, sample_is_depended_on and
// sample_has_redundancy values are defined as documented in the Independent and
// Disposable Samples Box.
//
// The flag sample_is_non_sync_sample provides the same information as the sync
// sample table [8.6.2]. When this value is set 0 for a sample, it is the same
// as if the sample were not in a movie fragment and marked with an entry in the
// sync sample table (or, if all samples are sync samples, the sync sample table
// were absent).
//
// The sample_padding_value is defined as for the padding bits table. The
// sample_degradation_priority is defined as for the degradation priority table.
type TrackExtendsBox struct {
	FullHeader
	NullContainer

	// identifies the track; this shall be the track ID of a track in the Movie
	// Box
	TrackID uint32

	// these fields set up defaults used in the track fragments.
	DefaultSampleDescrptionIndex uint32
	DefaultSampleDuration        uint32
	DefaultSampleSize            uint32
	DefaultSampleFlags           uint32
}

var _ Box = (*TrackExtendsBox)(nil)

func init() {
	BoxRegistry[TrexBoxType] = func() Box { return &TrackExtendsBox{} }
}

func (b TrackExtendsBox) Mp4BoxType() BoxType {
	return TrexBoxType
}

func (b *TrackExtendsBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 4 // unsigned int(32) track_ID;
	b.Size += 4 // unsigned int(32) default_sample_description_index;
	b.Size += 4 // unsigned int(32) default_sample_duration;
	b.Size += 4 // unsigned int(32) default_sample_size;
	b.Size += 4 // unsigned int(32) default_sample_flags;
	return b.Size
}

func (b *TrackExtendsBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.TrackID); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.DefaultSampleDescrptionIndex); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.DefaultSampleDuration); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.DefaultSampleSize); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.DefaultSampleFlags); err != nil {
		return
	}
	return
}

func (b *TrackExtendsBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.TrackID); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.DefaultSampleDescrptionIndex); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.DefaultSampleDuration); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.DefaultSampleSize); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.DefaultSampleFlags); err != nil {
		return
	}
	return
}
