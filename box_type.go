package mp4

import (
	"io"
)

type Box interface {
	BoxContainer

	// Basic methods
	Mp4BoxSize() uint32
	Mp4BoxType() BoxType
	Mp4BoxSetType(boxType BoxType)
	Mp4BoxUserType() UserType
	Mp4BoxSetUserType(userType UserType)

	// I/O methods
	Mp4BoxUpdate() uint32
	Mp4BoxRead(r io.Reader, header *Header) (err error)
	Mp4BoxWrite(w io.Writer) (err error)
}

type BoxContainer interface {
	Mp4BoxIsContainer() bool
	Mp4BoxUpdateChildren() uint32
	Mp4BoxReadChildren(r io.Reader, size uint32) (err error)
	Mp4BoxWriteChildren(w io.Writer) (err error)

	Mp4BoxAppend(box Box) (err error)
	Mp4BoxReplaceChildren(boxes []Box) (err error)
	Mp4BoxChildren() []Box
	Mp4BoxFirstChild() Box
	Mp4BoxLastChild() Box
	Mp4BoxFindAll(boxType BoxType) []Box
	Mp4BoxFindFirst(boxType BoxType) Box
	Mp4BoxFindLast(boxType BoxType) Box
	Mp4BoxRecursiveFindAll(boxType BoxType) []Box
	Mp4BoxRecursiveFindFirst(boxType BoxType) Box
}
