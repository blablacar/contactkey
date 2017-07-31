package deployers

import (
	"fmt"
	"reflect"
)

var registry = make(map[string]reflect.Type)

type Deployer interface {
	listUnits(env string) ([]string, error)
}

func register(deployer interface{}) {
	t := reflect.TypeOf(deployer).Elem()
	registry[fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())] = t
}

func makeInstance(name string) interface{} {
	return reflect.New(registry[name]).Elem().Interface()
}
