/*Package clone provides image cloning function.*/
package clone

import (
	"image"
	"image/draw"
	//"fmt"
	"github.com/phrfp/bild/parallel"
)

// PadMethod is the method used to fill padded pixels.
type PadMethod uint8

const (
	// NoFill leaves the padded pixels empty.
	NoFill = iota
	// EdgeExtend extends the closest edge pixel.
	EdgeExtend
	// EdgeWrap wraps around the pixels of an image.
	EdgeWrap
)

// AsRGBA returns an RGBA copy of the supplied image.
func AsRGBA(src image.Image) *image.RGBA {
	bounds := src.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, src, bounds.Min, draw.Src)
	return img
}

// AsGray16 returns a Gray16 copy of the supplied image.
func AsGray16(src image.Image) *image.Gray16 {
	bounds := src.Bounds()
	img := image.NewGray16(bounds)
	draw.Draw(img, bounds, src, bounds.Min, draw.Src)
	return img
}


// Pad returns an RGBA copy of the src image paramter with its edges padded
// using the supplied PadMethod.
// Parameter padX and padY correspond to the amount of padding to be applied
// on each side.
// Parameter m is the PadMethod to fill the new pixels.
//
// Usage example:
//
//		result := Pad(img, 5,5, EdgeExtend)
//
func Pad(src image.Image, padX, padY int, m PadMethod) *image.RGBA {

		var result *image.RGBA
		switch m {
		case EdgeExtend:
			result = extend(src, padX, padY)
		case NoFill:
			result = noFill(src, padX, padY)
		case EdgeWrap:
			result = wrap(src, padX, padY)
		default:
			result = extend(src, padX, padY)
		}
		return result

}

func PadG16(src image.Image, padX, padY int, m PadMethod) *image.Gray16{
	var result *image.Gray16
	switch m {
	case EdgeExtend:
		result = extendG16(src, padX, padY)
	case NoFill:
		result = noFillG16(src, padX, padY)
	default:
		result = extendG16(src, padX, padY)
	}
	return result
}


func noFill(img image.Image, padX, padY int) *image.RGBA {
	srcBounds := img.Bounds()
	paddedW, paddedH := srcBounds.Dx()+2*padX, srcBounds.Dy()+2*padY
	newBounds := image.Rect(0, 0, paddedW, paddedH)
	fillBounds := image.Rect(padX, padY, padX+srcBounds.Dx(), padY+srcBounds.Dy())

	dst := image.NewRGBA(newBounds)
	draw.Draw(dst, fillBounds, img, srcBounds.Min, draw.Src)

	return dst
}

func noFillG16(img image.Image, padX, padY int) *image.Gray16 {
	srcBounds := img.Bounds()
	paddedW, paddedH := srcBounds.Dx()+2*padX, srcBounds.Dy()+2*padY
	newBounds := image.Rect(0, 0, paddedW, paddedH)
	fillBounds := image.Rect(padX, padY, padX+srcBounds.Dx(), padY+srcBounds.Dy())

	dst := image.NewGray16(newBounds)
	draw.Draw(dst, fillBounds, img, srcBounds.Min, draw.Src)

	return dst
}

func extend(img image.Image, padX, padY int) *image.RGBA {
	dst := noFill(img, padX, padY)
	paddedW, paddedH := dst.Bounds().Dx(), dst.Bounds().Dy()

	parallel.Line(paddedH, func(start, end int) {
		for y := start; y < end; y++ {
			iy := y
			if iy < padY {
				iy = padY
			} else if iy >= paddedH-padY {
				iy = paddedH - padY - 1
			}

			for x := 0; x < paddedW; x++ {
				ix := x
				if ix < padX {
					ix = padX
				} else if x >= paddedW-padX {
					ix = paddedW - padX - 1
				} else if iy == y {
					// This only enters if we are not in a y-padded area or
					// x-padded area, so nothing to extend here.
					// So simply jump to the next padded-x index.
					x = paddedW - padX - 1
					continue
				}

				dstPos := y*dst.Stride + x*4
				edgePos := iy*dst.Stride + ix*4

				dst.Pix[dstPos+0] = dst.Pix[edgePos+0]
				dst.Pix[dstPos+1] = dst.Pix[edgePos+1]
				dst.Pix[dstPos+2] = dst.Pix[edgePos+2]
				dst.Pix[dstPos+3] = dst.Pix[edgePos+3]
			}
		}
	})

	return dst
}

func extendG16(img image.Image, padX, padY int) *image.Gray16 {
	dst := noFillG16(img, padX, padY)
	paddedW, paddedH := dst.Bounds().Dx(), dst.Bounds().Dy()

	parallel.Line(paddedH, func(start, end int) {
		for y := start; y < end; y++ {
			iy := y
			if iy < padY {
				iy = padY
			} else if iy >= paddedH-padY {
				iy = paddedH - padY - 1
			}

			for x := 0; x < paddedW; x++ {
				ix := x
				if ix < padX {
					ix = padX
				} else if x >= paddedW-padX {
					ix = paddedW - padX - 1
				} else if iy == y {
					// This only enters if we are not in a y-padded area or
					// x-padded area, so nothing to extend here.
					// So simply jump to the next padded-x index.
					x = paddedW - padX - 1
					continue
				}

				dstPos := y*dst.Stride + x*2
				edgePos := iy*dst.Stride + ix*2
			//	fmt.Println("Edge Pos: ", edgePos, " Dst Pos: ", dstPos)
				dst.Pix[dstPos+0] = dst.Pix[edgePos+0]
				dst.Pix[dstPos+1] = dst.Pix[edgePos+1]

			}
		}
	})

	return dst
}

func wrap(img image.Image, padX, padY int) *image.RGBA {
	dst := noFill(img, padX, padY)
	paddedW, paddedH := dst.Bounds().Dx(), dst.Bounds().Dy()

	parallel.Line(paddedH, func(start, end int) {
		for y := start; y < end; y++ {
			iy := y
			if iy < padY {
				iy = (paddedH - padY) - ((padY - y) % (paddedH - padY*2))
			} else if iy >= paddedH-padY {
				iy = padY - ((padY - y) % (paddedH - padY*2))
			}

			for x := 0; x < paddedW; x++ {
				ix := x
				if ix < padX {
					ix = (paddedW - padX) - ((padX - x) % (paddedW - padX*2))
				} else if ix >= paddedW-padX {
					ix = padX - ((padX - x) % (paddedW - padX*2))
				} else if iy == y {
					// This only enters if we are not in a y-padded area or
					// x-padded area, so nothing to extend here.
					// So simply jump to the next padded-x index.
					x = paddedW - padX - 1
					continue
				}

				dstPos := y*dst.Stride + x*4
				edgePos := iy*dst.Stride + ix*4

				dst.Pix[dstPos+0] = dst.Pix[edgePos+0]
				dst.Pix[dstPos+1] = dst.Pix[edgePos+1]
				dst.Pix[dstPos+2] = dst.Pix[edgePos+2]
				dst.Pix[dstPos+3] = dst.Pix[edgePos+3]
			}
		}
	})

	return dst
}
