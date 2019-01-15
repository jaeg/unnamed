package world

type Planet struct {
	Levels []*Level
}

func NewPlanet() (planet *Planet) {
	planet = &Planet{Levels: []*Level{}}
	planet.Levels = append(planet.Levels, NewOverworldSection(100, 100))
	return
}
