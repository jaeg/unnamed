package world

import (
	"fmt"
	"math"
	"math/rand"

	"unnamed/component"
	"unnamed/entity"
)

// Level .
type Level struct {
	data                  [][]TileSmall
	Entities              []*entity.Entity
	width, height         int
	id                    int
	left, right, up, down int
	theme                 string
}

//Tile .
type Tile struct {
	TileIndex   int
	SpriteIndex int
	Solid       bool
	StairsTo    int
	StairsX     int
	StairsY     int
	Water       bool
	Locked      bool
}

//Tile Types
//1 - open
//2 - solid
//3 - stairs [level id, to x, to y]
//4 - water
type TileSmall struct {
	TileIndex int
	Type      int
	Data      []int
	X         int
	Y         int
}

func newLevel(width int, height int) (level *Level) {
	level = &Level{width: width, height: height, left: -1, right: -1, up: -1, down: -1}

	data := make([][]TileSmall, width, height)
	for x := 0; x < width; x++ {
		col := []TileSmall{}
		for y := 0; y < height; y++ {
			col = append(col, TileSmall{TileIndex: 112, Type: 1, X: x, Y: y})
		}
		data[x] = append(data[x], col...)
	}

	level.data = data
	return
}

func NewOverworldSection(width int, height int) (level *Level) {
	fmt.Println("Creating new random level")
	level = newLevel(width, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if rand.Intn(1000) == 0 {
				level.GetTileAt(x, y).TileIndex = 13
				level.GetTileAt(x, y).Type = 3
				level.GetTileAt(x, y).Data = []int{-1, 0, 0}
			} else if rand.Intn(5) == 0 {
				level.GetTileAt(x, y).TileIndex = 121
			} else {
				level.GetTileAt(x, y).TileIndex = 122
			}
		}
	}

	//Generate Flower Medows
	for i := 0; i < 50; i++ {
		x := getRandom(1, width)
		y := getRandom(1, height)

		level.createCluster(x, y, 10, 123, 0, false, false)
	}

	//Generate Water
	for i := 0; i < 100; i++ {
		x := getRandom(1, width)
		y := getRandom(1, height)

		level.createCluster(x, y, 100, 181, 0, false, true)
	}

	return
}

func NewDungeon(width int, height int, stairsUpTo int, fromX int, fromY int) (level *Level, pX int, pY int) {
	fmt.Println("Creating dungeon")
	level = newLevel(width, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if rand.Intn(1000) == 0 {
				level.GetTileAt(x, y).TileIndex = 12
				level.GetTileAt(x, y).Type = 3
				level.GetTileAt(x, y).Data = []int{-1, 0, 0}
			} else if rand.Intn(5) == 0 {
				level.GetTileAt(x, y).TileIndex = 8
				level.GetTileAt(x, y).Type = 1
			} else {
				level.GetTileAt(x, y).TileIndex = 7
				level.GetTileAt(x, y).Type = 1
			}
		}
	}

	pX = getRandom(0, width)
	pY = getRandom(0, height)
	level.GetTileAt(pX, pY).Data = []int{stairsUpTo, fromX, fromY}
	level.GetTileAt(pX, pY).TileIndex = 13
	level.GetTileAt(pX, pY).Type = 3

	return
}

func (level *Level) GetTileAt(x int, y int) (tile *TileSmall) {
	if x < level.width && y < level.height && x >= 0 && y >= 0 {
		tile = &level.data[x][y]
	}
	return
}

//Get's the view frustum with the player in the center
func (level *Level) GetView(aX int, aY int, width int, height int, blind bool) (data [][]*TileSmall) {
	left := aX - width/2
	right := aX + width/2
	up := aY - height/2
	down := aY + height/2

	data = make([][]*TileSmall, width+1-width%2)

	cX := 0
	for x := left; x <= right; x++ {
		col := []*TileSmall{}
		for y := up; y <= down; y++ {
			currentTile := level.GetTileAt(x, y)
			if blind {
				if y < aY-height/4 || y > aY+height/4 || x > aX+width/4 || x < aX-width/4 {
					currentTile = nil
				}
			}

			if currentTile != nil {
				if los(aX, aY, x, y, level) == false {
					currentTile = nil
				}
			}

			col = append(col, currentTile)
		}
		data[cX] = append(data[cX], col...)
		cX++
	}
	return
}

