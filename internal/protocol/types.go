package protocol

import (
	"BlueIce/internal/nbt"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"strings"
)

// Definition of the VarInt type and read/write functions
type VarInt int32

func (v *VarInt) WriteTo(w io.Writer) (int64, error) {
	value := uint32(*v)
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

func (v *VarLong) WriteTo(w io.Writer) (int64, error) {
	value := uint64(*v)
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

type PrefixedArray[T any] struct {
	Length  VarInt
	Content []T
	Reader  func(io.Reader, *T) (int64, error)
}

func (p *PrefixedArray[T]) WriteTo(w io.Writer) (int64, error) {
	p.Length = VarInt(len(p.Content))

	size, err := p.Length.WriteTo(w)
	if err != nil {
		return size, err
	}

	for _, element := range p.Content {
		n, err := serialize(w, element)
		size += n
		if err != nil {
			return size, err
		}
	}

	return size, nil
}

func (p *PrefixedArray[T]) ReadFrom(r io.Reader) (int64, error) {
	size, err := p.Length.ReadFrom(r)

	if err != nil {
		return 0, err
	}

	p.Content = make([]T, p.Length)
	for i := range p.Content {
		var element T
		n, err := p.Reader(r, &element)
		size += n
		if err != nil {
			return size, err
		}
		p.Content[i] = element
	}

	return size, nil
}

type PrefixedOptional[T any] struct {
	Present bool
	Content T

	Reader func(io.Reader, *T) (int64, error)
}

func (p *PrefixedOptional[T]) WriteTo(w io.Writer) (int64, error) {
	if p.Present {
		size, err := w.Write([]byte{1})
		if err != nil {
			return int64(size), err
		}

		n, err := serialize(w, p.Content)
		size += int(n)
		return int64(size), err
	}

	size, err := w.Write([]byte{0})
	if err != nil {
		return int64(size), err
	}
	return int64(size), err
}

func (p *PrefixedOptional[T]) ReadFrom(r io.Reader) (int64, error) {
	size := 0

	err := binary.Read(r, binary.BigEndian, &p.Present)
	size += 1
	if err != nil {
		return int64(size), err
	}

	if p.Present {
		if p.Reader == nil {
			return int64(size), errors.New("reader is nil even though field is present")
		}

		n, err := p.Reader(r, &p.Content)
		size += int(n)
		if err != nil {
			return int64(size), err
		}
	}

	return int64(size), err
}

type GameProfile struct {
	UUID    [16]byte
	Name    string
	Options PrefixedArray[GameProfileOption]
}

func (p *GameProfile) WriteTo(w io.Writer) (int64, error) {
	size, err := w.Write(p.UUID[:])
	if err != nil {
		return int64(size), err
	}

	n, err := WriteString(w, p.Name)
	size += int(n)
	if err != nil {
		return int64(size), err
	}

	n, err = p.Options.WriteTo(w)
	size += int(n)
	if err != nil {
		return int64(size), err
	}

	return int64(size), nil
}

func (p *GameProfile) ReadFrom(r io.Reader) (int64, error) {
	size, err := r.Read(p.UUID[:])
	if err != nil {
		return int64(size), err
	}

	n, err := ReadString(r, &p.Name)
	size += int(n)
	if err != nil {
		return int64(size), err
	}

	p.Options.Reader = func(r io.Reader, gpo *GameProfileOption) (int64, error) {
		return gpo.ReadFrom(r)
	}

	n, err = p.Options.ReadFrom(r)
	size += int(n)
	if err != nil {
		return int64(size), err
	}

	return int64(size), nil
}

type GameProfileOption struct {
	Name      string
	Value     string
	Signature PrefixedOptional[string]
}

func (gpo *GameProfileOption) WriteTo(w io.Writer) (int64, error) {
	size, err := WriteString(w, gpo.Name)
	if err != nil {
		return size, err
	}

	n, err := WriteString(w, gpo.Value)
	size += n
	if err != nil {
		return size, err
	}

	n, err = gpo.Signature.WriteTo(w)
	size += n
	if err != nil {
		return size, err
	}

	return size, nil
}

func (gpo *GameProfileOption) ReadFrom(r io.Reader) (int64, error) {
	size, err := ReadString(r, &gpo.Name)
	if err != nil {
		return size, err
	}

	n, err := ReadString(r, &gpo.Value)
	size += n
	if err != nil {
		return size, err
	}

	gpo.Signature.Reader = func(r io.Reader, s *string) (int64, error) {
		return ReadString(r, s)
	}

	n, err = gpo.Signature.ReadFrom(r)
	size += n
	if err != nil {
		return size, err
	}

	return size, nil
}

type Identifier struct {
	Namespace string
	Path      string
}

func NewIdentifier(namespace, path string) Identifier {
	return Identifier{
		Namespace: namespace,
		Path:      path,
	}
}

func NewIdentifierFromPath(path string) Identifier {
	return NewIdentifier("minecraft", path)
}

func NewIdentifierFromString(str string) Identifier {
	splitted := strings.Split(str, ":")
	if len(splitted) < 2 {
		return NewIdentifier("minecraft", splitted[0])
	} else if len(splitted) == 2 {
		return NewIdentifier(splitted[0], splitted[1])
	}
	return Identifier{}
}

func (id *Identifier) String() string {
	return id.Namespace + ":" + id.Path
}

func (id *Identifier) WriteTo(w io.Writer) (int64, error) {
	str := id.String()

	length := VarInt(len(str))
	size, err := length.WriteTo(w)
	if err != nil {
		return size, err
	}

	n, err := io.WriteString(w, str)
	size += int64(n)
	return size, err
}

func (id *Identifier) ReadFrom(r io.Reader) (int64, error) {

	var length VarInt
	size, err := length.ReadFrom(r)
	if err != nil {
		return size, err
	}

	buffer := make([]byte, length)
	n, err := r.Read(buffer)
	size += int64(n)
	if err != nil {
		return size, err
	}

	str := string(buffer)
	splitted := strings.Split(str, ":")

	if len(splitted) < 2 {
		id.Namespace = "minecraft"
		id.Path = splitted[0]
	} else if len(splitted) == 2 {
		id.Namespace = splitted[0]
		id.Path = splitted[1]
	} else {
		return size, errors.New("invalid identifier")
	}

	return size, nil
}

type Position struct {
	X, Z int32
	Y    int16
}

func (pos *Position) WriteTo(w io.Writer) (int64, error) {
	total := int64(0)

	var encoded = uint64(0)
	x := uint64(uint32(pos.X)) // Casting twice because uint64 should be filled with 0 and not 1
	encoded |= (x >> 31) << 25 // Writing sign
	encoded |= x & 0x1FFFFFF   // Writing payload
	encoded <<= 26             // Shifting for next writing

	z := uint64(uint32(pos.Z))
	encoded |= (z >> 31) << 25
	encoded |= z & 0x1FFFFFF
	encoded <<= 26

	y := uint64(uint16(pos.Y))
	encoded |= (y >> 15) << 11
	encoded |= y & 0x7FF

	var buffer [8]byte
	binary.BigEndian.PutUint64(buffer[:], encoded)
	total += int64(len(buffer))
	_, err := w.Write(buffer[:])
	return total, err
}

func (pos *Position) ReadFrom(r io.Reader) (int64, error) {
	var buffer [8]byte
	n, err := r.Read(buffer[:])
	if err != nil {
		return int64(n), err
	}

	encoded := binary.BigEndian.Uint64(buffer[:])
	y := uint64(0)
	y |= encoded & 0x7FF
	encoded >>= 11

	if encoded&0x1 != 0 {
		y |= 0xF800
	}

	pos.Y = int16(y)

	encoded >>= 1

	z := uint64(0)
	z |= encoded & 0x1FFFFFF
	encoded >>= 25
	if encoded&0x1 != 0 {
		z |= 0xFC000000
	}

	pos.Z = int32(z)

	encoded >>= 1

	x := uint64(0)
	x |= encoded & 0x1FFFFFF
	encoded >>= 25
	if encoded&0x1 != 0 {
		x |= 0xFC000000
	}

	pos.X = int32(x)

	return int64(n), nil
}

type NBTValue struct {
	Value any
}

func (v *NBTValue) WriteTo(w io.Writer) (int64, error) {
	return nbt.WriteNBT(w, v.Value)
}

type BitSet struct {
	content []int64
}

func (b *BitSet) Set(index int32) {
	indexInLong := index % 64
	longIndex := index / 64
	if int32(len(b.content)) <= longIndex {
		b.content = append(b.content, 0)
	}
	b.content[longIndex] |= 1 << indexInLong
}

func (b *BitSet) SetRange(start, end int32) {
	for i := start; i <= end; i++ {
		b.Set(i)
	}
}

func (b *BitSet) Get(index int32) bool {
	indexInLong := index % 64
	longIndex := index / 64
	if int32(len(b.content)) <= longIndex {
		return false
	}
	return (b.content[longIndex]>>indexInLong)&1 == 1
}

func (b *BitSet) WriteTo(w io.Writer) (int64, error) {
	length := VarInt(len(b.content))
	total, err := length.WriteTo(w)
	if err != nil {
		return total, err
	}

	for _, long := range b.content {
		n, err := WriteInt64(w, long)
		total += n
		if err != nil {
			return total, err
		}
	}

	return total, nil
}

type Heightmap struct {
	Type        VarInt
	Data        [256]uint16
	WorldHeight uint16
}

func (h *Heightmap) Set(x, z int, height uint16) {
	h.Data[x+z*16] = height
}

func (h *Heightmap) Get(x, z int) uint16 {
	return h.Data[x+z*16]
}

func (h *Heightmap) WriteTo(w io.Writer) (int64, error) {
	total := int64(0)

	bitsPerEntry := int(math.Ceil(math.Log2(float64(h.WorldHeight) + 1)))
	entriesPerLong := 64 / bitsPerEntry
	numLongs := VarInt(int32(math.Ceil(float64(len(h.Data)) / float64(entriesPerLong))))

	n, err := h.Type.WriteTo(w)
	total += n
	if err != nil {
		return total, err
	}

	n, err = numLongs.WriteTo(w)
	total += n
	if err != nil {
		return total, err
	}

	var buffer int64
	mask := int64(1<<bitsPerEntry) - 1

	for i, height := range h.Data {
		bitIndex := (i % entriesPerLong) * bitsPerEntry

		buffer |= (int64(height) & mask) << bitIndex

		if ((i+1)%entriesPerLong) == 0 || i == len(h.Data)-1 {
			n, err := WriteInt64(w, buffer)
			total += n
			if err != nil {
				return total, err
			}
			buffer = 0
		}
	}

	return total, nil
}

type LightArray [2048]uint8

func (l *LightArray) WriteTo(w io.Writer) (int64, error) {
	total := int64(0)
	n, err := WriteInt32(w, int32(len(l)))
	total += n
	if err != nil {
		return total, err
	}
	for _, value := range l {
		n, err := w.Write([]byte{value})
		total += int64(n)
		if err != nil {
			return total, err
		}
	}
	return total, nil
}
