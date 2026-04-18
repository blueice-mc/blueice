package nbt

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"unsafe"
)

func WriteNBT(w io.Writer, v any) (int64, error) {
	f := deref(reflect.ValueOf(v))

	total := int64(0)

	n, err := w.Write([]byte{byte(TagCompound)})
	total += int64(n)
	if err != nil {
		return total, err
	}

	m, err := writeCompound(w, f)
	total += m
	return total, err
}

// follows pointers and interfaces until a concrete value is found
func deref(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}

// ensures the value can be addressed. if not, a copy is created and a new pointer will be returned
func addr(v reflect.Value) reflect.Value {
	if v.CanAddr() {
		return v.Addr()
	}

	ptr := reflect.New(v.Type())
	ptr.Elem().Set(v)
	return ptr
}

func tagTypeOf(v reflect.Value) TagType {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return TagEnd
		}
		v = v.Elem()
	}

	switch v.Kind() {
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
		switch v.Type().Elem().Kind() {
		case reflect.Int8:
			return TagByteArray
		case reflect.Int32:
			return TagIntArray
		case reflect.Int64:
			return TagLongArray
		default:
			return TagList
		}
	case reflect.Struct:
		return TagCompound
	case reflect.Map:
		return TagCompound
	default:
		return TagEnd
	}
}

// writeValue calls the right writer for the given reflection value.
// v should not be a pointer
func writeValue(w io.Writer, v reflect.Value) (int64, error) {
	if v.Kind() == reflect.Interface {
		v = deref(v)
	}

	switch tagTypeOf(v) {
	case TagByte:
		return writeByte(w, v)
	case TagShort:
		return writeShort(w, v)
	case TagInt:
		return writeInt(w, v)
	case TagLong:
		return writeLong(w, v)
	case TagFloat:
		return writeFloat(w, v)
	case TagDouble:
		return writeDouble(w, v)
	case TagByteArray:
		return writeByteArray(w, v)
	case TagString:
		return writeString(w, v)
	case TagIntArray:
		return writeIntArray(w, v)
	case TagLongArray:
		return writeLongArray(w, v)
	case TagList:
		return writeList(w, v)
	case TagCompound:
		if v.Kind() == reflect.Map {
			return writeMap(w, v)
		}
		return writeCompound(w, v)
	default:
		return 0, fmt.Errorf("unknown tag type: %v", v.Type())
	}
}

func writeList(w io.Writer, v reflect.Value) (int64, error) {
	v = deref(v)
	if !v.IsValid() || v.Kind() != reflect.Slice {
		return 0, fmt.Errorf("cannot write non-slice type to list")
	}

	total := int64(0)

	var header [5]byte

	tagType := getTagType(v.Type().Elem())
	header[0] = byte(tagType)

	binary.BigEndian.PutUint32(header[1:], uint32(v.Len()))
	n, err := w.Write(header[:])
	total += int64(n)
	if err != nil {
		return total, err
	}

	for i := 0; i < v.Len(); i++ {
		n, err := writeValue(w, v.Index(i))
		total += n
		if err != nil {
			return total, err
		}
	}

	return total, nil
}

func writeMap(w io.Writer, v reflect.Value) (int64, error) {
	total := int64(0)
	for _, key := range v.MapKeys() {
		val := v.MapIndex(key)
		n, err := writeNamedTagHeader(w, tagTypeOf(val), key.String())
		total += n
		if err != nil {
			return total, err
		}
		n, err = writeValue(w, val)
		total += n
		if err != nil {
			return total, err
		}
	}
	n, err := w.Write([]byte{byte(TagEnd)})
	total += int64(n)
	return total, err
}

func writeCompound(w io.Writer, v reflect.Value) (int64, error) {
	v = deref(v)
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return 0, fmt.Errorf("expected struct, got %v", v.Kind())
	}

	total := int64(0)
	t := v.Type()

	// iterate through all fields in the struct
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				continue
			}
			field = deref(field)
		}

		fieldType := t.Field(i)

		tagName, omitempty := parseTag(fieldType.Tag.Get("nbt"))

		// Skip if nbt tag name is empty
		if tagName == "" || tagName == "-" {
			continue
		}

		if omitempty && field.IsZero() {
			continue
		}

		n, err := writeNamedTagHeader(w, tagTypeOf(field), tagName)
		total += n
		if err != nil {
			return total, err
		}

		n, err = writeValue(w, field)
		total += n
		if err != nil {
			return total, err
		}
	}

	n, err := w.Write([]byte{byte(TagEnd)})
	total += int64(n)
	return total, err
}

