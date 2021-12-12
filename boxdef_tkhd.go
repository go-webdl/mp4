package mp4

import (
	"encoding/binary"
	"io"
)

// 8.3.2 Track Header Box

// Box Type: ‘tkhd’
// Container: Track Box (‘trak’)
// Mandatory: Yes
// Quantity: Exactly one

// This box specifies the characteristics of a single track. Exactly one Track
// Header Box is contained in a track.
//
// In the absence of an edit list, the presentation of a track starts at the
// beginning of the overall presentation. An empty edit is used to offset the
// start time of a track.
//
// The default value of the track header flags for media tracks is 7
// (track_enabled, track_in_movie, track_in_preview). If in a presentation all
// tracks have neither track_in_movie nor track_in_preview set, then all tracks
// shall be treated as if both flags were set on all tracks. Server hint tracks
// should have the track_in_movie and track_in_preview set to 0, so that they
// are ignored for local playback and preview.
//
// Under the ‘iso3’ brand or brands that share its requirements, the width and
// height in the track header are measured on a notional 'square' (uniform)
// grid. Track video data is normalized to these dimensions (logically) before
// any transformation or placement caused by a layup or composition system.
// Track (and movie) matrices, if used, also operate in this uniformly‐scaled
// space.
//
// The duration field here does not include the duration of following movie
// fragments, if any, but only of the media in the enclosing Movie Box. The
// Movie Extends Header box may be used to document the duration including movie
// fragments, when desired and possible.
type TrackHeaderBox struct {
	FullHeader
	NullContainer

	// is an integer that declares the creation time of this track (in seconds
	// since midnight, Jan. 1, 1904, in UTC time).
	CreationTime uint64

	// is an integer that declares the most recent time the track was modified
	// (in seconds since midnight, Jan. 1, 1904, in UTC time).
	ModificationTime uint64

	// is an integer that uniquely identifies this track over the entire
	// life‐time of this presentation. Track IDs are never re‐used and cannot be
	// zero.
	TrackID uint32

	// is an integer that indicates the duration of this track (in the timescale
	// indicated in the Movie Header Box). The value of this field is equal to
	// the sum of the durations of all of the track’s edits. If there is no edit
	// list, then the duration is the sum of the sample durations, converted
	// into the timescale in the Movie Header Box. If the duration of this track
	// cannot be determined then duration is set to all 1s.
	Duration uint64

	// specifies the front‐to‐back ordering of video tracks; tracks with lower
	// numbers are closer to the viewer. 0 is the normal value, and ‐1 would be
	// in front of track 0, and so on.
	Layer int16

	// is an integer that specifies a group or collection of tracks. If this
	// field is 0 there is no information on possible relations to other tracks.
	// If this field is not 0, it should be the same for tracks that contain
	// alternate data for one another and different for tracks belonging to
	// different such groups. Only one track within an alternate group should be
	// played or streamed at any one time, and must be distinguishable from
	// other tracks in the group via attributes such as bitrate, codec,
	// language, packet size etc. A group may have only one member.
	AlternateGroup int16

	// is a fixed 8.8 value specifying the track's relative audio volume. Full
	// volume is 1.0 (0x0100) and is the normal value. Its value is irrelevant
	// for a purely visual track. Tracks may be composed by combining them
	// according to their volume, and then using the overall Movie Header Box
	// volume setting; or more complex audio composition (e.g. MPEG‐4 BIFS) may
	// be used.
	Volume int16

	// provides a transformation matrix for the video; (u,v,w) are restricted
	// here to (0,0,1), hex (0,0,0x40000000).
	Matrix [9]int32

	// width and height fixed‐point 16.16 values are track‐dependent as follows:
	//
	// For text and subtitle tracks, they may, depending on the coding format,
	// describe the suggested size of the rendering area. For such tracks, the
	// value 0x0 may also be used to indicate that the data may be rendered at
	// any size, that no preferred size has been indicated and that the actual
	// size may be determined by the external context or by reusing the width
	// and height of another track. For those tracks, the flag
	// track_size_is_aspect_ratio may also be used.
	//
	// For non‐visual tracks (e.g. audio), they should be set to zero.
	//
	// For all other tracks, they specify the track's visual presentation size.
	// These need not be the same as the pixel dimensions of the images, which
	// is documented in the sample description(s); all images in the sequence
	// are scaled to this size, before any overall transformation of the track
	// represented by the matrix. The pixel dimensions of the images are the
	// default values.
	Width  uint32
	Height uint32
}

