package mp4

import (
	"encoding/binary"
	"io"

	"github.com/google/uuid"
)

// 8.2 Track Encryption box

// Box Type: `tenc’
// Container: Scheme Information Box (‘schi’)
// Mandatory: No (Yes, for protected tracks)
// Quantity: Zero or one

// The Track Encryption Box contains default values for the isProtected flag,
// Per_Sample_IV_Size, and KID for the entire track. In the case where
// pattern-based encryption is in effect, it supplies the pattern and when
// Constant IVs are in use, it supplies the Constant IV. These values are used
// as the encryption parameters for the samples in this track unless over-ridden
// by the sample group description associated with a group of samples. For files
// with only one key per track, this box allows the basic encryption parameters
// to be specified once per track instead of being repeated per sample.
//
// If both the value of default_isProtected is 1 and default_Per_Sample_IV_Size
// is 0, then the default_constant_IV_size for all samples that use these
// settings SHALL be present. A Constant IV SHALL NOT be used with counter-mode
// encryption. A sample group description may supply keys or keys and Constant
// IVs for sample groups that override these default values for those samples
// mapped to the group.
//
// > NOTE The version field of the Track Encryption Box is set to a value
// greater than zero when the pattern encryption defined in 9.6 is used and to
// zero otherwise.
type ProtectionSystemSpecificHeaderBox struct {
	FullHeader
	NullContainer

	// specifies a UUID that uniquely identifies the content protection system
	// that this header belongs to.
	SystemID uuid.UUID

	// identifies a key identifier that the Data field applies to. If not set,
	// then the Data array SHALL apply to all KIDs in the movie or movie
	// fragment containing this box.
	KIDList [][16]byte

	// holds the content protection system specific data.
	Data []byte
}

var _ Box = (*ProtectionSystemSpecificHeaderBox)(nil)

func init() {
	BoxRegistry[PsshBoxType] = func() Box { return &ProtectionSystemSpecificHeaderBox{} }
}

func (b ProtectionSystemSpecificHeaderBox) Mp4BoxType() BoxType {
	return PsshBoxType
}

func (b *ProtectionSystemSpecificHeaderBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 16 // unsigned int(8)[16] SystemID;
	if b.Version > 0 {
		b.Size += 4                           // unsigned int(32) KID_count;
		b.Size += 16 * uint32(len(b.KIDList)) // unsigned int(8)[16] KID [KID_count];
	}
	b.Size += 4                   // unsigned int(32) DataSize;
	b.Size += uint32(len(b.Data)) // unsigned int(8)[DataSize] Data;
	return b.Size
}

func (b *ProtectionSystemSpecificHeaderBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.SystemID); err != nil {
		return
	}
	var count uint32
	if b.Version > 0 {
		if err = binary.Read(r, binary.BigEndian, &count); err != nil {
			return
		}
		b.KIDList = make([][16]byte, count)
		if err = binary.Read(r, binary.BigEndian, b.KIDList); err != nil {
			return
		}
	}
	if err = binary.Read(r, binary.BigEndian, &count); err != nil {
		return
	}
	b.Data = make([]byte, count)
	if err = binary.Read(r, binary.BigEndian, b.Data); err != nil {
		return
	}
	return
}

func (b *ProtectionSystemSpecificHeaderBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.SystemID); err != nil {
		return
	}
	if b.Version > 0 {
		if err = binary.Write(w, binary.BigEndian, uint32(len(b.KIDList))); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, b.KIDList); err != nil {
			return
		}
	}
	if err = binary.Write(w, binary.BigEndian, uint32(len(b.Data))); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.Data); err != nil {
		return
	}
	return
}
