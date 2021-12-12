package mp4

import (
	"encoding/binary"
	"io"
)

// 5.3.2 Sample Encryption Box

// Box Type: ‘uuid’
// Container: Track Fragment Box (‘traf’)
// Mandatory: No
// Quantity: Zero or one

// The Sample Encryption box contains the sample specific encryption data. It is
// used when the sample data in the track or fragment is encrypted. The box MUST
// be present for Track Fragment Boxes or Sample Table Boxes that contain or
// refer to sample data for tracks containing encrypted data. It SHOULD be
// omitted for unencrypted content.
//
// For an AlgorithmID of Not Encrypted, no initialization vectors are needed and
// this table SHOULD be omitted.
//
// For an AlgorithmID of AES-CTR, if the IV_size field is 16 then the
// InitializationVector specifies the entire 128 bit IV value used as the
// counter value. If the InitializationVector field is 8, then its value is
// copied to bytes 0 to 7 of the 16 byte block passed to AES ECB and bytes 8 to
// 15 are set to zero. However the initial counter value is specified, bytes 8
// to 15 are used as a simple block counter that is incremented for each block
// of the sample processed and is kept in network byte order.
//
// Regardless of the length specified in the IV_size field, the initialization
// vectors for a given key MUST be unique for each sample in all Tracks. It is
// RECOMMENDED that the initial initialization vector be randomly generated and
// then incremented for each additional protected sample added. This provides
// entropy and ensures that the initialization vectors are unique.
//
// For an AlgorithmID of AES-CBC, initialization vectors must by 16 bytes long
// and MUST be constructed such that the IV for the first sample in a fragment
// is randomly generated and subsequent samples within the same fragment use the
// last block of ciphertext from the previous sample as their IV. Note that the
// IV for each sample is still added to the SampleEncryptionBox (even though it
// can be retrieved from the previous sample) to facilitate random sample
// access.
//
// The sub sample encryption entries SHALL NOT include an entry with a zero
// value in both the BytesOfClearData field and in the BytesOfEncryptedData
// field. Further, it is RECOMMENDED that the sub sample encryption entries be
// as compactly represented as possible. For example, instead two entries with
// {15 clear, 0 encrypted}, {17 clear, 500 encrypted} use one entry of {32
// clear, 500 encrypted}.
//
// See Section 6, Encryption of Track Level Data, for further details on how
// encryption is applied.
type SampleEncryptionBox struct {
	FullHeader
	NullContainer

	// is the identifier of the encryption algorithm used to encrypt the track.
	AlgorithmID PiffAlgorithmID

	// is the size in bytes of the InitializationVector field
	IVSize PiffIVSize

	// is a key identifier that uniquely identifies the key needed to decrypt samples referred to by this sample encryption box.
	KID [16]uint8

	Samples []SampleEncryptionSampleEntry
}

const (
	// If the Override TrackEncryptionBox parameters flag is set, then the
	// SampleEncryptionBox specifies the AlgorithmID, IV_size, and KID
	// parameters. If not present, then the default values from the
	// TrackEncryptionBox SHOULD be used for this fragment.
	FLAG_SENC_OVERRIDE_TRACK_ENCRYPTION_BOX_PARAMS uint32 = 0x01

	// If the UseSubSampleEncryption flag is set, then the track fragment that
	// contains this Sample Encryption Box SHALL use Subsample encryption as
	// described in 9.5. When this flag is set, Subsample mapping data follows
	// each InitializationVector. The Subsample mapping data consists of the
	// number of Subsamples for each sample, followed by an array of values
	// describing the number of bytes of clear data and the number of bytes of
	// encrypted data for each Subsample.
	FLAG_SENC_USE_SUBSAMPLE_ENCRYPTION uint32 = 0x02
)

type PiffAlgorithmID uint32

const (
	PiffNotEncrypted PiffAlgorithmID = 0x00
	PiffAES128CTR    PiffAlgorithmID = 0x01
	PiffAES128CBC    PiffAlgorithmID = 0x02
)

type PiffIVSize uint8

const (
	PiffIVSize64Bit  PiffIVSize = 8
	PiffIVSize128Bit PiffIVSize = 16
)

var _ Box = (*SampleEncryptionBox)(nil)

func init() {
	BoxRegistry[SencBoxType] = func() Box { return &SampleEncryptionBox{} }
	UUIDBoxRegistry[SampleEncryptionBoxUserType] = func() Box { return &SampleEncryptionBox{} }
}

type SampleEncryptionSampleEntry struct {
	// specifies the initialization vector required for decryption of the
	// sample.
	InitializationVector []byte
	Subsamples           []SampleEncryptionSubsampleEntry
}

