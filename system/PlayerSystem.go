package system

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unnamed/component"
	"unnamed/world"
)

// PlayerSystem .
func PlayerSystem(planets map[string]*world.Planet) map[string]*world.Planet {
	for _, planet := range planets {
		for currentLevel, level := range planet.Levels {
			//fmt.Println(t, len(level.Entities))
			for _, entity := range level.Entities {

				if entity.HasComponent("PlayerComponent") {
					playerComponent := entity.GetComponent("PlayerComponent").(*component.PlayerComponent)
					pc := entity.GetComponent("PositionComponent").(*component.PositionComponent)

					if entity.InteractingWith != nil {
						//Handle interactions
						switch entity.InteractingWith.GetType() {
						case "ShopComponent":
							sC := entity.InteractingWith.(*component.ShopComponent)
							if entity.Shown == false {
								playerComponent.AddMessage("Shop for all your shopping needs!")
								var items string
								for i, item := range sC.ItemsForSale {
									items += strconv.Itoa(i) + ": " + item
								}
								playerComponent.AddMessage(items)
								entity.Shown = true
							}

							input := getInput()
							if input == "Q" {
								entity.Shown = false
								entity.InteractingWith = nil
								playerComponent.AddMessage("Goodbye!")
							}

							if choice64, err := strconv.ParseInt(input, 10, 64); err == nil {
								choice := int(choice64)
								if len(sC.ItemsForSale) > choice {
									inventoryComponent := entity.GetComponent("InventoryComponent").(*component.InventoryComponent)
									inventoryComponent.AddItem(sC.ItemsForSale[choice])
									playerComponent.AddMessage("Thank you, here you go.")
								}
							}

						case "InteractComponent":
							ic := entity.InteractingWith.(*component.InteractComponent)
							playerComponent.AddMessage(ic.Message[getRandom(0, len(ic.Message))])
							entity.InteractingWith = nil
						}
					} else if entity.HasComponent("MyTurnComponent") {
						dc := entity.GetComponent("DirectionComponent").(*component.DirectionComponent)
						command := getInput()
						deltaX := 0
						deltaY := 0
						switch command {
						case "E": //Interact
							direction := getInput()
							if direction == "" {
								playerComponent.AddMessage("Wasn't given a direction to interact!")
							} else {
								x := pc.X
								y := pc.Y
								switch direction {
								case "W":
									y -= 1
								case "S":
									y += 1
								case "A":
									x -= 1
								case "D":
									x += 1
								}
								fmt.Println("GetAt", x, y)
								interactEntity := level.GetEntityAt(x, y)
								if interactEntity != nil {
									if interactEntity.HasComponent("ShopComponent") {
										sC := interactEntity.GetComponent("ShopComponent")
										entity.InteractingWith = sC
									} else if interactEntity.HasComponent("InteractComponent") {
										sC := interactEntity.GetComponent("InteractComponent")
										entity.InteractingWith = sC
									}
								} else {
									playerComponent.AddMessage("Nothing to interact with here")
								}
							}
						case "W":
							deltaY = -1
							dc.Direction = 2
						case "S":
							deltaY = 1
							dc.Direction = 1
						case "A":
							deltaX = -1
							dc.Direction = 3
						case "D":
							deltaX = 1
							dc.Direction = 0
						case "F":
							direction := getInput()
							if direction == "" {
								playerComponent.AddMessage("Wasn't given a direction to shoot!")
							} else {
								playerComponent.AddMessage("Shoot in the " + direction + " direction!")
							}
						}

						if deltaX != 0 || deltaY != 0 {
							tile := level.GetTileAt(pc.X+deltaX, pc.Y+deltaY)
							interactable := level.GetInteractableEntityAt(pc.X+deltaX, pc.Y+deltaY)
							if interactable != nil {
								ic := interactable.GetComponent("InteractComponent").(*component.InteractComponent)
								playerComponent.AddMessage(ic.Message[getRandom(0, len(ic.Message))])
							}
							if level.GetSolidEntityAt(pc.X+deltaX, pc.Y+deltaY) == nil {
								if tile == nil {
									playerComponent.AddMessage("You've hit the edge of the world!")
								} else if tile.Type == 2 || tile.Type == 4 {
									playerComponent.AddMessage("You can't walk that way!")
								} else {
									pc.Y += deltaY
									pc.X += deltaX

									//Stairs
									if tile.Type == 3 {
										if tile.Data[0] == -1 {
											levelIndex := len(planet.Levels)
											newLevel, tX, tY := world.NewDungeon(100, 100, currentLevel, pc.X, pc.Y)
											planet.Levels = append(planet.Levels, newLevel)
											tile.Type = 3
											tile.Data = []int{levelIndex, tX, tY}
										}
										level.RemoveEntity(entity)
										planet.Levels[tile.Data[0]].AddEntity(entity)
										pc.X = tile.Data[1]
										pc.Y = tile.Data[2]
									}
								}
							} else {
								playerComponent.AddMessage("Something blocks you in that direction!")
							}
						}
					}
				}
			}
		}
	}

	return planets
}

func checkStairs(x int, y int, tile world.TileSmall) {

}

func getInput() string {
	fmt.Print("Input:")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	return strings.ToUpper(strings.Trim(input, " \n"))
}
