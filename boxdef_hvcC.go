package mp4

import (
	"io"

	"github.com/go-webdl/media-codec/hevc"
)

type HEVCConfigurationBox struct {
	Header
	NullContainer
	HEVCConfig hevc.HEVCDecoderConfigurationRecord
}

var _ Box = (*HEVCConfigurationBox)(nil)

func init() {
	BoxRegistry[HvcCBoxType] = func() Box { return &HEVCConfigurationBox{} }
}

func (b HEVCConfigurationBox) Mp4BoxType() BoxType {
	return HvcCBoxType
}

func (b *HEVCConfigurationBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += b.HEVCConfig.RecordSize()
	return b.Size
}

func (b *HEVCConfigurationBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.HEVCConfig.RecordRead(r); err != nil {
		return
	}
	return
}

func (b *HEVCConfigurationBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.HEVCConfig.RecordWrite(w); err != nil {
		return
	}
	return
}
