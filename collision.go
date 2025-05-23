package kar

import (
	"image"
	"kar/items"
	"math"

	"github.com/setanarut/v"
)

const EPSILON = 1e-8

// const EPSILON float64 = 0.0

type AABB struct {
	Pos  Vec
	Half Vec
}

func (a AABB) Left() float64   { return a.Pos.X - a.Half.X }
func (a AABB) Right() float64  { return a.Pos.X + a.Half.X }
func (a AABB) Top() float64    { return a.Pos.Y - a.Half.Y }
func (a AABB) Bottom() float64 { return a.Pos.Y + a.Half.Y }
func (a AABB) TopLeft() Vec    { return Vec{a.Pos.X - a.Half.X, a.Pos.Y - a.Half.Y} }
func (a AABB) Size() Vec       { return Vec{a.Half.X * 2, a.Half.Y * 2} }

func (a *AABB) SetLeft(l float64)   { a.Pos.X = l + a.Half.X }
func (a *AABB) SetRight(r float64)  { a.Pos.X = r - a.Half.X }
func (a *AABB) SetTop(t float64)    { a.Pos.Y = t + a.Half.Y }
func (a *AABB) SetBottom(b float64) { a.Pos.Y = b - a.Half.Y }

type HitInfo struct {
	Pos    Vec
	Delta  Vec
	Normal Vec
	Time   float64
}

func (h *HitInfo) Reset() {
	*h = HitInfo{}
}

type HitInfo2 struct {
	Right, Bottom, Left, Top bool
	Delta                    Vec
}

func (h *HitInfo2) Reset() {
	*h = HitInfo2{}
}

// AABBPlatform moving platform b collision
func AABBPlatform(a, b *AABB, aVel, bVel *Vec, h *HitInfo2) bool {
	// Calculate old positions using velocities
	aOldPos := Vec{a.Pos.X - aVel.X, a.Pos.Y - aVel.Y}
	bOldPos := Vec{b.Pos.X - bVel.X, b.Pos.Y - bVel.Y}

	// Check collision at current positions using half dimensions
	xDist := math.Abs(a.Pos.X - b.Pos.X)
	yDist := math.Abs(a.Pos.Y - b.Pos.Y)

	// Combined half widths and heights
	combinedHalfW := a.Half.X + b.Half.X
	combinedHalfH := a.Half.Y + b.Half.Y

	// Early exit check
	if xDist > combinedHalfW || yDist > combinedHalfH {
		return false
	}

	// Calculate old distances using calculated old positions
	oldXDist := math.Abs(aOldPos.X - bOldPos.X)
	oldYDist := math.Abs(aOldPos.Y - bOldPos.Y)

	// Check collision direction and calculate position delta
	if yDist < combinedHalfH {
		if a.Pos.Y > b.Pos.Y && oldYDist >= combinedHalfH {
			h.Delta.Y = (b.Pos.Y + combinedHalfH + EPSILON) - a.Pos.Y
			h.Top = true
		} else if a.Pos.Y < b.Pos.Y && oldYDist >= combinedHalfH {
			h.Delta.Y = (b.Pos.Y - combinedHalfH - EPSILON) - a.Pos.Y
			h.Bottom = true
		}
	}

	if xDist < combinedHalfW {
		if a.Pos.X > b.Pos.X && oldXDist >= combinedHalfW {
			h.Delta.X = (b.Pos.X + combinedHalfW + EPSILON) - a.Pos.X
			h.Left = true
		} else if a.Pos.X < b.Pos.X && oldXDist >= combinedHalfW {
			h.Delta.X = (b.Pos.X - combinedHalfW - EPSILON) - a.Pos.X
			h.Right = true
		}
	}

	return true
}

func Segment(a *AABB, pos, delta, padding Vec, hit *HitInfo) bool {
	scale := v.One.Div(delta)
	signX := math.Copysign(1, scale.X)
	signY := math.Copysign(1, scale.Y)
	nearTimeX := (a.Pos.X - signX*(a.Half.X+padding.X) - pos.X) * scale.X
	nearTimeY := (a.Pos.Y - signY*(a.Half.Y+padding.Y) - pos.Y) * scale.Y
	farTimeX := (a.Pos.X + signX*(a.Half.X+padding.X) - pos.X) * scale.X
	farTimeY := (a.Pos.Y + signY*(a.Half.Y+padding.Y) - pos.Y) * scale.Y
	if math.IsNaN(nearTimeY) {
		nearTimeY = math.Inf(1)
	}
	if math.IsNaN(farTimeY) {
		farTimeY = math.Inf(1)
	}
	if nearTimeX > farTimeY || nearTimeY > farTimeX {
		return false
	}
	nearTime := max(nearTimeX, nearTimeY)
	farTime := min(farTimeX, farTimeY)
	if nearTime >= 1 || farTime <= 0 {
		return false
	}
	if hit == nil {
		return true
	}
	hit.Time = max(0, min(1, nearTime))

	if nearTimeX > nearTimeY {
		hit.Normal.X = -signX
		hit.Normal.Y = 0
	} else {
		hit.Normal.X = 0
		hit.Normal.Y = -signY
	}
	hit.Delta.X = (1.0 - hit.Time) * -delta.X
	hit.Delta.Y = (1.0 - hit.Time) * -delta.Y

	hit.Pos = pos.Add(delta.Scale(hit.Time))
	return true
}

