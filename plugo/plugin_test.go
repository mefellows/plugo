package plugo

import (
	"log"
	"testing"
)

var MiddlewareMockFunc = func() (interface{}, error) {
	return MockPluginThing{}, nil
}

func TestLookupFactory(t *testing.T) {
	PluginFactories.Register(MiddlewareMockFunc, "screwything")

	f, ok := PluginFactories.Lookup("screwything")

	if !ok {
		t.Fatalf("Expected lookup to be OK")
	}

	sym, err := f()

	if err != nil {
		t.Fatalf("Did not expect err: %v", err)
	}

	if sym == nil {
		t.Fatalf("Expected symptom not to be nil")
	}

	// Need to cast to specific type
	s := sym.(PluginThing)

	// Should not panic!
	s.DoSomething()
}

func TestLookupPluginAndConfigure(t *testing.T) {
	var data = []byte(`
port: 8080
name: Foo
things:
  - name: fireThing
    config:
      name: bar
      age: 21
`)

	PluginFactories.Register(func() (interface{}, error) {
		return &Thing{}, nil
	}, "fireThing")

	// Load Configuration
	var confLoader *ConfigLoader
	c := &MyPluginConfig{}
	confLoader.Load(data, c)

	// Load all plugins
	things := make([]*Thing, len(c.Things))
	plugins := LoadPluginsWithConfig(confLoader, c.Things)
	for i, p := range plugins {
		things[i] = p.(*Thing)
		things[i].WhoAmiI()
	}
}

type MyPluginConfig struct {
	Port   int
	Name   string
	Things []PluginConfig
}
type Thing struct {
	Name string
	Age  int
}

func (t *Thing) WhoAmiI() {
	log.Printf("My name is %s and I'm %d", t.Name, t.Age)
}

type PluginThing interface {
	DoSomething()
}
type MockPluginThing struct {
}

func (m MockPluginThing) DoSomething() {
	log.Println("Doing something!")
}
