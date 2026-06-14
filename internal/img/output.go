// Copyright © 2020 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package img

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"strings"
	"unicode"

	"github.com/esimov/stackblur-go"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/gonvenience/bunt"
	"github.com/gonvenience/font"
	"github.com/gonvenience/term"
	cjkfont "github.com/homeport/termshot/internal/font"
	imgfont "golang.org/x/image/font"
)

const (
	red    = "#ED655A"
	yellow = "#E1C04C"
	green  = "#71BD47"
)

const (
	defaultFontSize = 12
	defaultFontDPI  = 144
)

// commandIndicator is the string to be used to indicate the command in the screenshot
var commandIndicator = func() string {
	if val, ok := os.LookupEnv("TS_COMMAND_INDICATOR"); ok {
		return val
	}

	return "➜"
}()

type Scaffold struct {
	content bunt.String

	factor float64

	columns int

	defaultForegroundColor color.Color

	clipCanvas bool

	drawDecorations bool
	drawShadow      bool

	shadowBaseColor string
	shadowRadius    uint8
	shadowOffsetX   float64
	shadowOffsetY   float64

	padding float64
	margin  float64

	regular     imgfont.Face
	bold        imgfont.Face
	italic      imgfont.Face
	boldItalic  imgfont.Face
	lineSpacing float64
	tabSpaces   int

	// CJK font faces for Chinese character support
	cjkRegular    imgfont.Face
	cjkBold       imgfont.Face
	cjkItalic     imgfont.Face
	cjkBoldItalic imgfont.Face

	// Symbol font faces for Unicode symbol support (DejaVu Sans has wide coverage)
	symbolRegular    imgfont.Face
	symbolBold       imgfont.Face
	symbolItalic     imgfont.Face
	symbolBoldItalic imgfont.Face
}

func NewImageCreator() Scaffold {
	f := 2.0

	fontFaceOptions := &truetype.Options{
		Size: f * defaultFontSize,
		DPI:  defaultFontDPI,
	}

	cjkFont := cjkfont.NotoSansCJK{}
	symbolFont := cjkfont.DejaVuSans{}

	return Scaffold{
		defaultForegroundColor: bunt.LightGray,

		factor: f,

		margin:  f * 48,
		padding: f * 24,

		drawDecorations: true,
		drawShadow:      true,

		shadowBaseColor: "#10101066",
		shadowRadius:    uint8(math.Min(f*16, 255)),
		shadowOffsetX:   f * 16,
		shadowOffsetY:   f * 16,

		regular:    font.Hack.Regular(fontFaceOptions),
		bold:       font.Hack.Bold(fontFaceOptions),
		italic:     font.Hack.Italic(fontFaceOptions),
		boldItalic: font.Hack.BoldItalic(fontFaceOptions),

		cjkRegular:    cjkFont.Regular(fontFaceOptions),
		cjkBold:       cjkFont.Bold(fontFaceOptions),
		cjkItalic:     cjkFont.Italic(fontFaceOptions),
		cjkBoldItalic: cjkFont.BoldItalic(fontFaceOptions),

		symbolRegular:    symbolFont.Regular(fontFaceOptions),
		symbolBold:       symbolFont.Bold(fontFaceOptions),
		symbolItalic:     symbolFont.Italic(fontFaceOptions),
		symbolBoldItalic: symbolFont.BoldItalic(fontFaceOptions),

		lineSpacing: 1.2,
		tabSpaces:   2,
	}
}

func (s *Scaffold) SetFontFaceRegular(face imgfont.Face) { s.regular = face }

func (s *Scaffold) SetFontFaceBold(face imgfont.Face) { s.bold = face }

func (s *Scaffold) SetFontFaceItalic(face imgfont.Face) { s.italic = face }

func (s *Scaffold) SetFontFaceBoldItalic(face imgfont.Face) { s.boldItalic = face }

func (s *Scaffold) SetColumns(columns int) { s.columns = columns }

