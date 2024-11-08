package system

import (
	"kar"
	"kar/arc"
	"kar/engine/vectorg"
	"kar/items"
	"kar/world"

	"github.com/mlange-42/arche/ecs"
	"github.com/setanarut/anim"
	"github.com/setanarut/cm"
	"github.com/setanarut/vec"
)

var (
	playerEntity      ecs.Entity
	playerSpawnPos    vec2
	playerBody        *cm.Body
	playerPos         vec.Vec2
	playerInv         *arc.Inventory
	playerAnim        *anim.AnimationPlayer
	playerDrawOptions *arc.DrawOptions
)

// func EntityRemoveCallback(w *ecs.World, e ecs.EntityEvent) {
// 	fmt.Println(w.Alive(e.Entity))
// }

// var ls = listener.NewCallback(
// 	EntityRemoveCallback,
// 	event.EntityRemoved,
// )

type Spawn struct{}

func (s *Spawn) Init() {
	// kar.WorldECS.SetListener(&ls)
	gameWorld = world.NewWorld(
		kar.WorldSize.X,
		kar.WorldSize.Y,
		kar.ChunkSize,
		kar.BlockSize,
	)
	findSpawnPosition()
	// playerSpawnPos = vec.Vec2{400, 400}
	playerEntity = arc.SpawnMario(playerSpawnPos)
	playerBody = arc.MapBody.Get(playerEntity).Body
	_, playerDrawOptions, playerAnim, _, playerInv = arc.MapPlayer.Get(playerEntity)
	gameWorld.LoadChunks(playerSpawnPos)

}
func (s *Spawn) Update() {

	if debugDrawingEnabled {
		vectorg.GlobalTransform.Reset()
		cmDrawer.GeoM.Reset()
		Camera.ApplyCameraTransform(cmDrawer.GeoM)
		Camera.ApplyCameraTransform(vectorg.GlobalTransform)
	}

	playerPos = playerBody.Position()
	// playerBody.SetVelocityUpdateFunc(PlayerFlyVelocityFunc)
	if kar.WorldECS.Alive(playerEntity) {
		gameWorld.UpdateChunks(playerBody.Position())
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