const (
	// Indicates that the track is enabled. Flag value is 0x000001. A disabled
	// track (the low bit is zero) is treated as if it were not present.
	FLAG_TKHD_TRACK_ENABLED uint32 = 0x000001

	// Indicates that the track is used in the presentation. Flag value is
	// 0x000002.
	FLAG_TKHD_TRACK_IN_MOVIE uint32 = 0x000002

	// Indicates that the track is used when previewing the presentation. Flag
	// value is 0x000004.
	FLAG_TKHD_TRACK_IN_PREVIEW uint32 = 0x000004

	// Indicates that the width and height fields are not expressed in pixel
	// units. The values have the same units but these units are not specified.
	// The values are only an indication of the desired aspect ratio. If the
	// aspect ratios of this track and other related tracks are not identical,
	// then the respective positioning of the tracks is undefined, possibly
	// defined by external contexts. Flag value is 0x000008.
	FLAG_TKHD_TRACK_SIZE_IS_ASPECT_RATIO uint32 = 0x000008
)

var _ Box = (*TrackHeaderBox)(nil)

func init() {
	BoxRegistry[TkhdBoxType] = func() Box { return &TrackHeaderBox{} }
}

func (b TrackHeaderBox) Mp4BoxType() BoxType {
	return TkhdBoxType
}

func (b *TrackHeaderBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	if b.Version == 1 {
		b.Size += 8 // unsigned int(64) creation_time;
		b.Size += 8 // unsigned int(64) modification_time;
		b.Size += 4 // unsigned int(32) track_ID;
		b.Size += 4 // const unsigned int(32) reserved = 0;
		b.Size += 8 // unsigned int(64) duration;
	} else {
		b.Size += 4 // unsigned int(32) creation_time;
		b.Size += 4 // unsigned int(32) modification_time;
		b.Size += 4 // unsigned int(32) track_ID;
		b.Size += 4 // const unsigned int(32) reserved = 0;
		b.Size += 4 // unsigned int(32) duration;
	}
	b.Size += 4 * 2 // const unsigned int(32)[2] reserved = 0;
	b.Size += 2     // template int(16) layer = 0;
	b.Size += 2     // template int(16) alternate_group = 0;
	b.Size += 2     // template int(16) volume;
	b.Size += 2     // const unsigned int(16) reserved = 0;
	b.Size += 4 * 9 // template int(32)[9] matrix;
	b.Size += 4     // unsigned int(32) width;
	b.Size += 4     // unsigned int(32) height;
	return b.Size
}

func (b *TrackHeaderBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if b.Version == 1 {
		if err = binary.Read(r, binary.BigEndian, &b.CreationTime); err != nil {
			return
		}
		if err = binary.Read(r, binary.BigEndian, &b.ModificationTime); err != nil {
			return
		}
		if err = binary.Read(r, binary.BigEndian, &b.TrackID); err != nil {
			return
		}
		var reserved uint32
		if err = binary.Read(r, binary.BigEndian, &reserved); err != nil {
			return
		}
		if err = binary.Read(r, binary.BigEndian, &b.Duration); err != nil {
			return
		}
	} else {
		var tmp uint32
		if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
			return
		}
		b.CreationTime = uint64(tmp)
		if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
			return
		}
		b.ModificationTime = uint64(tmp)
		if err = binary.Read(r, binary.BigEndian, &b.TrackID); err != nil {
			return
		}
		var reserved uint32
		if err = binary.Read(r, binary.BigEndian, &reserved); err != nil {
			return
		}
		if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
			return
		}
		b.Duration = uint64(tmp)
	}
	var reserved [2]uint32
	if err = binary.Read(r, binary.BigEndian, &reserved); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.Layer); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.AlternateGroup); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.Volume); err != nil {
		return
	}
	var reserved2 uint16
	if err = binary.Read(r, binary.BigEndian, &reserved2); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.Matrix); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.Width); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.Height); err != nil {
		return
	}
	return
}

func (b *TrackHeaderBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}

	if b.Version == 1 {
		if err = binary.Write(w, binary.BigEndian, b.CreationTime); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, b.ModificationTime); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, b.TrackID); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, uint32(0)); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, b.Duration); err != nil {
			return
		}
	} else {
		if err = binary.Write(w, binary.BigEndian, uint32(b.CreationTime)); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, uint32(b.ModificationTime)); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, b.TrackID); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, uint32(0)); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, uint32(b.Duration)); err != nil {
			return
		}
	}
	if err = binary.Write(w, binary.BigEndian, [2]uint32{}); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.Layer); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.AlternateGroup); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.Volume); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint16(0)); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.Matrix); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.Width); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.Height); err != nil {
		return
	}
	return
}