// OverlapSweep returns hit info for b
func Overlap(a, b *AABB, hit *HitInfo) bool {
	dx := b.Pos.X - a.Pos.X
	px := b.Half.X + a.Half.X - math.Abs(dx)
	if px <= 0 {
		return false
	}

	dy := b.Pos.Y - a.Pos.Y
	py := b.Half.Y + a.Half.Y - math.Abs(dy)
	if py <= 0 {
		return false
	}

	if hit == nil {
		return true
	}

	// hit reset here
	if px < py {
		sx := math.Copysign(1, dx)
		hit.Delta.X = px * sx
		hit.Normal.X = sx
		hit.Pos.X = a.Pos.X + a.Half.X*sx
		hit.Pos.Y = b.Pos.Y
	} else {
		sy := math.Copysign(1, dy)
		hit.Delta.Y = py * sy
		hit.Normal.Y = sy
		hit.Pos.X = b.Pos.X
		hit.Pos.Y = a.Pos.Y + a.Half.Y*sy
	}
	return true
}

// OverlapSweep returns hit info for b
func OverlapSweep(staticA, b *AABB, bDelta Vec, hit *HitInfo) bool {
	if bDelta.IsZero() {
		return Overlap(staticA, b, hit)
	}
	result := Segment(staticA, b.Pos, bDelta, b.Half, hit)
	if result {
		hit.Time = max(0, min(1, hit.Time))
		direction := bDelta.Unit()
		hit.Pos.X = max(staticA.Pos.X-staticA.Half.X, min(staticA.Pos.X+staticA.Half.X, hit.Pos.X+direction.X*b.Half.X))
		hit.Pos.Y = max(hit.Pos.Y+direction.Y*b.Half.Y, min(staticA.Pos.Y-staticA.Half.Y, staticA.Pos.Y+staticA.Half.Y))
	}
	return result
}

// OverlapSweep2 returns hit info for b
func OverlapSweep2(a, b *AABB, aDelta, bDelta Vec, hit *HitInfo) bool {
	delta := bDelta.Sub(aDelta)
	isCollide := OverlapSweep(a, b, delta, hit)
	if isCollide {
		hit.Pos = hit.Pos.Add(aDelta.Scale(hit.Time))
		if hit.Normal.X != 0 {
			hit.Pos.X = b.Pos.X + (bDelta.X * hit.Time) - (hit.Normal.X * b.Half.X)
		} else {
			hit.Pos.Y = b.Pos.Y + (bDelta.Y * hit.Time) - (hit.Normal.Y * b.Half.Y)
		}
	}
	return isCollide
}

// HitTileInfo stores information about a collision with a tile
type HitTileInfo struct {
	TileCoords image.Point // X,Y coordinates of the tile in the tilemap
	Normal     Vec         // Normal vector of the collision (-1/0/1)
}

// Collider handles collision detection between rectangles and a 2D tilemap
type Collider struct {
	Collisions []HitTileInfo // List of collisions from last check
	TileSize   image.Point   // Width and height of tiles
	TileMap    [][]uint8     // 2D grid of tile IDs
}

// CollisionCallback is called when collisions occur, receiving collision info and final movement
type CollisionCallback func(hitInfos []HitTileInfo, delta Vec)

// Collide checks for collisions when moving a rectangle and returns the allowed movement
func (c *Collider) Collide(rect AABB, delta Vec, onCollide CollisionCallback) Vec {
	c.Collisions = c.Collisions[:0]

	if delta.IsZero() {
		return delta
	}

	if math.Abs(delta.X) > math.Abs(delta.Y) {
		if delta.X != 0 {
			delta.X = c.CollideX(&rect, delta.X)
		}
		if delta.Y != 0 {
			rect.Pos.X += delta.X
			delta.Y = c.CollideY(&rect, delta.Y)
		}
	} else {
		if delta.Y != 0 {
			delta.Y = c.CollideY(&rect, delta.Y)
		}
		if delta.X != 0 {
			rect.Pos.Y += delta.Y
			delta.X = c.CollideX(&rect, delta.X)
		}
	}

	if onCollide != nil {
		onCollide(c.Collisions, delta)
	}
	return delta
}

