package main

import (
	"fmt"

	handles "github.com/royalmcpe/golang-handles-map"
)

type EntityHandle handles.Handle

type Entity struct {
	handle *EntityHandle

	Name string
}

func (e *Entity) SetHandle(handle *handles.Handle) {
	e.handle = (*EntityHandle)(handle)
}

func (e Entity) Handle() *handles.Handle {
	return (*handles.Handle)(e.handle)
}

func main() {
	m := handles.NewMap[Entity](1024)
	handle1, ok1 := m.Add(Entity{Name: "Test"})
	handle2, ok2 := m.Add(Entity{Name: "Test2"})
	if !ok1 || !ok2 {
		panic(fmt.Sprintf("handle1: %v, handle2: %v", handle1, handle2))
	}

	if h2e := m.Get(handle2); h2e != nil {
		h2e.Name = "Test3"
	}

	m.Remove(handle1)

	_, _ = m.Add(Entity{Name: "Test4"})
}