func (s *Scaffold) SetMargin(margin float64) { s.margin = margin * s.factor }

func (s *Scaffold) SetPadding(padding float64) { s.padding = padding * s.factor }

func (s *Scaffold) DrawDecorations(value bool) { s.drawDecorations = value }

func (s *Scaffold) DrawShadow(value bool) { s.drawShadow = value }

func (s *Scaffold) ClipCanvas(value bool) { s.clipCanvas = value }

func (s *Scaffold) GetFixedColumns() int {
	if s.columns != 0 {
		return s.columns
	}

	columns, _ := term.GetTerminalSize()
	return columns
}

func (s *Scaffold) AddCommand(args ...string) error {
	return s.AddContent(strings.NewReader(
		bunt.Sprintf("Lime{%s} DimGray{%s}\n",
			commandIndicator,
			strings.Join(args, " "),
		),
	))
}

func (s *Scaffold) AddContent(in io.Reader) error {
	parsed, err := bunt.ParseStream(in)
	if err != nil {
		return fmt.Errorf("failed to parse input stream: %w", err)
	}

	var tmp bunt.String
	var counter int
	for _, cr := range *parsed {
		counter++

		if cr.Symbol == '\n' {
			counter = 0
		}

		// Add an additional newline in case the column
		// count is reached and line wrapping is needed
		if counter > s.GetFixedColumns() {
			counter = 0
			tmp = append(tmp, bunt.ColoredRune{
				Settings: cr.Settings,
				Symbol:   '\n',
			})
		}

		tmp = append(tmp, cr)
	}

	s.content = append(s.content, tmp...)

	return nil
}

func (s *Scaffold) fontHeight() float64 {
	return float64(s.regular.Metrics().Height >> 6)
}

// isCJKChar checks if a rune is a CJK (Chinese, Japanese, Korean) character
func isCJKChar(r rune) bool {
	return unicode.Is(unicode.Han, r) ||
		unicode.Is(unicode.Hiragana, r) ||
		unicode.Is(unicode.Katakana, r) ||
		unicode.Is(unicode.Hangul, r) ||
		(r >= 0x3000 && r <= 0x303F) || // CJK Symbols and Punctuation
		(r >= 0xFF00 && r <= 0xFFEF) // Halfwidth and Fullwidth Forms
}

// isSymbolChar checks if a rune is a Unicode symbol character
func isSymbolChar(r rune) bool {
	// Common Unicode symbol ranges that may not be in Hack font
	return (r >= 0x2190 && r <= 0x21FF) || // Arrows
		(r >= 0x2200 && r <= 0x22FF) || // Mathematical Operators
		(r >= 0x2300 && r <= 0x23FF) || // Miscellaneous Technical
		(r >= 0x2460 && r <= 0x24FF) || // Enclosed Alphanumerics
		(r >= 0x2500 && r <= 0x257F) || // Box Drawing
		(r >= 0x2580 && r <= 0x259F) || // Block Elements
		(r >= 0x25A0 && r <= 0x25FF) || // Geometric Shapes
		(r >= 0x2600 && r <= 0x26FF) || // Miscellaneous Symbols
		(r >= 0x2700 && r <= 0x27BF) // Dingbats (✓ ✘ ✔ etc.)
}

func (s *Scaffold) measureContent() (width float64, height float64) {
	var tmp = make([]rune, len(s.content))
	for i, cr := range s.content {
		tmp[i] = cr.Symbol
	}

	lines := strings.Split(
		strings.TrimSuffix(
			string(tmp),
			"\n",
		),
		"\n",
	)

	// temporary drawer for reference calucation
	tmpDrawer := &imgfont.Drawer{Face: s.regular}

	// width, either by using longest line, or by fixed column value
	switch s.columns {
	case 0: // unlimited: max width of all lines
		for _, line := range lines {
			advance := tmpDrawer.MeasureString(line)
			if lineWidth := float64(advance >> 6); lineWidth > width {
				width = lineWidth
			}
		}

	default: // fixed: max width based on column count
		width = float64(tmpDrawer.MeasureString(strings.Repeat("a", s.GetFixedColumns())) >> 6)
	}

	// height, lines times font height and line spacing
	height = float64(len(lines)) * s.fontHeight() * s.lineSpacing

	return width, height
}

