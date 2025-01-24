package tilemap

import (
	"image"
	"kar/items"
	"math/rand/v2"

	"github.com/setanarut/fastnoise"
)

type WorldOpts struct {
	Seed                int
	SurfaceFlatness     float64
	HighestSurfaceLevel float64
	LowestSurfaceLevel  float64
}

type Generator struct {
	NoiseState *fastnoise.State[float64]
	PCG        *rand.PCG
	Rand       *rand.Rand
	Opts       WorldOpts
	Tilemap    *TileMap
}

func DefaultWorldOpts() WorldOpts {
	return WorldOpts{
		SurfaceFlatness:     0,
		HighestSurfaceLevel: 0,
		LowestSurfaceLevel:  50,
	}
}

func NewGenerator(t *TileMap) *Generator {
	g := &Generator{
		NoiseState: fastnoise.New[float64](),
		Opts:       DefaultWorldOpts(),
		Tilemap:    t,
	}
	g.PCG = rand.NewPCG(0, 1024)
	g.Rand = rand.New(g.PCG)
	g.NoiseState.Seed = 0
	return g
}

func (g *Generator) BlockState(x, y int) uint16 {
	val := mapRange(g.NoiseState.Noise2D(x, 0), -1, 1, g.Opts.LowestSurfaceLevel, g.Opts.HighestSurfaceLevel)
	if y > int(val) {
		return items.Stone
	} else {
		return items.Air
	}
}

func (g *Generator) SetSeed(seed int) {
	g.PCG.Seed(uint64(seed), 1024)
	g.NoiseState.Seed = seed
}

func (g *Generator) Generate(tm *TileMap) {
	for y := range tm.H {
		for x := range tm.W {
			tm.Grid[y][x] = g.BlockState(x, y)
		}
	}
	g.MakeSurface(tm)
}

func (g *Generator) MakeSurface(tm *TileMap) {
	rect := g.GetSurfaceBounds()
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			upperBlockID := tm.Get(x, y-1)
			currentBlockID := tm.Get(x, y)
			if upperBlockID == items.Air && currentBlockID == items.Stone {
				tm.Set(x, y, items.GrassBlock)

				if g.Rand.Float64() < 0.15 {
					if tm.Get(x-1, y-1) != items.OakLog && tm.Get(x+1, y-1) != items.OakLog {
						g.MakeTree(x, y-1)
					}
				}

				tm.Set(x, y+1, items.Dirt)
				tm.Set(x, y+2, items.Dirt)
				if g.Rand.Float64() < 0.9 {
					tm.Set(x, y+3, items.Dirt)
				}
				if g.Rand.Float64() < 0.7 {
					tm.Set(x, y+4, items.Dirt)
				}
				if g.Rand.Float64() < 0.5 {
					tm.Set(x, y+5, items.Dirt)
				}

				if g.Rand.Float64() < 0.3 {
					tm.Set(x, y+6, items.Dirt)
				}
				if g.Rand.Float64() < 0.1 {
					tm.Set(x, y+7, items.Dirt)
				}
			}
		}
	}

}

func (g *Generator) MakeTree(x, y int) {

	if g.Rand.Float64() < 0.8 {
		g.Tilemap.Set(x, y-7, items.OakLeaves)
		g.Tilemap.Set(x+1, y-6, items.OakLeaves)
		g.Tilemap.Set(x, y-6, items.OakLeaves)
		g.Tilemap.Set(x-1, y-6, items.OakLeaves)
		g.Tilemap.Set(x+1, y-5, items.OakLeaves)
		g.Tilemap.Set(x, y-5, items.OakLeaves)
		g.Tilemap.Set(x-1, y-5, items.OakLeaves)
		g.Tilemap.Set(x+1, y-4, items.OakLeaves)
		g.Tilemap.Set(x, y-4, items.OakLeaves)
		g.Tilemap.Set(x-1, y-4, items.OakLeaves)
		g.Tilemap.Set(x, y-3, items.OakLog)
		g.Tilemap.Set(x, y-2, items.OakLog)
		g.Tilemap.Set(x, y-1, items.OakLog)
		g.Tilemap.Set(x, y, items.OakLog)
	} else {
		g.Tilemap.Set(x, y-3, items.OakLeaves)
		g.Tilemap.Set(x-1, y-2, items.OakLeaves)
		g.Tilemap.Set(x+1, y-2, items.OakLeaves)
		g.Tilemap.Set(x, y-2, items.OakLeaves)
		g.Tilemap.Set(x, y-1, items.OakLog)
		g.Tilemap.Set(x, y, items.OakLog)

	}

}

func mapRange(v, a, b, c, d float64) float64 {
	return (v-a)/(b-a)*(d-c) + c
}

func (g *Generator) GetSurfaceBounds() image.Rectangle {
	return image.Rect(
		0,
		int(g.Opts.HighestSurfaceLevel),
		g.Tilemap.W,
		int(g.Opts.LowestSurfaceLevel)+1,
	)
}
