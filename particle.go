package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"math/rand"
	"os"
)

const ()

type (
	Particle struct {
		X     float64
		Y     float64
		Color color.RGBA
		Type  int
	}
)

var (
	particles  []Particle
	canvas     *image.RGBA
	imageCount = 0
	stepCount  = 0

	ParticleCount = 1920 * 768
	ImageWidth    = 1920
	ImageHeight   = 1080
	Density       = 50.0
	Iterations    = 16384
	OutDir        = "out"

	jpegOptions = jpeg.Options{Quality: 100}
)

func (p *Particle) Move() bool {
	x := int(math.Round(p.X))
	y := int(math.Round(p.Y))

	bgColor := color.RGBA{255, 255, 255, 255}

	moved := false
	//
	if p.Type == 0 && stepCount%2 == 0 {
		c := canvas.At((x+1)%ImageWidth, y)
		if c == bgColor {
			p.X += 1
			moved = true
		}
		if x >= ImageWidth {
			p.X = 0
		}
	}
	if p.Type == 1 && stepCount%2 == 1 {
		c := canvas.At(x, (y+1)%ImageHeight)
		if c == bgColor {
			p.Y += 1
			moved = true
		}

		if y >= ImageHeight {
			p.Y = 0
		}
	}
	return moved

}

func main() {
	log.Println("hi")

	w := flag.Int("width", 1920, "the width of the rendered image")
	h := flag.Int("height", 1080, "the height of the rendered image")
	d := flag.Float64("density", 50.0, "the density of the particles")
	i := flag.Int("iterations", 16384, "number of iterations")
	o := flag.String("out", "out", "The directory to save the files")
	// configurable colors?
	flag.Parse()

	ImageWidth = *w
	ImageHeight = *h
	Density = *d
	Iterations = *i
	OutDir = *o

	ParticleCount = int(Density / 100.0 * float64(ImageWidth*ImageHeight))

	r := image.Rect(0, 0, ImageWidth, ImageHeight)
	canvas = image.NewRGBA(r)
	particles = make([]Particle, ParticleCount)
	for i := 0; i < ParticleCount; i += 1 {
		partColor := color.RGBA{255, 0, 0, 255}
		if i%2 == 0 {
			partColor = color.RGBA{0, 0, 255, 255}
		}
		particles[i] = Particle{
			X:     rand.Float64() * float64(ImageWidth),
			Y:     rand.Float64() * float64(ImageHeight),
			Color: partColor,
			Type:  i % 2,
		}
	}

	for i := 0; i < Iterations; i += 1 {
		log.Printf("%d finished", i)
		moveCount := moveParticles()
		if moveCount < 1 && i > 0 {
			break
		}
		renderBackground()
		renderParticles()
		saveImage()
		stepCount += 1
	}

}

func moveParticles() int {
	moveCount := 0
	for i := 0; i < len(particles); i += 1 {
		moved := particles[i].Move()
		if moved {
			moveCount += 1
		}
	}
	return moveCount
}

func renderParticles() {
	for i := 0; i < len(particles); i += 1 {
		x := int(math.Round(particles[i].X))
		y := int(math.Round(particles[i].Y))
		canvas.Set(x, y, particles[i].Color)
	}
}

func renderBackground() {
	for x := 0; x < ImageWidth; x += 1 {
		for y := 0; y < ImageHeight; y += 1 {
			canvas.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}
}

func saveImage() {
	f, err := os.Create(fmt.Sprintf("%s/%dx%d-at-%0.2f-%06d-image.jpg", OutDir, ImageWidth, ImageHeight, Density, imageCount))
	if err != nil {
		log.Fatal(err)
	}

	if err := jpeg.Encode(f, canvas, &jpegOptions); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	imageCount += 1
}
