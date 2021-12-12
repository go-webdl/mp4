package mp4

import (
	"fmt"
	"io"
)

type Container struct {
	Children        []Box
	TypeChildrenMap map[BoxType][]Box
}

var _ BoxContainer = (*Container)(nil)

func (b *Container) Mp4BoxIsContainer() bool {
	return true
}

func (b *Container) Mp4BoxUpdateChildren() (size uint32) {
	for _, child := range b.Children {
		size += child.Mp4BoxUpdate()
	}
	return
}

func (b *Container) Mp4BoxReadChildren(r io.Reader, size uint32) (err error) {
	remainingSize := int64(size)
	for remainingSize > 0 {
		var child Box
		if child, err = ReadBox(r); err != nil {
			return
		}
		remainingSize -= int64(child.Mp4BoxSize())
		if remainingSize < 0 {
			err = fmt.Errorf("child box %s exceeds parent boundary: %w", child.Mp4BoxType(), ErrInvalidFormat)
			return
		}
		b.Mp4BoxAppend(child)
	}
	return
}

func (b *Container) Mp4BoxWriteChildren(w io.Writer) (err error) {
	for _, child := range b.Children {
		if err = child.Mp4BoxWrite(w); err != nil {
			return
		}
	}
	return
}

func (b *Container) Mp4BoxAppend(box Box) (err error) {
	b.Children = append(b.Children, box)
	if b.TypeChildrenMap == nil {
		b.TypeChildrenMap = make(map[BoxType][]Box)
	}
	b.TypeChildrenMap[box.Mp4BoxType()] = append(b.TypeChildrenMap[box.Mp4BoxType()], box)
	return
}

func (b *Container) Mp4BoxReplaceChildren(boxes []Box) (err error) {
	b.Children = boxes
	b.TypeChildrenMap = make(map[BoxType][]Box)
	for _, box := range boxes {
		b.TypeChildrenMap[box.Mp4BoxType()] = append(b.TypeChildrenMap[box.Mp4BoxType()], box)
	}
	return
}

func (b *Container) Mp4BoxChildren() []Box {
	if len(b.Children) == 0 {
		return nil
	}
	return b.Children[:]
}

func (b *Container) Mp4BoxFirstChild() Box {
	if len(b.Children) == 0 {
		return nil
	}
	return b.Children[0]
}

func (b *Container) Mp4BoxLastChild() Box {
	if len(b.Children) == 0 {
		return nil
	}
	return b.Children[len(b.Children)]
}

func (b *Container) Mp4BoxFindAll(boxType BoxType) []Box {
	if b.TypeChildrenMap == nil {
		return nil
	}
	boxes := b.TypeChildrenMap[boxType]
	if len(boxes) == 0 {
		return nil
	}
	return boxes[:]
}

func (b *Container) Mp4BoxFindFirst(boxType BoxType) Box {
	if b.TypeChildrenMap == nil {
		return nil
	}
	boxes := b.TypeChildrenMap[boxType]
	if len(boxes) == 0 {
		return nil
	}
	return boxes[0]
}

func (b *Container) Mp4BoxFindLast(boxType BoxType) Box {
	if b.TypeChildrenMap == nil {
		return nil
	}
	boxes := b.TypeChildrenMap[boxType]
	if len(boxes) == 0 {
		return nil
	}
	return boxes[len(boxes)]
}

func (b *Container) Mp4BoxRecursiveFindAll(boxType BoxType) []Box {
	boxes := b.Mp4BoxFindAll(boxType)
	for _, child := range b.Children {
		boxes = append(boxes, child.Mp4BoxRecursiveFindAll(boxType)...)
	}
	return boxes
}

func (b *Container) Mp4BoxRecursiveFindFirst(boxType BoxType) (box Box) {
	box = b.Mp4BoxFindFirst(boxType)
	if box != nil {
		return
	}
	for _, child := range b.Children {
		box = child.Mp4BoxRecursiveFindFirst(boxType)
		if box != nil {
			return
		}
	}
	return
}

type NullContainer struct{}

var _ BoxContainer = (*NullContainer)(nil)

func (b *NullContainer) Mp4BoxIsContainer() bool {
	return false
}

func (b *NullContainer) Mp4BoxUpdateChildren() (size uint32) {
	return
}

func (b *NullContainer) Mp4BoxReadChildren(r io.Reader, size uint32) (err error) {
	return
}

func (b *NullContainer) Mp4BoxWriteChildren(w io.Writer) (err error) {
	return
}

func (b *NullContainer) Mp4BoxAppend(box Box) error {
	return ErrChildBoxNotSupported
}

func (b *NullContainer) Mp4BoxReplaceChildren(boxes []Box) error {
	return ErrChildBoxNotSupported
}

func (b *NullContainer) Mp4BoxChildren() []Box {
	return nil
}

func (b *NullContainer) Mp4BoxFirstChild() Box {
	return nil
}

func (b *NullContainer) Mp4BoxLastChild() Box {
	return nil
}

func (b *NullContainer) Mp4BoxFindAll(boxType BoxType) []Box {
	return nil
}

func (b *NullContainer) Mp4BoxFindFirst(boxType BoxType) Box {
	return nil
}

func (b *NullContainer) Mp4BoxFindLast(boxType BoxType) Box {
	return nil
}

func (b *NullContainer) Mp4BoxRecursiveFindAll(boxType BoxType) []Box {
	return nil
}

func (b *NullContainer) Mp4BoxRecursiveFindFirst(boxType BoxType) (box Box) {
	return nil
}
