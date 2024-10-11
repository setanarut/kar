package res

import (
	"image"
	"kar/engine/util"
	"kar/itm"

	"github.com/anthonynsimon/bild/blend"
	"github.com/setanarut/anim"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	AtlasPlayer     = util.LoadEbitenImageFromFS(assets, "assets/player.png")
	Border          = util.LoadEbitenImageFromFS(assets, "assets/border64.png")
	Hotbar          = util.LoadEbitenImageFromFS(assets, "assets/hotbar.png")
	HotbarSelection = util.LoadEbitenImageFromFS(assets, "assets/hotbarselect.png")
	SpriteFrames    = make(map[uint16][]*ebiten.Image)
)

func init() {
	blockAtlas := util.LoadEbitenImageFromFS(assets, "assets/blocks.png")
	itemAtlas := util.LoadEbitenImageFromFS(assets, "assets/items.png")
	s := 16

	// blocks
	SpriteFrames[itm.Bedrock] = anim.SubImages(blockAtlas, 176, 0, s, s, 1, false)
	SpriteFrames[itm.Grass] = anim.SubImages(blockAtlas, 0, 0, s, s, 11, false)
	SpriteFrames[itm.Snow] = anim.SubImages(blockAtlas, 0, s, s, s, 11, false)
	SpriteFrames[itm.Dirt] = anim.SubImages(blockAtlas, 0, 2*s, s, s, 11, false)
	SpriteFrames[itm.Sand] = anim.SubImages(blockAtlas, 0, 3*s, s, s, 11, false)
	// SpriteFrames[itm.Stone] = anim.SubImages(blockAtlas, 0, 4*s, s, s, 11, false)
	SpriteFrames[itm.Stone] = getFrames(1)
	SpriteFrames[itm.CoalOre] = anim.SubImages(blockAtlas, 0, 5*s, s, s, 11, false)
	SpriteFrames[itm.GoldOre] = anim.SubImages(blockAtlas, 0, 6*s, s, s, 11, false)
	SpriteFrames[itm.IronOre] = anim.SubImages(blockAtlas, 0, 7*s, s, s, 11, false)
	SpriteFrames[itm.DiamondOre] = anim.SubImages(blockAtlas, 0, 8*s, s, s, 11, false)
	SpriteFrames[itm.CopperOre] = anim.SubImages(blockAtlas, 0, 9*s, s, s, 11, false)
	SpriteFrames[itm.EmeraldOre] = anim.SubImages(blockAtlas, 0, 10*s, s, s, 11, false)
	SpriteFrames[itm.LapisOre] = anim.SubImages(blockAtlas, 0, 11*s, s, s, 11, false)
	SpriteFrames[itm.RedstoneOre] = anim.SubImages(blockAtlas, 0, 12*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateStone] = anim.SubImages(blockAtlas, 0, 13*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateCoalOre] = anim.SubImages(blockAtlas, 0, 14*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateGoldOre] = anim.SubImages(blockAtlas, 0, 15*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateIronOre] = anim.SubImages(blockAtlas, 0, 16*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateDiamondOre] = anim.SubImages(blockAtlas, 0, 17*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateCopperOre] = anim.SubImages(blockAtlas, 0, 18*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateEmeraldOre] = anim.SubImages(blockAtlas, 0, 19*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateLapisOre] = anim.SubImages(blockAtlas, 0, 20*s, s, s, 11, false)
	SpriteFrames[itm.DeepslateRedStoneOre] = anim.SubImages(blockAtlas, 0, 21*s, s, s, 11, false)
	SpriteFrames[itm.Log] = anim.SubImages(blockAtlas, 0, 22*s, s, s, 11, false)
	SpriteFrames[itm.Leaves] = anim.SubImages(blockAtlas, 0, 23*s, s, s, 11, false)
	SpriteFrames[itm.Planks] = anim.SubImages(blockAtlas, 0, 24*s, s, s, 11, false)
	SpriteFrames[itm.Cobblestone] = anim.SubImages(blockAtlas, 0, 25*s, s, s, 11, false)
	SpriteFrames[itm.CobbledDeepslate] = anim.SubImages(blockAtlas, 0, 26*s, s, s, 11, false)
	// items
	SpriteFrames[itm.Sapling] = anim.SubImages(itemAtlas, 0, 0, s, s, 1, false)
	SpriteFrames[itm.Torch] = anim.SubImages(itemAtlas, s, 0, s, s, 1, false)
	SpriteFrames[itm.Coal] = anim.SubImages(itemAtlas, 2*s, 0, s, s, 1, false)
	SpriteFrames[itm.RawGold] = anim.SubImages(itemAtlas, 3*s, 0, s, s, 1, false)
	SpriteFrames[itm.RawIron] = anim.SubImages(itemAtlas, 4*s, 0, s, s, 1, false)
	SpriteFrames[itm.Diamond] = anim.SubImages(itemAtlas, 5*s, 0, s, s, 1, false)
	SpriteFrames[itm.RawCopper] = anim.SubImages(itemAtlas, 6*s, 0, s, s, 1, false)
	SpriteFrames[itm.CharCoal] = anim.SubImages(itemAtlas, 7*s, 0, s, s, 1, false)
	SpriteFrames[itm.Emerald] = anim.SubImages(itemAtlas, 8*s, 0, s, s, 1, false)
	SpriteFrames[itm.LapisLazuli] = anim.SubImages(itemAtlas, 9*s, 0, s, s, 1, false)
	SpriteFrames[itm.Redstone] = anim.SubImages(itemAtlas, 10*s, 0, s, s, 1, false)
	SpriteFrames[itm.WoodAxe] = anim.SubImages(itemAtlas, 11*s, 0, s, s, 1, false)
	SpriteFrames[itm.WoodHoe] = anim.SubImages(itemAtlas, 12*s, 0, s, s, 1, false)
	SpriteFrames[itm.WoodPickaxe] = anim.SubImages(itemAtlas, 13*s, 0, s, s, 1, false)
	SpriteFrames[itm.WoodShovel] = anim.SubImages(itemAtlas, 14*s, 0, s, s, 1, false)
	SpriteFrames[itm.WoodSword] = anim.SubImages(itemAtlas, 15*s, 0, s, s, 1, false)
	SpriteFrames[itm.StoneAxe] = anim.SubImages(itemAtlas, 16*s, 0, s, s, 1, false)
	SpriteFrames[itm.StoneHoe] = anim.SubImages(itemAtlas, 17*s, 0, s, s, 1, false)
	SpriteFrames[itm.StonePickaxe] = anim.SubImages(itemAtlas, 18*s, 0, s, s, 1, false)
	SpriteFrames[itm.StoneShovel] = anim.SubImages(itemAtlas, 19*s, 0, s, s, 1, false)
	SpriteFrames[itm.StoneSword] = anim.SubImages(itemAtlas, 20*s, 0, s, s, 1, false)
	SpriteFrames[itm.GoldenAxe] = anim.SubImages(itemAtlas, 21*s, 0, s, s, 1, false)
	SpriteFrames[itm.GoldenHoe] = anim.SubImages(itemAtlas, 22*s, 0, s, s, 1, false)
	SpriteFrames[itm.GoldenPickaxe] = anim.SubImages(itemAtlas, 23*s, 0, s, s, 1, false)
	SpriteFrames[itm.GoldenShovel] = anim.SubImages(itemAtlas, 24*s, 0, s, s, 1, false)
	SpriteFrames[itm.GoldenSword] = anim.SubImages(itemAtlas, 25*s, 0, s, s, 1, false)
	SpriteFrames[itm.IronAxe] = anim.SubImages(itemAtlas, 26*s, 0, s, s, 1, false)
	SpriteFrames[itm.IronHoe] = anim.SubImages(itemAtlas, 27*s, 0, s, s, 1, false)
	SpriteFrames[itm.IronPickaxe] = anim.SubImages(itemAtlas, 28*s, 0, s, s, 1, false)
	SpriteFrames[itm.IronShovel] = anim.SubImages(itemAtlas, 29*s, 0, s, s, 1, false)
	SpriteFrames[itm.IronSword] = anim.SubImages(itemAtlas, 30*s, 0, s, s, 1, false)
	SpriteFrames[itm.DiamondAxe] = anim.SubImages(itemAtlas, 31*s, 0, s, s, 1, false)
	SpriteFrames[itm.DiamondHoe] = anim.SubImages(itemAtlas, 0, s, s, s, 1, false)
	SpriteFrames[itm.DiamondPickaxe] = anim.SubImages(itemAtlas, s, s, s, s, 1, false)
	SpriteFrames[itm.DiamondShovel] = anim.SubImages(itemAtlas, 2*s, s, s, s, 1, false)
	SpriteFrames[itm.DiamondSword] = anim.SubImages(itemAtlas, 3*s, s, s, s, 1, false)
	SpriteFrames[itm.NetheriteAxe] = anim.SubImages(itemAtlas, 4*s, s, s, s, 1, false)
	SpriteFrames[itm.NetheriteHoe] = anim.SubImages(itemAtlas, 5*s, s, s, s, 1, false)
	SpriteFrames[itm.NetheritePickaxe] = anim.SubImages(itemAtlas, 6*s, s, s, s, 1, false)
	SpriteFrames[itm.NetheriteShovel] = anim.SubImages(itemAtlas, 7*s, s, s, s, 1, false)
	SpriteFrames[itm.NetheriteSword] = anim.SubImages(itemAtlas, 8*s, s, s, s, 1, false)
	SpriteFrames[itm.NetheriteScrap] = anim.SubImages(itemAtlas, 9*s, s, s, s, 1, false)
	SpriteFrames[itm.NetheriteIngot] = anim.SubImages(itemAtlas, 10*s, s, s, s, 1, false)
	SpriteFrames[itm.GoldIngot] = anim.SubImages(itemAtlas, 11*s, s, s, s, 1, false)
	SpriteFrames[itm.IronIngot] = anim.SubImages(itemAtlas, 12*s, s, s, s, 1, false)
	SpriteFrames[itm.CopperIngot] = anim.SubImages(itemAtlas, 13*s, s, s, s, 1, false)
	SpriteFrames[itm.Stick] = anim.SubImages(itemAtlas, 14*s, s, s, s, 1, false)

}

