package protocol

import (
	"errors"
	"io"
)

// Definition of the VarInt type and read/write functions
type VarInt int32

func (v VarInt) WriteTo(w io.Writer) (int64, error) {
	value := uint32(v)
	var size int64
	var buf [1]byte

	for {
		currentByte := value & 0x7f
		value >>= 7

		if value != 0 {
			currentByte |= 0x80
		}

		buf[0] = byte(currentByte)
		_, err := w.Write(buf[:])

		if err != nil {
			return 0, err
		}

		size++

		if size > 5 {
			return 0, errors.New("VarInt is too big (max 5 bytes)")
		}

		if value == 0 {
			break
		}
	}

	return size, nil
}

func (v *VarInt) ReadFrom(r io.Reader) (int64, error) {
	var value int32
	var size int64
	var buf [1]byte

	for {
		_, err := r.Read(buf[:])

		if err != nil {
			return 0, err
		}

		currentByte := buf[0]

		value |= int32(currentByte&0x7F) << (size * 7)
		size++

		if size > 5 {
			return 0, errors.New("VarInt is too big (max 5 bytes)")
		}

		if currentByte&0x80 == 0 {
			break
		}
	}

	*v = VarInt(value)
	return size, nil
}

// Definition of the VarLong type and read/write functions
type VarLong int64

func (v VarLong) WriteTo(w io.Writer) (int64, error) {
	value := uint64(v)
	var size int64
	var buf [1]byte

	for {
		currentByte := value & 0x7f
		value >>= 7

		if value != 0 {
			currentByte |= 0x80
		}

		buf[0] = byte(currentByte)
		_, err := w.Write(buf[:])

		if err != nil {
			return 0, err
		}

		size++

		if size > 10 {
			return 0, errors.New("VarLong is too big (max 10 bytes)")
		}

		if value == 0 {
			break
		}
	}

	return size, nil
}

func (v *VarLong) ReadFrom(r io.Reader) (int64, error) {
	var value int64
	var size int64
	var buf [1]byte

	for {
		_, err := r.Read(buf[:])

		if err != nil {
			return 0, err
		}

		currentByte := buf[0]

		value |= int64(currentByte&0x7F) << (size * 7)
		size++

		if size > 10 {
			return 0, errors.New("VarLong is too big (max 10 bytes)")
		}

		if currentByte&0x80 == 0 {
			break
		}
	}

	*v = VarLong(value)
	return size, nil
}

type String []byte

func (s String) WriteTo(w io.Writer) (int64, error) {
	str := []byte(s)

	length := int32(len(str))

	n, err := VarInt(length).WriteTo(w)
	if err != nil {
		return 0, err
	}

	m, err := w.Write(str[:])

	return n + int64(m), nil
}

func (s *String) ReadFrom(r io.Reader) (int64, error) {
	var length VarInt

	n, err := length.ReadFrom(r)

	if err != nil {
		return 0, err
	}

	if length > 32768 {
		return 0, errors.New("String is too long (max 32768 bytes)")
	}

	buf := make([]byte, length)
	m, err := io.ReadFull(r, buf)

	if err != nil {
		return 0, err
	}

	*s = String(buf)

	return n + int64(m), nil
}
