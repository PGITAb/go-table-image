package tableimage

import (
	"image"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

func defaultFaceOptions() *opentype.FaceOptions {
	return &opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingNone,
	}
}

var incosolataface font.Face = inconsolata.Bold8x16

//https://cs.opensource.google/go/x/image/+/062f8c9f:font/opentype/opentype_test.go;l=27;bpv=1;bpt=1
func initface() font.Face {
	font, err := sfnt.Parse(gomedium.TTF)
	if err != nil {
		panic(err)
	}
	regular, err := opentype.NewFace(font, defaultFaceOptions())
	if err != nil {
		panic(err)
	}
	return regular
}

var gootherface = initface()

var goboldface = basicfont.Face{
	Advance: 7,
	Width:   6,
	Height:  13,
	Ascent:  11,
	Descent: 2,
	//Left:    -1,
	Mask: &image.Alpha{
		Stride: 6,
		Rect:   image.Rectangle{Max: image.Point{6, 96 * 13}},
		Pix:    gomono.TTF, // this is incorrect, refer to https://cs.opensource.google/go/x/image/+/062f8c9f:font/inconsolata/bold8x16.go;bpv=1;bpt=1
		//Pix: *(inconsolata.Bold8x16).Mask.Pix,
	},
	//Ranges: []basicfont.Range{
	//	{'\u0020', '\u007f', 0},
	//	{'\ufffd', '\ufffe', 95},
	//},
}

func (ti *TableImage) setRgba() {
	img := image.NewRGBA(image.Rect(0, 0, ti.width, ti.height))
	//set image background
	draw.Draw(img, img.Bounds(), &image.Uniform{getColorByHex(ti.backgroundColor)}, image.ZP, draw.Src)
	ti.img = img
}

func (ti *TableImage) setRgbaTH(hexcolor string) {
	th_row := image.Rect(0, 0, ti.width, 1*rowSpace+separatorPadding)

	//set image background
	draw.Draw(ti.img, th_row, &image.Uniform{getColorByHex(hexcolor)}, image.ZP, draw.Src)

	//draw.Draw(img, img.Bounds(), &image.Uniform{getColorByHex(hexcolor)}, image.ZP, draw.Src)
	//ti.img = ti.img
}

func (ti *TableImage) addString(x, y int, label string, color string) {

	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst: ti.img,
		Src: image.NewUniform(getColorByHex(color)),
		//Face: basicfont.Face7x13, //original
		//Face: gootherface,  //working face with ttf
		Face: incosolataface, //additional option
		//Face: &goboldface,  //not working
		Dot: point,
	}
	d.DrawString(label)
}

//Thx to https://github.com/StephaneBunel/bresenham
func (ti *TableImage) addLine(x1, y1, x2, y2 int, color string) {

	var dx, dy, e, slope int
	col := getColorByHex(color)
	// Because drawing p1 -> p2 is equivalent to draw p2 -> p1,
	// I sort points in x-axis order to handle only half of possible cases.
	if x1 > x2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}

	dx, dy = x2-x1, y2-y1
	// Because point is x-axis ordered, dx cannot be negative
	if dy < 0 {
		dy = -dy
	}

	switch {

	// Is line a point ?
	case x1 == x2 && y1 == y2:
		ti.img.Set(x1, y1, col)

	// Is line an horizontal ?
	case y1 == y2:
		for ; dx != 0; dx-- {
			ti.img.Set(x1, y1, col)
			x1++
		}
		ti.img.Set(x1, y1, col)

	// Is line a vertical ?
	case x1 == x2:
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for ; dy != 0; dy-- {
			ti.img.Set(x1, y1, col)
			y1++
		}
		ti.img.Set(x1, y1, col)

	// Is line a diagonal ?
	case dx == dy:
		if y1 < y2 {
			for ; dx != 0; dx-- {
				ti.img.Set(x1, y1, col)
				x1++
				y1++
			}
		} else {
			for ; dx != 0; dx-- {
				ti.img.Set(x1, y1, col)
				x1++
				y1--
			}
		}
		ti.img.Set(x1, y1, col)

	// wider than high ?
	case dx > dy:
		if y1 < y2 {
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				ti.img.Set(x1, y1, col)
				x1++
				e -= dy
				if e < 0 {
					y1++
					e += slope
				}
			}
		} else {
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				ti.img.Set(x1, y1, col)
				x1++
				e -= dy
				if e < 0 {
					y1--
					e += slope
				}
			}
		}
		ti.img.Set(x2, y2, col)

	// higher than wide.
	default:
		if y1 < y2 {
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				ti.img.Set(x1, y1, col)
				y1++
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		} else {
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				ti.img.Set(x1, y1, col)
				y1--
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		}
		ti.img.Set(x2, y2, col)
	}
}
