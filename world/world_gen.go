package world

import (
	"image"
	"image/color"
	"kar/engine/mathutil"
	"kar/items"
	"math/rand/v2"

	"github.com/setanarut/fastnoise"
)

var ns = fastnoise.New[float64]()
var random *rand.Rand

func GenerateWorld(w, h int) *image.Gray16 {

	// seeds
	seed := mathutil.RandRangeInt(0, 1000000)
	ns.Seed = seed
	random = rand.New(rand.NewPCG(42, uint64(seed)))

	ns.Frequency = 0.0003
	ns.Octaves = 3
	ns.Lacunarity = 5
	ns.NoiseType(fastnoise.OpenSimplex2)
	ns.Gain = 1

	img := image.NewGray16(image.Rect(0, 0, w, h))
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.SetGray16(x, y, color.Gray16{BlockState(x, y)})
		}
	}

	// ikinci geçiş
	for y := 0; y < 200; y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if img.Gray16At(x, y).Y == items.Dirt {
				upperBlock := img.Gray16At(x, y-1).Y
				if upperBlock == items.Air {
					// rn := mathutil.RandRangeInt(1, 26)
					// img.SetGray16(x, y, color.Gray16{uint16(rn)})
					img.SetGray16(x, y, color.Gray16{items.GrassBlock})
				}
			}
		}
	}
	return img
}

func BlockState(x, y int) uint16 {
	depth := float64(y)
	amp := 60.0
	ns.FractalType(fastnoise.FractalFBm)
	n1 := 10 + amp + ns.Noise2D(x, 0)*amp
	n2 := 20 + amp + ns.Noise2D(x, y)*amp

	if depth >= n1 && depth <= n2 {
		grad := mathutil.MapRange(float64(y), n1, n2, 0, 1)

		if (random.Float64()+grad)/2 > 0.8 {
			return items.Stone
		} else {
			return items.Dirt
		}
		// return items.Dirt
	}
	if depth >= n2 && depth <= amp*4 {
		return items.Stone
	}

	if depth >= n2 && depth <= 250 {
		return items.Stone
	}
	if depth >= 251 && depth <= 509 {
		if depth >= 251 && depth < 251+3 {
			if random.Float64() > 0.5 {
				return items.Stone
			} else {
				return items.Deepslate
			}
		}
		return items.Deepslate
	}
	if depth == 510 {
		if random.Float64() > 0.5 {
			return items.Deepslate
		} else {
			return items.Bedrock
		}
	}
	if depth == 511 {
		return items.Bedrock
	}
	return items.Air
}
