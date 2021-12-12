package mp4

import "io"

func ReadHeader(r io.Reader) (header *Header, err error) {
	header = &Header{}
	if err = header.ReadHeader(r, nil); err != nil {
		return
	}
	return
}

func ReadBox(r io.Reader) (box Box, err error) {
	var header *Header
	if header, err = ReadHeader(r); err != nil {
		return
	}
	return ReadBoxAfterHeader(r, header)
}

func ReadBoxAfterHeader(r io.Reader, header *Header) (box Box, err error) {
	if header.Type == UuidBoxType {
		box = NewUUIDBox(header.UserType)
	} else {
		box = NewBox(header.Type)
	}
	if err = box.Mp4BoxRead(r, header); err != nil {
		return
	}
	return
}
