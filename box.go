package mp4

var BoxRegistry = make(map[BoxType]func() Box)
var UUIDBoxRegistry = make(map[UserType]func() Box)

func NewBox(boxType BoxType) (box Box) {
	if boxFn := BoxRegistry[boxType]; boxFn != nil {
		box = boxFn()
	} else {
		box = &UnknownBox{}
	}
	return
}

func NewUUIDBox(userType UserType) (box Box) {
	if boxFn := UUIDBoxRegistry[userType]; boxFn != nil {
		box = boxFn()
	} else {
		box = &UnknownBox{}
	}
	return
}
