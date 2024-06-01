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
	TerrainImg   *image.Gray
	NoiseOptions NoiseOptions
	Noise        opensimplex.Noise
	ChunkSize    float64
}

func NewTerrain(seed int64, chunkSize float64) *Terrain {
	terr := &Terrain{
		NoiseOptions: DefaultNoiseOptions(),
		Noise:        opensimplex.NewNormalized(seed),
		ChunkSize:    chunkSize,
	}
	return terr
}

func (tr *Terrain) Generate(threshold bool) {
	tr.TerrainImg = image.NewGray(image.Rect(0, 0, 1024, 1024))
	for y := 0; y < tr.TerrainImg.Bounds().Dy(); y++ {
		for x := 0; x < tr.TerrainImg.Bounds().Dx(); x++ {
			v := tr.Eval2WithOptions(x, y)
			if threshold {
				tr.TerrainImg.SetGray(x, y, color.Gray{Y: uint8(math.Round(v)) * 255})
			} else {
				tr.TerrainImg.SetGray(x, y, color.Gray{Y: uint8(v * 255)})
			}
		}
	}
}

func (tr *Terrain) WriteChunkImage(chunkCoord image.Point, filename string) {
	img := tr.ChunkImage(chunkCoord)
	WriteImage(img, filename)
}

func (tr *Terrain) ChunkCoord(pos cm.Vec2, blockSize float64) image.Point {
	return pos.Div(tr.ChunkSize).Floor().Div(blockSize).Point()
}

func (tr *Terrain) WriteMapImage() {
	WriteImage(tr.MapImageInvertY(), "map.png")
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
func (tr *Terrain) MapImageInvertY() *image.Gray {
	size := tr.TerrainImg.Bounds().Dx()
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

func (tr *Terrain) SpawnChunk(chunkCoord image.Point, callback func(pos cm.Vec2)) {
	chunksize := int(tr.ChunkSize)
	for y := 0; y < chunksize; y++ {
		for x := 0; x < chunksize; x++ {
			blockX := x + (chunksize * chunkCoord.X)
			blockY := y + (chunksize * chunkCoord.Y)
			blockNumber := tr.TerrainImg.GrayAt(blockX, blockY)
			blockPos := cm.Vec2{float64(blockX), float64(blockY)}
			blockPos = blockPos.Mult(50) // blok boyutu
			if blockNumber.Y > 128 {
				callback(blockPos)
			}
		}
	}
}

func (tr *Terrain) Eval2WithOptions(x, y int) float64 {
	return eval2WithOpts(x, y, tr.Noise, tr.NoiseOptions)
}

func WriteImage(img *image.Gray, filename string) {
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
