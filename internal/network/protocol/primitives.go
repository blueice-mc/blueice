package protocol

import (
	"encoding/binary"
	"io"
	"math"
)

// ReadBool reads a single byte from the given reader and interprets it as a boolean
func ReadBool(r io.Reader, v *bool) (int64, error) {
	var read int8
	total, err := ReadInt8(r, &read)
	*v = read == 1
	return total, err
}

// WriteBool writes a single byte to the given writer, representing the boolean value
func WriteBool(w io.Writer, v bool) (int64, error) {
	if v {
		return WriteInt8(w, 1)
	}
	return WriteInt8(w, 0)
}

// ReadUint8 reads a single unsigned byte from the given reader
func ReadUint8(r io.Reader, v *uint8) (int64, error) {
	var buffer [1]byte
	total, err := r.Read(buffer[:])
	*v = buffer[0]
	return int64(total), err
}

// WriteUint8 writes a single unsigned byte to the given writer
func WriteUint8(w io.Writer, v uint8) (int64, error) {
	total, err := w.Write([]byte{v})
	return int64(total), err
}

// ReadInt8 reads a single signed byte from the given reader
func ReadInt8(r io.Reader, v *int8) (int64, error) {
	var buffer [1]byte
	total, err := r.Read(buffer[:])
	*v = int8(buffer[0])
	return int64(total), err
}

// WriteInt8 writes a single signed byte to the given writer
func WriteInt8(w io.Writer, v int8) (int64, error) {
	total, err := w.Write([]byte{byte(v)})
	return int64(total), err
}

// ReadInt16 reads a signed 16-bit integer from the given reader
func ReadInt16(r io.Reader, v *int16) (int64, error) {
	var buffer [2]byte
	_, err := r.Read(buffer[:])
	*v = int16(binary.BigEndian.Uint16(buffer[:]))
	return 2, err
}

// WriteInt16 writes a signed 16-bit integer to the given writer
func WriteInt16(w io.Writer, v int16) (int64, error) {
	var buffer [2]byte
	binary.BigEndian.PutUint16(buffer[:], uint16(v))
	total, err := w.Write(buffer[:])
	return int64(total), err
}

// ReadInt32 reads a signed 32-bit integer from the given reader
func ReadInt32(r io.Reader, v *int32) (int64, error) {
	var buffer [4]byte
	_, err := r.Read(buffer[:])
	*v = int32(binary.BigEndian.Uint32(buffer[:]))
	return 4, err
}

// WriteInt32 writes a signed 32-bit integer to the given writer
func WriteInt32(w io.Writer, v int32) (int64, error) {
	var buffer [4]byte
	binary.BigEndian.PutUint32(buffer[:], uint32(v))
	total, err := w.Write(buffer[:])
	return int64(total), err
}

// ReadInt64 reads a signed 64-bit integer from the given reader
func ReadInt64(r io.Reader, v *int64) (int64, error) {
	var buffer [8]byte
	_, err := r.Read(buffer[:])
	*v = int64(binary.BigEndian.Uint64(buffer[:]))
	return 8, err
}

// WriteInt64 writes a signed 64-bit integer to the given writer
func WriteInt64(w io.Writer, v int64) (int64, error) {
	var buffer [8]byte
	binary.BigEndian.PutUint64(buffer[:], uint64(v))
	total, err := w.Write(buffer[:])
	return int64(total), err
}

// ReadFloat32 reads a 32-bit floating point number from the given reader
func ReadFloat32(r io.Reader, v *float32) (int64, error) {
	var buffer [4]byte
	_, err := r.Read(buffer[:])
	*v = math.Float32frombits(binary.BigEndian.Uint32(buffer[:]))
	return 4, err
}

// WriteFloat32 writes a 32-bit floating point number to the given writer
func WriteFloat32(w io.Writer, v float32) (int64, error) {
	var buffer [4]byte
	binary.BigEndian.PutUint32(buffer[:], math.Float32bits(v))
	total, err := w.Write(buffer[:])
	return int64(total), err
}

// ReadFloat64 reads a 64-bit floating point number from the given reader
func ReadFloat64(r io.Reader, v *float64) (int64, error) {
	var buffer [8]byte
	_, err := r.Read(buffer[:])
	*v = math.Float64frombits(binary.BigEndian.Uint64(buffer[:]))
	return 8, err
}

// WriteFloat64 writes a 64-bit floating point number to the given writer
func WriteFloat64(w io.Writer, v float64) (int64, error) {
	var buffer [8]byte
	binary.BigEndian.PutUint64(buffer[:], math.Float64bits(v))
	total, err := w.Write(buffer[:])
	return int64(total), err
}

// ReadString reads a string from the given reader as a prefixed byte array
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

// WriteString writes a string to the given writer as a prefixed byte array
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
