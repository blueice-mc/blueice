package api

type SerializedPlayer struct {
	UUID     [16]byte
	Name     string
	EntityID int32

	X, Y, Z    float64
	Yaw, Pitch float32

	Gamemode uint8
	Health   float32
	Food     int32
}

type SerializedLoginEvent struct {
	UUID          [16]byte
	Name          string
	Cancelled     bool
	CancelMessage string
}

type SerializedJoinEvent struct {
	Player      SerializedPlayer
	JoinMessage string
}
