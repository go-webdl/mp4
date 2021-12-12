package mp4

import (
	"encoding/binary"
	"io"
)

// 4.3 File Type Box

// Box Type: `ftyp’
// Container: File
// Mandatory: Yes
// Quantity: Exactly one (but see below)

// Files written to this version of this specification must contain a file‐type
// box. For compatibility with an earlier version of this specification, files
// may be conformant to this specification and not contain a filetype box. Files
// with no file‐type box should be read as if they contained an FTYP box with
// Major_brand='mp41', minor_version=0, and the single compatible brand 'mp41'.
//
// A media‐file structured to this part of this specification may be compatible
// with more than one detailed specification, and it is therefore not always
// possible to speak of a single ‘type’ or ‘brand’ for the file. This means that
// the utility of the file name extension and Multipurpose Internet Mail
// Extension (MIME) type are somewhat reduced.
//
// This box must be placed as early as possible in the file (e.g. after any
// obligatory signature, but before any significant variable‐size boxes such as
// a Movie Box, Media Data Box, or Free Space). It identifies which
// specification is the ‘best use’ of the file, and a minor version of that
// specification; and also a set of other specifications to which the file
// complies. Readers implementing this format should attempt to read files that
// are marked as compatible with any of the specifications that the reader
// implements. Any incompatible change in a specification should therefore
// register a new ‘brand’ identifier to identify files conformant to the new
// specification.
//
// The minor version is informative only. It does not appear for
// compatible‐brands, and must not be used to determine the conformance of a
// file to a standard. It may allow more precise identification of the major
// specification, for inspection, debugging, or improved decoding.
//
// Files would normally be externally identified (e.g. with a file extension or
// mime type) that identifies the‘best use’ (major brand), or the brand that the
// author believes will provide the greatest compatibility.
type FileTypeBox struct {
	Header
	NullContainer
	MajorBrand       FourCC   // brand identifier
	MinorVersion     uint32   // an informative integer for the minor version of the major brand
	CompatibleBrands []FourCC // a list of brands
}

var _ Box = (*FileTypeBox)(nil)

func init() {
	BoxRegistry[FtypBoxType] = func() Box { return &FileTypeBox{} }
}

func (b FileTypeBox) Mp4BoxType() BoxType {
	return FtypBoxType
}

func (b *FileTypeBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += 4                                   // unsigned int(32) major_brand;
	b.Size += 4                                   // unsigned int(32) minor_brand;
	b.Size += 4 * uint32(len(b.CompatibleBrands)) // unsigned int(32) compatible_brands[];
	return b.Size
}

func (b *FileTypeBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.MajorBrand); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.MinorVersion); err != nil {
		return
	}
	b.CompatibleBrands = make([]FourCC, (b.Size-b.HeaderSize()-4-4)/4)
	if err = binary.Read(r, binary.BigEndian, b.CompatibleBrands); err != nil {
		return
	}
	return
}

func (b *FileTypeBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.MajorBrand); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.MinorVersion); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.CompatibleBrands); err != nil {
		return
	}
	return
}
