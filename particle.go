package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
)

const (
	ParticleCount = 1920 * 768
	IMAGE_WIDTH   = 1920
	IMAGE_HEIGHT  = 1080
)

type (
	Particle struct {
		X     float64
		Y     float64
		Vx    float64
		Vy    float64
		Color color.RGBA
		Type  int
	}
)

var (
	particles  []Particle
	canvas     *image.RGBA
	imageCount = 0
	stepCount  = 0
)

func (p *Particle) Move() {
	if p.X+p.Vx < 0 || p.X+p.Vx > IMAGE_WIDTH {
		p.Vx = p.Vx * -1.0
	}
	if p.Y+p.Vy < 0 || p.Y+p.Vy > IMAGE_HEIGHT {
		p.Vy = p.Vy * -1.0
	}

	// p.X += p.Vx
	// p.Y += p.Vy

	x := int(math.Round(p.X))
	y := int(math.Round(p.Y))

	bgColor := color.RGBA{255, 255, 255, 255}

	if p.Type == 0 && stepCount%2 == 0 {
		c := canvas.At(x+1, y)
		if c == bgColor {
			p.X += 1
		}
	}
	if p.Type == 1 && stepCount%2 == 1 {
		c := canvas.At(x, y+1)
		if c == bgColor {
			p.Y += 1
		}
	}
	if y+1 >= IMAGE_HEIGHT {
		p.Y = 0
	}
	if x+1 >= IMAGE_WIDTH {
		p.X = 0
	}
}

func main() {
	log.Println("hi")
	r := image.Rect(0, 0, IMAGE_WIDTH, IMAGE_HEIGHT)
	canvas = image.NewRGBA(r)
	particles = make([]Particle, ParticleCount)
	for i := 0; i < ParticleCount; i += 1 {
		partColor := color.RGBA{255, 0, 0, 255}
		if i%2 == 0 {
			partColor = color.RGBA{0, 0, 255, 255}
		}
		particles[i] = Particle{
			X:     rand.Float64() * float64(IMAGE_WIDTH),
			Y:     rand.Float64() * float64(IMAGE_HEIGHT),
			Vx:    (rand.Float64() - 0.5) * 1,
			Vy:    (rand.Float64() - 0.5) * 1,
			Color: partColor,
			Type:  i % 2,
		}
	}
	for i := 0; i < 16384; i += 1 {
		moveParticles()
		renderBackground()
		renderParticles()
		saveImage()
		stepCount += 1
	}

}

func moveParticles() {
	for i := 0; i < len(particles); i += 1 {
		particles[i].Move()
	}
}

func renderParticles() {
	for i := 0; i < len(particles); i += 1 {
		x := int(math.Round(particles[i].X))
		y := int(math.Round(particles[i].Y))
		canvas.Set(x, y, particles[i].Color)
	}
}

func renderBackground() {
	for x := 0; x < IMAGE_WIDTH; x += 1 {
		for y := 0; y < IMAGE_HEIGHT; y += 1 {
			canvas.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}
}

func saveImage() {
	f, err := os.Create(fmt.Sprintf("out/%05d-image.png", imageCount))
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, canvas); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	imageCount += 1
}
