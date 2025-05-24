package main

import (
	"log"

	handles "github.com/royalmcpe/handle-map"
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
	handle, ok := m.Add(Entity{Name: "Test"})
	if !ok {
		panic("no entity")
	}

	e := m.Get(handle)
	log.Println(e.Name)
}
