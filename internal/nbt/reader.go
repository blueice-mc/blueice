package nbt

import (
	"encoding/binary"
	"errors"
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
	var length uint16
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
	var typeId int8
	if err := binary.Read(r, binary.BigEndian, &typeId); err != nil {
		return 0, err
	}
	size := int64(1)

	var length int32
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return 0, err
	}
	size += 4

	field := v.Elem()
	data := reflect.MakeSlice(field.Type(), int(length), int(length))

	for i := 0; i < int(length); i++ {
		n, err := read(r, typeId, data.Index(i).Addr())
		size += n
		if err != nil {
			return size, err
		}
	}

	return size, nil
}

func readCompound(r io.Reader, vPtr reflect.Value) (int64, error) {
	size := int64(0)

	v := vPtr.Elem()
	fieldMap := mapFields(v.Type())

	for {
		var typeId int8
		if err := binary.Read(r, binary.BigEndian, &typeId); err != nil {
			return size, err
		}
		size++

		if typeId == 0x00 {
			break // Type 0x00 is TAG_End
		}

		var nameLength uint16
		if err := binary.Read(r, binary.BigEndian, &nameLength); err != nil {
			return size, err
		}
		size += 2

		name := make([]byte, nameLength)
		if err := binary.Read(r, binary.BigEndian, name); err != nil {
			return size, err
		}

		fieldNr, ok := fieldMap[string(name)]

		if !ok {
			n, err := skip(r, typeId)
			size += n
			if err != nil {
				return size, err
			}
			continue
		}

		field := v.Field(fieldNr)
		n, err := read(r, typeId, field)
		size += n
		if err != nil {
			return size, err
		}
	}

	return size, nil
}

func skip(r io.Reader, typeId int8) (int64, error) {
	switch typeId {
	case 0x01: // TAG_Byte
		n, err := io.CopyN(io.Discard, r, 1)
		return n, err
	case 0x02: // TAG_Short
		n, err := io.CopyN(io.Discard, r, 2)
		return n, err
	case 0x03: // TAG_Int
		n, err := io.CopyN(io.Discard, r, 4)
		return n, err
	case 0x04: // TAG_Long
		n, err := io.CopyN(io.Discard, r, 8)
		return n, err
	case 0x05: // TAG_Float
		n, err := io.CopyN(io.Discard, r, 4)
		return n, err
	case 0x06: // TAG_Double
		n, err := io.CopyN(io.Discard, r, 8)
		return n, err
	case 0x07: // TAG_Byte_Array
		var length int32
		if err := binary.Read(r, binary.BigEndian, &length); err != nil {
			return 0, err
		}
		n, err := io.CopyN(io.Discard, r, int64(length))
		return 4 + int64(n), err
	case 0x08: // TAG_String
		var length uint16
		if err := binary.Read(r, binary.BigEndian, &length); err != nil {
			return 0, err
		}
		n, err := io.CopyN(io.Discard, r, int64(length))
		return 2 + int64(n), err
	case 0x09: // TAG_List
		size := int64(0)

		var listTypeId int8
		if err := binary.Read(r, binary.BigEndian, &listTypeId); err != nil {
			return 0, err
		}
		size += 1

		var length int32
		if err := binary.Read(r, binary.BigEndian, &length); err != nil {
			return 1, err
		}
		size += 4

		for i := 0; i < int(length); i++ {
			n, err := skip(r, listTypeId)
			size += n
			if err != nil {
				return size, err
			}
		}

		return size, nil
	case 0x0A: // TAG_Compound
		size := int64(0)

		for {
			var fieldTypeId int8
			if err := binary.Read(r, binary.BigEndian, &fieldTypeId); err != nil {
				return 0, err
			}
			size++

			if fieldTypeId == 0x00 { // TAG_End
				return size, nil
			}

			var nameLength uint16
			if err := binary.Read(r, binary.BigEndian, &nameLength); err != nil {
				return size, err
			}
			size += 2

			n, err := io.CopyN(io.Discard, r, int64(nameLength))
			size += n

			if err != nil {
				return size, err
			}

			n, err = skip(r, fieldTypeId)
			size += n
			if err != nil {
				return size, err
			}
		}
	}

	return 0, fmt.Errorf("unsupported type: %v", typeId)
}

func read(r io.Reader, typeId int8, field reflect.Value) (int64, error) {
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
		n, err := readCompound(r, field.Addr())
		size += n
		if err != nil {
			return size, err
		}
		break
	case 0x0B: // TAG_Int_Array
		n, err := readIntArray(r, field.Addr())
		size += n
		if err != nil {
			return size, err
		}
		break
	case 0x0C: // TAG_Long_Array
		n, err := readLongArray(r, field.Addr())
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

func ReadNBT[T any](r io.Reader, s *T) (int64, error) {
	v := reflect.ValueOf(s)

	var size int64

	var typeId int8
	if err := binary.Read(r, binary.BigEndian, &typeId); err != nil {
		return size, err
	}
	size += 1

	if typeId != int8(0x0A) {
		return size, errors.New("nbt root must be compound")
	}

	n, err := readCompound(r, v)
	size += n
	if err != nil {
		return size, err
	}

	return size, nil
}

func WriteNBT[T any](w io.Writer, s T) {

}
