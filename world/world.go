package world

import (
	"image"
	"image/color"
	"kar"
	"kar/arc"
	"kar/items"
	"math"
	"slices"

	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

var toRemove []ecs.Entity

var mooreNeighbors = [8]image.Point{
	{1, 0},
	{1, 1},
	{0, 1},
	{-1, 1},
	{-1, 0},
	{-1, -1},
	{0, -1},
	{1, -1},
}

type World struct {
	W, H         float64
	BlockSize    float64
	ChunkSize    vec.Vec2
	Image        *image.Gray16
	LoadedChunks []image.Point
	PlayerChunk  image.Point
}

func NewWorld(w, h float64, chunkSize vec.Vec2, blockSize float64) *World {
	terr := &World{
		Image:     GenerateWorld(int(w), int(h)),
		ChunkSize: chunkSize,
		BlockSize: blockSize,
		W:         w,
		H:         h,
	}

	return terr
}

func (tr *World) SpawnChunk(chunkCoord image.Point) {
	ChunkSizeX, ChunkSizeY := int(tr.ChunkSize.X), int(tr.ChunkSize.Y)
	for y := 0; y < ChunkSizeY; y++ {
		for x := 0; x < ChunkSizeX; x++ {
			blockX := x + (ChunkSizeX * chunkCoord.X)
			blockY := y + (ChunkSizeY * chunkCoord.Y)
			blockType := tr.Image.Gray16At(blockX, blockY).Y
			blockPos := vec.Vec2{float64(blockX), float64(blockY)}
			blockPos = blockPos.Scale(tr.BlockSize)
			if blockType != items.Air {
				arc.SpawnBlock(blockPos, blockType)
			}
		}
	}
}

func (tr *World) LoadChunks(playerPos vec.Vec2) {
	tr.LoadedChunks = GetPlayerNeighborChunks(WorldToChunk(playerPos))
	for _, coord := range tr.LoadedChunks {
		tr.SpawnChunk(coord)
	}
}

// Spawn/Destroy chunks
func (tr *World) UpdateChunks(playerPos vec.Vec2) {

	playerChunkTemp := WorldToChunk(playerPos)
	if tr.PlayerChunk != playerChunkTemp {
		tr.PlayerChunk = playerChunkTemp
		neighborChunks := GetPlayerNeighborChunks(tr.PlayerChunk)
		intersectionChunks := findChunkCoordsIntersection(neighborChunks, tr.LoadedChunks)

		for _, toLoad := range neighborChunks {
			if !slices.Contains(intersectionChunks, toLoad) {
				tr.SpawnChunk(toLoad)
			}
		}

		for _, toUnload := range tr.LoadedChunks {
			if !slices.Contains(intersectionChunks, toUnload) {
				tr.DeSpawnChunk(toUnload)
			}
		}

		tr.LoadedChunks = neighborChunks

	}

}

func (tr *World) DeSpawnChunk(chunkCoord image.Point) {

	q := arc.ItemFilter.Query(&kar.WorldECS)
	for q.Next() {
		itm := q.Get()
		e := q.Entity()
		if itm.Chunk == chunkCoord {
			b := arc.BodyMapper.Get(e)
			kar.Space.AddPostStepCallback(removeBodyPostStep, b, nil)
			toRemove = append(toRemove, e)
		}
	}

	// The world is unlocked again.
	// Actually remove the collected entities.
	for _, e := range toRemove {
		kar.WorldECS.RemoveEntity(e)
	}

	// Empty the slice, so we can reuse it in the next time step.
	toRemove = toRemove[:0]
}

func (tr *World) ChunkImage(chunkCoord image.Point) *image.Gray16 {
	chunkSizeX, chunkSizeY := int(tr.ChunkSize.X), int(tr.ChunkSize.Y)
	img := image.NewGray16(image.Rect(0, 0, chunkSizeX, chunkSizeY))
	for y := 0; y < chunkSizeY; y++ {
		for x := 0; x < chunkSizeX; x++ {
			gray := tr.Image.Gray16At(
				x+(chunkSizeX*chunkCoord.X),
				y+(chunkSizeY*chunkCoord.Y))
			img.Set(x, y, gray)
		}
	}
	return img
}
func (tr *World) EmptyBlockInChunk(chunk image.Point) (bool, image.Point) {
	chunkSizeX, chunkSizeY := int(tr.ChunkSize.X), int(tr.ChunkSize.Y)
	for y := 0; y < chunkSizeY; y++ {
		for x := 0; x < chunkSizeX; x++ {
			blockX := x + (chunkSizeX * chunk.X)
			blockY := y + (chunkSizeY * chunk.Y)
			if tr.Image.Gray16At(blockX, blockY).Y == 0 {
				return true, image.Point{blockX, blockY}
			}

		}
	}
	return false, image.Point{}
}

func (tr *World) BlockTypeAtPixelCoord(p image.Point) uint16 {
	return uint16(tr.Image.Gray16At(p.X, p.Y).Y)
}

func (tr *World) IsAir(x, y int) bool {
	if tr.Image.Gray16At(x, y).Y == items.Air {
		return true
	} else {
		return false
	}
}

func PixelToWorld(x, y int) vec.Vec2 {
	// worldPos := vec.Vec2{float64(x) * kar.BlockSize, float64(y) * kar.BlockSize}
	// return worldPos.Sub(kar.BlockCenterOffset)
	return vec.Vec2{float64(x) * kar.BlockSize, float64(y) * kar.BlockSize}
}
func WorldToPixel(pos vec.Vec2) image.Point {
	// pos = pos.Add(kar.BlockCenterOffset)
	x := math.Floor(pos.X / kar.BlockSize)
	y := math.Floor(pos.Y / kar.BlockSize)
	return image.Point{int(x), int(y)}
}
func WorldToChunk(pos vec.Vec2) image.Point {
	// pos = pos.Add(kar.BlockCenterOffset)
	return image.Point{
		int(math.Floor((pos.X / kar.ChunkSize.X) / kar.BlockSize)),
		int(math.Floor((pos.Y / kar.ChunkSize.Y) / kar.BlockSize))}
}

func ApplyColorMap(img *image.Gray16, clr map[uint16]color.RGBA) *image.RGBA {
	colored := image.NewRGBA(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			colored.SetRGBA(x, y, clr[img.Gray16At(x, y).Y])
		}
	}
	return colored
}

func findChunkCoordsIntersection(s1, s2 []image.Point) []image.Point {
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

func removeBodyPostStep(space *cm.Space, body, data interface{}) {
	space.RemoveBodyWithShapes(body.(*cm.Body))
}
