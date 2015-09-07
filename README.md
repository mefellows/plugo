# Plugo

Simple Plugin abstraction with built-in configuration capability. Plugo provides a simple API to register and lookup plugins, 
and then inject configuration into them from a YAML/JSON configuration file.

## Creating a Plugin

First, create your plugin `struct`:

thing.go:

```go
type Thing struct {
	Name string
	Age  int
}

// Register with plugo with name 'fireThing'
func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &Thing{}, nil
	}, "fireThing")
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
```

Then in your main class, load the configuration from file (or in this case, hardcoded) and unmarshall the plugins,
complete with configuration applied:

main.go:

```go
type MyPluginConfig struct {
	Port   int
	Name   string
	Things []PluginConfig
}

func main() {
	var data = []byte(`
port: 8080
name: Foo
things:
  - name: fireThing
    config:
      name: bar
      age: 21
`)

	// Load Configuration
	var confLoader *plugo.ConfigLoader
	c := &MyPluginConfig{}
	confLoader.Load(data, c)

	// Load all plugins
	things := make([]*Thing, len(c.Things))
	plugins := plugo.LoadPluginsWithConfig(confLoader, c.Things)
	for i, p := range plugins {
		things[i] = p.(*Thing)
		things[i].WhoAmiI()
	}
}

```