func (level *Level) GetEntitiesAround(x int, y int, width int, height int) (entities []*entity.Entity) {
	left := x - width/2
	right := x + width/2
	up := y - height/2
	down := y + height/2

	entitiesLen := len(level.Entities)

	for i := 0; i < entitiesLen; i++ {
		entity := level.Entities[i]

		if entity.HasComponent("PositionComponent") {
			pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
			if pc.X >= left && pc.X <= right && pc.Y >= up && pc.Y <= down {
				entities = append(entities, entity)
			}
		}
	}

	return
}

func (level *Level) GetPlayersAround(x int, y int, width int, height int) (entities []*entity.Entity) {
	left := x - width/2
	right := x + width/2
	up := y - height/2
	down := y + height/2

	entitiesLen := len(level.Entities)

	for i := 0; i < entitiesLen; i++ {
		entity := level.Entities[i]

		if entity.HasComponent("PositionComponent") {
			pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
			if pc.X >= left && pc.X <= right && pc.Y >= up && pc.Y <= down {
				if entity.HasComponent("PlayerComponent") {
					entities = append(entities, entity)
				}
			}
		}
	}

	return
}

func (level *Level) GetEntityAt(x int, y int) (entity *entity.Entity) {
	for i := 0; i < len(level.Entities); i++ {
		entity = level.Entities[i]
		if entity.HasComponent("PositionComponent") {
			pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
			if pc.X == x && pc.Y == y {
				return
			}
		}
	}
	entity = nil
	return
}

func (level *Level) GetSolidEntityAt(x int, y int) (entity *entity.Entity) {
	for i := 0; i < len(level.Entities); i++ {
		entity = level.Entities[i]
		if entity.HasComponent("PositionComponent") {
			if entity.HasComponent("SolidComponent") {
				pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
				if pc.X == x && pc.Y == y {
					return
				}
			}
		}
	}
	entity = nil
	return
}

func (level *Level) GetInteractableEntityAt(x int, y int) (entity *entity.Entity) {
	for i := 0; i < len(level.Entities); i++ {
		entity = level.Entities[i]
		if entity.HasComponent("PositionComponent") {
			if entity.HasComponent("InteractComponent") {
				pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
				if pc.X == x && pc.Y == y {
					return
				}
			}
		}
	}
	entity = nil
	return
}

func (level *Level) AddEntity(entity *entity.Entity) {
	level.Entities = append(level.Entities, entity)
}

func (level *Level) RemoveEntity(entity *entity.Entity) {
	for i := 0; i < len(level.Entities); i++ {
		if level.Entities[i] == entity {
			level.Entities = append(level.Entities[:i], level.Entities[i+1:]...)

		}
	}
}

func (level *Level) createCluster(x int, y int, size int, tileIndex int, spriteIndex int, solid bool, water bool) {
	for i := 0; i < 200; i++ {
		n := getRandom(1, 6)
		e := getRandom(1, 6)
		s := getRandom(1, 6)
		w := getRandom(1, 6)

		if n == 1 {
			x += 1
		}

		if s == 1 {
			x--
		}

		if e == 1 {
			y++
		}

		if w == 1 {
			y--
		}

		if level.GetTileAt(x, y) != nil {
			tile := level.GetTileAt(x, y)
			tile.TileIndex = tileIndex

			if solid {
				tile.Type = 2
			} else if water {
				tile.Type = 4
			} else {
				tile.Type = 1
			}

		}
	}
}

func getRandom(low int, high int) int {
	if low == high {
		return low
	}
	return (rand.Intn((high - low))) + low
}

func Sgn(a int) int {
	switch {
	case a < 0:
		return -1
	case a > 0:
		return +1
	}
	return 0
}

