package mp4

import (
	"encoding/binary"
	"io"

	"golang.org/x/text/language"
)

// 8.4.1 Media Box

// Box Type: ‘mdhd’
// Container: Media Box (‘mdia’)
// Mandatory: Yes
// Quantity: Exactly one

// The media header declares overall information that is media‐independent, and
// relevant to characteristics of the media in a track.
type MediaHeaderBox struct {
	FullHeader
	NullContainer

	// is an integer that declares the creation time of the media in this track
	// (in seconds since midnight, Jan. 1, 1904, in UTC time).
	CreationTime uint64

	// is an integer that declares the most recent time the media in this track was
	// modified (in seconds since midnight, Jan. 1, 1904, in UTC time).
	ModificationTime uint64

	// is an integer that specifies the time‐scale for this media; this is the
	// number of time units that pass in one second. For example, a time
	// coordinate system that measures time in sixtieths of a second has a time
	// scale of 60.
	Timescale uint32

	// is an integer that declares the duration of this media (in the scale of
	// the timescale). If the duration cannot be determined then duration is set
	// to all 1s.
	Duration uint64

	// declares the language code for this media. See ISO 639‐2/T for the set of three
	// character codes. Each character is packed as the difference between its ASCII
	// value and 0x60. Since the code is confined to being three lower‐case letters,
	// these values are strictly positive.
	Language language.Base
}

var _ Box = (*MediaHeaderBox)(nil)

func init() {
	BoxRegistry[MdhdBoxType] = func() Box { return &MediaHeaderBox{} }
}

func (b MediaHeaderBox) Mp4BoxType() BoxType {
	return MdhdBoxType
}

func (b *MediaHeaderBox) Mp4BoxUpdate() uint32 {
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
	// bit(1) pad = 0;
	// unsigned int(5)[3] language; // ISO-639-2/T language code
	b.Size += 2
	b.Size += 2 // unsigned int(16) pre_defined = 0;
	return b.Size
}

func (b *MediaHeaderBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
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
	var lang uint16
	if err = binary.Read(r, binary.BigEndian, &lang); err != nil {
		return
	}
	if b.Language, err = language.ParseBase(string([]byte{
		(byte(lang>>10) & 0x1F) + 0x60,
		(byte(lang>>5) & 0x1F) + 0x60,
		(byte(lang) & 0x1F) + 0x60,
	})); err != nil {
		return
	}
	var tmp uint16
	if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
		return
	}
	return
}

func (b *MediaHeaderBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}

	if b.Version == 1 {
		if err = binary.Write(w, binary.BigEndian, &b.CreationTime); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, &b.ModificationTime); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, &b.Timescale); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, &b.Duration); err != nil {
			return
		}
	} else {
		tmp := uint32(b.CreationTime)
		if err = binary.Write(w, binary.BigEndian, tmp); err != nil {
			return
		}
		tmp = uint32(b.ModificationTime)
		if err = binary.Write(w, binary.BigEndian, tmp); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, b.Timescale); err != nil {
			return
		}
		tmp = uint32(b.Duration)
		if err = binary.Write(w, binary.BigEndian, tmp); err != nil {
			return
		}
	}
	iso3 := b.Language.ISO3()
	lang := (uint16(iso3[0]) << 10) | (uint16(iso3[0]) << 5) | uint16(iso3[0])
	if err = binary.Write(w, binary.BigEndian, lang); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint16(0)); err != nil {
		return
	}
	return
}
