package system

import (
	"kar"
	"kar/arche"
	"kar/comp"
	"kar/controller"
	"kar/items"
	"kar/world"
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
	playerBody.SetVelocityUpdateFunc(controller.VelocityFunc)
}
func (s *Spawn) Update() {
	if playerEntry.Valid() {
		playerPos = playerBody.Position()
		playerVel = playerBody.Velocity()
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
