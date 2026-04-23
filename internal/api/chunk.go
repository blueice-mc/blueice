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

type SerializedChunkGenerationEvent struct {
	Chunk     *SerializedChunk // the generated chunk
	Generated bool             // if the event handler generated chunk data. if this turns out to be false, a fallback chunk generator will be used.
}
