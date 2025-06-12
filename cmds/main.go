package main

import (
	"log"

	handles "github.com/royalmcpe/handle-map"
)

type Entity struct {
	Name string
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
