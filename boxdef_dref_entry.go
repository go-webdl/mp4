package mp4

import (
	"bytes"
	"fmt"
	"io"
)

// 8.7.2 Data Reference Box

// Box Types: ‘url ‘, ‘urn ‘
// Container: Data Information Box (‘dref’)
// Mandatory: Yes (at least one of ‘url ‘ or ‘urn ‘ shall be present)
// Quantity: One or more

// The data reference object contains a table of data references (normally URLs)
// that declare the location(s) of the media data used within the presentation.
// The data reference index in the sample description ties entries in this table
// to the samples in the track. A track may be split over several sources in
// this way.
//
// If the flag is set indicating that the data is in the same file as this box,
// then no string (not even an empty one) shall be supplied in the entry field.
//
// The entry_count in the DataReferenceBox shall be 1 or greater; each
// DataEntryBox within the DataReferenceBox shall be either a DataEntryUrnBox or
// a DataEntryUrlBox.
//
// > NOTE Though the count is 32 bits, the number of items is usually much
// fewer, and is restricted by the fact that the reference index in the sample
// table is only 16 bits
//
// When a file that has data entries with the flag set indicating that the media
// data is in the same file, is split into segments for transport, the value of
// this flag does not change, as the file is (logically) reassembled after the
// transport operation.
type DataEntryBox struct {
	FullHeader
	NullContainer

	// Name is a URN, and is required in a URN entry.
	Name NullTerminatedString

	// Location is a URL, and is required in a URL entry and optional in a URN
	// entry, where it gives a location to find the resource with the given
	// name. Each is a null‐terminated string using UTF‐8 characters. If the
	// self‐contained flag is set, the URL form is used and no string is
	// present; the box terminates with the entry‐flags field. The URL type
	// should be of a service that delivers a file (e.g. URLs of type file,
	// http, ftp etc.), and which services ideally also permit random access.
	// Relative URLs are permissible and are relative to the file containing the
	// Movie Box that contains this data reference.
	Location NullTerminatedString
}

var _ Box = (*DataEntryBox)(nil)

func init() {
	BoxRegistry[UrnBoxType] = func() Box { return &DataEntryBox{} }
	BoxRegistry[UrlBoxType] = func() Box { return &DataEntryBox{} }
}

func (b DataEntryBox) Mp4BoxType() BoxType {
	if b.Name != "" {
		return UrnBoxType
	}
	return UrlBoxType
}

func (b *DataEntryBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	if b.Mp4BoxFlags()&FLAG_DREF_SAME_FILE == 0 {
		if b.Type == UrnBoxType {
			b.Size += b.Name.Size() // string name;
		}
		b.Size += b.Location.Size() // string location;
	}
	return b.Size
}

func (b *DataEntryBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if b.Mp4BoxFlags()&FLAG_DREF_SAME_FILE == 0 {
		buf := make([]byte, b.Size-b.headerSize())
		if _, err = io.ReadFull(r, buf); err != nil {
			return
		}
		parts := bytes.Split(buf, []byte{0})
		if len(parts) < 2 {
			err = fmt.Errorf("dref entry missing string data: %w", ErrInvalidFormat)
			return
		}
		if len(parts) > 3 || len(parts[2]) > 0 {
			err = fmt.Errorf("dref entry has too many string data: %w", ErrInvalidFormat)
			return
		}
		if len(parts) == 2 {
			b.Type = UrlBoxType
			b.Location = NullTerminatedString(parts[0])
		} else {
			b.Type = UrnBoxType
			b.Name = NullTerminatedString(parts[0])
			b.Location = NullTerminatedString(parts[1])
		}
	}
	return
}

func (b *DataEntryBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if b.Mp4BoxFlags()&FLAG_DREF_SAME_FILE == 0 {
		if b.Type == UrnBoxType {
			if err = b.Name.Write(w); err != nil {
				return
			}
		}
		if err = b.Location.Write(w); err != nil {
			return
		}
	}
	return
}
