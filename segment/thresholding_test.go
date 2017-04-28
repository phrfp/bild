package segment

import (
	"image"
	"testing"

	"github.com/phrfp/bild/util"
)

func TestThreshold(t *testing.T) {
	cases := []struct {
		level    uint8
		img      image.Image
		expected *image.Gray
	}{
		{
			level: 0,
			img: &image.RGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x80, 0x80, 0x80,
				},
			},
			expected: &image.Gray{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2,
				Pix: []uint8{
					0xFF, 0xFF,
					0xFF, 0xFF,
				},
			},
		},
		{
			level: 128,
			img: &image.RGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x80, 0x80, 0x80,
				},
			},
			expected: &image.Gray{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2,
				Pix: []uint8{
					0x00, 0xFF,
					0xFF, 0x00,
				},
			},
		},
		{
			level: 255,
			img: &image.RGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x80, 0x80, 0x80,
				},
			},
			expected: &image.Gray{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2,
				Pix: []uint8{
					0x00, 0xFF,
					0xFF, 0x00,
				},
			},
		},
		{
			level: 127,
			img: &image.RGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0xFF, 0xC0, 0xC0, 0xC0, 0xFF,
					0x40, 0x40, 0x40, 0x40, 0x80, 0x80, 0x80, 0x80,
				},
			},
			expected: &image.Gray{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2,
				Pix: []uint8{
					0xFF, 0xFF,
					0x00, 0xFF,
				},
			},
		},
	}

	for _, c := range cases {
		actual := Threshold(c.img, c.level)
		if !util.GrayImageEqual(actual, c.expected) {
			t.Errorf("%s: expected: %v actual: %v", "Threshold", c.expected, actual)
		}
	}
}

func TestThresholdG16(t *testing.T) {
	cases := []struct {
		level_up  uint16
		level_lwr uint16
		img      image.Image
		expected *image.Gray
	}{
		{
			level_lwr: 0,
			level_up: 0xffff,
			img: &image.Gray16{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2*2,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF,
				},
			},
			expected: &image.Gray{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2,
				Pix: []uint8{
					0x00, 0x00,
					0x00, 0x00,
				},
			},
		},
		{
			level_lwr: 128,
			level_up: 0xffff,
			img: &image.Gray16{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 2,
				Pix: []uint8{
					0x00, 0x08, 0x80, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF,
				},
			},
			expected: &image.Gray{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2,
				Pix: []uint8{
					0xFF, 0x00,
					0x00, 0x00,
				},
			},
		},
	}

	for _, c := range cases {
		actual := ThresholdG16(c.img, c.level_lwr, c.level_up)
		if !util.GrayImageEqual(actual, c.expected) {
		t.Errorf("%s: \nexpected:%v\nactual:%v\n", "Trheshold", util.GrayToString(c.expected), util.GrayToString(actual))
		}
	}
}
