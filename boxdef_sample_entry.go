package mp4

import (
	"encoding/binary"
	"io"
)

// 8.5.2 Sample Description Box
type SampleEntry struct {
	Header
	Container

	// is an integer that contains the index of the data reference to use to
	// retrieve data associated with samples that use this sample description.
	// Data references are stored in Data Reference Boxes. The index ranges from
	// 1 to the number of data references.
	DataReferenceIndex uint16
}

func (b *SampleEntry) SampleEntrySize() (size uint32) {
	size = b.HeaderSize()
	size += 6 // const unsigned int(8)[6] reserved = 0;
	size += 2 // unsigned int(16) data_reference_index;
	return
}

func (b *SampleEntry) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	var tmp [4]uint16
	if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
		return
	}
	b.DataReferenceIndex = tmp[3]
	return
}

func (b *SampleEntry) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, [4]uint16{0, 0, 0, b.DataReferenceIndex}); err != nil {
		return
	}
	return
}