//Ported from http://www.roguebasin.com/index.php?title=Simple_Line_of_Sight
func los(pX int, pY int, tX int, tY int, level *Level) bool {
	deltaX := pX - tX
	deltaY := pY - tY

	absDeltaX := math.Abs(float64(deltaX))
	absDeltaY := math.Abs(float64(deltaY))

	signX := Sgn(deltaX)
	signY := Sgn(deltaY)

	if absDeltaX > absDeltaY {
		t := absDeltaY*2 - absDeltaX
		for {
			if t >= 0 {
				tY += signY
				t -= absDeltaX * 2
			}

			tX += signX
			t += absDeltaY * 2

			if tX == pX && tY == pY {
				return true
			}
			if level.GetTileAt(tX, tY).Type == 2 {
				break
			}
		}
		return false
	}

	t := absDeltaX*2 - absDeltaY

	for {
		if t >= 0 {
			tX += signX
			t -= absDeltaY * 2
		}
		tY += signY
		t += absDeltaX * 2
		if tX == pX && tY == pY {
			return true
		}

		if level.GetTileAt(tX, tY).Type == 2 {
			break
		}
	}

	return false

}

// Ported from here: https://github.com/tome2/tome2/blob/master/src/generate.cc
func (level *Level) buildRecursiveRoom(x1 int, y1 int, x2 int, y2 int, power int) {
	xSize := x2 - x1
	ySize := y2 - y1

	if xSize < 0 || ySize < 0 {
		return
	}

	var choice int
	if power < 3 && xSize > 12 && ySize > 12 {
		choice = 1
	} else {
		if power < 10 {
			if getRandom(0, 10) > 2 && xSize < 8 && ySize < 8 {
				choice = 4
			} else {
				choice = getRandom(0, 2) + 1
			}
		} else {
			choice = getRandom(0, 3) + 1
		}
	}

	if choice == 1 {
		//Outer walls
		for x := x1; x <= x2; x++ {
			level.GetTileAt(x, y1).TileIndex = 0
			level.GetTileAt(x, y2).TileIndex = 0
		}

		for y := y1 + 1; y < y2; y++ {
			level.GetTileAt(x1, y).TileIndex = 6
			level.GetTileAt(x2, y).TileIndex = 6
		}

		if getRandom(0, 2) == 0 {
			y := getRandom(0, ySize) + y1
			level.GetTileAt(x1, y).TileIndex = 121
			level.GetTileAt(x2, y).TileIndex = 121
		} else {
			x := getRandom(0, xSize) + x1
			level.GetTileAt(x, y1).TileIndex = 121
			level.GetTileAt(x, y2).TileIndex = 121
		}

		//Size of keep
		t1 := getRandom(0, ySize/3) + y1
		t2 := y2 - getRandom(0, ySize/3)
		t3 := getRandom(0, xSize/3) + x1
		t4 := x2 - getRandom(0, xSize/3)

		//Above and below
		level.buildRecursiveRoom(x1+1, y1+1, x2-1, t1, power+1)
		level.buildRecursiveRoom(x1+1, t2, x2-1, y2, power+1)

		//Left and right
		level.buildRecursiveRoom(x1+1, t1+1, t3, t2-1, power+3)
		level.buildRecursiveRoom(t4, t1+1, x2-1, t2-1, power+3)

		x1 = t3
		x2 = t4
		y1 = t1
		y2 = t2
		xSize = x2 - x1
		ySize = y2 - y1
		power += 2
	}

	if choice == 4 || choice == 1 {
		if xSize < 3 || ySize < 3 {
			for y := y1; y < y2; y++ {
				for x := x1; x < x2; x++ {
					level.GetTileAt(x, y).TileIndex = 165
				}
			}

			return
		}

		//make outside walls
		for x := x1 + 1; x <= x2-1; x++ {
			level.GetTileAt(x, y1+1).TileIndex = 165
			level.GetTileAt(x, y2-1).TileIndex = 165
		}

		for y := y1 + 1; y < y2-1; y++ {
			level.GetTileAt(x1+1, y).TileIndex = 165
			level.GetTileAt(x2-1, y).TileIndex = 165
		}

		//Make door
		y := getRandom(0, ySize-3) + y1 + 2
		if getRandom(0, 2) == 0 {
			level.GetTileAt(x1+1, y).TileIndex = 123
		} else {
			level.GetTileAt(x2-1, y).TileIndex = 123
		}

		level.buildRecursiveRoom(x1+2, y1+2, x2-2, y2-2, power+3)
	}

	if choice == 2 {
		if xSize < 3 {
			for y := y1; y < y2; y++ {
				for x := x1; x < x2; x++ {
					level.GetTileAt(x, y).TileIndex = 165
				}
			}

			return
		}

		t1 := getRandom(0, xSize-2) + x1 + 1
		level.buildRecursiveRoom(x1, y1, t1, y2, power-2)
		level.buildRecursiveRoom(t1+1, y1, x2, y2, power-2)
	}

	if choice == 3 {
		if ySize < 3 {
			for y := y1; y < y2; y++ {
				for x := x1; x < x2; x++ {
					level.GetTileAt(x, y).TileIndex = 165
				}
			}

			return
		}

		t1 := getRandom(0, ySize-2) + y1 + 1
		level.buildRecursiveRoom(x1, y1, x2, t1, power-2)
		level.buildRecursiveRoom(x1, t1+1, x2, y2, power-2)
	}
}

