package tilemap

import (
	"kar/items"
	"math/rand"
	"time"

	"github.com/setanarut/fastnoise"
)

var ns = fastnoise.New[float64]()
var random *rand.Rand

var high float64 = 50
var low float64 = 70

func BlockState(x, y, w, h int) uint16 {
	val := mapRange(ns.Noise2D(x, 0), -1, 1, low, high)
	if y > int(val) {
		return items.Stone
	} else {
		return items.Air
	}
}

func Generate(tm *TileMap) {
	// seeds
	ns.Seed = int(time.Now().UnixNano())
	random = rand.New(rand.NewSource(int64(ns.Seed)))

	ns.Frequency = 0.01
	// ns.Octaves = 3
	// ns.Lacunarity = 5
	ns.NoiseType(fastnoise.OpenSimplex2)
	// ns.Gain = 1

	for y := 0; y < tm.H; y++ {
		for x := 0; x < tm.W; x++ {
			tm.Grid[y][x] = BlockState(x, y, tm.W, tm.H)
		}
	}

	// MakeGrass(tm)
	// MakeTrees(tm)
}

func MakeGrass(tm *TileMap) {
	for y := 1; y < tm.H/4; y++ {
		for x := 0; x < tm.W; x++ {
			if tm.Grid[y][x] == items.Dirt {
				upperBlock := tm.Grid[y-1][x]
				if upperBlock == items.Air {
					tm.Grid[y][x] = items.GrassBlock
				}
			}
		}
	}
}

// Make trees
func MakeTrees(tm *TileMap) {
	min := 10
	max := tm.W - 10
	treeCount := 40 // max try count
	for i := 0; i < treeCount; i++ {
		x := min + random.Intn(max-min+1)
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
