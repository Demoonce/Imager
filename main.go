package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"fmt"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	cli "github.com/urfave/cli/v2"
)

func ToText(filename string, width int, chars []rune) string {
	data, err := imgio.Open(filename)
	if err != nil {
		log.Fatalln("Can't open this file")
	}
	w := data.Bounds().Dx()
	h := data.Bounds().Dy()
	scale := float64(w) / float64(h)
	if width == 0 {
		width = w
	}
	var res_height = int(float64(width) / (scale + 1))
	data = transform.Resize(data, width, res_height, transform.NearestNeighbor)
	if data == nil {
		log.Fatalln("Can't resize image")
	}
	var result = make([]rune, 0, width*res_height+h)
	var base int
	if 256%len(chars) == 0 {
		base = 256 / len(chars)
	} else {
		base = 256/len(chars) + 1
	}
	for a := 0; a < res_height; a++ {
		for b := 0; b < width; b++ {
			r, g, b, _ := data.At(b, a).RGBA()
			gray := int((r + g + b) / 768)
			result = append(result, chars[int(gray/base)])
		}
		result = append(result, '\n')
	}
	return string(result)
}

func Execute(ctx *cli.Context) error {
	results := []string{}

	str := []rune(ctx.String("chars"))
	w := ctx.Int("width")
	outfile := ctx.String("outfile")
	for _, file := range ctx.Args().Slice() {
		results = append(results, ToText(file, w, str))
	}

	if outfile == "" {
		for _, res := range results {
			fmt.Println(res)
		}
	} else {
		for n, res := range results {
			var name, ext string
			split := strings.Split(outfile, ".")
			if len(split) > 1 {
				name = strings.Join(split[:len(split)-1], ".")
				ext = split[len(split)-1]
			}
			file, err := os.Create(name + strconv.Itoa(n+1) + "." + ext)
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
				Usage: "set your own `symbols` to print an image using them. They would be better ordered by 'brightness' or, simply put, an area they occupy within a text, from the lower to the higher values",
			},
			&cli.StringFlag{
				Name:        "outfile",
				Value:       "",
				Usage:       "specify filename to write text to `filename`",
				DefaultText: "stdout",
			},
			&cli.IntFlag{
				Name:        "width",
				Value:       0,
				DefaultText: "original size",
				Usage:       "set result text line `length` (image scales to keep aspect ratio)",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
