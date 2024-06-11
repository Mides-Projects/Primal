package pubsub

type ServerDownPacket struct {
	ServerId string `json:"_id"`
}

func NewServerDownPacket(serverId string) ServerDownPacket {
	return ServerDownPacket{
		ServerId: serverId,
	}
}
