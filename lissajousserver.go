// Lissajous Server is a server that produces lissajous images.
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
	"net/url"
	"strconv"
)

var palette = []color.Color{color.Black, color.RGBA{0x00, 0x99, 0x00, 0xff}}

func main() {

	handler := func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		lissajous(w, params)
	}
	http.HandleFunc("/", handler)

    log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func lissajous(out io.Writer, params url.Values) {

	// number of complete x oscillator revolutions
	cycles, err := strconv.Atoi(params.Get("cycles"))
	if err != nil || cycles < 1 {
		cycles  = 5     
	}

	// image canvas covers [-size..+size]
	size, err := strconv.Atoi(params.Get("size"))
	if err != nil || size < 1 {
		size  = 300     
	}

	// number of animation frames
	nframes, err := strconv.Atoi(params.Get("nframes"))
	if err != nil || nframes < 1 {
		nframes = 64     
	}

	// delay between frames in 10ms units
	delay, err := strconv.Atoi(params.Get("delay"))
	if err != nil || delay < 1 {
		delay = 8     
	}
	
    const res     = 0.001 // angular resolution 

    freq := rand.Float64() * 3.0 // relative frequency of y oscillator
    anim := gif.GIF{LoopCount: nframes}
    phase := 0.0 // phase difference
    for i := 0; i < nframes; i++ {
        rect := image.Rect(0, 0, 2*size+1, 2*size+1)
        img := image.NewPaletted(rect, palette)
        for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
            x := math.Sin(t)
            y := math.Sin(t*freq + phase)
            img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5),
                1)
        }
        phase += 0.1
        anim.Delay = append(anim.Delay, delay)
        anim.Image = append(anim.Image, img)
    }
    gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
