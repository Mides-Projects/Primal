package protocol

type PlayerChangedServer struct {
    Username string
    XUID     string

    OldServerName string
    NewServerName string
}

// ShieldId returns the packet ID.
func (pk *PlayerChangedServer) ShieldId() int32 {
    return 0x01
}

// Unmarshal unmarshals the object from the given IO.
func (pk *PlayerChangedServer) Unmarshal(io IO) {
    io.String(&pk.Username)
    io.String(&pk.XUID)

    io.String(&pk.OldServerName)
    io.String(&pk.NewServerName)
}

// Marshal marshals the object into the given IO.
func (pk *PlayerChangedServer) Marshal(io IO) {
    io.String(&pk.Username)
    io.String(&pk.XUID)

    io.String(&pk.OldServerName)
    io.String(&pk.NewServerName)
}
