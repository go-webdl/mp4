package mp4

import (
	"encoding/binary"
	"io"
)

// 8.7.2 Data Reference Box

// Box Types: ‘dref’
// Container: Data Information Box (‘dinf’)
// Mandatory: Yes
// Quantity: Exactly one

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
type DataReferenceBox struct {
	FullHeader
	Container
}

const (
	FLAG_DREF_SAME_FILE uint32 = 0x000001
)

var _ Box = (*DataReferenceBox)(nil)

func init() {
	BoxRegistry[DrefBoxType] = func() Box { return &DataReferenceBox{} }
}

func (b DataReferenceBox) Mp4BoxType() BoxType {
	return DrefBoxType
}

func (b *DataReferenceBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.headerSize()
	b.Size += 4 // unsigned int(32) entry_count;
	// for (i=1; i <= entry_count; i++) {
	//     DataEntryBox(entry_version, entry_flags) data_entry;
	// }
	b.Size += b.Mp4BoxUpdateChildren()
	return b.Size
}

func (b *DataReferenceBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
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
	return
}

func (b *DataReferenceBox) Mp4BoxWrite(w io.Writer) (err error) {
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