// CollideX checks for collisions along the X axis and returns the allowed X movement
//
// reset list before -> tileCollider.Collisions = tileCollider.Collisions[:0]
func (c *Collider) CollideX(rect *AABB, deltaX float64) float64 {
	checkLimit := max(1, int(math.Ceil(math.Abs(deltaX)/float64(c.TileSize.Y)))+1)
	rectTileTopCoord := int(math.Floor(rect.Top() / float64(c.TileSize.Y)))
	rectTileBottomCoord := int(math.Ceil(rect.Bottom()/float64(c.TileSize.Y))) - 1
	if deltaX > 0 {
		rectRight := rect.Right()
		startRightX := int(math.Floor(rectRight / float64(c.TileSize.X)))
		endX := startRightX + checkLimit
		endX = min(endX, len(c.TileMap[0]))
		for y := rectTileTopCoord; y <= rectTileBottomCoord; y++ {
			if y < 0 || y >= len(c.TileMap) {
				continue
			}
			for x := startRightX; x < endX; x++ {
				if x < 0 || x >= len(c.TileMap[0]) {
					continue
				}
				if !items.HasTag(c.TileMap[y][x], items.NonSolidBlock) {
					tileLeft := float64(x * c.TileSize.X)
					collision := tileLeft - rectRight
					if collision <= deltaX {
						deltaX = collision
						c.Collisions = append(c.Collisions, HitTileInfo{
							TileCoords: image.Point{x, y},
							Normal:     v.Left,
						})
					}
				}
			}
		}
	}
	if deltaX < 0 {
		rectLeft := rect.Left()
		endX := int(math.Floor(rectLeft / float64(c.TileSize.X)))
		startX := endX - checkLimit
		startX = max(startX, 0)
		for y := rectTileTopCoord; y <= rectTileBottomCoord; y++ {
			if y < 0 || y >= len(c.TileMap) {
				continue
			}
			for x := startX; x <= endX; x++ {
				if x < 0 || x >= len(c.TileMap[0]) {
					continue
				}
				if !items.HasTag(c.TileMap[y][x], items.NonSolidBlock) {
					tileRight := float64((x + 1) * c.TileSize.X)
					collision := tileRight - rectLeft
					if collision >= deltaX {
						deltaX = collision
						c.Collisions = append(c.Collisions, HitTileInfo{
							TileCoords: image.Point{x, y},
							Normal:     v.Right,
						})
					}
				}
			}
		}
	}
	return deltaX
}

// CollideY checks for collisions along the Y axis and returns the allowed Y movement
//
// reset list before -> tileCollider.Collisions = tileCollider.Collisions[:0]
func (c *Collider) CollideY(rect *AABB, deltaY float64) float64 {
	checkLimit := max(1, int(math.Ceil(math.Abs(deltaY)/float64(c.TileSize.Y)))+1)
	rectTileLeftCoord := int(math.Floor(rect.Left() / float64(c.TileSize.X)))
	rectTileRightCoord := int(math.Ceil(rect.Right()/float64(c.TileSize.X))) - 1
	if deltaY > 0 {
		rectBottom := rect.Bottom()
		startBottomY := int(math.Floor(rectBottom / float64(c.TileSize.Y)))
		endY := startBottomY + checkLimit
		endY = min(endY, len(c.TileMap))

		for x := rectTileLeftCoord; x <= rectTileRightCoord; x++ {
			if x < 0 || x >= len(c.TileMap[0]) {
				continue
			}
			for y := startBottomY; y < endY; y++ {
				if y < 0 || y >= len(c.TileMap) {
					continue
				}
				if !items.HasTag(c.TileMap[y][x], items.NonSolidBlock) {
					tileTop := float64(y * c.TileSize.Y)
					collision := tileTop - rectBottom
					if collision <= deltaY {
						deltaY = collision
						c.Collisions = append(c.Collisions, HitTileInfo{
							TileCoords: image.Point{x, y},
							Normal:     v.Up,
						})
					}
				}
			}
		}
	}
	if deltaY < 0 {
		rectTop := rect.Top()
		endY := int(math.Floor(rectTop / float64(c.TileSize.Y)))
		startY := endY - checkLimit
		startY = max(startY, 0)
		for x := rectTileLeftCoord; x <= rectTileRightCoord; x++ {
			if x < 0 || x >= len(c.TileMap[0]) {
				continue
			}
			for y := startY; y <= endY; y++ {
				if y < 0 || y >= len(c.TileMap) {
					continue
				}
				if !items.HasTag(c.TileMap[y][x], items.NonSolidBlock) {
					tileBottom := float64((y + 1) * c.TileSize.Y)
					collision := tileBottom - rectTop
					if collision >= deltaY {
						deltaY = collision
						c.Collisions = append(c.Collisions, HitTileInfo{
							TileCoords: image.Point{x, y},
							Normal:     v.Down,
						})
					}
				}
			}
		}
	}
	return deltaY
}
