package system

import (
	"unnamed/component"
	"unnamed/world"
)

var statusConditions = []string{"Poisoned"}

// StatusConditionSystem .
func StatusConditionSystem(planets map[string]*world.Planet) {
	for _, planet := range planets {
		for _, level := range planet.Levels {
			for _, entity := range level.Entities {
				for _, statusCondition := range statusConditions {
					if entity.HasComponent(statusCondition + "Component") {
						pc := entity.GetComponent(statusCondition + "Component").(component.DecayingComponent)

						if pc.Decay() {
							entity.RemoveComponent(statusCondition + "Component")
							if entity.HasComponent("PlayerComponent") {
								playerComponent := entity.GetComponent("PlayerComponent").(*component.PlayerComponent)
								playerComponent.AddMessage(statusCondition + " has cleared!")
							}
						}
					}
				}
			}
		}
	}
}
