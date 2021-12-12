package mp4

import (
	"encoding/binary"
	"io"
)

// 12.1.4 Pixel Aspect Ratio and Clean Aperture

// The pixel aspect ratio and clean aperture of the video may be specified using
// the ‘pasp’ and ‘clap’ sample entry boxes, respectively. These are both
// optional; if present, they over‐ride the declarations (if any) in structures
// specific to the video codec, which structures should be examined if these
// boxes are absent. For maximum compatibility, these boxes should follow, not
// precede, any boxes defined in or required by derived specifications.
//
// In the PixelAspectRatioBox, hSpacing and vSpacing have the same units, but
// those units are unspecified: only the ratio matters. hSpacing and vSpacing
// may or may not be in reduced terms, and they may reduce to 1/1. Both of them
// must be positive.
//
// They are defined as the aspect ratio of a pixel, in arbitrary units. If a
// pixel appears H wide and V tall, then hSpacing/vSpacing is equal to H/V. This
// means that a square on the display that is n pixels tall needs to be
// n*vSpacing/hSpacing pixels wide to appear square.
//
// > NOTE When adjusting pixel aspect ratio, normally, the horizontal dimension
// of the video is scaled, if needed (i.e. if the final display system has a
// different pixel aspect ratio from the video source).
//
// > NOTE It is recommended that the original pixels, and the composed
// transform, be carried through the pipeline as far as possible. If the
// transformation resulting from ‘correcting’ pixel aspect ratio to a square
// grid, normalizing to the track dimensions, composition or placement (e.g.
// track and/or movie matrix), and normalizing to the display characteristics,
// is a unity matrix, then no re‐sampling need be done. In particular, video
// should not be resampled more than once in the process of rendering, if at all
// possible.
//
// There are notionally four values in the CleanApertureBox. These parameters
// are represented as a fraction N/D. The fraction may or may not be in reduced
// terms. We refer to the pair of parameters fooN and fooD as foo. For horizOff
// and vertOff, D must be positive and N may be positive or negative. For
// cleanApertureWidth and cleanApertureHeight, both N and D must be positive.
//
// > NOTE These are fractional numbers for several reasons. First, in some
// systems the exact width after pixel aspect ratio correction is integral, not
// the pixel count before that correction. Second, if video is resized in the
// full aperture, the exact expression for the clean aperture may not be
// integral. Finally, because this is represented using centre and offset, a
// division by two is needed, and so half‐values can occur.
//
// Considering the pixel dimensions as defined by the VisualSampleEntry width
// and height. If picture centre of the image is at pcX and pcY, then horizOff
// and vertOff are defined as follows:
//
//     pcX = horizOff + (width - 1)/2
//     pcY = vertOff + (height - 1)/2;
//
// Typically, horizOff and vertOff are zero, so the image is centred about the
// picture centre. The leftmost/rightmost pixel and the topmost/bottommost line
// of the clean aperture fall at:
//
//     pcX ± (cleanApertureWidth - 1)/2
//     pcY ± (cleanApertureHeight - 1)/2;
type PixelAspectRatioBox struct {
	Header
	NullContainer

	// define the relative width and height of a pixel
	HSpacing uint32
	VSpacing uint32
}

var _ Box = (*PixelAspectRatioBox)(nil)

func init() {
	BoxRegistry[PaspBoxType] = func() Box { return &PixelAspectRatioBox{} }
}

func (b PixelAspectRatioBox) Mp4BoxType() BoxType {
	return PaspBoxType
}

func (b *PixelAspectRatioBox) Mp4BoxUpdate() uint32 {
	b.Type = b.Mp4BoxType()
	b.Size = b.HeaderSize()
	b.Size += 4 // unsigned int(32) hSpacing;
	b.Size += 4 // unsigned int(32) vSpacing;
	return b.Size
}

func (b *PixelAspectRatioBox) Mp4BoxRead(r io.Reader, header *Header) (err error) {
	if err = b.ReadHeader(r, header); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.HSpacing); err != nil {
		return
	}
	if err = binary.Read(r, binary.BigEndian, &b.VSpacing); err != nil {
		return
	}
	return
}

func (b *PixelAspectRatioBox) Mp4BoxWrite(w io.Writer) (err error) {
	if err = b.WriteHeader(w); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.HSpacing); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, b.VSpacing); err != nil {
		return
	}
	return
}
