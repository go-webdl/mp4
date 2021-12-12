package mp4

import (
	"io"

	"github.com/go-webdl/media-codec/dovi"
)

type DOVIConfigurationBox struct {
	Header
	NullContainer
	DOVIConfig dovi.DOVIDecoderConfigurationRecord
}

var _ Box = (*DOVIConfigurationBox)(nil)

func init() {
	BoxRegistry[DvcCBoxType] = func() Box { return &DOVIConfigurationBox{} }
	BoxRegistry[DvvCBoxType] = func() Box { return &DOVIConfigurationBox{} }
	BoxRegistry[DvwCBoxType] = func() Box { return &DOVIConfigurationBox{} }
}

func (b DOVIConfigurationBox) Mp4BoxType() BoxType {
	return b.Header.Mp4BoxType()
}

func (b *DOVIConfigurationBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += b.DOVIConfig.RecordSize()
	return b.Size
}

func (b *DOVIConfigurationBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = b.DOVIConfig.RecordRead(r); err != nil {
		return
	}
	return
}

func (b *DOVIConfigurationBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = b.DOVIConfig.RecordWrite(w); err != nil {
		return
	}
	return
}
