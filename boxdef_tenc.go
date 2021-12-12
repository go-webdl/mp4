package mp4

import (
	"encoding/binary"
	"io"
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
type TrackEncryptionBox struct {
	FullHeader
	NullContainer

	// specifies the count of the encrypted Blocks in the protection pattern,
	// where each Block is of size 16-bytes. See 9.1 for further details.
	DefaultCryptByteBlock uint8

	// specifies the count of the unencrypted Blocks in the protection pattern.
	// See the skip_byte_block field in 9.1 for further details.
	DefaultSkipByteBlock uint8

	// is the protection flag which indicates the default protection state of
	// the samples in the track. See the isProtected field in 9.1 for further
	// details.
	DefaultIsProtected uint8

	// is the default Initialization Vector size in bytes. See the
	// Per_Sample_IV_Size field in 9.1 for further details.
	DefaultPerSampleIVSize uint8

	// is the default key identifier used for samples in this track. See the KID
	// field in 9.1 for further details.
	DefaultKID [16]byte

	// is the size of a possible default Initialization Vector for all samples.
	DefaultConstantIVSize uint8

	// if present, is the default Initialization Vector for all samples. See the
	// constant_IV field in 9.1 for further details.
	DefaultConstantIV []byte
}

var _ Box = (*TrackEncryptionBox)(nil)

func init() {
	BoxRegistry[TencBoxType] = func() Box { return &TrackEncryptionBox{} }
}

func (b TrackEncryptionBox) Mp4BoxType() BoxType {
	return TencBoxType
}

func (b *TrackEncryptionBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	// unsigned int(8) reserved = 0;
	// if (version==0) {
	// unsigned int(8) reserved = 0;
	// }
	// else { // version is 1 or greater
	// unsigned int(4) default_crypt_byte_block;
	// unsigned int(4) default_skip_byte_block;
	// }
	// unsigned int(8) default_isProtected;
	// unsigned int(8) default_Per_Sample_IV_Size;
	b.Size += 4
	b.Size += 16 // unsigned int(8)[16] default_KID;
	if b.DefaultIsProtected == 1 && b.DefaultPerSampleIVSize == 0 {
		b.Size += 1                                // unsigned int(8) default_constant_IV_size;
		b.Size += uint32(len(b.DefaultConstantIV)) // unsigned int(8)[default_constant_IV_size] default_constant_IV;
	}
	return b.Size
}

func (b *TrackEncryptionBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	var tmp uint32
	if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
		return
	}
	if b.Version != 0 {
		b.DefaultCryptByteBlock = uint8(tmp >> 24)
		b.DefaultSkipByteBlock = uint8(tmp >> 16 & 0xff)
	}
	b.DefaultIsProtected = uint8(tmp >> 8 & 0xff)
	b.DefaultPerSampleIVSize = uint8(tmp & 0xff)
	if err = binary.Read(r, binary.BigEndian, &b.DefaultKID); err != nil {
		return
	}
	if b.DefaultIsProtected == 1 && b.DefaultPerSampleIVSize == 0 {
		if err = binary.Read(r, binary.BigEndian, &b.DefaultConstantIVSize); err != nil {
			return
		}
		b.DefaultConstantIV = make([]byte, b.DefaultConstantIVSize)
		if err = binary.Read(r, binary.BigEndian, b.DefaultConstantIV); err != nil {
			return
		}
	}
	return
}

func (b *TrackEncryptionBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	var tmp uint32
	if b.Version != 0 {
		tmp |= uint32(b.DefaultCryptByteBlock)<<24 | uint32(b.DefaultSkipByteBlock)<<16
	}
	tmp |= uint32(b.DefaultIsProtected)<<8 | uint32(b.DefaultPerSampleIVSize)
	if err = binary.Write(w, binary.BigEndian, tmp); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.DefaultKID); err != nil {
		return
	}
	if b.DefaultIsProtected == 1 && b.DefaultPerSampleIVSize == 0 {
		if err = binary.Write(w, binary.BigEndian, b.DefaultConstantIVSize); err != nil {
			return
		}
		b.DefaultConstantIV = make([]byte, b.DefaultConstantIVSize)
		if err = binary.Write(w, binary.BigEndian, b.DefaultConstantIV); err != nil {
			return
		}
	}
	return
}
