package system

import (
	"kar"
	"kar/arche"
	"kar/comp"
	"kar/items"
	"kar/types"
	"kar/world"

	"github.com/setanarut/cm"
	"github.com/setanarut/kamera/v2"
	"github.com/setanarut/vec"
	"github.com/yohamta/donburi"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	gameWorld   *world.World
	playerEntry *donburi.Entry
	playerBody  *cm.Body
	inventory   *types.Inventory
	ecsWorld    donburi.World
	cmSpace     *cm.Space
	camera      *kamera.Camera
)

var (
	playerPos        vec.Vec2
	playerSpawnPoint vec.Vec2
)

type Spawn struct{}

func (s *Spawn) Init() {

	ecsWorld = donburi.NewWorld()

	cmSpace = cm.NewSpace()

	cmSpace.SetGravity(vec.Vec2{0, kar.Gravity})
	cmSpace.Damping = kar.Damping

	cmSpace.CollisionBias = kar.CollisionBias
	cmSpace.CollisionSlop = kar.CollisionSlop

	if kar.UseSpatialHash {
		cmSpace.UseSpatialHash(kar.SpatialHashDim, kar.SpatialHashCount)

	}
	cmSpace.Iterations = kar.Iterations

	gameWorld = world.NewWorld(
		kar.WorldSize.X,
		kar.WorldSize.Y,
		kar.ChunkSize,
		kar.BlockSize,
	)
	playerSpawnPoint = FindAirSpawnPoint()
	gameWorld.LoadChunks(cmSpace, ecsWorld, playerSpawnPoint)
	playerEntry = arche.SpawnPlayer(playerSpawnPoint, cmSpace, ecsWorld)

	inventory = comp.Inventory.Get(playerEntry)
	playerBody = comp.Body.Get(playerEntry)
	gameWorld.UpdateChunks(cmSpace, ecsWorld, playerBody.Position())

}
func (s *Spawn) Update() {

	if playerEntry.Valid() {
		playerPos = playerBody.Position()
		gameWorld.UpdateChunks(cmSpace, ecsWorld, playerPos)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if ebiten.IsKeyPressed(ebiten.KeyC) {
			x, y := camera.ScreenToWorld(ebiten.CursorPosition())
			arche.SpawnDebugBox(cmSpace, ecsWorld, vec.Vec2{x, y})
		}
	}
}

func (s *Spawn) Draw() {}

func FindAirSpawnPoint() vec.Vec2 {
	for y := range 200 {
		posUp := gameWorld.Image.Gray16At(10, y).Y
		posDown := gameWorld.Image.Gray16At(10, y+1).Y
		if posDown != items.Air && posUp == items.Air {
			return world.PixelCoordToWorldPos(10, y)
		}
	}
	return vec.Vec2{}
}
