package terr

import (
	"image"
	"image/color"
	"kar/comp"
	"kar/engine/cm"
	"kar/engine/io"
	"kar/engine/vec"
	"kar/res"
	"kar/types"
	"math"
	"slices"

	"github.com/ojrac/opensimplex-go"
	"github.com/yohamta/donburi"
)

var mooreNeighbours = [8]image.Point{{1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}, {0, -1}, {1, -1}}

type Terrain struct {
	TerrainImg           *image.Gray
	NoiseOptions         NoiseOptions
	Noise                opensimplex.Noise
	ChunkSize, BlockSize float64
	LoadedChunks         []image.Point
	threshold            bool
	MapSize              int
}

func NewTerrain(seed int64, mapSize int, chunkSize float64, blockSize float64) *Terrain {
	terr := &Terrain{
		NoiseOptions: DefaultNoiseOptions(),
		Noise:        opensimplex.NewNormalized(seed),
		ChunkSize:    chunkSize,
		BlockSize:    blockSize,
		MapSize:      mapSize,

		threshold: true,
	}
	return terr
}

func (tr *Terrain) Generate() {
	tr.TerrainImg = image.NewGray(image.Rect(0, 0, tr.MapSize, tr.MapSize))
	for y := 0; y < tr.TerrainImg.Bounds().Dy(); y++ {
		for x := 0; x < tr.TerrainImg.Bounds().Dx(); x++ {
			v := tr.Eval2WithOptions(x, y)
			if tr.threshold {
				tr.TerrainImg.SetGray(x, y, color.Gray{Y: uint8(math.Round(v)) * 255})
			} else {
				tr.TerrainImg.SetGray(x, y, color.Gray{Y: uint8(v * 255)})
			}
		}
	}
}

func (tr *Terrain) SpawnChunk(chunkCoord image.Point, blockSpawnCallbackFunc func(pos vec.Vec2, chunkCoord image.Point, blockType types.BlockType)) {
	chunksize := int(tr.ChunkSize)
	for y := 0; y < chunksize; y++ {
		for x := 0; x < chunksize; x++ {
			blockX := x + (chunksize * chunkCoord.X)
			blockY := y + (chunksize * chunkCoord.Y)
			blockNumber := tr.TerrainImg.GrayAt(blockX, blockY)
			blockPos := vec.Vec2{float64(blockX), float64(blockY)}
			blockPos = blockPos.Scale(tr.BlockSize)

			if tr.threshold {
				if blockNumber.Y != 0 {
					blockPos.X += res.BlockSize / 2
					blockSpawnCallbackFunc(blockPos.NegY(), chunkCoord, res.BlockStone)
				}
			} else {
				if blockNumber.Y < 85 {
					blockSpawnCallbackFunc(blockPos, chunkCoord, res.BlockStone)
				}
				if blockNumber.Y > 85 && blockNumber.Y < 100 {
					blockSpawnCallbackFunc(blockPos, chunkCoord, res.BlockDirt)
				}

			}

		}
	}
}
func (tr *Terrain) EmptyBlockInChunk(chunkCoord image.Point) (bool, image.Point) {
	chunksize := int(tr.ChunkSize)
	for y := 0; y < chunksize; y++ {
		for x := 0; x < chunksize; x++ {
			blockX := x + (chunksize * chunkCoord.X)
			blockY := y + (chunksize * chunkCoord.Y)
			if tr.TerrainImg.GrayAt(blockX, blockY).Y == 0 {
				return true, image.Point{blockX, blockY}
			}

		}
	}
	return false, image.Point{}
}

// Spawn/Destroy chunks
func (tr *Terrain) UpdateChunks(playerChunk image.Point, blockSpawnCallbackFunc func(pos vec.Vec2, chunkCoord image.Point, blockType types.BlockType)) {
	playerChunks := GetPlayerChunks(playerChunk)
	intersectionChunks := FindIntersection(playerChunks, tr.LoadedChunks)

	for _, toLoad := range playerChunks {
		if !slices.Contains(intersectionChunks, toLoad) {
			flip := toLoad
			flip.Y *= -1
			tr.SpawnChunk(flip, blockSpawnCallbackFunc)
		}
	}

	for _, toUnload := range tr.LoadedChunks {
		if !slices.Contains(intersectionChunks, toUnload) {
			flip := toUnload
			flip.Y *= -1
			tr.DeSpawnChunk(flip)
		}
	}

	tr.LoadedChunks = playerChunks
}

