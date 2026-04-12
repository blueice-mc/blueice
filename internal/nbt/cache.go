package nbt

import "reflect"

var structCache = make(map[reflect.Type]map[string]int)

func mapFields(t reflect.Type) map[string]int {
	m, ok := structCache[t]

	if ok {
		return m
	}

	m = make(map[string]int)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("nbt")

		if len(tag) == 0 {
			tag = field.Name
		}

		m[tag] = i
	}

	structCache[t] = m
	return m
}
