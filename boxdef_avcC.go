package mp4

import (
	"io"

	"github.com/go-webdl/media-codec/avc"
)

type AVCConfigurationBox struct {
	Header
	NullContainer
	AVCConfig avc.AVCDecoderConfigurationRecord
}

var _ Box = (*AVCConfigurationBox)(nil)

func init() {
	BoxRegistry[AvcCBoxType] = func() Box { return &AVCConfigurationBox{} }
}

func (b AVCConfigurationBox) Mp4BoxType() BoxType {
	return AvcCBoxType
}

func (b *AVCConfigurationBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += b.AVCConfig.RecordSize()
	return b.Size
}

func (b *AVCConfigurationBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.AVCConfig.RecordRead(r); err != nil {
		return
	}
	return
}

func (b *AVCConfigurationBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.AVCConfig.RecordWrite(w); err != nil {
		return
	}
	return
}
