package res

import (
	"kar/engine/util"
	"kar/types"

	"github.com/setanarut/anim"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	AtlasPlayer  = util.LoadEbitenImageFromFS(assets, "assets/player.png")
	Border48     = util.LoadEbitenImageFromFS(assets, "assets/border48.png")
	Border32     = util.LoadEbitenImageFromFS(assets, "assets/border32.png")
	Slot16       = util.LoadEbitenImageFromFS(assets, "assets/slot16.png")
	SpriteFrames = make(map[types.ItemID][]*ebiten.Image)
)

func init() {
	blockAtlas := util.LoadEbitenImageFromFS(assets, "assets/blocks.png")
	itemAtlas := util.LoadEbitenImageFromFS(assets, "assets/items.png")
	s := 16

	// blocks
	for y := range 22 {
		SpriteFrames[types.ItemID(y+1)] = anim.SubImages(blockAtlas, 0, y*s, s, s, 11, false)
	}

	// items
	for i := range 7 {
		id := i + 26
		x := i * 16
		SpriteFrames[types.ItemID(id)] = anim.SubImages(itemAtlas, x, 0, s, s, 1, false)
	}
}
