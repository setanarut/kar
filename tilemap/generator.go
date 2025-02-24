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
	HighestGoldLevel    float64
	LowestGoldLevel     float64
	HighestIronLevel    float64
	LowestIronLevel     float64
	HighestCoalLevel    float64
	LowestCoalLevel     float64
	HighestDiamondLevel float64
	LowestDiamondLevel  float64
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
		HighestSurfaceLevel: 10,
		LowestSurfaceLevel:  30,
		HighestCoalLevel:    32,
		LowestCoalLevel:     100,
		HighestGoldLevel:    35,
		LowestGoldLevel:     45,
		HighestIronLevel:    50,
		LowestIronLevel:     70,
		HighestDiamondLevel: 100,
		LowestDiamondLevel:  110,
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

func (g *Generator) Generate() {
	g.StoneLayer() // Base stone layer
	g.Surface()    // Fill surface (tree/dirt)
	g.CoalOreLayer(20)
	g.GoldOreLayer(20)
	g.IronOreLayer(150)
	g.DiamondOreLayer(20)
	g.FillBedrockLayer()
}

func (g *Generator) FillBedrockLayer() {
	for x := range g.Tilemap.W {
		g.Tilemap.Set(x, g.Tilemap.H-1, items.Bedrock)
	}
}
func (g *Generator) DiamondOreLayer(n int) {
	rect := image.Rect(0, int(g.Opts.HighestDiamondLevel), g.Tilemap.W, int(g.Opts.LowestDiamondLevel))
	for range n {
		x, y := g.RandomPointInRect(rect)
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x, y, items.DiamondOre)
		}
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x, y-1, items.DiamondOre)
		}
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x-1, y, items.DiamondOre)
		}
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x-1, y-1, items.DiamondOre)
		}

	}
}
func (g *Generator) CoalOreLayer(n int) {
	rect := image.Rect(0, int(g.Opts.HighestCoalLevel), g.Tilemap.W, int(g.Opts.LowestCoalLevel))
	for range n {
		x, y := g.RandomPointInRect(rect)
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x, y, items.CoalOre)
		}
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x, y-1, items.CoalOre)
		}
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x-1, y, items.CoalOre)
		}
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x-1, y-1, items.CoalOre)
		}

	}
}
func (g *Generator) GoldOreLayer(n int) {
	rect := image.Rect(0, int(g.Opts.HighestGoldLevel), g.Tilemap.W, int(g.Opts.LowestGoldLevel))
	for range n {
		x, y := g.RandomPointInRect(rect)
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x, y, items.GoldOre)
		}
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x, y-1, items.GoldOre)
		}
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x-1, y, items.GoldOre)
		}
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x-1, y-1, items.GoldOre)
		}

	}
}
func (g *Generator) IronOreLayer(n int) {
	rect := image.Rect(0, int(g.Opts.HighestIronLevel), g.Tilemap.W, int(g.Opts.LowestIronLevel))
	for range n {
		x, y := g.RandomPointInRect(rect)
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x, y, items.IronOre)
		}
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x, y-1, items.IronOre)
		}
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x-1, y, items.IronOre)
		}
		if g.Rand.Float64() < 0.5 {
			g.Tilemap.Set(x-1, y-1, items.IronOre)
		}

	}
}
func (g *Generator) StoneLayer() {
	for y := range g.Tilemap.H {
		for x := range g.Tilemap.W {
			val := mapRange(g.NoiseState.Noise2D(x, 0), -1, 1, g.Opts.LowestSurfaceLevel, g.Opts.HighestSurfaceLevel)
			if y > int(val) {
				g.Tilemap.Grid[y][x] = items.Stone
			} else {
				g.Tilemap.Grid[y][x] = items.Air
			}
		}
	}
}

func (g *Generator) Surface() {
	// surface bounds
	for y := int(g.Opts.HighestSurfaceLevel); y <= int(g.Opts.LowestSurfaceLevel); y++ {
		for x := range g.Tilemap.W {
			upperBlockID := g.Tilemap.GetID(x, y-1)
			currentBlockID := g.Tilemap.GetID(x, y)
			if upperBlockID == items.Air && currentBlockID == items.Stone {
				g.Tilemap.Set(x, y, items.GrassBlock)

				if g.Rand.Float64() < 0.15 {
					if g.Tilemap.GetID(x-1, y-1) != items.OakLog && g.Tilemap.GetID(x+1, y-1) != items.OakLog {
						g.makeTree(x, y-1)
					}
				}

				g.Tilemap.Set(x, y+1, items.Dirt)
				g.Tilemap.Set(x, y+2, items.Dirt)
				if g.Rand.Float64() < 0.9 {
					g.Tilemap.Set(x, y+3, items.Dirt)
				}
				if g.Rand.Float64() < 0.7 {
					g.Tilemap.Set(x, y+4, items.Dirt)
				}
				if g.Rand.Float64() < 0.5 {
					g.Tilemap.Set(x, y+5, items.Dirt)
				}

				if g.Rand.Float64() < 0.3 {
					g.Tilemap.Set(x, y+6, items.Dirt)
				}
				if g.Rand.Float64() < 0.1 {
					g.Tilemap.Set(x, y+7, items.Dirt)
				}
			}
		}
	}

}

func (g *Generator) makeTree(x, y int) {

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

func (g *Generator) SetSeed(seed int) {
	g.PCG.Seed(uint64(seed), 1024)
	g.NoiseState.Seed = seed
}

func (g *Generator) RandomPointInRect(rect image.Rectangle) (int, int) {
	return rect.Min.X + g.Rand.IntN(rect.Dx()), rect.Min.Y + g.Rand.IntN(rect.Dy())
}
