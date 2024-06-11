package server

import (
	"github.com/gorilla/mux"
	"github.com/holypvp/primal/server/routes"
	"sync"
)

var (
	mutex     sync.Mutex
	portMutex sync.Mutex

	servers       = map[string]*ServerInfo{}
	serversByPort = map[int64]string{}

	instance *ServerModule
)

type ServerModule struct{}

func (m *ServerModule) LookupById(id string) *ServerInfo {
	mutex.Lock()
	defer mutex.Unlock()

	serverInfo, ok := servers[id]
	if !ok {
		return nil
	}

	return serverInfo
}

func (m *ServerModule) LookupByPort(port int64) *ServerInfo {
	portMutex.Lock()
	defer portMutex.Unlock()

	id, ok := serversByPort[port]
	if !ok {
		return nil
	}

	return m.LookupById(id)
}

func (m *ServerModule) Append(serverInfo *ServerInfo) {
	mutex.Lock()
	servers[serverInfo.Id()] = serverInfo
	mutex.Unlock()

	portMutex.Lock()
	serversByPort[serverInfo.Port()] = serverInfo.Id()
	portMutex.Unlock()
}

func (m *ServerModule) Destroy(serverId string) {
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

func (m *ServerModule) Values() []*ServerInfo {
	mutex.Lock()
	defer mutex.Unlock()

	values := make([]*ServerInfo, 0, len(servers))
	for _, serverInfo := range servers {
		values = append(values, serverInfo)
	}

	return values
}

func LoadAll(router *mux.Router) {
	router.HandleFunc("/api/v2/servers/lookup", routes.LookupServers).Methods("GET")
	router.HandleFunc("/api/v2/servers/{id}/down", routes.ServerDownRoute).Methods("POST")
	router.HandleFunc("/api/v2/servers/{id}/tick", routes.ServerTickRoute).Methods("PATCH")
}

func Module() *ServerModule {
	if instance == nil {
		instance = &ServerModule{}
	}

	return instance
}
