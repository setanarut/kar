package system

import (
	"kar"
	"kar/arche"
	"kar/comp"
	"kar/items"
	"kar/types"
	"kar/world"

	"github.com/setanarut/anim"
	"github.com/setanarut/cm"
	"github.com/yohamta/donburi"
)

var (
	playerEntry       *donburi.Entry
	playerSpawnPos    vec2
	playerBody        *cm.Body
	playerInv         *types.Inventory
	playerAnim        *anim.AnimationPlayer
	playerDrawOptions *types.DrawOptions
)

type Spawn struct{}

func (s *Spawn) Init() {
	gameWorld = world.NewWorld(
		kar.WorldSize.X,
		kar.WorldSize.Y,
		kar.ChunkSize,
		kar.BlockSize,
	)
	findSpawnPosition()
	gameWorld.LoadChunks(Space, ecsWorld, playerSpawnPos)
	playerEntry = arche.SpawnPlayer(Space, ecsWorld, playerSpawnPos, 1, 0, 0)
	playerInv = comp.Inventory.Get(playerEntry)
	playerBody = comp.Body.Get(playerEntry)
	playerAnim = comp.AnimPlayer.Get(playerEntry)
	playerDrawOptions = comp.DrawOptions.Get(playerEntry)
	playerBody.SetVelocityUpdateFunc(VelocityFunc)
}
func (s *Spawn) Update() {
	if playerEntry.Valid() {
		playerPos = playerBody.Position()
		gameWorld.UpdateChunks(Space, ecsWorld, playerPos)

	}
}

func (s *Spawn) Draw() {
}

func findSpawnPosition() {
	for y := range 300 {
		posUp := gameWorld.Image.Gray16At(10, y).Y
		posDown := gameWorld.Image.Gray16At(10, y+1).Y
		if posDown != items.Air && posUp == items.Air {
			playerSpawnPos = world.PixelToWorld(10, y)
			break
		}
	}
}
