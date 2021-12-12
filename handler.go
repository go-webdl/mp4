package mp4

type Handler struct {
	HandleMp4Box func(box Box) error
}
