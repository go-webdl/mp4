package mp4

import (
	"encoding/binary"
	"io"
)

// 8.5.2 Sample Description Box
type BitRateBox struct {
	Header
	NullContainer
	BufferSizeDB uint32
	MaxBitrate   uint32
	AvgBitrate   uint32
}

var _ Box = (*BitRateBox)(nil)

func init() {
	BoxRegistry[BtrtBoxType] = func() Box { return &BitRateBox{} }
}

func (b BitRateBox) Mp4BoxType() BoxType {
	return BtrtBoxType
}

func (b *BitRateBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += 4 // unsigned int(32) bufferSizeDB;
	b.Size += 4 // unsigned int(32) maxBitrate;
	b.Size += 4 // unsigned int(32) avgBitrate;
	return b.Size
}

func (b *BitRateBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.BufferSizeDB); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.MaxBitrate); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.AvgBitrate); err != nil {
		return
	}
	return
}

func (b *BitRateBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.BufferSizeDB); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.MaxBitrate); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.AvgBitrate); err != nil {
		return
	}
	return
}
