package mp4

import (
	"encoding/binary"
	"io"

	"github.com/google/uuid"
)

type Header struct {
	Size     uint32
	Type     BoxType
	UserType UserType
}

type UserType uuid.UUID

func BoxTypeToUserType(boxType BoxType) UserType {
	return UserType{boxType[0], boxType[1], boxType[2], boxType[3], 0x00, 0x11, 0x00, 0x10, 0x80, 0x00, 0x00, 0xAA, 0x00, 0x38, 0x9B, 0x71}
}

func (h Header) HeaderSize() uint32 {
	if h.Type == UuidBoxType {
		return 24
	}
	return 8
}

func (h Header) Mp4BoxSize() uint32 {
	return h.Size
}

func (h Header) Mp4BoxType() BoxType {
	return h.Type
}

func (h *Header) Mp4BoxSetType(boxType BoxType) {
	h.Type = boxType
}

func (h Header) Mp4BoxUserType() UserType {
	if h.Type == UuidBoxType {
		return h.UserType
	}
	return BoxTypeToUserType(h.Type)
}

func (h *Header) Mp4BoxSetUserType(userType UserType) {
	h.UserType = userType
}

func (h *Header) ReadHeader(r io.Reader, header *Header) (err error) {
	if header == nil {
		if err = binary.Read(r, binary.BigEndian, &h.Size); err != nil {
			return
		}
		if err = binary.Read(r, binary.BigEndian, &h.Type); err != nil {
			return
		}
		if h.Type == UuidBoxType {
			if err = binary.Read(r, binary.BigEndian, &h.UserType); err != nil {
				return
			}
		}
	} else {
		*h = *header
	}
	return
}

func (h *Header) WriteHeader(w io.Writer) (err error) {
	if err = binary.Write(w, binary.BigEndian, h.Size); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, h.Type); err != nil {
		return
	}
	if h.Type == UuidBoxType {
		if err = binary.Write(w, binary.BigEndian, h.UserType); err != nil {
			return
		}
	}
	return
}

type FullHeader struct {
	Header
	Version uint8
	Flags   [3]uint8
}

func (h FullHeader) Mp4BoxFlags() uint32 {
	return uint32(h.Flags[0])<<16 | uint32(h.Flags[1])<<8 | uint32(h.Flags[2])
}

func (h *FullHeader) Mp4BoxSetFlags(flags uint32) {
	h.Flags[0] = uint8((flags >> 16) & 0xff)
	h.Flags[1] = uint8((flags >> 8) & 0xff)
	h.Flags[2] = uint8(flags & 0xff)
}

func (h FullHeader) headerSize() uint32 {
	return h.Header.HeaderSize() + 4
}

func (h *FullHeader) ReadHeader(r io.Reader, header *Header) (err error) {
	if header == nil {
		if err = h.Header.ReadHeader(r, header); err != nil {
			return
		}
	} else {
		h.Header = *header
	}
	if err = binary.Read(r, binary.BigEndian, &h.Version); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &h.Flags); err != nil {
		return
	}
	return
}

func (h *FullHeader) WriteHeader(w io.Writer) (err error) {
	if err = h.Header.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, h.Version); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, h.Flags); err != nil {
		return
	}
	return
}