func (s *Scaffold) image() (image.Image, error) {
	var f = func(value float64) float64 { return s.factor * value }

	var (
		corner   = f(6)
		radius   = f(9)
		distance = f(25)
	)

	contentWidth, contentHeight := s.measureContent()

	// Make sure the output window is big enough in case no content or very few
	// content will be rendered
	contentWidth = math.Max(contentWidth, 3*distance+3*radius)

	marginX, marginY := s.margin, s.margin
	paddingX, paddingY := s.padding, s.padding

	xOffset := marginX
	yOffset := marginY

	var titleOffset float64
	if s.drawDecorations {
		titleOffset = f(40)
	}

	width := contentWidth + 2*marginX + 2*paddingX
	height := contentHeight + 2*marginY + 2*paddingY + titleOffset

	dc := gg.NewContext(int(width), int(height))

	// Optional: Apply blurred rounded rectangle to mimic the window shadow
	//
	if s.drawShadow {
		xOffset -= s.shadowOffsetX / 2
		yOffset -= s.shadowOffsetY / 2

		bc := gg.NewContext(int(width), int(height))
		bc.DrawRoundedRectangle(xOffset+s.shadowOffsetX, yOffset+s.shadowOffsetY, width-2*marginX, height-2*marginY, corner)
		bc.SetHexColor(s.shadowBaseColor)
		bc.Fill()

		src := bc.Image()
		dst := image.NewNRGBA(src.Bounds())
		if err := stackblur.Process(dst, src, uint32(s.shadowRadius)); err != nil {
			return nil, err
		}

		dc.DrawImage(dst, 0, 0)
	}

	// Draw rounded rectangle with outline to produce impression of a window
	//
	dc.DrawRoundedRectangle(xOffset, yOffset, width-2*marginX, height-2*marginY, corner)
	dc.SetHexColor("#151515")
	dc.Fill()

	dc.DrawRoundedRectangle(xOffset, yOffset, width-2*marginX, height-2*marginY, corner)
	dc.SetHexColor("#404040")
	dc.SetLineWidth(f(1))
	dc.Stroke()

	// Optional: Draw window decorations (i.e. three buttons) to produce the
	// impression of an actional window
	//
	if s.drawDecorations {
		for i, color := range []string{red, yellow, green} {
			dc.DrawCircle(xOffset+paddingX+float64(i)*distance+f(4), yOffset+paddingY+f(4), radius)
			dc.SetHexColor(color)
			dc.Fill()
		}
	}

	// Apply the actual text into the prepared content area of the window
	//
	var x, y = xOffset + paddingX, yOffset + paddingY + titleOffset + s.fontHeight()
	for _, cr := range s.content {
		// Determine font face based on style settings and character type
		var fontFace imgfont.Face
		useCJK := isCJKChar(cr.Symbol)
		useSymbol := isSymbolChar(cr.Symbol)

		// Priority: Symbol > CJK > Regular
		if useSymbol {
			// Use DejaVu Sans for Unicode symbols
			switch cr.Settings & 0x1C {
			case 4: // Bold
				fontFace = s.symbolBold
			case 8: // Italic
				fontFace = s.symbolItalic
			case 12: // Bold Italic
				fontFace = s.symbolBoldItalic
			default: // Regular
				fontFace = s.symbolRegular
			}
		} else if useCJK {
			// Use Noto Sans CJK for CJK characters
			switch cr.Settings & 0x1C {
			case 4: // Bold
				fontFace = s.cjkBold
			case 8: // Italic
				fontFace = s.cjkItalic
			case 12: // Bold Italic
				fontFace = s.cjkBoldItalic
			default: // Regular
				fontFace = s.cjkRegular
			}
		} else {
			// Use Hack for regular ASCII/Latin characters
			switch cr.Settings & 0x1C {
			case 4: // Bold
				fontFace = s.bold
			case 8: // Italic
				fontFace = s.italic
			case 12: // Bold Italic
				fontFace = s.boldItalic
			default: // Regular
				fontFace = s.regular
			}
		}

		dc.SetFontFace(fontFace)

		str := string(cr.Symbol)
		w, h := dc.MeasureString(str)

		// Use consistent line height based on the regular font for uniform backgrounds
		// This ensures all characters on the same line have the same background height
		baseMetrics := s.regular.Metrics()
		baseAscent := float64(baseMetrics.Ascent >> 6)
		baseDescent := float64(baseMetrics.Descent >> 6)
		lineHeight := baseAscent + baseDescent

		// background color
		switch cr.Settings & 0x02 { //nolint:gocritic
		case 2:
			dc.SetRGB255(
				int((cr.Settings>>32)&0xFF), // #nosec G115
				int((cr.Settings>>40)&0xFF), // #nosec G115
				int((cr.Settings>>48)&0xFF), // #nosec G115
			)

			// Draw background rectangle with consistent height
			dc.DrawRectangle(x, y-baseAscent, w, lineHeight)
			dc.Fill()
		}

		// foreground color
		switch cr.Settings & 0x01 {
		case 1:
			dc.SetRGB255(
				int((cr.Settings>>8)&0xFF),  // #nosec G115
				int((cr.Settings>>16)&0xFF), // #nosec G115
				int((cr.Settings>>24)&0xFF), // #nosec G115
			)

		default:
			dc.SetColor(s.defaultForegroundColor)
		}

		switch str {
		case "\n":
			x = xOffset + paddingX
			y += h * s.lineSpacing
			continue

		case "\t":
			x += w * float64(s.tabSpaces)
			continue

		case "✗", "ˣ": // mitigate issue #1 by replacing it with a similar character
			str = "×"
		}

		dc.DrawString(str, x, y)

		// There seems to be no font face based way to do an underlined
		// string, therefore manually draw a line under each character
		if cr.Settings&0x1C == 16 {
			dc.DrawLine(x, y+f(4), x+w, y+f(4))
			dc.SetLineWidth(f(1))
			dc.Stroke()
		}

		x += w
	}

	return dc.Image(), nil
}

