package api

type SerializedChunk struct {
	X, Z      int32
	MinY      int16
	Height    uint16
	Sections  []SerializedChunkSection
	HeightMap [256]uint16
}

type SerializedChunkSection struct {
	Y      int16
	Blocks [4096]uint32
}
