package component

// MyTurnComponent .
type AppearanceComponent struct {
	SpriteIndex int
	Resource    string
	Char        string
}

func (pc AppearanceComponent) GetType() string {
	return "AppearanceComponent"
}
