package server

import (
	"github.com/holypvp/primal/server/object"
	"golang.org/x/exp/maps"
	"sync"
)

var inst = &service{
	servers: make(map[string]*object.ServerInfo),
	groups:  make(map[string]*object.ServerGroup),
}

// Service is a service that provides server information.
// It is used to cache server information and look up servers by their ID or port.
// It is also used to cache server groups and look up server groups by their ID.
// The service is thread-safe.
type service struct {
	servers   map[string]*object.ServerInfo
	serversMu sync.Mutex

	groups   map[string]*object.ServerGroup
	groupsMu sync.Mutex
}

// LookupById looks up a server by its ID.
func (s *service) LookupById(id string) *object.ServerInfo {
	s.serversMu.Lock()
	defer s.serversMu.Unlock()

	server, ok := s.servers[id]
	if !ok {
		return nil
	}

	return server
}

// LookupByPort looks up a server by its port.
func (s *service) LookupByPort(port int64) *object.ServerInfo {
	s.serversMu.Lock()
	defer s.serversMu.Unlock()

	for _, server := range s.servers {
		if server.Port() == port {
			return server
		}
	}

	return nil
}

// CacheServer caches a server in the service.
func (s *service) CacheServer(server *object.ServerInfo) {
	s.serversMu.Lock()
	s.servers[server.Id()] = server
	s.serversMu.Unlock()
}

// DestroyServer removes a server from the cache.
func (s *service) DestroyServer(id string) {
	s.serversMu.Lock()
	delete(s.servers, id)
	s.serversMu.Unlock()
}

func (s *service) Servers() []*object.ServerInfo {
	s.serversMu.Lock()
	defer s.serversMu.Unlock()

	return maps.Values(s.servers)
}

// LookupGroupById looks up a server group by its ID.
func (s *service) LookupGroupById(id string) *object.ServerGroup {
	s.groupsMu.Lock()
	defer s.groupsMu.Unlock()

	group, ok := s.groups[id]
	if !ok {
		return nil
	}

	return group
}

// CacheGroup caches a server group in the service.
func (s *service) CacheGroup(group *object.ServerGroup) {
	s.groupsMu.Lock()
	s.groups[group.Id()] = group
	s.groupsMu.Unlock()
}

// DestroyGroup removes a server group from the cache.
func (s *service) DestroyGroup(id string) {
	s.groupsMu.Lock()
	delete(s.groups, id)
	s.groupsMu.Unlock()
}

func (s *service) Groups() []*object.ServerGroup {
	s.groupsMu.Lock()
	defer s.groupsMu.Unlock()

	return maps.Values(s.groups)
}

// Service returns the server service instance.
func Service() *service {
	return inst
}
