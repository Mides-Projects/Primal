package pubsub

type server_create_packet struct {
	Id   string `json:"_id"`
	Port int64  `json:"port"`
}

func NewServerCreatePacket(id string, port int64) *server_create_packet {
	return &server_create_packet{
		Id:   id,
		Port: port,
	}
}
