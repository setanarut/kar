package terr

import (
	"image"
	"image/png"
	"math"
	"os"
	"strconv"

	"github.com/mazznoer/colorgrad"
	"github.com/ojrac/opensimplex-go"
)

type Terrain struct {
	Seed              int64
	MapSize, GridSize int
	NoiseOptions      NoiseOptions
	Gradient          colorgrad.Gradient
	Noise             opensimplex.Noise
}

func NewTerrain(seed int64, mapSize, gridSize int) *Terrain {
	terr := &Terrain{
		Seed:         seed,
		MapSize:      mapSize,
		GridSize:     gridSize,
		NoiseOptions: DefaultNoiseOptions(),
		Noise:        opensimplex.NewNormalized(seed),
	}
	terr.Gradient, _ = colorgrad.NewGradient().Build()
	return terr
}

func (tr *Terrain) SetSeed(seed int64) {
	tr.Noise = opensimplex.NewNormalized(seed)
}
func (tr *Terrain) ChunkSize() int {
	return tr.MapSize / tr.GridSize
}
func (tr *Terrain) ChunkCount() int {
	return tr.GridSize * tr.GridSize
}

func (tr *Terrain) BlockCount() int {
	return tr.MapSize * tr.MapSize
}

func (tr *Terrain) GetMapImage() *image.RGBA {
	offset := -(tr.MapSize / 2)
	img := image.NewRGBA(image.Rect(0, 0, tr.MapSize, tr.MapSize))
	for y := 0; y < tr.MapSize; y++ {
		for x := 0; x < tr.MapSize; x++ {
			t := Eval2WithOpts(
				x+offset,
				y+offset,
				tr.Noise,
				tr.NoiseOptions,
			)
			img.Set(x, tr.MapSize-y, tr.Gradient.At(t))
		}
	}
	return img
}

func (tr *Terrain) WriteChunkImage(chunkCoordX, chunkCoordY int) {
	img := tr.GetChunkImage(chunkCoordX, chunkCoordY)
	WriteImage(img, "chunk ("+strconv.Itoa(chunkCoordX)+", "+strconv.Itoa(chunkCoordY)+").png")
}
func (tr *Terrain) GetChunkForWorldPosition(x, y int) (int, int) {
	cs := float64(tr.ChunkSize())
	return int(math.Floor(float64(x) / 128)), int(math.Floor(float64(y) / cs))
}

func (tr *Terrain) WriteMapImage() {
	img := tr.GetMapImage()
	WriteImage(img, "map.png")
}

func (tr *Terrain) GetChunkImage(chunkCoordX, chunkCoordY int) *image.RGBA {
	chunkSize := tr.MapSize / tr.GridSize
	img := image.NewRGBA(image.Rect(0, 0, chunkSize, chunkSize))
	for y := 0; y < chunkSize; y++ {
		for x := 0; x < chunkSize; x++ {
			t := tr.GetBlockValue(x+(chunkSize*chunkCoordX), y+(chunkSize*chunkCoordY))
			img.Set(x, chunkSize-y, tr.Gradient.At(t))
		}
	}
	return img
}

func Eval2WithOpts(x, y int, nois opensimplex.Noise, o NoiseOptions) float64 {
	var total = 0.0
	for i := 0; i < o.Octaves; i++ {
		total += nois.Eval2(float64(x)*o.Frequency, float64(y)*o.Frequency) * o.Amplitude
		o.Amplitude *= o.Persistence
		o.Frequency *= o.Lacunarity
	}
	return total
}
func (tr *Terrain) GetBlockValue(x, y int) float64 {
	return Eval2WithOpts(x, y, tr.Noise, tr.NoiseOptions)
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
