package bits

import (
	"encoding/binary"
	"io"
)

type PackedLongArray struct {
	BitsPerEntry uint8
	Content      []int64
	Length       int32

	entriesPerLong int32
	mask           int64
}

// NewPackedLongArray creates a new PackedLongArray. It panics if bitsPerEntry is > 64.
func NewPackedLongArray(bitsPerEntry uint8, length int32) *PackedLongArray {
	if bitsPerEntry > 64 {
		panic("bits per entry must be <= 64")
	}

	entriesPerLong := 64 / int32(bitsPerEntry)
	longCount := (length + entriesPerLong - 1) / entriesPerLong

	return &PackedLongArray{
		BitsPerEntry:   bitsPerEntry,
		Content:        make([]int64, longCount),
		Length:         length,
		entriesPerLong: entriesPerLong,
		mask:           (1 << bitsPerEntry) - 1,
	}
}

// Set sets the value of the packed long array at the given index.
func (p *PackedLongArray) Set(index int32, value int64) {
	if index >= p.Length {
		panic("index out of bounds")
	}

	longIndex := index / p.entriesPerLong
	bitIndex := (index % p.entriesPerLong) * int32(p.BitsPerEntry)

	// clear current entry
	p.Content[longIndex] &= ^(p.mask << bitIndex)
	// set new value
	p.Content[longIndex] |= (value & p.mask) << bitIndex
}

// Get returns the value of the packed long array at the given index.
func (p *PackedLongArray) Get(index int32) int64 {
	if index >= p.Length {
		panic("index out of bounds")
	}

	longIndex := index / p.entriesPerLong
	bitIndex := (index % p.entriesPerLong) * int32(p.BitsPerEntry)

	// return value at position
	return (p.Content[longIndex] >> bitIndex) & p.mask
}

// WriteTo writes the packed long array to the given writer.
func (p *PackedLongArray) WriteTo(w io.Writer) (int64, error) {
	var buffer [8]byte
	total := int64(0)
	for _, long := range p.Content {
		binary.BigEndian.PutUint64(buffer[:], uint64(long))
		n, err := w.Write(buffer[:])
		total += int64(n)
		if err != nil {
			return total, err
		}
	}
	return total, nil
}