type SampleEncryptionSubsampleEntry struct {
	// specifies the number of bytes of clear data at the beginning of this sub
	// sample encryption entry. Note that this value may be zero if no clear
	// bytes exist for this entry.
	BytesOfClearData uint16

	// specifies the number of bytes of encrypted data following the clear data.
	// Note that this value may be zero if no encrypted bytes exist for this
	// entry.
	BytesOfProtectedData uint32
}

func (b SampleEncryptionBox) Mp4BoxType() BoxType {
	if b.Type == UuidBoxType || b.Type == SencBoxType {
		return b.Type
	}
	return SencBoxType
}

func (b SampleEncryptionBox) Mp4BoxUserType() UserType {
	if b.Type == UuidBoxType || b.UserType == SampleEncryptionBoxUserType {
		return SampleEncryptionBoxUserType
	}
	return b.Header.Mp4BoxUserType()
}

func (b *SampleEncryptionBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.UserType = b.Mp4BoxUserType()
	b.Size = b.headerSize()
	flags := b.Mp4BoxFlags()
	ivSize := PiffIVSize(8)
	if flags&FLAG_SENC_OVERRIDE_TRACK_ENCRYPTION_BOX_PARAMS > 0 {
		// unsigned int(24) AlgorithmID;
		// unsigned int(8) IV_size;
		// unsigned int(8)[16] KID;
		b.Size += 3 + 1 + 16
		ivSize = b.IVSize
	}
	b.Size += 4                                       // unsigned int(32) sample_count;
	b.Size += uint32(ivSize) * uint32(len(b.Samples)) // unsigned int(Per_Sample_IV_Size*8) InitializationVector;
	if flags&FLAG_SENC_USE_SUBSAMPLE_ENCRYPTION > 0 {
		b.Size += 2 * uint32(len(b.Samples)) // unsigned int(16) subsample_count;
		var subsampleTotal uint32
		for _, sample := range b.Samples {
			subsampleTotal += uint32(len(sample.Subsamples))
		}
		// {
		//     unsigned int(16) BytesOfClearData;
		//     unsigned int(32) BytesOfProtectedData;
		// } [ subsample_count ]
		b.Size += 6 * subsampleTotal
	}
	return b.Size
}

func (b *SampleEncryptionBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	flags := b.Mp4BoxFlags()
	ivSize := PiffIVSize(8)
	if flags&FLAG_SENC_OVERRIDE_TRACK_ENCRYPTION_BOX_PARAMS > 0 {
		var tmp uint32
		if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
			return
		}
		b.AlgorithmID = PiffAlgorithmID(tmp >> 8)
		b.IVSize = PiffIVSize(tmp & 0xFF)
		if err = binary.Read(r, binary.BigEndian, &b.KID); err != nil {
			return
		}
	}
	var sampleCount uint32
	if err = binary.Read(r, binary.BigEndian, &sampleCount); err != nil {
		return
	}
	b.Samples = make([]SampleEncryptionSampleEntry, sampleCount)
	for i := uint32(0); i < sampleCount; i++ {
		b.Samples[i].InitializationVector = make([]byte, ivSize)
		if _, err = io.ReadFull(r, b.Samples[i].InitializationVector); err != nil {
			return
		}
		if b.Mp4BoxFlags()&FLAG_SENC_USE_SUBSAMPLE_ENCRYPTION > 0 {
			var subsampleCount uint16
			if err = binary.Read(r, binary.BigEndian, &subsampleCount); err != nil {
				return
			}
			b.Samples[i].Subsamples = make([]SampleEncryptionSubsampleEntry, subsampleCount)
			if err = binary.Read(r, binary.BigEndian, b.Samples[i].Subsamples); err != nil {
				return
			}
		}
	}
	return
}

func (b *SampleEncryptionBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	flags := b.Mp4BoxFlags()
	if flags&FLAG_SENC_OVERRIDE_TRACK_ENCRYPTION_BOX_PARAMS > 0 {
		if err = binary.Write(w, binary.BigEndian, uint32(b.AlgorithmID)<<8|uint32(b.IVSize)); err != nil {
			return
		}
		if err = binary.Write(w, binary.BigEndian, b.KID); err != nil {
			return
		}
	}
	if err = binary.Write(w, binary.BigEndian, uint32(len(b.Samples))); err != nil {
		return
	}
	for _, sample := range b.Samples {
		if err = binary.Write(w, binary.BigEndian, sample.InitializationVector); err != nil {
			return
		}
		if b.Mp4BoxFlags()&FLAG_SENC_USE_SUBSAMPLE_ENCRYPTION > 0 {
			if err = binary.Write(w, binary.BigEndian, uint16(len(sample.Subsamples))); err != nil {
				return
			}
			for _, subsample := range sample.Subsamples {
				if err = binary.Write(w, binary.BigEndian, subsample); err != nil {
					return
				}
			}
		}
	}
	return
}
