package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// 8.5.2 Sample Description Box

// Box Types: ‘stsd’
// Container: Sample Table Box (‘stbl’)
// Mandatory: Yes
// Quantity: Exactly one

// The sample description table gives detailed information about the coding type
// used, and any initialization information needed for that coding.
//
// The information stored in the sample description box after the entry‐count is
// both track‐type specific as documented here, and can also have variants
// within a track type (e.g. different codings may use different specific
// information after some common fields, even within a video track).
//
// Which type of sample entry form is used is determined by the media handler,
// using a suitable form, such as one defined in clause 12, or defined in a
// derived specification, or registration.
//
// Multiple descriptions may be used within a track.
//
// > Note Though the count is 32 bits, the number of items is usually much
// fewer, and is restricted by the fact that the reference index in the sample
// table is only 16 bits
//
// If the ‘format’ field of a SampleEntry is unrecognized, neither the sample
// description itself, nor the associated media samples, shall be decoded.
//
// > Note The definition of sample entries specifies boxes in a particular
// order, and this is usually also followed in derived specifications. For
// maximum compatibility, writers should construct files respecting the order
// both within specifications and as implied by the inheritance, whereas readers
// should be prepared to accept any box order.
//
// All string fields shall be null‐terminated, even if unused. “Optional” means
// there is at least one null byte.
//
// Entries that identify the format by MIME type, such as a
// TextSubtitleSampleEntry, TextMetaDataSampleEntry, or SimpleTextSampleEntry,
// all of which contain a MIME type, may be used to identify the format of
// streams for which a MIME type applies. A MIME type applies if the contents of
// the string in the optional configuration box (without its null termination),
// followed by the contents of a set of samples, starting with a sync sample and
// ending at the sample immediately preceding a sync sample, are concatenated in
// their entirety, and the result meets the decoding requirements for documents
// of that MIME type. Non‐sync samples should be used only if that format
// specifies the behaviour of ‘progressive decoding’, and then the sample times
// indicate when the results of such progressive decoding should be presented
// (according to the media type).
//
// > Note The samples in a track that is all sync samples are therefore each a
// valid document for that MIME type.
//
// In some classes derived from SampleEntry, namespace and schema_location are
// used both to identify the XML document content and to declare “brand” or
// profile compatibility. Multiple namespace identifiers indicate that the track
// conforms to the specification represented by each of the identifiers, some of
// which may identify supersets of the features present. A decoder should be
// able to decode all the namespaces in order to be able to decode and present
// correctly the media associated with this sample entry.
//
// > Note Additionally, namespace identifiers may represent performance
// constraints, such as limits on document size, font size, drawing rate, etc.,
// as well as syntax constraints such as features that are not permitted or
// ignored.
type SampleDescriptionBox struct {
	FullHeader
	Container
}

var _ Box = (*SampleDescriptionBox)(nil)

func init() {
	BoxRegistry[StsdBoxType] = func() Box { return &SampleDescriptionBox{} }
}

func (b SampleDescriptionBox) Mp4BoxType() BoxType {
	return StsdBoxType
}

func (b *SampleDescriptionBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 4 // unsigned int(32) entry_count;
	b.Size += b.Mp4BoxUpdateChildren()
	return b.Size
}

func (b *SampleDescriptionBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	var entryCount uint32
	if err = binary.Read(r, binary.BigEndian, &entryCount); err != nil {
		return
	}
	if err = b.Mp4BoxReadChildren(r, b.Size-b.headerSize()-4); err != nil {
		return
	}
	if len(b.Children) != int(entryCount) {
		err = fmt.Errorf("stsd entry count mismatch: %w", ErrInvalidFormat)
		return
	}
	return
}

func (b *SampleDescriptionBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint32(len(b.Children))); err != nil {
		return
	}
	if err = b.Mp4BoxWriteChildren(w); err != nil {
		return
	}
	return
}