func writeNamedTagHeader(w io.Writer, t TagType, name string) (int64, error) {
	var buf [3]byte
	buf[0] = byte(t)
	binary.BigEndian.PutUint16(buf[1:], uint16(len(name)))
	total, err := w.Write(buf[:])

	if err != nil {
		return int64(total), err
	}

	n, err := w.Write([]byte(name))
	total += n
	if err != nil {
		return int64(total), err
	}

	return int64(total), err
}

func parseTag(tag string) (string, bool) {
	name, _, _ := strings.Cut(tag, ",")
	return name, strings.Contains(tag, ",omitempty")
}

// Primitive writers

func writeByte(w io.Writer, v reflect.Value) (int64, error) {
	n, err := w.Write([]byte{byte(v.Int())})
	return int64(n), err
}

func writeShort(w io.Writer, v reflect.Value) (int64, error) {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], uint16(v.Int()))
	n, err := w.Write(buf[:])
	return int64(n), err
}

func writeInt(w io.Writer, v reflect.Value) (int64, error) {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(v.Int()))
	n, err := w.Write(buf[:])
	return int64(n), err
}

func writeLong(w io.Writer, v reflect.Value) (int64, error) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(v.Int()))
	n, err := w.Write(buf[:])
	return int64(n), err
}

func writeFloat(w io.Writer, v reflect.Value) (int64, error) {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], math.Float32bits(float32(v.Float())))
	n, err := w.Write(buf[:])
	return int64(n), err
}

func writeDouble(w io.Writer, v reflect.Value) (int64, error) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(v.Float()))
	n, err := w.Write(buf[:])
	return int64(n), err
}

func writeByteArray(w io.Writer, v reflect.Value) (int64, error) {
	var header [4]byte
	length := v.Len()
	binary.BigEndian.PutUint32(header[:], uint32(length))
	total, err := w.Write(header[:])
	if err != nil {
		return int64(total), err
	}

	if length > 0 {
		rawBytes := unsafe.Slice((*byte)(unsafe.Pointer(v.Pointer())), length)
		n, err := w.Write(rawBytes)
		total += n
		if err != nil {
			return int64(total), err
		}
	}

	return int64(total), nil
}

func writeString(w io.Writer, v reflect.Value) (int64, error) {
	var header [2]byte
	str := v.String()
	binary.BigEndian.PutUint16(header[:], uint16(len(str)))
	total, err := w.Write(header[:])
	if err != nil {
		return int64(total), err
	}

	if len(str) > 0 {
		n, err := io.WriteString(w, str)
		total += n
		if err != nil {
			return int64(total), err
		}
	}

	return int64(total), nil
}

func writeIntArray(w io.Writer, v reflect.Value) (int64, error) {
	var header [4]byte
	length := v.Len()
	binary.BigEndian.PutUint32(header[:], uint32(length))
	total, err := w.Write(header[:])
	if err != nil {
		return int64(total), err
	}

	slice := v.Interface().([]int32)

	var buffer [256]byte
	for i := range length {
		pos := i % 64
		binary.BigEndian.PutUint32(buffer[4*pos:4*(pos+1)], uint32(slice[i]))

		// Flush data if the buffer is full
		if pos == 63 || i+1 == length {
			n, err := w.Write(buffer[:4*(pos+1)])
			total += n
			if err != nil {
				return int64(total), err
			}
		}
	}

	return int64(total), nil
}

func writeLongArray(w io.Writer, v reflect.Value) (int64, error) {
	var header [4]byte
	length := v.Len()
	binary.BigEndian.PutUint32(header[:], uint32(length))
	total, err := w.Write(header[:])
	if err != nil {
		return int64(total), err
	}

	slice := v.Interface().([]int64)

	var buffer [512]byte
	for i := range length {
		pos := i % 64
		binary.BigEndian.PutUint64(buffer[8*pos:8*(pos+1)], uint64(slice[i]))

		// Flush data if the buffer is full
		if pos == 63 || i+1 == length {
			n, err := w.Write(buffer[:8*(pos+1)])
			total += n
			if err != nil {
				return int64(total), err
			}
		}
	}

	return int64(total), nil
}
