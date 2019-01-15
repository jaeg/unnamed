package main

import (
	"math/rand"
	"os/exec"
	"time"

	"unnamed/component"
	"unnamed/entity"
	"unnamed/system"
	"unnamed/world"

	"github.com/nsf/termbox-go"
)

var planets map[string]*world.Planet

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	//start := time.Now()
	planets = make(map[string]*world.Planet)
	planets["hub"] = world.NewPlanet()
	//elapsed := time.Since(start)

	//Player
	newPlayerEntity := entity.Entity{}
	playerComponent := &component.PlayerComponent{}
	newPlayerEntity.AddComponent(playerComponent)
	initiativeComponent := &component.InitiativeComponent{DefaultValue: 10, Ticks: 1}
	newPlayerEntity.AddComponent(initiativeComponent)
	positionComponent := &component.PositionComponent{X: 0, Y: 0, Level: 0}
	newPlayerEntity.AddComponent(positionComponent)
	newPlayerEntity.AddComponent(&component.AppearanceComponent{SpriteIndex: 0, Resource: "npc"})
	newPlayerEntity.AddComponent(&component.DirectionComponent{Direction: 0})
	newPlayerEntity.AddComponent(&component.SolidComponent{})
	newPlayerEntity.AddComponent(&component.InventoryComponent{})
	//entities = append(entities, &newPlayerEntity)
	planets["hub"].Levels[0].AddEntity(&newPlayerEntity)

	for i := 0; i < 10; i++ {
		entity := entity.Entity{}
		x := 1
		y := 1
		if i != 0 {
			x = rand.Intn(30)
			y = rand.Intn(30)
			message := []string{"Hello there!", "Like my hat?", "It's dangerous out here at night."}
			entity.AddComponent(&component.InteractComponent{Message: message})
			entity.AddComponent(&component.AppearanceComponent{Char: "P"})
		} else {
			entity.AddComponent(&component.ShopComponent{ItemsForSale: []string{"Sword", "Bow", "Shield", "Meat"}})
			entity.AddComponent(&component.AppearanceComponent{Char: "S"})
		}

		//entity.AddComponent(&component.WanderAIComponent{})
		entity.AddComponent(&component.InitiativeComponent{DefaultValue: 4, Ticks: 1})
		entity.AddComponent(&component.PositionComponent{X: x, Y: y, Level: 0})
		entity.AddComponent(&component.DirectionComponent{Direction: 0})
		entity.AddComponent(&component.SolidComponent{})

		//entities = append(entities, &newPlayerEntity)
		planets["hub"].Levels[0].AddEntity(&entity)

	}

	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// restore the echoing state when exiting
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	ticker := time.NewTicker(time.Second / 16)

	running := true
	system.RenderSystem(planets)
	for _ = range ticker.C {
		if !running {
			break
		}

		//start := time.Now()
		system.InitiativeSystem(planets)
		planets = system.PlayerSystem(planets)
		system.AISystem(planets)
		system.RenderSystem(planets)
		system.StatusConditionSystem(planets)
		planets = system.CleanUpSystem(planets)
		//elapsed := time.Since(start)
		//log.Printf("Game loop took %s", elapsed)
	}

}
