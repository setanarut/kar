package world

import (
	"image"
	"image/color"
	"kar/comp"
	"kar/engine/util"
	"kar/items"
	"kar/res"
	"kar/types"
	"slices"

	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
	"github.com/yohamta/donburi"
)

var mooreNeighbors = [8]image.Point{{1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}, {0, -1}, {1, -1}}

type World struct {
	W, H                 float64
	ChunkSize, BlockSize float64
	Image                *image.Gray
	LoadedChunks         []image.Point
}

func NewWorld(w, h, chunkSize, blockSize float64) *World {
	terr := &World{
		Image:     GenerateWorld(int(w), int(h)),
		ChunkSize: chunkSize,
		BlockSize: blockSize,
		W:         w,
		H:         h,
	}
	return terr
}

func (tr *World) BlockTypeAtPixelCoord(p image.Point) types.ItemType {
	return types.ItemType(tr.Image.GrayAt(p.X, p.Y).Y)
}

func (tr *World) SpawnChunk(chunkCoord image.Point, spawnBlock types.BlockSpawnFunc) {
	chunksize := int(tr.ChunkSize)
	for y := 0; y < chunksize; y++ {
		for x := 0; x < chunksize; x++ {
			blockX := x + (chunksize * chunkCoord.X)
			blockY := y + (chunksize * chunkCoord.Y)
			blockType := types.ItemType(tr.Image.GrayAt(blockX, blockY).Y)
			blockPos := vec.Vec2{float64(blockX), float64(blockY)}
			blockPos = blockPos.Scale(tr.BlockSize)
			if blockType != items.Air {
				spawnBlock(blockPos, chunkCoord, blockType)
			}
		}
	}
}

func (tr *World) EmptyBlockInChunk(chunkCoord image.Point) (bool, image.Point) {
	chunksize := int(tr.ChunkSize)
	for y := 0; y < chunksize; y++ {
		for x := 0; x < chunksize; x++ {
			blockX := x + (chunksize * chunkCoord.X)
			blockY := y + (chunksize * chunkCoord.Y)
			if tr.Image.GrayAt(blockX, blockY).Y == 0 {
				return true, image.Point{blockX, blockY}
			}

		}
	}
	return false, image.Point{}
}

// Spawn/Destroy chunks
func (tr *World) UpdateChunks(playerChunk image.Point, spawnBlock types.BlockSpawnFunc) {
	playerChunks := GetPlayerNeighborChunks(playerChunk)
	intersectionChunks := FindChunkCoordsIntersection(playerChunks, tr.LoadedChunks)

	for _, toLoad := range playerChunks {
		if !slices.Contains(intersectionChunks, toLoad) {
			tr.SpawnChunk(toLoad, spawnBlock)
		}
	}

	for _, toUnload := range tr.LoadedChunks {
		if !slices.Contains(intersectionChunks, toUnload) {
			tr.DeSpawnChunk(toUnload)
		}
	}

	tr.LoadedChunks = playerChunks
}

func (tr *World) DeSpawnChunk(chunkCoord image.Point) {
	comp.Item.Each(res.ECSWorld, func(e *donburi.Entry) {
		if comp.Item.Get(e).ChunkCoord == chunkCoord {
			b := comp.Body.Get(e)
			destroyBodyWithEntry(b)
		}
	})
}

func (tr *World) WorldPosToChunkCoord(worldPos vec.Vec2) image.Point {
	return util.Vec2ToPoint(worldPos.Div(tr.ChunkSize).Div(tr.BlockSize))
}

func (tr *World) PixelSpaceToWorldSpace(worldPos vec.Vec2) image.Point {
	return util.Vec2ToPoint(worldPos.Div(res.BlockSize))
}

func (tr *World) WorldSpaceToPixelSpace(pos vec.Vec2) image.Point {
	mapPos := util.Vec2ToPoint(pos).Div(int(res.BlockSize))
	return mapPos
}

func (tr *World) ChunkImage(chunkCoord image.Point) *image.Gray {
	chunksize := int(tr.ChunkSize)
	img := image.NewGray(image.Rect(0, 0, chunksize, chunksize))
	for y := 0; y < chunksize; y++ {
		for x := 0; x < chunksize; x++ {
			gray := tr.Image.GrayAt(x+(chunksize*chunkCoord.X), y+(chunksize*chunkCoord.Y))
			img.Set(x, y, gray)
		}
	}
	return img
}

func ApplyColorMap(img *image.Gray, colors map[types.ItemType]color.RGBA) *image.RGBA {
	colored := image.NewRGBA(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			blockType := types.ItemType(img.GrayAt(x, y).Y)
			colored.SetRGBA(x, y, colors[blockType])
		}
	}
	return colored
}

func FindChunkCoordsIntersection(s1, s2 []image.Point) []image.Point {
	intersection := make([]image.Point, 0)
	for _, point := range s2 {
		if slices.Contains(s1, point) {
			intersection = append(intersection, point)
		}
	}
	return intersection
}

func GetPlayerNeighborChunks(playerChunk image.Point) []image.Point {
	playerChunks := make([]image.Point, 0)
	playerChunks = append(playerChunks, playerChunk)

	for _, neighbor := range mooreNeighbors {
		playerChunks = append(playerChunks, playerChunk.Add(neighbor))
	}
	return playerChunks
}

func destroyBodyWithEntry(b *cm.Body) {
	s := b.FirstShape().Space()
	if s.ContainsBody(b) {
		e := b.UserData.(*donburi.Entry)
		e.Remove()
		s.AddPostStepCallback(removeBodyPostStep, b, false)
	}
}

func removeBodyPostStep(space *cm.Space, body, data interface{}) {
	space.RemoveBodyWithShapes(body.(*cm.Body))
}
