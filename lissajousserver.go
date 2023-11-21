package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

var palette = []color.Color{color.Black, color.RGBA{0x00, 0x99, 0x00, 0xff}}

func main() {

	handler := func(w http.ResponseWriter, r *http.Request) {
		lissajous(w)
	}
	http.HandleFunc("/", handler)

    log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func lissajous(out io.Writer) {
    const (
        cycles  = 5     // number of complete x oscillator revolutions
        res     = 0.001 // angular resolution
        size    = 100   // image canvas covers [-size..+size]
        nframes = 64    // number of animation frames
        delay   = 8     // delay between frames in 10ms units
    )
	if len(os.Args) > 1 {
		changePalette(strings.ToLower(os.Args[1]), 1)
		if len(os.Args) > 2 {
			changePalette(strings.ToLower(os.Args[2]), 0)
		}
	}
    freq := rand.Float64() * 3.0 // relative frequency of y oscillator
    anim := gif.GIF{LoopCount: nframes}
    phase := 0.0 // phase difference
    for i := 0; i < nframes; i++ {
        rect := image.Rect(0, 0, 2*size+1, 2*size+1)
        img := image.NewPaletted(rect, palette)
        for t := 0.0; t < cycles*2*math.Pi; t += res {
            x := math.Sin(t)
            y := math.Sin(t*freq + phase)
            img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
                1)
        }
        phase += 0.1
        anim.Delay = append(anim.Delay, delay)
        anim.Image = append(anim.Image, img)
    }
    gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

func changePalette(option string, index int) {
	switch option {
	case "black":
		palette[index] = color.Black
	case "blue":
		palette[index] = color.RGBA{0x00, 0x00, 0x99, 0xff}
	case "cyan":
		palette[index] = color.RGBA{0x00, 0x99, 0x99, 0xff}
	case "green":
		palette[index] = color.RGBA{0x00, 0x99, 0x00, 0xff}
	case "magenta":
		palette[index] = color.RGBA{0x99, 0x00, 0x99, 0xff}
	case "red":
		palette[index] = color.RGBA{0x99, 0x00, 0x00, 0xff}
	case "yellow":
		palette[index] = color.RGBA{0x99, 0x99, 0x00, 0xff}
	case "white":
		palette[index] = color.White
	}

}
