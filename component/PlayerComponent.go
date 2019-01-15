package component

// PlayerComponent
type PlayerComponent struct {
	MessageLog      []string
	InteractingWith int
}

//GetType get the type
func (PlayerComponent) GetType() string {
	return "PlayerComponent"
}

func (pc *PlayerComponent) AddMessage(x string) {
	pc.MessageLog = append(pc.MessageLog, x)
}
