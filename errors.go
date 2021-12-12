package mp4

import "errors"

var ErrInvalidFormat = errors.New("mp4 format error")
var ErrChildBoxNotSupported = errors.New("this box cannot have child boxes")
var ErrUnsupportedSerialization = errors.New("serialization not supported")