// Write writes the scaffold content as PNG into the provided writer
//
// Deprecated: Use [Scaffold.WritePNG] instead.
func (s *Scaffold) Write(w io.Writer) error {
	return s.WritePNG(w)
}

// WritePNG writes the scaffold content as PNG into the provided writer
func (s *Scaffold) WritePNG(w io.Writer) error {
	img, err := s.image()
	if err != nil {
		return err
	}

	// Optional: Clip image to minimum size by removing all surrounding transparent pixels
	//
	if s.clipCanvas {
		if imgRGBA, ok := img.(*image.RGBA); ok {
			var minX, minY = math.MaxInt, math.MaxInt
			var maxX, maxY = 0, 0

			var bounds = imgRGBA.Bounds()
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
					r, g, b, a := imgRGBA.At(x, y).RGBA()
					isTransparent := r == 0 && g == 0 && b == 0 && a == 0

					if !isTransparent {
						if x < minX {
							minX = x
						}

						if y < minY {
							minY = y
						}

						if x > maxX {
							maxX = x
						}

						if y > maxY {
							maxY = y
						}
					}
				}
			}

			img = imgRGBA.SubImage(image.Rect(minX, minY, maxX, maxY))
		}
	}

	return png.Encode(w, img)
}

// WriteRaw writes the scaffold content as-is into the provided writer
func (s *Scaffold) WriteRaw(w io.Writer) error {
	_, err := w.Write([]byte(s.content.String()))
	return err
}
