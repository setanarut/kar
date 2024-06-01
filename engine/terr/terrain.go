package terr

import (
	"image"
	"image/color"
	"image/png"
	"kar/engine/cm"
	"os"

	"github.com/ojrac/opensimplex-go"
)

type Terrain struct {
	TerrainData1024 [1024][1024]uint8
	NoiseOptions    NoiseOptions
	Noise           opensimplex.Noise
	ChunkSize       float64
	TerrainSize     float64
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
	for y := 0; y < int(tr.TerrainSize); y++ {
		for x := 0; x < int(tr.TerrainSize); x++ {
			tr.TerrainData1024[x][y] = uint8(tr.Eval2WithOptions(x, y) * 255)
		}
	}
}

func (tr *Terrain) MapImage() *image.RGBA {
	terraSize := int(tr.TerrainSize)
	img := image.NewRGBA(image.Rect(0, 0, int(tr.TerrainSize), int(tr.TerrainSize)))
	for y := 0; y < terraSize; y++ {
		for x := 0; x < terraSize; x++ {
			value := tr.TerrainData1024[x][y]
			img.Set(x, terraSize-y, color.Gray{value})
		}
	}
	return img
}

func (tr *Terrain) WriteChunkImage(chunkCoord image.Point, filename string) {
	img := tr.ChunkImage(chunkCoord)
	// WriteImage(img, "chunk ("+strconv.Itoa(chunkCoordX)+", "+strconv.Itoa(chunkCoordY)+").png")
	WriteImage(img, filename)
}

func (tr *Terrain) ChunkCoord(pos cm.Vec2, blockSize float64) image.Point {

	return pos.Div(tr.ChunkSize).Floor().Div(blockSize).Point()
}

func (tr *Terrain) WriteMapImage() {
	img := tr.MapImage()
	WriteImage(img, "map.png")
}

func (tr *Terrain) ChunkImage(chunkCoord image.Point) *image.RGBA {
	chunksize := int(tr.ChunkSize)
	img := image.NewRGBA(image.Rect(0, 0, chunksize, chunksize))
	for y := 0; y < chunksize; y++ {
		for x := 0; x < chunksize; x++ {
			value := tr.TerrainData1024[x+(chunksize*chunkCoord.X)][y+(chunksize*chunkCoord.Y)]
			// resim için y eksenini ters çevir
			img.Set(x, chunksize-y, color.Gray{value})
		}
	}
	return img
}

func (tr *Terrain) SpawnChunk(chunkCoord image.Point, callback func(pos cm.Vec2)) {
	chunksize := int(tr.ChunkSize)
	for y := 0; y < chunksize; y++ {
		for x := 0; x < chunksize; x++ {
			blockCoordX := x + (chunksize * chunkCoord.X)
			blockCoordY := y + (chunksize * chunkCoord.Y)
			blockNumber := tr.TerrainData1024[blockCoordX][blockCoordY]
			blockPos := cm.Vec2{float64(blockCoordX), float64(blockCoordY)}
			blockPos = blockPos.Mult(50) // blok boyutu
			if blockNumber > 128 {
				callback(blockPos)
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
