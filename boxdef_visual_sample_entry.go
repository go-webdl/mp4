package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// 12.1.3 Visual Sample entry

// Video tracks use VisualSampleEntryBox.
//
// In video tracks, the frame_count field must be 1 unless the specification for
// the media format explicitly documents this template field and permits larger
// values. That specification must document both how the individual frames of
// video are found (their size information) and their timing established. That
// timing might be as simple as dividing the sample duration by the frame count
// to establish the frame duration.
//
// The width and height in the video sample entry document the pixel counts that
// the codec will deliver; this enables the allocation of buffers. Since these
// are counts they do not take into account pixel aspect ratio.
type VisualSampleEntryBox struct {
	SampleEntry

	Width           uint16
	Height          uint16
	HorizResolution uint32
	VertResolution  uint32
	FrameCount      uint16
	CompressorName  string
	Depth           uint16

	Clap *CleanApertureBox
	Pasp *PixelAspectRatioBox
}

var _ Box = (*VisualSampleEntryBox)(nil)

func init() {
	BoxRegistry[Avc1BoxType] = func() Box { return &VisualSampleEntryBox{} }
	BoxRegistry[Avc2BoxType] = func() Box { return &VisualSampleEntryBox{} }
	BoxRegistry[Avc3BoxType] = func() Box { return &VisualSampleEntryBox{} }
	BoxRegistry[Avc4BoxType] = func() Box { return &VisualSampleEntryBox{} }
	BoxRegistry[Dva1BoxType] = func() Box { return &VisualSampleEntryBox{} }
	BoxRegistry[DvavBoxType] = func() Box { return &VisualSampleEntryBox{} }
	BoxRegistry[Dvh1BoxType] = func() Box { return &VisualSampleEntryBox{} }
	BoxRegistry[DvheBoxType] = func() Box { return &VisualSampleEntryBox{} }
	BoxRegistry[Hev1BoxType] = func() Box { return &VisualSampleEntryBox{} }
	BoxRegistry[Hvc1BoxType] = func() Box { return &VisualSampleEntryBox{} }
}

func (b *VisualSampleEntryBox) VisualSampleEntrySize() (size uint32) {
	size = b.SampleEntrySize()
	size += 2     // unsigned int(16) pre_defined = 0;
	size += 2     // const unsigned int(16) reserved = 0;
	size += 4 * 3 // unsigned int(32)[3] pre_defined = 0;
	size += 2     // unsigned int(16) width;
	size += 2     // unsigned int(16) height;
	size += 4     // template unsigned int(32) horizresolution = 0x00480000; // 72 dpi
	size += 4     // template unsigned int(32) vertresolution = 0x00480000; // 72 dpi
	size += 4     // const unsigned int(32) reserved = 0;
	size += 2     // template unsigned int(16) frame_count = 1;
	size += 32    // string[32] compressorname;
	size += 2     // template unsigned int(16) depth = 0x0018;
	size += 2     // int(16) pre_defined = -1;
	if b.Clap != nil {
		size += b.Clap.Mp4BoxUpdate() // CleanApertureBox clap; // optional
	}
	if b.Pasp != nil {
		size += b.Pasp.Mp4BoxUpdate() // PixelAspectRatioBox pasp; // optional
	}
	return
}

func (b *VisualSampleEntryBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.VisualSampleEntrySize()
	b.Size += b.Mp4BoxUpdateChildren()
	return b.Size
}

func (b *VisualSampleEntryBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.SampleEntry.Mp4BoxRead(r, header); err != nil {
		return
	}
	var reserved [4]uint32
	if err = binary.Read(r, binary.BigEndian, &reserved); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.Width); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.Height); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.HorizResolution); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.VertResolution); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, reserved[:1]); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.FrameCount); err != nil {
		return
	}
	var compressorname [32]byte
	if err = binary.Read(r, binary.BigEndian, &compressorname); err != nil {
		return
	}
	if compressorname[0] > 31 {
		err = fmt.Errorf("visual sample entry got compressor name with length field exceeds 31: %w", ErrInvalidFormat)
		return
	} else if compressorname[0] > 0 {
		b.CompressorName = string(compressorname[1 : compressorname[0]+1])
	}
	if err = binary.Read(r, binary.BigEndian, &b.Depth); err != nil {
		return
	}
	var reserved2 int16
	if err = binary.Read(r, binary.BigEndian, &reserved2); err != nil {
		return
	}
	if err = b.Mp4BoxReadChildren(r, b.Size-b.VisualSampleEntrySize()); err != nil {
		return
	}
	return
}

func (b *VisualSampleEntryBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.SampleEntry.Mp4BoxWrite(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, [4]uint32{}); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.Width); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.Height); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.HorizResolution); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.VertResolution); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint32(0)); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.FrameCount); err != nil {
		return
	}
	var compressorname [32]byte
	if len(b.CompressorName) > 31 {
		err = fmt.Errorf("visual sample entry got compressor name length exceeds 31: %w", ErrInvalidFormat)
		return
	} else if len(b.CompressorName) > 0 {
		compressorname[0] = byte(len(b.CompressorName))
		copy(compressorname[1:32], []byte(b.CompressorName)[:])
		if err = binary.Write(w, binary.BigEndian, compressorname); err != nil {
			return
		}
	}
	if err = binary.Write(w, binary.BigEndian, b.Depth); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, int16(-1)); err != nil {
		return
	}
	if err = b.Mp4BoxWriteChildren(w); err != nil {
		return
	}
	return
}
