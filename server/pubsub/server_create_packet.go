package pubsub

type ServerCreatePacket struct {
	Id   string `json:"_id"`
	Port int64  `json:"port"`
}

func NewServerCreatePacket(id string, port int64) ServerCreatePacket {
	return ServerCreatePacket{
		Id:   id,
		Port: port,
	}
}
