// Copyright © 2024 The Homeport Team
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

package font

import (
	_ "embed"
	"fmt"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed assets/DejaVuSans.ttf
var dejaVuSansRegular []byte

//go:embed assets/DejaVuSans-Bold.ttf
var dejaVuSansBold []byte

// DejaVuSans provides embedded DejaVu Sans font for Unicode symbol support
type DejaVuSans struct{}

// Regular returns the regular font face with the given options
func (DejaVuSans) Regular(options *truetype.Options) font.Face {
	return parseTTFFont(dejaVuSansRegular, options)
}

// Bold returns the bold font face with the given options
func (DejaVuSans) Bold(options *truetype.Options) font.Face {
	return parseTTFFont(dejaVuSansBold, options)
}

// Italic returns the regular font face (DejaVu Sans italic not embedded)
func (DejaVuSans) Italic(options *truetype.Options) font.Face {
	return parseTTFFont(dejaVuSansRegular, options)
}

// BoldItalic returns the bold font face (DejaVu Sans bold italic not embedded)
func (DejaVuSans) BoldItalic(options *truetype.Options) font.Face {
	return parseTTFFont(dejaVuSansBold, options)
}

// parseTTFFont parses a TrueType Font and returns a font face
func parseTTFFont(data []byte, options *truetype.Options) font.Face {
	f, err := truetype.Parse(data)
	if err != nil {
		panic(fmt.Sprintf("failed to parse TTF font: %v", err))
	}

	face := truetype.NewFace(f, options)
	return face
}
