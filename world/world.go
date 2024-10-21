package world

import (
	"image"
	"image/color"
	"kar"
	"kar/arche"
	"kar/comp"
	"kar/engine/util"
	"kar/items"
	"slices"

	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
	"github.com/yohamta/donburi"
)

var mooreNeighbors = [8]image.Point{{1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}, {0, -1}, {1, -1}}

type World struct {
	W, H         float64
	BlockSize    float64
	ChunkSize    vec.Vec2
	Image        *image.Gray16
	LoadedChunks []image.Point
}

func NewWorld(worldW, worldH float64, chunkSize vec.Vec2, blockSize float64) *World {
	terr := &World{
		Image:     GenerateWorld(int(worldW), int(worldH)),
		ChunkSize: chunkSize,
		BlockSize: blockSize,
		W:         worldW,
		H:         worldH,
	}
	return terr
}

func (tr *World) SpawnChunk(s *cm.Space, w donburi.World, chunkCoord image.Point) {
	ChunkSizeX, ChunkSizeY := int(tr.ChunkSize.X), int(tr.ChunkSize.Y)
	for y := 0; y < ChunkSizeY; y++ {
		for x := 0; x < ChunkSizeX; x++ {
			blockX := x + (ChunkSizeX * chunkCoord.X)
			blockY := y + (ChunkSizeY * chunkCoord.Y)
			blockType := tr.Image.Gray16At(blockX, blockY).Y
			blockPos := vec.Vec2{float64(blockX), float64(blockY)}
			blockPos = blockPos.Scale(tr.BlockSize)
			if blockType != items.Air {
				arche.SpawnBlock(s, w, blockPos, blockType)

				// if resources.Cam != nil {
				// x, y := resources.Cam.Target()
				// camPos := vec.Vec2{x, y}
				// if blockPos.Distance(camPos) < 10000 {
				// arche.SpawnBlock(blockPos, blockType)
				// }
				// }
			}
		}
	}
}

// Spawn/Destroy chunks
func (tr *World) UpdateChunks(s *cm.Space, w donburi.World, playerChunk image.Point) {
	playerChunks := GetPlayerNeighborChunks(playerChunk)
	intersectionChunks := findChunkCoordsIntersection(playerChunks, tr.LoadedChunks)

	for _, toLoad := range playerChunks {
		if !slices.Contains(intersectionChunks, toLoad) {
			tr.SpawnChunk(s, w, toLoad)
		}
	}

	for _, toUnload := range tr.LoadedChunks {
		if !slices.Contains(intersectionChunks, toUnload) {
			tr.DeSpawnChunk(w, toUnload)
		}
	}

	tr.LoadedChunks = playerChunks
}

func (tr *World) DeSpawnChunk(w donburi.World, chunkCoord image.Point) {
	comp.Item.Each(w, func(e *donburi.Entry) {
		if comp.Item.Get(e).Chunk == chunkCoord {
			b := comp.Body.Get(e)
			destroyBodyWithEntry(b)
		}
	})
}

func (tr *World) WorldPosToChunkCoord(worldPos vec.Vec2) image.Point {
	return image.Point{
		int((worldPos.X / tr.ChunkSize.X) / tr.BlockSize),
		int((worldPos.Y / tr.ChunkSize.Y) / tr.BlockSize)}
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

func PixelToWorldSpace(x, y int) vec.Vec2 {
	return vec.Vec2{float64(x) * kar.BlockSize, float64(y) * kar.BlockSize}
}
func WorldSpaceToPixelSpace(pos vec.Vec2) image.Point {
	return util.Vec2ToPoint(pos).Div(int(kar.BlockSize))
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

func destroyBodyWithEntry(b *cm.Body) {
	s := b.Shapes[0].Space
	if s.ContainsBody(b) {
		e := b.UserData.(*donburi.Entry)
		e.Remove()
		s.AddPostStepCallback(removeBodyPostStep, b, false)
	}
}

func removeBodyPostStep(space *cm.Space, body, data interface{}) {
	space.RemoveBodyWithShapes(body.(*cm.Body))
}
