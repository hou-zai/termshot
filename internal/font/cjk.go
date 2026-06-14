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
	"golang.org/x/image/font/opentype"
)

//go:embed assets/NotoSansCJK-Regular.ttc
var notoSansCJKRegular []byte

//go:embed assets/NotoSansCJK-Bold.ttc
var notoSansCJKBold []byte

// NotoSansCJK provides embedded Noto Sans CJK (思源黑体) font faces with full Chinese character support
type NotoSansCJK struct{}

// Regular returns the regular font face with the given options
func (NotoSansCJK) Regular(options *truetype.Options) font.Face {
	return parseTTCFont(notoSansCJKRegular, 0, options)
}

// Bold returns the bold font face with the given options
func (NotoSansCJK) Bold(options *truetype.Options) font.Face {
	return parseTTCFont(notoSansCJKBold, 0, options)
}

// Italic returns the regular font face (CJK fonts typically don't have true italic variants)
func (NotoSansCJK) Italic(options *truetype.Options) font.Face {
	return parseTTCFont(notoSansCJKRegular, 0, options)
}

// BoldItalic returns the bold font face (CJK fonts typically don't have true italic variants)
func (NotoSansCJK) BoldItalic(options *truetype.Options) font.Face {
	return parseTTCFont(notoSansCJKBold, 0, options)
}

// parseTTCFont parses a TrueType Collection (TTC) font and returns a specific font face
func parseTTCFont(data []byte, index int, options *truetype.Options) font.Face {
	// Parse as OpenType Collection
	collection, err := opentype.ParseCollection(data)
	if err != nil {
		panic(fmt.Sprintf("failed to parse TTC font collection: %v", err))
	}

	if index >= collection.NumFonts() {
		panic(fmt.Sprintf("font index %d out of range (collection has %d fonts)", index, collection.NumFonts()))
	}

	// Get the font at the specified index
	f, err := collection.Font(index)
	if err != nil {
		panic(fmt.Sprintf("failed to get font at index %d: %v", index, err))
	}

	// Create font face with the provided options
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    options.Size,
		DPI:     options.DPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to create font face: %v", err))
	}

	return face
}
