package fgray


import "github.com/phrfp/bild/math/f64"

type GRAYF64 struct {
	GY float64
}


// NewRGBAF64 returns a new RGBAF64 color based on the provided uint8 values.
// uint8 value 0 maps to 0, 128 to 0.5 and 255 to 1.0.
func NewGrayF64( gr uint16) GRAYF64 {
	return  GRAYF64{float64(gr)/65535}
}

func NewGrayG8F64( gr uint8) GRAYF64 {
	return  GRAYF64{float64(gr)/255}
}


// Clamp limits the channel values of the Gray16 color to the range 0.0 to 1.0.
func (ch *GRAYF64) ClampG16() {
	ch.GY = f64.Clamp(ch.GY, 0, 1)
}
