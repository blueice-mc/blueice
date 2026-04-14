package nbt

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"unsafe"
)

func writeByte(w io.Writer, v reflect.Value) (int64, error) {
	field := v.Elem()
	b := int8(field.Int())

	n, err := w.Write([]byte{byte(b)})
	return int64(n), err
}

func writeShort(w io.Writer, v reflect.Value) (int64, error) {
	field := v.Elem()
	s := int16(field.Int())

	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], uint16(s))
	n, err := w.Write(buf[:])
	return int64(n), err
}

func writeInt(w io.Writer, v reflect.Value) (int64, error) {
	field := v.Elem()
	i := int32(field.Int())

	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(i))
	n, err := w.Write(buf[:])
	return int64(n), err
}

func writeLong(w io.Writer, v reflect.Value) (int64, error) {
	field := v.Elem()
	l := field.Int()

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(l))
	n, err := w.Write(buf[:])
	return int64(n), err
}

func writeFloat(w io.Writer, v reflect.Value) (int64, error) {
	field := v.Elem()
	f := float32(field.Float())

	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], math.Float32bits(f))
	n, err := w.Write(buf[:])
	return int64(n), err
}

func writeDouble(w io.Writer, v reflect.Value) (int64, error) {
	field := v.Elem()
	d := field.Float()

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(d))
	n, err := w.Write(buf[:])
	return int64(n), err
}

func writeByteArray(w io.Writer, v reflect.Value) (int64, error) {
	field := v.Elem()
	length := int32(field.Len())
	size := int64(0)

	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(length))
	n, err := w.Write(buf[:])
	size += int64(n)
	if err != nil {
		return size, err
	}

	if length > 0 {
		rawBytes := unsafe.Slice((*byte)(unsafe.Pointer(field.Pointer())), length)
		n, err = w.Write(rawBytes[:])
		size += int64(n)
	}

	return size, err
}

func writeIntArray(w io.Writer, v reflect.Value) (int64, error) {
	field := v.Elem()
	size := int64(0)

	slice := field.Interface().([]int32)

	var header [4]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(slice)))
	n, err := w.Write(header[:])
	size += int64(n)
	if err != nil {
		return size, err
	}

	var buffer [256]byte
	for i := range slice {
		pos := i % 64
		binary.BigEndian.PutUint32(buffer[4*pos:4*(pos+1)], uint32(slice[i]))

		if pos == 63 || i == len(slice)-1 {
			n, err := w.Write(buffer[:4*(pos+1)])
			size += int64(n)
			if err != nil {
				return size, err
			}
		}
	}

	return size, nil
}

func writeLongArray(w io.Writer, v reflect.Value) (int64, error) {
	field := v.Elem()
	size := int64(0)

	slice := field.Interface().([]int64)

	var header [4]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(slice)))
	n, err := w.Write(header[:])
	size += int64(n)
	if err != nil {
		return size, err
	}

	var buffer [512]byte
	for i := range slice {
		pos := i % 64
		binary.BigEndian.PutUint64(buffer[8*pos:8*(pos+1)], uint64(slice[i]))

		if pos == 63 || i == len(slice)-1 {
			n, err := w.Write(buffer[:8*(pos+1)])
			size += int64(n)
			if err != nil {
				return size, err
			}
		}
	}

	return size, nil
}

func writeString(w io.Writer, v reflect.Value) (int64, error) {
	field := v.Elem()
	str := field.String() // TODO MUTF-8
	size := int64(0)

	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], uint16(len(str)))
	n, err := w.Write(buf[:])
	size += int64(n)
	if err != nil {
		return size, err
	}

	n, err = io.WriteString(w, str)
	size += int64(n)
	if err != nil {
		return size, err
	}

	return size, nil
}

func writeList(w io.Writer, v reflect.Value) (int64, error) {
	size := int64(0)
	field := v.Elem()
	elementType := field.Type().Elem()
	tagType := getTagType(elementType)

	// Writing the type ID
	n, err := w.Write([]byte{byte(tagType)})
	size += int64(n)
	if err != nil {
		return size, err
	}

	// Writing the length
	length := int32(field.Len())
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(length))
	n, err = w.Write(buf[:])
	size += int64(n)
	if err != nil {
		return size, err
	}

	// Writing the content
	for i := 0; i < int(length); i++ {
		element := field.Index(i)
		n, err := write(w, element.Addr())
		size += int64(n)
		if err != nil {
			return size, err
		}
	}

	return size, nil
}

func writeCompound(w io.Writer, v reflect.Value) (int64, error) {
	size := int64(0)

	field := v

	for field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface {
		if field.IsNil() {
			return 0, nil
		}
		field = field.Elem()
	}

	for i := 0; i < field.NumField(); i++ {
		f := field.Field(i)
		fieldPtr := f.Addr()
		actualType := f.Type()

		if actualType.Kind() == reflect.Ptr {
			actualType = actualType.Elem()
		}

		tagType := getTagType(actualType)

		fieldName := field.Type().Field(i).Tag.Get("nbt")
		if fieldName == "" || fieldName == "-" {
			continue
		}

		if f.Kind() == reflect.Ptr && f.IsNil() {
			continue
		}

		// Writing the type

		n, err := w.Write([]byte{byte(tagType)})
		size += int64(n)
		if err != nil {
			return size, err
		}

		// Writing the name

		// Length of the name
		var buf [2]byte
		binary.BigEndian.PutUint16(buf[:], uint16(len(fieldName)))
		n, err = w.Write(buf[:])
		size += int64(n)
		if err != nil {
			return size, err
		}

		// Name as String
		n, err = io.WriteString(w, fieldName)
		size += int64(n)
		if err != nil {
			return size, err
		}

		// Writing the field content
		m, err := write(w, fieldPtr)
		size += m
		if err != nil {
			return size, err
		}
	}

	// Write TAG_End
	n, err := w.Write([]byte{byte(TagEnd)})
	size += int64(n)
	if err != nil {
		return size, err
	}
	return size, nil
}

func write(w io.Writer, v reflect.Value) (int64, error) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return 0, nil
		}
	}

	tagType := getTagType(v.Elem().Type())

	switch tagType {
	case TagEnd:
		return 0, nil
	case TagByte:
		n, err := writeByte(w, v)
		return n, err
	case TagShort:
		n, err := writeShort(w, v)
		return n, err
	case TagInt:
		n, err := writeInt(w, v)
		return n, err
	case TagLong:
		n, err := writeLong(w, v)
		return n, err
	case TagFloat:
		n, err := writeFloat(w, v)
		return n, err
	case TagDouble:
		n, err := writeDouble(w, v)
		return n, err
	case TagByteArray:
		n, err := writeByteArray(w, v)
		return n, err
	case TagString:
		n, err := writeString(w, v)
		return n, err
	case TagList:
		n, err := writeList(w, v)
		return n, err
	case TagCompound:
		n, err := writeCompound(w, v)
		return n, err
	case TagIntArray:
		n, err := writeIntArray(w, v)
		return n, err
	case TagLongArray:
		n, err := writeLongArray(w, v)
		return n, err
	default:
		return 0, fmt.Errorf("unknown tag type: %v", tagType)
	}
}

func WriteNBT(w io.Writer, s any) (int64, error) {
	ptr := reflect.ValueOf(s)

	if ptr.Kind() == reflect.Interface {
		ptr = ptr.Elem()
	}

	// Writing type first
	n, err := w.Write([]byte{byte(TagCompound)})
	if err != nil {
		return int64(n), err
	}

	m, err := writeCompound(w, ptr)
	return int64(n) + m, err
}
