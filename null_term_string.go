package mp4

import (
	"fmt"
	"io"
)

type NullTerminatedString string

func (s NullTerminatedString) Size() uint32 {
	return uint32(len(s)) + 1
}

func (s *NullTerminatedString) ReadOfSize(r io.Reader, size uint32) (err error) {
	if size < 1 {
		err = fmt.Errorf("null-terminated cannot have size of 0: %w", ErrInvalidFormat)
		return
	}
	b := make([]byte, size)
	if _, err = io.ReadFull(r, b); err != nil {
		return
	}
	if b[size-1] != 0 {
		err = fmt.Errorf("string not null-terminated: %w", ErrInvalidFormat)
		return
	}
	*s = NullTerminatedString(b[:size-1])
	return
}

func (s NullTerminatedString) Write(w io.Writer) (err error) {
	if _, err = w.Write([]byte(s)); err != nil {
		return
	}
	if _, err = w.Write([]byte{0}); err != nil {
		return
	}
	return
}
