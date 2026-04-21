package protocol

import (
	"encoding/binary"
	"io"
	"math"
)

func WriteBool(w io.Writer, v bool) (int64, error) {
	if v {
		return WriteInt8(w, 1)
	}
	return WriteInt8(w, 0)
}

// WriteUint8 writes a single unsigned byte to the given writer
func WriteUint8(w io.Writer, v uint8) (int64, error) {
	total, err := w.Write([]byte{v})
	return int64(total), err
}

// WriteInt8 writes a single signed byte to the given writer
func WriteInt8(w io.Writer, v int8) (int64, error) {
	total, err := w.Write([]byte{byte(v)})
	return int64(total), err
}

// WriteInt16 writes a signed 16-bit integer to the given writer
func WriteInt16(w io.Writer, v int16) (int64, error) {
	var buffer [2]byte
	binary.BigEndian.PutUint16(buffer[:], uint16(v))
	total, err := w.Write(buffer[:])
	return int64(total), err
}

// WriteInt32 writes a signed 32-bit integer to the given writer
func WriteInt32(w io.Writer, v int32) (int64, error) {
	var buffer [4]byte
	binary.BigEndian.PutUint32(buffer[:], uint32(v))
	total, err := w.Write(buffer[:])
	return int64(total), err
}

// WriteInt64 writes a signed 64-bit integer to the given writer
func WriteInt64(w io.Writer, v int64) (int64, error) {
	var buffer [8]byte
	binary.BigEndian.PutUint64(buffer[:], uint64(v))
	total, err := w.Write(buffer[:])
	return int64(total), err
}

// WriteFloat32 writes a 32-bit floating point number to the given writer
func WriteFloat32(w io.Writer, v float32) (int64, error) {
	var buffer [4]byte
	binary.BigEndian.PutUint32(buffer[:], math.Float32bits(v))
	total, err := w.Write(buffer[:])
	return int64(total), err
}

// WriteFloat64 writes a 64-bit floating point number to the given writer
func WriteFloat64(w io.Writer, v float64) (int64, error) {
	var buffer [8]byte
	binary.BigEndian.PutUint64(buffer[:], math.Float64bits(v))
	total, err := w.Write(buffer[:])
	return int64(total), err
}

func WriteString(w io.Writer, v string) (int64, error) {
	total := int64(0)
	length := VarInt(len(v))
	n, err := length.WriteTo(w)
	total += n
	if err != nil {
		return total, err
	}
	m, err := io.WriteString(w, v)
	total += int64(m)
	return total, nil
}

func ReadString(r io.Reader, v *string) (int64, error) {
	total := int64(0)
	var length VarInt
	n, err := length.ReadFrom(r)
	total += n
	if err != nil {
		return total, err
	}

	buffer := make([]byte, length)
	m, err := io.ReadFull(r, buffer)
	total += int64(m)
	if err != nil {
		return total, err
	}
	*v = string(buffer)
	return total, nil
}
