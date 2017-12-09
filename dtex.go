package main

import (
	"errors"
	//"fmt"
	"github.com/docopt/docopt-go"
	"image"
	_ "image/png"
	"io/ioutil"
	"os"
)

type textureFormat int

const (
	bpp2 textureFormat = iota
	bpp4
	bpp8
	bpp16
	a3i5
	a5i3
	c4x4
)

var (
	formatStrings = map[string]textureFormat{
		"2bpp":  bpp2,
		"4bpp":  bpp4,
		"8bpp":  bpp8,
		"16bpp": bpp16,
		"a3i5":  a3i5,
		"a5i3":  a5i3,
		"4x4c":  c4x4,
	}
	bitsPerPixel = map[textureFormat]int{
		bpp2:  2,
		bpp4:  4,
		bpp8:  8,
		bpp16: 16,
		a3i5:  8,
		a5i3:  8,
		c4x4:  2,
	}
)

func main() {
	usage := `dtex, a texture converter for Nintendo DS homebrew

Usage:
    dtex <input_filename> to (2bpp | 4bpp | 8bpp | 16bpp | a3i5 | a5i3 |
      4x4c) at <output_filename>
    dtex <input_filename> to (2bpp | 4bpp | 8bpp | 16bpp | a3i5 | a5i3 |
      4x4c) palette at <output_filename>
    dtex --help
    dtex --version

Options:
    --help     Print this message and exit
    --version  Show version number and exit
    --format <format>  Output texture format; one of 2bpp, 4bpp, 8bpp, 16bpp,
                       a3i5, a5i3, compressed`
	args, _ := docopt.Parse(usage, nil, true, "dtex 0.1.0", false)
	infile := args["<input_filename>"].(string)
	outfile := args["<output_filename>"].(string)
	format := bpp2
	for _, f := range []string{"2bpp", "4bpp", "8bpp", "16bpp", "a3i5", "a5i3", "4x4c"} {
		if args[f].(bool) {
			format = formatStrings[f]
			break
		}
	}
	converter := convertPalettedImage
	if args["palette"].(bool) {
		converter = convertPalette
	}
	convert_file(infile, outfile, format, converter)
}

func convert_file(infile, outfile string, format textureFormat, converter func(*image.Paletted, textureFormat) ([]byte, error)) error {
	src, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return err
	}
	if palimg := img.(*image.Paletted); palimg != nil {
		output, err := converter(palimg, format)
		if err != nil {
			return err
		}
		ioutil.WriteFile(outfile, output, 0644)
	} else {
		return errors.New("non-paletted image formats are not yet supported")
	}
	return nil
}

// TODO: Write the palette converter
func convertPalette(img *image.Paletted, format textureFormat) ([]byte, error) {
	return nil, errors.New("palette conversion not yet implemented")
}

func convertPalettedImage(img *image.Paletted, format textureFormat) ([]byte, error) {
	if format == bpp16 || format == c4x4 {
		return nil, errors.New("requested format conversion not implemented")
	}
	return convertToNBitImage(bitsPerPixel[format], img.Pix), nil
}

func pack(values []byte, shift int) byte {
	var mask byte = (1 << uint(shift)) - 1
	var b byte = 0
	for i, value := range values {
		b |= (value & mask) << uint(i*shift)
	}
	return b
}

func convertToNBitImage(bpp int, indexes []byte) []byte {
	pixelsPerByte := 8 / bpp
	outimg := make([]byte, len(indexes)/pixelsPerByte)
	for outpix := 0; outpix < len(indexes)/pixelsPerByte; outpix++ {
		outimg[outpix] = pack(indexes[outpix*pixelsPerByte:outpix*pixelsPerByte+pixelsPerByte], bpp)
	}
	return outimg
}
