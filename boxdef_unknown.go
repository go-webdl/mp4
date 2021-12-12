package mp4

import (
	"io"
)

type UnknownBox struct {
	Header
	NullContainer
	Data []byte
}

var _ Box = (*UnknownBox)(nil)

func (b UnknownBox) Mp4BoxType() BoxType {
	return b.Type
}

func (b *UnknownBox) Mp4BoxUpdate() uint32 {
	return b.Size
}

func (b *UnknownBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	b.Data = make([]byte, b.Size-b.HeaderSize())
	readSize, err := io.ReadFull(r, b.Data)
	if err != nil {
		return
	}
	if readSize < len(b.Data) {
		err = io.EOF
	}
	return
}

func (b *UnknownBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if _, err = w.Write(b.Data); err != nil {
		return
	}
	return
}
