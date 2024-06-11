package server

import (
	"sync"
)

var (
	mutex     sync.Mutex
	portMutex sync.Mutex

	servers       = map[string]*ServerInfo{}
	serversByPort = map[int64]string{}

	instance *ServerService
)

type ServerService struct{}

func (m *ServerService) LookupById(id string) *ServerInfo {
	mutex.Lock()
	defer mutex.Unlock()

	serverInfo, ok := servers[id]
	if !ok {
		return nil
	}

	return serverInfo
}

func (m *ServerService) LookupByPort(port int64) *ServerInfo {
	portMutex.Lock()
	defer portMutex.Unlock()

	id, ok := serversByPort[port]
	if !ok {
		return nil
	}

	return m.LookupById(id)
}

func (m *ServerService) Append(serverInfo *ServerInfo) {
	mutex.Lock()
	servers[serverInfo.Id()] = serverInfo
	mutex.Unlock()

	portMutex.Lock()
	serversByPort[serverInfo.Port()] = serverInfo.Id()
	portMutex.Unlock()
}

func (m *ServerService) Destroy(serverId string) {
	mutex.Lock()
	defer mutex.Unlock()

	serverInfo, ok := servers[serverId]
	if !ok {
		return
	}

	portMutex.Lock()
	defer portMutex.Unlock()

	delete(servers, serverId)
	delete(serversByPort, serverInfo.Port())
}

func (m *ServerService) Values() []*ServerInfo {
	mutex.Lock()
	defer mutex.Unlock()

	values := make([]*ServerInfo, 0, len(servers))
	for _, serverInfo := range servers {
		values = append(values, serverInfo)
	}

	return values
}

func Service() *ServerService {
	if instance == nil {
		instance = &ServerService{}
	}

	return instance
}
