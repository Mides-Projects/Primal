package pubsub

type ServerStatusPacket struct {
	ServerId string `json:"_id"`
}

func NewServerStatusPacket(serverId string) ServerStatusPacket {
	return ServerStatusPacket{
		ServerId: serverId,
	}
}
