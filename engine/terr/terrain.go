package terr

import (
	"image"
	"image/color"
	"image/png"
	"kar/engine/cm"
	"math"
	"os"

	"github.com/ojrac/opensimplex-go"
)

type Terrain struct {
	TerrainData1024 [1024][1024]uint8
	NoiseOptions    NoiseOptions
	Noise           opensimplex.Noise
	ChunkSize       int
	TerrainSize     int
}

func NewTerrain(seed int64) *Terrain {
	terr := &Terrain{
		NoiseOptions: DefaultNoiseOptions(),
		Noise:        opensimplex.NewNormalized(seed),
		ChunkSize:    16,
		TerrainSize:  1024,
	}
	return terr
}

func (tr *Terrain) Generate() {
	for y := 0; y < tr.TerrainSize; y++ {
		for x := 0; x < tr.TerrainSize; x++ {
			tr.TerrainData1024[x][y] = uint8(tr.Eval2WithOptions(x, y) * 255)
		}
	}
}

func (tr *Terrain) MapImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, tr.TerrainSize, tr.TerrainSize))
	for y := 0; y < tr.TerrainSize; y++ {
		for x := 0; x < tr.TerrainSize; x++ {
			value := tr.TerrainData1024[x][y]
			img.Set(x, tr.TerrainSize-y, color.Gray{value})
		}
	}
	return img
}

func (tr *Terrain) WriteChunkImage(chunkCoordX, chunkCoordY int, filename string) {
	img := tr.ChunkImage(chunkCoordX, chunkCoordY)
	// WriteImage(img, "chunk ("+strconv.Itoa(chunkCoordX)+", "+strconv.Itoa(chunkCoordY)+").png")
	WriteImage(img, filename)
}

func (tr *Terrain) ChunkCoord(pos cm.Vec2) image.Point {
	cs := float64(tr.ChunkSize)
	return image.Point{int(math.Floor(float64(pos.X) / cs)), int(math.Floor(float64(pos.Y) / cs))}
}

func (tr *Terrain) WriteMapImage() {
	img := tr.MapImage()
	WriteImage(img, "map.png")
}

func (tr *Terrain) ChunkImage(chunkCoordX, chunkCoordY int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, tr.ChunkSize, tr.ChunkSize))
	for y := 0; y < tr.ChunkSize; y++ {
		for x := 0; x < tr.ChunkSize; x++ {
			value := tr.TerrainData1024[x+(tr.ChunkSize*chunkCoordX)][y+(tr.ChunkSize*chunkCoordY)]
			// resim için y eksenini ters çevir
			img.Set(x, tr.ChunkSize-y, color.Gray{value})
		}
	}
	return img
}
func (tr *Terrain) SpawnChunk(coord image.Point, callback func(pos cm.Vec2)) {
	for y := 0; y < tr.ChunkSize; y++ {
		for x := 0; x < tr.ChunkSize; x++ {
			blockCoordX := x + (tr.ChunkSize * coord.X)
			blockCoordY := y + (tr.ChunkSize * coord.Y)
			blockNumber := tr.TerrainData1024[blockCoordX][blockCoordY]
			p := cm.Vec2{float64(blockCoordX), float64(blockCoordY)}
			p = p.Mult(50) // blok boyutu

			if blockNumber > 128 {
				callback(p)
			}
		}
	}
}

func (tr *Terrain) Eval2WithOptions(x, y int) float64 {
	return eval2WithOpts(x, y, tr.Noise, tr.NoiseOptions)
}

func WriteImage(img *image.RGBA, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	png.Encode(file, img)
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
		Frequency:   0.03,
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