func distance(x1 int, y1 int, x2 int, y2 int) int {
	var dy int
	if y1 > y2 {
		dy = y1 - y2
	} else {
		dy = (y2 - y1)
	}

	var dx int
	if x1 > x2 {
		dx = x1 - x2
	} else {
		dx = x2 - x1
	}

	var d int
	if dy > dx {
		d = dy + (dx >> 1)
	} else {
		d = dx + (dy >> 1)
	}

	return d
}

type coords struct {
	x, y int
}

func (level *Level) buildStoreCircle(qx int, qy int, xx int, yy int) coords {
	rad := 2 + getRandom(0, 2)

	y0 := qy + yy*9 + 6
	x0 := qx + xx*14 + 12

	//Building
	for y := y0 - rad; y <= y0+rad; y++ {
		for x := x0 - rad; x <= x0+rad; x++ {
			if distance(x0, y0, x, y) > rad {
				continue
			}
			level.GetTileAt(x, y).TileIndex = 165
		}
	}

	//Door location
	tmp := getRandom(0, 4)

	if ((tmp == 0) && (yy == 1)) ||
		((tmp == 1) && (yy == 0)) ||
		((tmp == 2) && (xx == 3)) ||
		((tmp == 3) && (xx == 0)) {
		/* Pick a new direction */
		tmp = getRandom(0, 4)
	}

	coords := coords{}
	switch tmp {
	case 0:
		for y := y0; y <= y0+rad; y++ {
			level.GetTileAt(x0, y).TileIndex = 123
		}
		coords.x = x0
		coords.y = y0 + rad

	case 1:
		for y := y0 - rad; y <= y0; y++ {
			level.GetTileAt(x0, y).TileIndex = 123
		}
		coords.x = x0
		coords.y = y0 - rad

	case 2:
		for x := x0; x <= x0+rad; x++ {
			level.GetTileAt(x, y0).TileIndex = 123
		}
		coords.x = x0 + rad
		coords.y = y0

	case 3:
		for x := x0 - rad; x <= x0; x++ {
			level.GetTileAt(x, y0).TileIndex = 123
		}
		coords.x = x0 - rad
		coords.y = y0
	}

	return coords
}

func (level *Level) townGenCircle(qx int, qy int, width int, height int) {
	var buildings []coords
	//buildings
	for y := 0; y < 2; y++ {
		for x := 0; x < 4; x++ {
			buildings = append(buildings, level.buildStoreCircle(qx, qy, x, y))
		}
	}
}
