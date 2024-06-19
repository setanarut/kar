package terr

import (
	"image"
	"image/color"
	"kar/comp"
	"kar/engine/cm"
	"kar/engine/io"
	"kar/engine/mathutil"
	"kar/engine/vec"
	"kar/items"
	"kar/res"
	"kar/types"
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
	MapSize              float64
}

func NewTerrain(seed int64, mapSize float64, chunkSize float64, blockSize float64) *Terrain {
	terr := &Terrain{
		NoiseOptions: DefaultNoiseOptions(),
		Noise:        opensimplex.NewNormalized(seed),
		ChunkSize:    chunkSize,
		BlockSize:    blockSize,
		MapSize:      mapSize,
	}
	return terr
}

func (tr *Terrain) Generate() {
	tr.TerrainImg = image.NewGray(image.Rect(0, 0, int(tr.MapSize), int(tr.MapSize)))
	for y := 0; y < tr.TerrainImg.Bounds().Dy(); y++ {
		for x := 0; x < tr.TerrainImg.Bounds().Dx(); x++ {
			noiseValue := tr.Eval2WithOptions(x, y)
			if noiseValue < 0.33 {
				tr.TerrainImg.SetGray(x, y, color.Gray{Y: uint8(items.Air)})
			}
			if mathutil.InRange(noiseValue, 0.33, 0.66) {
				tr.TerrainImg.SetGray(x, y, color.Gray{Y: uint8(items.Dirt)})
			}
			if noiseValue > 0.66 {
				tr.TerrainImg.SetGray(x, y, color.Gray{Y: uint8(items.IronOre)})
			}
		}
	}
}

func (tr *Terrain) SpawnChunk(chunkCoord image.Point, spawnBlock types.BlockSpawnFunc) {
	chunksize := int(tr.ChunkSize)
	for y := 0; y < chunksize; y++ {
		for x := 0; x < chunksize; x++ {
			blockX := x + (chunksize * chunkCoord.X)
			blockY := y + (chunksize * chunkCoord.Y)
			blockType := types.ItemType(tr.TerrainImg.GrayAt(blockX, blockY).Y)
			blockPos := vec.Vec2{float64(blockX), float64(blockY)}
			blockPos = blockPos.Scale(tr.BlockSize)
			if blockType != items.Air {
				spawnBlock(blockPos, chunkCoord, blockType)
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
func (tr *Terrain) UpdateChunks(playerChunk image.Point, spawnBlock types.BlockSpawnFunc) {
	playerChunks := GetPlayerChunks(playerChunk)
	intersectionChunks := FindIntersection(playerChunks, tr.LoadedChunks)

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

func (tr *Terrain) DeSpawnChunk(chunkCoord image.Point) {
	comp.Item.Each(res.World, func(e *donburi.Entry) {
		if comp.Item.Get(e).ChunkCoord == chunkCoord {
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
	return worldPos.Div(tr.ChunkSize).Div(tr.BlockSize).Point()
}

func (tr *Terrain) InTerrainBounds(worldPos vec.Vec2) bool {
	s := tr.MapSize * tr.BlockSize
	return worldPos.X > s || worldPos.Y > s
}

func (tr *Terrain) WorldSpaceToMapSpace(worldPos vec.Vec2) image.Point {
	// return worldPos.Div(tr.BlockSize).Point()
	return worldPos.Point().Div(int(res.BlockSize))
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
