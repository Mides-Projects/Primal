package protocol

type PlayerJoinedNetwork struct {
	Username string
	XUID     string

	ServerName string
}

// ShieldId returns the packet ID.
func (pk *PlayerJoinedNetwork) ShieldId() int32 {
	return 0x00
}

// Unmarshal unmarshals the object from the given IO.
func (pk *PlayerJoinedNetwork) Unmarshal(io IO) {
	io.String(&pk.Username)
	io.String(&pk.XUID)

	io.String(&pk.ServerName)
}

// Marshal marshals the object into the given IO.
func (pk *PlayerJoinedNetwork) Marshal(io IO) {
	io.String(&pk.Username)
	io.String(&pk.XUID)

	io.String(&pk.ServerName)
}
