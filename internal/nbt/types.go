package nbt

import "reflect"

type TagType byte

const (
	TagEnd       TagType = 0
	TagByte      TagType = 1
	TagShort     TagType = 2
	TagInt       TagType = 3
	TagLong      TagType = 4
	TagFloat     TagType = 5
	TagDouble    TagType = 6
	TagByteArray TagType = 7
	TagString    TagType = 8
	TagList      TagType = 9
	TagCompound  TagType = 10
	TagIntArray  TagType = 11
	TagLongArray TagType = 12
)

func getTagType(t reflect.Type) TagType {
	switch t.Kind() {
	case reflect.Bool:
		return TagByte
	case reflect.Int8:
		return TagByte
	case reflect.Int16:
		return TagShort
	case reflect.Int32:
		return TagInt
	case reflect.Int64:
		return TagLong
	case reflect.Float32:
		return TagFloat
	case reflect.Float64:
		return TagDouble
	case reflect.String:
		return TagString
	case reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.Int8, reflect.Uint8:
			return TagByteArray
		case reflect.Int, reflect.Int32:
			return TagIntArray
		case reflect.Int64:
			return TagLongArray
		default:
			return TagList
		}
	case reflect.Struct:
		return TagCompound
	default:
		return TagEnd
	}
}
