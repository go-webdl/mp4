package mp4

import (
	"io"

	"github.com/go-webdl/media-codec/hevc"
)

type DolbyVisionELHEVCConfigurationBox struct {
	Header
	NullContainer
	HEVCConfig hevc.HEVCDecoderConfigurationRecord
}

var _ Box = (*DolbyVisionELHEVCConfigurationBox)(nil)

func init() {
	BoxRegistry[HvcEBoxType] = func() Box { return &DolbyVisionELHEVCConfigurationBox{} }
}

func (b DolbyVisionELHEVCConfigurationBox) Mp4BoxType() BoxType {
	return HvcEBoxType
}

func (b *DolbyVisionELHEVCConfigurationBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += b.HEVCConfig.RecordSize()
	return b.Size
}

func (b *DolbyVisionELHEVCConfigurationBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.HEVCConfig.RecordRead(r); err != nil {
		return
	}
	return
}

func (b *DolbyVisionELHEVCConfigurationBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.HEVCConfig.RecordWrite(w); err != nil {
		return
	}
	return
}
