package nbt

import (
	"encoding/binary"
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
	l := int64(field.Int())

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

func writeCompound(w io.Writer, v reflect.Value) (int64, error) {
	field := v.Elem()
	
}

func write(w io.Writer, v reflect.Value) (int64, error) {
	t := v.Elem().Type()
	tagType := getTagType(t)
	size := int64(0)

	n, err := w.Write([]byte{byte(tagType)})
	size += int64(n)
	if err != nil {
		return size, err
	}

	switch tagType {
	case TagEnd:
		return size, nil

	}
}

func WriteNBT[T any](w io.Writer, s T) {

}
