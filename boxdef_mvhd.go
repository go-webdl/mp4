package mp4

import (
	"encoding/binary"
	"io"
)

// 8.2.2 Movie Header Box

// Box Type: ‘mvhd’
// Container: Movie Box (‘moov’)
// Mandatory: Yes
// Quantity: Exactly one

// This box defines overall information which is media‐independent, and relevant
// to the entire presentation considered as a whole.
type MovieHeaderBox struct {
	FullHeader
	NullContainer

	// is an integer that declares the creation time of the presentation (in
	// seconds since midnight, Jan. 1, 1904, in UTC time)
	CreationTime uint64

	// is an integer that declares the most recent time the presentation was
	// modified (in seconds since midnight, Jan. 1, 1904, in UTC time)
	ModificationTime uint64

	// is an integer that specifies the time‐scale for the entire presentation;
	// this is the number of time units that pass in one second. For example, a
	// time coordinate system that measures time in sixtieths of a second has a
	// time scale of 60.
	Timescale uint32

	// is an integer that declares length of the presentation (in the indicated
	// timescale). This property is derived from the presentation’s tracks: the
	// value of this field corresponds to the duration of the longest track in
	// the presentation. If the duration cannot be determined then duration is
	// set to all 1s.
	Duration uint64

	// is a fixed point 16.16 number that indicates the preferred rate to play
	// the presentation; 1.0 (0x00010000) is normal forward playback
	Rate int32

	// is a fixed point 8.8 number that indicates the preferred playback volume.
	// 1.0 (0x0100) is full volume.
	Volume int16

	// provides a transformation matrix for the video; (u,v,w) are restricted
	// here to (0,0,1), hex values (0,0,0x40000000).
	Matrix [9]int32

	// is a non‐zero integer that indicates a value to use for the track ID of
	// the next track to be added to this presentation. Zero is not a valid
	// track ID value. The value of next_track_ID shall be larger than the
	// largest track‐ID in use. If this value is equal to all 1s (32‐bit
	// maxint), and a new media track is to be added, then a search must be made
	// in the file for an unused track identifier.
	NextTrackID uint32
}

var _ Box = (*MovieHeaderBox)(nil)

func init() {
	BoxRegistry[MvhdBoxType] = func() Box { return &MovieHeaderBox{} }
}

func (b MovieHeaderBox) Mp4BoxType() BoxType {
	return MvhdBoxType
}

func (b *MovieHeaderBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	if b.Version == 1 {
		b.Size += 8 // unsigned int(64) creation_time;
		b.Size += 8 // unsigned int(64) modification_time;
		b.Size += 4 // unsigned int(32) timescale;
		b.Size += 8 // unsigned int(64) duration;
	} else {
		b.Size += 4 // unsigned int(32) creation_time;
		b.Size += 4 // unsigned int(32) modification_time;
		b.Size += 4 // unsigned int(32) timescale;
		b.Size += 4 // unsigned int(32) duration;
	}
	b.Size += 4     // template int(32) rate = 0x00010000;
	b.Size += 2     // template int(16) volume = 0x0100;
	b.Size += 2     // const bit(16) reserved = 0;
	b.Size += 4 * 2 // const unsigned int(32)[2] reserved = 0;
	b.Size += 4 * 9 // template int(32)[9] matrix;
	b.Size += 4 * 6 // bit(32)[6] pre_defined = 0;
	b.Size += 4     // unsigned int(32) next_track_ID;
	return b.Size
}

func (b *MovieHeaderBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
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
		if err = binary.Read(r, binary.BigEndian, &b.Timescale); err != nil {
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
		if err = binary.Read(r, binary.BigEndian, &b.Timescale); err != nil {
			return
		}
		if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
			return
		}
		b.Duration = uint64(tmp)
	}
	if err = binary.Read(r, binary.BigEndian, &b.Rate); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.Volume); err != nil {
		return
	}
	var reserved [5]uint16
	if err = binary.Read(r, binary.BigEndian, &reserved); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.Matrix); err != nil {
		return
	}
	var reserved2 [6]uint32
	if err = binary.Read(r, binary.BigEndian, &reserved2); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.NextTrackID); err != nil {
		return
	}
	return
}

func (b *MovieHeaderBox) Mp4BoxWrite(w io.Writer) (err error) {
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
		if err = binary.Write(w, binary.BigEndian, b.Timescale); err != nil {
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
		if err = binary.Write(w, binary.BigEndian, b.Timescale); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, uint32(b.Duration)); err != nil {
			return
		}
	}
	if err = binary.Write(w, binary.BigEndian, b.Rate); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.Volume); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, [5]uint16{}); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.Matrix); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, [6]uint32{}); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.NextTrackID); err != nil {
		return
	}
	return
}
