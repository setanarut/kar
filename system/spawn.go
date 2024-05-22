package system

import (
	"kar/arche"
	"kar/comp"
	"kar/engine"
	"kar/engine/cm"
	"kar/res"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

type SpawnSystem struct {
}

func NewSpawnSystem() *SpawnSystem {

	return &SpawnSystem{}
}

func (sys *SpawnSystem) Init() {
	res.CurrentRoom = res.ScreenRect

	res.Rooms = make([]cm.BB, 0)
	res.Rooms = append(res.Rooms, res.CurrentRoom)                                        // middle 0
	res.Rooms = append(res.Rooms, res.CurrentRoom.Offset(cm.Vec2{0, res.CurrentRoom.T}))  // top 1
	res.Rooms = append(res.Rooms, res.CurrentRoom.Offset(cm.Vec2{0, -res.CurrentRoom.T})) // bottom 2
	res.Rooms = append(res.Rooms, res.CurrentRoom.Offset(cm.Vec2{-res.CurrentRoom.R, 0})) // left 3
	res.Rooms = append(res.Rooms, res.CurrentRoom.Offset(cm.Vec2{res.CurrentRoom.R, 0}))  // right 4

	arche.SpawnRoom(res.Rooms[0], arche.RoomOptions{true, true, true, true, 1, 2, 3, 4})
	arche.SpawnRoom(res.Rooms[1], arche.RoomOptions{true, false, true, true, 5, -1, 6, 7})
	arche.SpawnRoom(res.Rooms[2], arche.RoomOptions{false, true, true, true, -1, 8, 9, 10})
	arche.SpawnRoom(res.Rooms[3], arche.RoomOptions{true, true, true, false, 11, 12, 13, -1})
	arche.SpawnRoom(res.Rooms[4], arche.RoomOptions{true, true, false, true, 14, 15, -1, 16})

	ResetLevel()

	// res.World.OnRemove(func(world donburi.World, entity donburi.Entity) {
	// 	e := world.Entry(entity)
	// 	if e.HasComponent(comp.EnemyTag) {
	// 		p := comp.Body.Get(e).Position()
	// 		i := comp.Inventory.Get(e)
	// 		for _, v := range i.Keys {
	// 			arche.SpawnDefaultKeyCollectible(v, p)
	// 		}
	// 		for range i.Bombs {
	// 			arche.SpawnDefaultBomb(p)
	// 		}
	// 		for range i.Snowballs {
	// 			arche.SpawnDefaultSnowball(p)
	// 		}
	// 	}

	// })

}

func (sys *SpawnSystem) Update() {

	// Reset Level
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		ResetLevel()
	}

	// worldPos := res.Camera.ScreenToWorld(ebiten.CursorPosition())
	// cursor := engine.InvPosVectY(worldPos, res.CurrentRoom.T)

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {

		for range 4 {
			// arche.SpawnDefaultEnemy(engine.RandomPointInBB(res.CurrentRoom, 64))
			arche.SpawnDefaultMob(engine.RandomPointInBB(res.CurrentRoom, 64))
		}

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {

		for range 10 {
			arche.SpawnDefaultBomb(engine.RandomPointInBB(res.CurrentRoom, 64))
		}

	}

}

func (sys *SpawnSystem) Draw() {
}

func ResetLevel() {
	res.CurrentRoom = res.ScreenRect
	res.Camera.LookAt(res.CurrentRoom.Center())

	player, ok := comp.PlayerTag.First(res.World)
	if ok {
		DestroyEntryWithBody(player)
		arche.SpawnDefaultPlayer(res.CurrentRoom.Center().Add(cm.Vec2{0, -100}))
	} else {
		arche.SpawnDefaultPlayer(res.CurrentRoom.Center().Add(cm.Vec2{0, -100}))
	}

	comp.EnemyTag.Each(res.World, func(e *donburi.Entry) {
		DestroyEntryWithBody(e)
	})

	comp.BombTag.Each(res.World, func(e *donburi.Entry) {
		DestroyEntryWithBody(e)
	})

	// // top room
	// for i := 5; i < 8; i++ {
	// 	arche.SpawnDefaultMob(engine.RandomPointInBB(res.Rooms[1], 20))
	// }
	// center room

	for i := 1; i < 5; i++ {
		arche.SpawnDefaultMob(engine.RandomPointInBB(res.Rooms[0], 20))
	}
	// // bottom room
	// for i := 8; i < 11; i++ {
	// 	arche.SpawnDefaultMob(engine.RandomPointInBB(res.Rooms[2], 20))
	// }

}
