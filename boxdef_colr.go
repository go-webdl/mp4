package mp4

import (
	"encoding/binary"
	"io"
)

// 12.1.5 Colour information
type ColourInformationBox struct {
	Header
	NullContainer

	// An indication of the type of colour information supplied. For colour_type
	// ‘nclx’: these fields are exactly the four bytes defined for
	// PTM_COLOR_INFO( ) in A.7.2 of ISO/IEC 29199‐2 but note that the full
	// range flag is here in a different bit position
	ColourType FourCC

	ColourPrimaries         uint16
	TransferCharacteristics uint16
	MatrixCoefficients      uint16
	FullRange               bool

	ICCProfile []byte

	UnknownData []byte
}

var _ Box = (*ColourInformationBox)(nil)

func init() {
	BoxRegistry[ColrBoxType] = func() Box { return &ColourInformationBox{} }
}

func (b ColourInformationBox) Mp4BoxType() BoxType {
	return ColrBoxType
}

func (b *ColourInformationBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += 4 // unsigned int(32) colour_type;
	if b.ColourType == NclxFourCC {
		b.Size += 2 // unsigned int(16) colour_primaries;
		b.Size += 2 // unsigned int(16) transfer_characteristics;
		b.Size += 2 // unsigned int(16) matrix_coefficients;
		b.Size += 1 // unsigned int(1) full_range_flag;
		// unsigned int(7) reserved = 0;
	} else if b.ColourType == NclcFourCC {
		b.Size += 2 // unsigned int(16) colour_primaries;
		b.Size += 2 // unsigned int(16) transfer_characteristics;
		b.Size += 2 // unsigned int(16) matrix_coefficients;
	} else if b.ColourType == RiccFourCC || b.ColourType == ProfFourCC {
		b.Size += uint32(len(b.ICCProfile))
	}
	return b.Size
}

func (b *ColourInformationBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.ColourType); err != nil {
		return
	}
	if b.ColourType == NclxFourCC || b.ColourType == NclcFourCC {
		if err = binary.Read(r, binary.BigEndian, &b.ColourPrimaries); err != nil {
			return
		}
		if err = binary.Read(r, binary.BigEndian, &b.TransferCharacteristics); err != nil {
			return
		}
		if err = binary.Read(r, binary.BigEndian, &b.MatrixCoefficients); err != nil {
			return
		}
		if b.ColourType == NclxFourCC {
			var tmp uint8
			if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
				return
			}
			b.FullRange = (tmp >> 7) > 0
		}
	} else {
		tmp := make([]byte, b.Size-b.HeaderSize()-4)
		if _, err = io.ReadFull(r, tmp); err != nil {
			return
		}
		if b.ColourType == RiccFourCC || b.ColourType == ProfFourCC {
			b.ICCProfile = tmp
		} else {
			b.UnknownData = tmp
		}
	}
	return
}

func (b *ColourInformationBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.ColourType); err != nil {
		return
	}
	if b.ColourType == NclxFourCC || b.ColourType == NclcFourCC {
		if err = binary.Write(w, binary.BigEndian, b.ColourPrimaries); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, b.TransferCharacteristics); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, b.MatrixCoefficients); err != nil {
			return
		}
		if b.ColourType == NclxFourCC {
			var tmp uint8
			if b.FullRange {
				tmp = 1 << 7
			}
			if err = binary.Write(w, binary.BigEndian, tmp); err != nil {
				return
			}
		}
	} else if b.ColourType == RiccFourCC || b.ColourType == ProfFourCC && len(b.ICCProfile) > 0 {
		if err = binary.Write(w, binary.BigEndian, b.ICCProfile); err != nil {
			return
		}
	} else if len(b.UnknownData) > 0 {
		if err = binary.Write(w, binary.BigEndian, b.UnknownData); err != nil {
			return
		}
	}
	return
}
