package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"fmt"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	term "github.com/gookit/color"
	cli "github.com/urfave/cli/v2"
)

// `.-~_,:!^+<>=;*/()?{}[]7123#4$5%980@
// 712345908
// `.-,:!^+<=*/(?{[#$%@

func ToText(filename string, width int, chars []rune, color bool, use_gaussian bool) string {
	data, err := imgio.Open(filename)
	if err != nil {
		log.Fatalln("Can't open this file", err)
	}
	w := data.Bounds().Dx()
	h := data.Bounds().Dy()
	scale := float64(w) / float64(h)
	if width == 0 {
		width = w
	}
	var res_height = int(float64(width) / (scale + 1))
	if use_gaussian {
		data = transform.Resize(data, width, res_height, transform.Gaussian)
	} else {
		data = transform.Resize(data, width, res_height, transform.Linear)
	}

	if data == nil {
		log.Fatalln("Can't resize image")
	}
	var result string
	var base int
	var string_chars = make([]string, len(chars))
	for a := range chars {
		string_chars[a] = string(chars[a])
	}
	if 256%len(chars) == 0 {
		base = 256 / len(chars)
	} else {
		base = 256/len(chars) + 1
	}
	for a := 0; a < res_height; a++ {
		for b := 0; b < width; b++ {
			r, g, b, _ := data.At(b, a).RGBA()
			gray := int((r + g + b) / 768)
			red, green, blue := uint8(r), uint8(g), uint8(b)
			ch := string_chars[gray/base]
			if color {
				result += term.RGB(red, green, blue).Sprint(ch)
			} else {
				result += ch
			}
		}
		result += "\n"
	}
	return result
}

func CreateFile() {
}

func Execute(ctx *cli.Context) error {
	results := make([]string, 0)

	str := []rune(ctx.String("chars"))
	w := ctx.Int("width")
	outfile := ctx.String("outfile")
	color := ctx.Bool("color")
	if outfile != "" && color {
		fmt.Println("Can't use color with output files")
		os.Exit(1)
	}
	if ctx.Args().Len() > 1 {
		for _, file := range ctx.Args().Slice() {
			results = append(results, ToText(file, w, str, color, false))
		}
	} else {
		results = append(results, ToText(ctx.Args().Slice()[0], w, str, color, true))
	}

	if outfile == "" {
		for _, res := range results {
			fmt.Println(res)
		}
	} else {
		for n, res := range results {
			var name, ext string
			var file *os.File
			var err error
			split := strings.Split(outfile, ".")
			if len(split) > 1 {
				name = strings.Join(split[:len(split)-1], ".")
				ext = split[len(split)-1]
			}
			if n == 0 {
				file, err = os.Create(name + "." + ext)
			} else {
				file, err = os.Create(name + strconv.Itoa(n+1) + "." + ext)
			}
			if err != nil {
				return err
			}
			_, err = file.WriteString(res)
			if err != nil {
				return err
			}
			file.Close()
		}
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:  "imager",
		Usage: "this tool can be used to convert an image to text art",
		Action: func(ctx *cli.Context) error {
			return Execute(ctx)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "chars",
				Value: ".,^!&%$@",
				Usage: "set your own `symbols` to make an image using them. They would be better ordered by 'brightness' or, simply put, the size of area they occupy within the text, from the lower to the higher values",
			},
			&cli.StringFlag{
				Name:        "outfile",
				Value:       "",
				Usage:       "specify `filename` in which the text should be written to",
				DefaultText: "stdout",
			},
			&cli.IntFlag{
				Name:        "width",
				Value:       0,
				DefaultText: "original size",
				Usage:       "set the text line `length` (the image scales to keep aspect ratio)",
			},
			&cli.BoolFlag{
				Name:  "color",
				Value: false,
				Usage: "set this flag to use colors (it's only supported with stdout output of course)",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
