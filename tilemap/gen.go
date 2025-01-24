package tilemap

import (
	"image"
	"kar/items"
	"math/rand"

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
	Rand       *rand.Rand
	Opts       WorldOpts
}

func DefaultWorldOpts() WorldOpts {
	return WorldOpts{
		SurfaceFlatness:     0,
		HighestSurfaceLevel: 0,
		LowestSurfaceLevel:  50,
	}
}

func NewGenerator() *Generator {
	g := &Generator{
		NoiseState: fastnoise.New[float64](),
		Rand:       rand.New(rand.NewSource(int64(0))),
		Opts:       DefaultWorldOpts(),
	}
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
	g.Rand.Seed(int64(seed))
	g.NoiseState.Seed = seed
}

func (g *Generator) Generate(tm *TileMap) {
	for y := 0; y < tm.H; y++ {
		for x := 0; x < tm.W; x++ {
			tm.Grid[y][x] = g.BlockState(x, y)
		}
	}
	g.MakeDirt(tm)
	// g.MakeTrees(tm)
}

func (g *Generator) MakeDirt(tm *TileMap) {
	rect := image.Rect(0, int(g.Opts.HighestSurfaceLevel), tm.W, int(g.Opts.LowestSurfaceLevel)+1)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			upperBlockID := tm.Get(x, y-1)
			currentBlockID := tm.Get(x, y)
			if upperBlockID == items.Air && currentBlockID == items.Stone {
				tm.Set(x, y, items.Dirt)
				tm.Set(x, y+1, items.Dirt)
				tm.Set(x, y+2, items.Dirt)
				tm.Set(x, y+3, items.Dirt)
			}
		}
	}

}

// Make trees
func (g *Generator) MakeTrees(tm *TileMap) {
	min := 10
	max := tm.W - 10
	treeCount := 40 // max try count
	for i := 0; i < treeCount; i++ {
		x := min + g.Rand.Intn(max-min+1)
		for y := 0; y < tm.H; y++ {
			up := tm.Get(x, y)
			down := tm.Get(x, y+1)
			if up == items.Air {
				if down == items.Dirt || down == items.GrassBlock {
					MakeTree(x, y, tm)
					break
				}
			}
		}
	}
}

func MakeTree(x, y int, tm *TileMap) {
	tm.Set(x, y, items.OakLog)
	tm.Set(x, y-1, items.OakLog)
	tm.Set(x, y-2, items.OakLog)
	tm.Set(x, y-3, items.OakLog)
	tm.Set(x-1, y-4, items.OakLeaves)
	tm.Set(x, y-4, items.OakLeaves)
	tm.Set(x+1, y-4, items.OakLeaves)
	tm.Set(x-1, y-5, items.OakLeaves)
	tm.Set(x, y-5, items.OakLeaves)
	tm.Set(x+1, y-5, items.OakLeaves)
	tm.Set(x-1, y-6, items.OakLeaves)
	tm.Set(x, y-6, items.OakLeaves)
	tm.Set(x+1, y-6, items.OakLeaves)
}

func mapRange(v, a, b, c, d float64) float64 {
	return (v-a)/(b-a)*(d-c) + c
}
