package mp4

import (
	"io"
)

// 8.12.1 Protection Scheme Information Box

// Box Types: ‘sinf’
// Container: Protected Sample Entry, or Item Protection Box (‘ipro’)
// Mandatory: Yes
// Quantity: One or More

// The Protection Scheme Information Box contains all the information required
// both to understand the encryption transform applied and its parameters, and
// also to find other information such as the kind and location of the key
// management system. It also documents the original (unencrypted) format of the
// media. The Protection Scheme Information Box is a container Box. It is
// mandatory in a sample entry that uses a code indicating a protected stream.
//
// When used in a protected sample entry, this box must contain the original
// format box to document the original format. At least one of the following
// signalling methods must be used to identify the protection applied:
//
//     a) MPEG‐4 systems with IPMP: no other boxes, when IPMP descriptors in
//        MPEG‐4 systems streams are used;
//
//     b) Scheme signalling: a SchemeTypeBox and SchemeInformationBox, when
//        these are used (either both must occur, or neither).
//
// At least one protection scheme information box must occur in a protected
// sample entry. When more than one occurs, they are equivalent, alternative,
// descriptions of the same protection. Readers should choose one to process.
type ProtectionSchemeInfoBox struct {
	Header
	Container
}

var _ Box = (*ProtectionSchemeInfoBox)(nil)

func init() {
	BoxRegistry[SinfBoxType] = func() Box { return &ProtectionSchemeInfoBox{} }
}

func (b ProtectionSchemeInfoBox) Mp4BoxType() BoxType {
	return SinfBoxType
}

func (b *ProtectionSchemeInfoBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += b.Mp4BoxUpdateChildren()
	return b.Size
}

func (b *ProtectionSchemeInfoBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.Mp4BoxReadChildren(r, b.Size-b.HeaderSize()); err != nil {
		return
	}
	return
}

func (b *ProtectionSchemeInfoBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.Mp4BoxWriteChildren(w); err != nil {
		return
	}
	return
}
