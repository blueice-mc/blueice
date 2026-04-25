package bits

type PackedLongArray struct {
	BitsPerEntry uint8
	Content      []int64
}

// NewPackedLongArray creates a new PackedLongArray. It panics if bitsPerEntry is > 64.
func NewPackedLongArray(bitsPerEntry uint8) *PackedLongArray {

	if bitsPerEntry > 64 {
		panic("bits per entry must be <= 64")
	}

	return &PackedLongArray{
		BitsPerEntry: bitsPerEntry,
		Content:      make([]int64, 0),
	}
}
