package tilemap

import (
	"kar/engine/mathutil"
	"kar/items"
	"math/rand/v2"
	"time"

	"github.com/setanarut/fastnoise"
)

var ns = fastnoise.New[float64]()
var random *rand.Rand

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

func Generate(tm *TileMap) {

	// seeds
	ns.Seed = int(time.Now().UnixNano())
	random = rand.New(rand.NewPCG(42, uint64(ns.Seed)))

	ns.Frequency = 0.0003
	ns.Octaves = 3
	ns.Lacunarity = 5
	ns.NoiseType(fastnoise.OpenSimplex2)
	ns.Gain = 1

	for y := 0; y < tm.H; y++ {
		for x := 0; x < tm.W; x++ {
			tm.Grid[y][x] = BlockState(x, y, tm.W, tm.H)
		}
	}

	// Make grass blocks
	for y := 0; y < tm.H/4; y++ {
		for x := 0; x < tm.W; x++ {
			if tm.Grid[y][x] == items.Dirt {
				upperBlock := tm.Grid[y-1][x]
				if upperBlock == items.Air {
					tm.Grid[y][x] = items.GrassBlock
				}
			}
		}
	}

	// Make trees
	min := 10
	max := tm.W - 10
	treeCount := 40 // max try count
	for range treeCount {
		x := min + random.IntN(max-min+1)
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

func BlockState(x, y, w, h int) uint16 {
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
				return items.Obsidian
			}
		}
		return items.Obsidian
	}
	if depth == 510 {
		if random.Float64() > 0.5 {
			return items.Obsidian
		} else {
			return items.Bedrock
		}
	}
	if depth == float64(h) {
		return items.Bedrock
	}
	return items.Air
}