func getFrames(id int) []*ebiten.Image {
	return convert2ebiten(makeStages(base64ToImage(textures[id].Base64)))
}

func convert2ebiten(st []image.Image) []*ebiten.Image {
	l := make([]*ebiten.Image, 0)
	for _, v := range st {
		l = append(l, ebiten.NewImageFromImage(v))
	}
	return l
}

func makeStages(block image.Image) []image.Image {
	cracks := base64ToImage(stagesBase64Img)
	frames := make([]image.Image, 0)
	frames = append(frames, block)
	for i := range 10 {
		x := i * 16
		rec := image.Rect(x, 0, x+16, x+16)
		si := cracks.(*image.NRGBA).SubImage(rec)
		frames = append(frames, blend.Overlay(block, si))
	}
	return frames
}

var stagesBase64Img string = "iVBORw0KGgoAAAANSUhEUgAAAKAAAAAQCAYAAACRBXRYAAABgGlDQ1BzUkdCIElFQzYxOTY2LTIuMQAAKJF1kc8rRFEUxz8zaISJQllYTMJq/Bo1sVFmEkrSGGWwmXnezKh5M6/3ZpJsle0UJTZ+LfgL2CprpYiUrK2JDXrO89RMMud27vnc773ndO+54I5mFM2s7gctmzci4yHffGzB53mmDg+t9OGNK6Y+OjMzRUV7v8Nlx5seu1blc/9a/bJqKuCqFR5RdCMvPCE8tZrXbd4WblHS8WXhU2G/IRcUvrX1hMPPNqcc/rTZiEbC4G4S9qXKOFHGStrQhOXldGqZgvJ7H/slDWp2blZih3g7JhHGCeFjkjHCBBlgWOYgPQTolRUV8vt/8qfJSa4is84aBiukSJPHL2pBqqsSk6KrMjKs2f3/21czORhwqjeEoObJsl67wLMFX0XL+ji0rK8jqHqEi2wpP3cAQ2+iF0ta5z40bsDZZUlL7MD5JrQ96HEj/iNVibuTSXg5AW8Mmq+hbtHp2e8+x/cQXZevuoLdPeiW841L31QKZ90VQ2gVAAAACXBIWXMAAAsTAAALEwEAmpwYAAACdklEQVRoge2ZUXLEIAiG/2R6LXMUn7yVefIoezDtQ8uUUoxCtt10Z5nxIZHPGEUgBK215dVeTbYQQmutLTnnpvXL+6Tf4+la3sejX/TV+pvvZb08sSGE1jOw3tg55+bh3/CSQ9m2rd5ut9XLbdtWAcAyRozRzccYzc8nXWIBoJSClBKIpz4+N8kDqFb+qQ3Qazyc55tiEVpoD59SWgFUC68YWpeX60JzPeI/54SU0lpKqXy8U/yjw81vhCHie/mLle/lN2f43hg8VFn4nHOjMKjxXF+79vI85Hr4X/OAZ73PPYROHZ+LZV7Ecxl5JT6+h/c+n3S5d+E8D80yTMcY3bz0cFb+34RgSy5EuiQ8RIw2X7IarxmG5fmP4EspKKVUwVfgp6Fb+N4hmeWX1tpy9CJeL7bvey2lgBJh6ziSlydtxAIfibCFl97lajxdX4nnjsH1/KM8hsd0TyPem4v9NS/1LLymd2/+aBwtZ5T80V56eZ77WfkQQht+hJxJ5Ecv/Ux8b51m+V6dzMPzuXD+yKFY+d48LXzOuQ3D2UzIOxJvGeO/CV8nLY8ciayT8b6Z8ThfSvnWR3xKaZX7odXpRrz8kPDyKaV1aa2fAz5a+EQ9LDCXtPfYq/CAzRFoPBmFzMF4rj3L85z+LH/pr2Cv8REL+Dy49kdgVrRE/Cw/+oCTZaYOr3pRWRie4Xm/lycjvbQBngn/d2RN4VQ5NHW2CtDx+JV7FG2cQZnmG68dBisv19bJVwC4dAh+Ftn3Xa2XzepRWQnoHyxZeqL7sx6Vl0qAL0PXeN5/lr+0B3wWmfXGPb0ZXhaXSTSPpPFM70e/xmte3cO/A/sae4F/B2fXAAAAAElFTkSuQmCC"
