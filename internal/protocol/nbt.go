package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"reflect"
	"unsafe"
)

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

func readByte(r io.Reader, v reflect.Value) (int64, error) {
	var b int8

	if err := binary.Read(r, binary.BigEndian, &b); err != nil {
		return 0, err
	}

	field := v.Elem()
	if field.CanSet() {
		field.SetInt(int64(b))
	}

	return 1, nil
}

func readShort(r io.Reader, v reflect.Value) (int64, error) {
	var s int16

	if err := binary.Read(r, binary.BigEndian, &s); err != nil {
		return 0, err
	}

	field := v.Elem()
	if field.CanSet() {
		field.SetInt(int64(s))
	}

	return 2, nil
}

func readInt(r io.Reader, v reflect.Value) (int64, error) {
	var i int32

	if err := binary.Read(r, binary.BigEndian, &i); err != nil {
		return 0, err
	}

	field := v.Elem()
	if field.CanSet() {
		field.SetInt(int64(i))
	}

	return 4, nil
}

func readLong(r io.Reader, v reflect.Value) (int64, error) {
	var l int64

	if err := binary.Read(r, binary.BigEndian, &l); err != nil {
		return 0, err
	}

	field := v.Elem()
	if field.CanSet() {
		field.SetInt(int64(l))
	}

	return 8, nil
}

func readFloat(r io.Reader, v reflect.Value) (int64, error) {
	var f float32

	if err := binary.Read(r, binary.BigEndian, &f); err != nil {
		return 0, err
	}

	field := v.Elem()
	if field.CanSet() {
		field.SetFloat(float64(f))
	}

	return 4, nil
}

func readDouble(r io.Reader, v reflect.Value) (int64, error) {
	var d float64

	if err := binary.Read(r, binary.BigEndian, &d); err != nil {
		return 0, err
	}

	field := v.Elem()
	if field.CanSet() {
		field.SetFloat(d)
	}

	return 8, nil
}

func readByteArray(r io.Reader, v reflect.Value) (int64, error) {
	var length int32
	size := int64(0)

	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return 0, err
	}
	size += 4

	data := make([]int8, length)
	rawBytes := *(*[]byte)(unsafe.Pointer(&data))

	n, err := io.ReadFull(r, rawBytes[:])
	size += int64(n)
	if err != nil {
		return size, err
	}

	field := v.Elem()
	if field.CanSet() {
		field.Set(reflect.ValueOf(data))
	}

	return size, nil
}

func readIntArray(r io.Reader, v reflect.Value) (int64, error) {
	var length int32
	size := int64(0)

	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return 0, err
	}
	size += 4

	field := v.Elem()

	data := reflect.MakeSlice(field.Type(), int(length), int(length))
	if err := binary.Read(r, binary.BigEndian, data.Interface()); err != nil {
		return size, err
	}

	size += int64(length) * 4

	if field.CanSet() {
		field.Set(data)
	}

	return size, nil
}

func readLongArray(r io.Reader, v reflect.Value) (int64, error) {
	var length int32
	size := int64(0)

	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return 0, err
	}
	size += 4

	field := v.Elem()

	data := reflect.MakeSlice(field.Type(), int(length), int(length))
	if err := binary.Read(r, binary.BigEndian, data.Interface()); err != nil {
		return size, err
	}

	size += int64(length) * 8

	if field.CanSet() {
		field.Set(data)
	}

	return size, nil
}

func readString(r io.Reader, v reflect.Value) (int64, error) {
	var length int16
	size := int64(0)

	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return 0, err
	}
	size += 2

	field := v.Elem()

	data := make([]byte, length)
	n, err := io.ReadFull(r, data[:])
	size += int64(n)
	if err != nil {
		return size, err
	}

	if field.CanSet() {
		field.SetString(string(data)) // TODO MUTF-8
	}

	return size, nil
}

func readList(r io.Reader, v reflect.Value) (int64, error) {

}

func readCompound(r io.Reader, reflectType reflect.Type, vPtr reflect.Value) (int64, error) {
	size := int64(0)

	fieldMap := mapFields(reflectType)
	v := vPtr.Elem()

	for {
		var typeId int8
		if err := binary.Read(r, binary.BigEndian, &typeId); err != nil {
			return size, err
		}
		size++

		if typeId == 0x00 {
			break // Type 0x00 is TAG_End
		}

		var name String
		n, err := name.ReadFrom(r)
		size += n
		if err != nil {
			return size, err
		}

		fieldNr, ok := fieldMap[string(name.Content)]

		if !ok {
			return size, fmt.Errorf(`field "%s" not found`, string(name.Content))
		}

		field := v.Field(fieldNr)
		n, err = read(r, typeId, reflectType, field)
		size += n
		if err != nil {
			return size, err
		}
	}

	return size, nil
}

func read(r io.Reader, typeId int8, reflectType reflect.Type, field reflect.Value) (int64, error) {
	size := int64(0)

	switch typeId {
	case 0x01: // TAG_Byte
		n, err := readByte(r, field.Addr())
		size += n
		if err != nil {
			return size, err
		}
		break
	case 0x02: // TAG_Short
		n, err := readShort(r, field.Addr())
		size += n
		if err != nil {
			return size, err
		}
		break
	case 0x03: // TAG_Int
		n, err := readInt(r, field.Addr())
		size += n
		if err != nil {
			return size, err
		}
	case 0x04: // TAG_Long
		n, err := readLong(r, field.Addr())
		size += n
		if err != nil {
			return size, err
		}
	case 0x05: // TAG_Float
		n, err := readFloat(r, field.Addr())
		size += n
		if err != nil {
			return size, err
		}
		break
	case 0x06: // TAG_Double
		n, err := readDouble(r, field.Addr())
		size += n
		if err != nil {
			return size, err
		}
		break
	case 0x07: // TAG_Byte_Array
		n, err := readByteArray(r, field.Addr())
		size += n
		if err != nil {
			return size, err
		}
		break
	case 0x08: // TAG_String
		n, err := readString(r, field.Addr())
		size += n
		if err != nil {
			return size, err
		}
		break
	case 0x09: // TAG_List
		n, err := readList(r, field.Addr())
		size += n
		if err != nil {
			return size, err
		}
		break
	case 0x0A: // TAG_Compound
		n, err := readCompound(r, field.Type(), field.Addr())
		size += n
		if err != nil {
			return size, err
		}
		break
	default:
		log.Println("Unknown NBT typeId:", typeId)
	}

	return size, nil
}

func ReadNBT[T any](r io.Reader, s *T) {

}

func WriteNBT[T any](w io.Writer, s T) {

}