func (tr *Terrain) DeSpawnChunk(chunkCoord image.Point) {
	comp.Block.Each(res.World, func(e *donburi.Entry) {
		if comp.Block.Get(e).ChunkCoord == chunkCoord {
			b := comp.Body.Get(e)
			destroyBodyWithEntry(b)
		}
	})
}

func (tr *Terrain) WriteChunkImage(chunkCoord image.Point, filename string) {
	img := tr.ChunkImage(chunkCoord)
	io.WriteImage(img, filename)
}

func (tr *Terrain) WorldPosToChunkCoord(worldPos vec.Vec2) image.Point {
	return worldPos.Div(tr.ChunkSize).Floor().Div(tr.BlockSize).Point()
}

func (tr *Terrain) InTerrainBounds(worldPos vec.Vec2) bool {
	s := float64(tr.MapSize) * tr.BlockSize
	return worldPos.X > s || worldPos.Y > s
}

func (tr *Terrain) WorldSpaceToMapSpace(worldPos vec.Vec2) image.Point {
	return worldPos.Div(tr.BlockSize).Point()
}

func (tr *Terrain) WriteTerrainImage(flipVertical bool) {
	if flipVertical {
		io.WriteImage(tr.TerrainImageFlipVertical(), "map.png")
	} else {
		io.WriteImage(tr.TerrainImg, "map.png")
	}
}

func (tr *Terrain) ChunkImage(chunkCoord image.Point) *image.Gray {
	chunksize := int(tr.ChunkSize)
	img := image.NewGray(image.Rect(0, 0, chunksize, chunksize))
	for y := 0; y < chunksize; y++ {
		for x := 0; x < chunksize; x++ {
			gray := tr.TerrainImg.GrayAt(x+(chunksize*chunkCoord.X), y+(chunksize*chunkCoord.Y))
			// resim için y eksenini ters çevir
			img.Set(x, chunksize-y, gray)
		}
	}
	return img
}

// Chipmunk coordinate system conversion
func (tr *Terrain) TerrainImageFlipVertical() *image.Gray {
	size := tr.TerrainImg.Bounds().Dy()
	img := image.NewGray(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			gray := tr.TerrainImg.GrayAt(x, y)
			// resim için y eksenini ters çevir
			img.Set(x, size-y, gray)
		}
	}
	return img
}

func (tr *Terrain) Eval2WithOptions(x, y int) float64 {
	return eval2WithOpts(x, y, tr.Noise, tr.NoiseOptions)
}

type NoiseOptions struct {
	Amplitude float64
	// Distance to view the noisemap.
	Frequency float64
	// How much each octave contributes to the overall shape (adjusts amplitude).
	Persistence float64
	// How much detail is added or removed at each octave (adjusts frequency).
	Lacunarity float64
	// The number of levels of detail you want you perlin noise to have.
	Octaves int
}

func DefaultNoiseOptions() NoiseOptions {
	return NoiseOptions{
		Amplitude:   1,
		Frequency:   0.2,
		Persistence: 0.,
		Lacunarity:  0.,
		Octaves:     1,
	}
}

func eval2WithOpts(x, y int, nois opensimplex.Noise, o NoiseOptions) float64 {
	var total = 0.0
	for i := 0; i < o.Octaves; i++ {
		total += nois.Eval2(float64(x)*o.Frequency, float64(y)*o.Frequency) * o.Amplitude
		o.Amplitude *= o.Persistence
		o.Frequency *= o.Lacunarity
	}
	return total
}

func FindIntersection(s1, s2 []image.Point) []image.Point {
	intersection := make([]image.Point, 0)
	for _, point := range s2 {
		if slices.Contains(s1, point) {
			intersection = append(intersection, point)
		}
	}
	return intersection
}

func GetPlayerChunks(playerChunk image.Point) []image.Point {
	playerChunks := make([]image.Point, 0)
	playerChunks = append(playerChunks, playerChunk)

	for _, neighbor := range mooreNeighbours {
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
