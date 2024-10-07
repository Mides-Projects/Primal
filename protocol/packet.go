package protocol

type Packet interface {
	// ShieldId returns the packet ID.
	ShieldId() int32

	// Unmarshal unmarshals the object from the given IO.
	Unmarshal(io IO)
	// Marshal marshals the object into the given IO.
	Marshal(io IO)
}
