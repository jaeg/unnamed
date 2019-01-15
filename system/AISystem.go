package system

import (
	"math/rand"

	"unnamed/world"

	"unnamed/component"
)

func getRandom(low int, high int) int {
	return (rand.Intn((high - low))) + low
}

// PlayerSystem .
func AISystem(planets map[string]*world.Planet) {
	for _, planet := range planets {
		for _, level := range planet.Levels {
			//fmt.Println(t, len(level.Entities))
			for _, entity := range level.Entities {

				if entity.HasComponent("WanderAIComponent") {
					if entity.HasComponent("MyTurnComponent") {
						pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
						dc := entity.GetComponent("DirectionComponent").(*component.DirectionComponent)

						deltaX := getRandom(-1, 2)
						deltaY := 0
						if deltaX == 0 {
							deltaY = getRandom(-1, 2)
						}
						if level.GetSolidEntityAt(pc.X+deltaX, pc.Y+deltaY) == nil {
							tile := level.GetTileAt(pc.X+deltaX, pc.Y+deltaY)
							if tile == nil {
							} else if tile.Type != 2 && tile.Type != 4 {
								pc.X += deltaX
								pc.Y += deltaY
							}
						}
						if deltaY > 0 {
							dc.Direction = 1
						}
						if deltaY < 0 {
							dc.Direction = 2
						}
						if deltaX < 0 {
							dc.Direction = 3
						}
						if deltaX > 0 {
							dc.Direction = 0
						}
					}
				}
			}
		}
	}
}
