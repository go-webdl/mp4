package mp4

import (
	"io"
)

// 8.5.1 Sample Table Box

// Box Type: ‘stbl’
// Container: Media Information Box (‘minf’)
// Mandatory: Yes
// Quantity: Exactly one

// The sample table contains all the time and data indexing of the media samples
// in a track. Using the tables here, it is possible to locate samples in time,
// determine their type (e.g. I‐frame or not), and determine their size,
// container, and offset into that container.
//
// If the track that contains the Sample Table Box references no data, then the
// Sample Table Box does not need to contain any sub‐boxes (this is not a very
// useful media track).
//
// If the track that the Sample Table Box is contained in does reference data,
// then the following sub‐boxes are required: Sample Description, Sample Size,
// Sample To Chunk, and Chunk Offset. Further, the Sample Description Box shall
// contain at least one entry. A Sample Description Box is required because it
// contains the data reference index field which indicates which Data Reference
// Box to use to retrieve the media samples. Without the Sample Description, it
// is not possible to determine where the media samples are stored. The Sync
// Sample Box is optional. If the Sync Sample Box is not present, all samples
// are sync samples.
//
// A.7 provides a narrative description of random access using the structures
// defined in the Sample Table Box.
type SampleTableBox struct {
	Header
	Container
}

var _ Box = (*SampleTableBox)(nil)

func init() {
	BoxRegistry[StblBoxType] = func() Box { return &SampleTableBox{} }
}

func (b SampleTableBox) Mp4BoxType() BoxType {
	return StblBoxType
}

func (b *SampleTableBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += b.Mp4BoxUpdateChildren()
	return b.Size
}

func (b *SampleTableBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.Mp4BoxReadChildren(r, b.Size-b.HeaderSize()); err != nil {
		return
	}
	return
}

func (b *SampleTableBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.Mp4BoxWriteChildren(w); err != nil {
		return
	}
	return
}
