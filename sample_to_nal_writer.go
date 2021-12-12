package mp4

import "io"

func NewSampleToNALWriter(w io.Writer, lengthSizeMinusOne uint32) io.Writer {
	return &sampleToNALWriterT{w: w, lengthSizeMinusOne: lengthSizeMinusOne}
}

type sampleToNALWriterT struct {
	w                  io.Writer
	lengthSizeMinusOne uint32
	nalLength          uint32
	nalBytesRead       uint32
}

func (w *sampleToNALWriterT) Write(p []byte) (n int, err error) {
	if p == nil {
		return
	}
	for n < len(p) {
		for w.nalBytesRead < (w.lengthSizeMinusOne+1) && n < len(p) {
			w.nalLength |= uint32(p[n]) << (8 * (w.lengthSizeMinusOne - w.nalBytesRead))
			n += 1
			w.nalBytesRead += 1
			if w.nalBytesRead == w.lengthSizeMinusOne+1 {
				// adjust nalLength to include the length itself
				w.nalLength += w.lengthSizeMinusOne + 1
				// write NAL start code
				if _, err = w.w.Write([]byte{0, 0, 0, 1}); err != nil {
					return
				}
			}
		}
		length := int(w.nalLength - w.nalBytesRead)
		if len(p[n:]) < length {
			length = len(p[n:])
		}
		if length > 0 {
			if _, err = w.w.Write(p[n : n+length]); err != nil {
				return
			}
			w.nalBytesRead += uint32(length)
			n += length
		}
		if w.nalLength == w.nalBytesRead {
			// finish current NALU
			w.nalLength = 0
			w.nalBytesRead = 0
		}
	}
	return
}
