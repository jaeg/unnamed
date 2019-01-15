package system

import (
	"unnamed/component"
	"unnamed/world"

	"github.com/nsf/termbox-go"
)

type entityView struct {
	X, Y int
	Char string
}

// RenderSystem .
func RenderSystem(planets map[string]*world.Planet) {
	//	os.Stdout.Write([]byte{0x1B, 0x5B, 0x33, 0x3B, 0x4A, 0x1B, 0x5B, 0x48, 0x1B, 0x5B, 0x32, 0x4A})
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	for _, planet := range planets {
		for _, level := range planet.Levels {
			var seeableEntities []entityView
			for _, entity := range level.Entities {
				if entity.HasComponent("AppearanceComponent") {
					ac := entity.GetComponent("AppearanceComponent").(*component.AppearanceComponent)
					pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)
					ev := entityView{X: pc.X, Y: pc.Y, Char: ac.Char}
					seeableEntities = append(seeableEntities, ev)
				}
			}

			for _, entity := range level.Entities {
				if entity.HasComponent("PlayerComponent") {
					viewWidth := 20
					viewHeight := 20
					pc := entity.GetComponent("PlayerComponent").(*component.PlayerComponent)
					positionComponent := entity.GetComponent("PositionComponent").(*component.PositionComponent)
					view := level.GetView(positionComponent.X, positionComponent.Y, viewWidth, viewHeight, false)
					/*fmt.Println("Pos:", positionComponent.X, positionComponent.Y)
					fmt.Println("View size", len(view), len(view[0]))*/
					if len(pc.MessageLog) > 0 {
						print_tb(0, viewHeight+5, termbox.ColorWhite, termbox.ColorBlack, pc.MessageLog[len(pc.MessageLog)-1])
					}
					for y := 0; y < len(view[0]); y++ {
						for x := 0; x < len(view); x++ {
							tile := view[x][y]
							if tile == nil {
								print_tb(x, y, termbox.ColorWhite, termbox.ColorBlack, "-")
							} else {
								if positionComponent.X == tile.X && positionComponent.Y == tile.Y {
									print_tb(x, y, termbox.ColorWhite, termbox.ColorBlack, "@")
								} else {
									drawTile := true
									for _, entity := range seeableEntities {
										if entity.X == tile.X && entity.Y == tile.Y {
											if drawTile {
												print_tb(x, y, termbox.ColorWhite, termbox.ColorBlack, entity.Char)
												drawTile = false
											}
										}
									}
									if drawTile {
										print_tb(x, y, termbox.ColorWhite, termbox.ColorBlack, "#")
									}
								}
							}
						}
					}
				}
			}
		}
	}

	termbox.Flush()
}

func print_tb(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}
