package world

import (
	"image"
	"image/color"
	"kar/engine/mathutil"
	"kar/items"
	"kar/types"
	"math/rand/v2"

	"github.com/setanarut/fastnoise"
)

var ns = fastnoise.New[float64]()
var random *rand.Rand

func GenerateWorld(w, h int) *image.Gray {

	// seeds
	seed := mathutil.RandRangeInt(0, 1000000)
	ns.Seed = seed
	random = rand.New(rand.NewPCG(42, uint64(seed)))

	ns.Frequency = 0.0003
	ns.Octaves = 3
	ns.Lacunarity = 5
	ns.NoiseType(fastnoise.OpenSimplex2)
	ns.Gain = 1

	img := image.NewGray(image.Rect(0, 0, w, h))
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.SetGray(x, y, gray(BlockState(x, y)))
		}
	}

	// ikinci geçiş
	for y := 0; y < 200; y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			blockType := types.ItemType(img.GrayAt(x, y).Y)
			if blockType == items.Dirt {
				upperBlock := types.ItemType(img.GrayAt(x, y-1).Y)
				if upperBlock == items.Air {
					img.SetGray(x, y, gray(items.Grass))
				}
			}
		}
	}

	return img
}

func BlockState(x, y int) types.ItemType {
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
				return items.DeepSlateStone
			}
		}
		return items.DeepSlateStone
	}
	if depth == 510 {
		if random.Float64() > 0.5 {
			return items.DeepSlateStone
		} else {
			return items.DeepSlateGoldOre
		}
	}
	if depth == 511 {
		return items.GoldOre
	}
	return items.Air
}

func gray(item types.ItemType) color.Gray {
	return color.Gray{uint8(item)}
}
