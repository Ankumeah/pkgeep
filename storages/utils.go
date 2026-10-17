package storages

import (
	"fmt"
	"os"
	"reflect"
)

func get[T any](field string, config map[string]any) T {
	_field, ok := config[field]
	if !ok {
		fmt.Println(field + " not provided in config")
		os.Exit(1)
	}

	f, ok := _field.(T)
	if !ok {
		fmt.Println(
			field + " in config must be " +
				reflect.TypeFor[T]().String() +
				" got " + reflect.TypeOf(_field).String(),
		)
		os.Exit(1)
	}

	return f
}
