package system

import (
	"unnamed/component"
	"unnamed/world"
)

// InitiativeSystem .
func InitiativeSystem(planets map[string]*world.Planet) {
	for _, planet := range planets {
		for _, level := range planet.Levels {
			for _, entity := range level.Entities {
				if entity.HasComponent("InitiativeComponent") {
					ic := entity.GetComponent("InitiativeComponent").(*component.InitiativeComponent)
					ic.Ticks--

					if ic.Ticks <= 0 {
						ic.Ticks = ic.DefaultValue
						if ic.OverrideValue > 0 {
							ic.Ticks = ic.OverrideValue
						}

						if entity.HasComponent("MyTurnComponent") == false {
							mTC := &component.MyTurnComponent{}
							entity.AddComponent(mTC)
						}
					}
				}
			}
		}
	}
}
