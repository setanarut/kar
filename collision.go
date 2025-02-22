package kar

import (
	"image"
	"kar/items"
	"kar/v"
	"math"
)

const EPSILON = 1e-8

type AABB struct {
	Pos  Vec
	Half Vec
}

type HitInfo struct {
	Pos    Vec
	Delta  Vec
	Normal Vec
	Time   float64
}

// HitTileInfo stores information about a collision with a tile
type HitTileInfo struct {
	TileCoords image.Point // X,Y coordinates of the tile in the tilemap
	Normal     Vec         // Normal vector of the collision (-1/0/1)
}

func (a AABB) Segment(pos, delta, padding Vec, hit *HitInfo) bool {
	scaleX := 1.0 / delta.X
	scaleY := 1.0 / delta.Y
	signX := sign(scaleX)
	signY := sign(scaleY)
	nearTimeX := (a.Pos.X - signX*(a.Half.X+padding.X) - pos.X) * scaleX
	nearTimeY := (a.Pos.Y - signY*(a.Half.Y+padding.Y) - pos.Y) * scaleY
	farTimeX := (a.Pos.X + signX*(a.Half.X+padding.X) - pos.X) * scaleX
	farTimeY := (a.Pos.Y + signY*(a.Half.Y+padding.Y) - pos.Y) * scaleY
	if math.IsNaN(nearTimeY) {
		nearTimeY = math.Inf(1)
	}
	if math.IsNaN(farTimeY) {
		farTimeY = math.Inf(1)
	}
	if nearTimeX > farTimeY || nearTimeY > farTimeX {
		return false
	}
	nearTime := math.Max(nearTimeX, nearTimeY)
	farTime := math.Min(farTimeX, farTimeY)
	if nearTime >= 1 || farTime <= 0 {
		return false
	}
	if hit == nil {
		return true
	}
	hit.Time = clamp(nearTime, 0, 1)
	if nearTimeX > nearTimeY {
		hit.Normal.X = -signX
		hit.Normal.Y = 0
	} else {
		hit.Normal.X = 0
		hit.Normal.Y = -signY
	}
	hit.Delta.X = (1.0 - hit.Time) * -delta.X
	hit.Delta.Y = (1.0 - hit.Time) * -delta.Y
	hit.Pos.X = pos.X + delta.X*hit.Time
	hit.Pos.Y = pos.Y + delta.Y*hit.Time
	return true
}

func (a AABB) Overlap(a2 AABB, hit *HitInfo) bool {
	dx := a2.Pos.X - a.Pos.X
	px := a2.Half.X + a.Half.X - math.Abs(dx)
	if px <= 0 {
		return false
	}

	dy := a2.Pos.Y - a.Pos.Y
	py := a2.Half.Y + a.Half.Y - math.Abs(dy)
	if py <= 0 {
		return false
	}

	if hit == nil {
		return true
	}

	// hit.Collider = box1
	hit.Delta = Vec{}
	hit.Normal = Vec{}
	hit.Time = 0 // boxes overlap
	if px < py {
		sx := sign(dx)
		hit.Delta.X = px * sx
		hit.Normal.X = sx
		hit.Pos.X = a.Pos.X + a.Half.X*sx
		hit.Pos.Y = a2.Pos.Y
	} else {
		sy := sign(dy)
		hit.Delta.Y = py * sy
		hit.Normal.Y = sy
		hit.Pos.X = a2.Pos.X
		hit.Pos.Y = a.Pos.Y + a.Half.Y*sy
	}
	return true
}

func (a AABB) OverlapSweep(a2 AABB, delta Vec, hit *HitInfo) bool {
	if delta.IsZero() {
		return a.Overlap(a2, hit)
	}
	result := a.Segment(a2.Pos, delta, a2.Half, hit)
	if result {
		// hit.Time = 1.0
		hit.Time = clamp(hit.Time-EPSILON, 0, 1)
		direction := delta.Unit()
		hit.Pos.X = clamp(
			hit.Pos.X+direction.X*a2.Half.X,
			a.Pos.X-a.Half.X,
			a.Pos.X+a.Half.X,
		)
		hit.Pos.Y = clamp(
			hit.Pos.Y+direction.Y*a2.Half.Y,
			a.Pos.Y-a.Half.Y,
			a.Pos.Y+a.Half.Y,
		)
	}
	return result
}

// Collider handles collision detection between rectangles and a 2D tilemap
type Collider struct {
	Collisions []HitTileInfo // List of collisions from last check
	TileSize   image.Point   // Width and height of tiles
	TileMap    [][]uint8     // 2D grid of tile IDs
}

// NewCollider creates a new tile collider with the given tilemap and tile dimensions
func NewCollider(tileMap [][]uint8, tileWidth, tileHeight int) *Collider {
	return &Collider{
		TileMap:  tileMap,
		TileSize: image.Point{tileWidth, tileHeight},
	}
}

// CollisionCallback is called when collisions occur, receiving collision info and final movement
type CollisionCallback func([]HitTileInfo, Vec)

// Collide checks for collisions when moving a rectangle and returns the allowed movement
func (c *Collider) Collide(rect AABB, delta Vec, onCollide CollisionCallback) Vec {
	c.Collisions = c.Collisions[:0]

	if delta.IsZero() {
		return delta
	}

	if math.Abs(delta.X) > math.Abs(delta.Y) {
		if delta.X != 0 {
			delta.X = c.CollideX(rect, delta.X)
		}
		if delta.Y != 0 {
			rect.Pos.X += delta.X
			delta.Y = c.CollideY(rect, delta.Y)
		}
	} else {
		if delta.Y != 0 {
			delta.Y = c.CollideY(rect, delta.Y)
		}
		if delta.X != 0 {

			rect.Pos.Y += delta.Y
			delta.X = c.CollideX(rect, delta.X)
		}
	}

	if onCollide != nil {
		onCollide(c.Collisions, delta)
	}
	return delta
}

// CollideX checks for collisions along the X axis and returns the allowed X movement
func (c *Collider) CollideX(rect AABB, deltaX float64) float64 {
	checkLimit := max(1, int(math.Ceil(math.Abs(deltaX)/float64(c.TileSize.Y)))+1)

	rectTop := rect.Pos.Y - rect.Half.Y
	rectBottom := rect.Pos.Y + rect.Half.Y

	rectTileTopCoord := int(math.Floor(rectTop / float64(c.TileSize.Y)))
	rectTileBottomCoord := int(math.Ceil((rectBottom)/float64(c.TileSize.Y))) - 1

	if deltaX > 0 {
		startRightX := int(math.Floor((rect.Pos.X + rect.Half.X) / float64(c.TileSize.X)))
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
					collision := tileLeft - (rect.Pos.X + rect.Half.X)
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
		rectLeft := rect.Pos.X - rect.Half.X

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
func (c *Collider) CollideY(rect AABB, deltaY float64) float64 {

	checkLimit := max(1, int(math.Ceil(math.Abs(deltaY)/float64(c.TileSize.Y)))+1)

	rectLeft := rect.Pos.X - rect.Half.X
	rectRight := rect.Pos.X + rect.Half.X

	rectTileLeftCoord := int(math.Floor(rectLeft / float64(c.TileSize.X)))
	rectTileRightCoord := int(math.Ceil(rectRight/float64(c.TileSize.X))) - 1

	if deltaY > 0 {
		rectBottom := rect.Pos.Y + rect.Half.Y
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
		rectTop := rect.Pos.Y - rect.Half.Y
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

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

func sign(value float64) float64 {
	if value < 0 {
		return -1
	}
	return 1
}
