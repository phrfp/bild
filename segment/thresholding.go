/*Package segment provides basic image segmentation and clusterring methods.*/
package segment

import (
	"image"
	"image/color"

	"github.com/phrfp/bild/clone"
	"github.com/phrfp/bild/util"
	"github.com/phrfp/bild/parallel"

)

// Threshold returns a grayscale image in which values from the param img that are
// smaller than the param level are set to black and values larger than or equal to
// it are set to white.
// Level must be of the range 0 to 255.
func Threshold(img image.Image, level uint8) *image.Gray {
	src := clone.AsRGBA(img)
	bounds := src.Bounds()

	dst := image.NewGray(bounds)

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			srcPos := y*src.Stride + x*4
			dstPos := y*dst.Stride + x

			c := src.Pix[srcPos : srcPos+4]
			r := util.Rank(color.RGBA{c[0], c[1], c[2], c[3]})

			if uint8(r) >= level {
				dst.Pix[dstPos] = 0xFF
			} else {
				dst.Pix[dstPos] = 0x00
			}
		}
	}

	return dst
}

func ThresholdG16( img image.Image, level_lower, level_upper uint16 ) *image.Gray {

	src := clone.AsGray16(img)
	bounds := src.Bounds()
	w := src.Bounds().Dx()
	h := src.Bounds().Dy()
	dst := image.NewGray(bounds)

		parallel.Line(w, func(start, end int){
			for x := start; x < end; x++ {
				for y:= 0; y < h; y++{
						ipos := y*src.Stride + x*2
						dstPos := y*dst.Stride + x
						tpix := uint16(src.Pix[ipos+0])<<8 | uint16(src.Pix[ipos+1])
						if (tpix >= level_lower) && (tpix <= level_upper) {
							dst.Pix[dstPos] = 0x00
						} else {
							dst.Pix[dstPos] = 0xFF
						}
				}
			}
		})
		return dst
}
