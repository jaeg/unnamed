package system

import (
	"unnamed/component"
	"fmt"
)

// HelloWorldSystem .
type HelloWorldSystem struct {
}

// Update .
func (HelloWorldSystem) Update(a *component.HelloWorldComponent) {
	fmt.Println("hello world")
}
