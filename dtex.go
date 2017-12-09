package main

import (
	"errors"
	//"fmt"
	"github.com/docopt/docopt-go"
	"image"
	//"image/color"
	_ "image/png"
	"io/ioutil"
	"os"
)

func main() {
	usage := `dtex, a texture converter for Nintendo DS homebrew

Usage:
    dtex <input_filename> to (2bpp | 4bpp | 8bpp | 16bpp | a3i5 | a5i3 |
      4x4c ) at <output_filename>
    dtex <input_filename> to (2bpp | 4bpp | 8bpp | 16bpp | a3i5 | a5i3 |
      4x4c ) palette at <output_filename>
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
	format := ""
	for _, f := range []string{"2bpp", "4bpp", "8bpp", "16bpp", "a3i5", "a5i3", "4x4c"} {
		if args[f].(bool) {
			format = f
			break
		}
	}
	convert_pal := args["palette"].(bool)
	convert_file(infile, outfile, format, convert_pal)
}

func convert_file(infile, outfile, format string, convert_pal bool) error {
	src, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err !=  nil {
		return err
	}
	if palimg := img.(image.PalettedImage); palimg != nil {
		if convert_pal {
			ioutil.WriteFile(outfile, convert_palette(palimg, format), 0644)
		} else {
			ioutil.WriteFile(outfile, convert_image(palimg, format), 0644)
		}
	} else {
		return errors.New("non-paletted image formats are not yet supported")
	}
	return nil
}

// TODO: Write the palette converter
func convert_palette(img image.PalettedImage, format string) []byte {
	return make([]byte, 0)
	//return make([]byte, len(img.ColorModel().(color.Palette)))
}

// TODO: Write the image converter
func convert_image(img image.PalettedImage, format string) []byte {
	return make([]byte, 0)
}
