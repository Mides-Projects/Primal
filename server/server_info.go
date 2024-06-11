package server

type ServerInfo struct {
	id   string
	port int64

	groups []string
}

func (i *ServerInfo) Id() string {
	return i.id
}

func (i *ServerInfo) Port() int64 {
	return i.port
}

func (i *ServerInfo) Groups() []string {
	return i.groups
}

func (i *ServerInfo) AddGroup(group string) {
	i.groups = append(i.groups, group)
}

func (i *ServerInfo) ClearGroups() {
	i.groups = []string{}
}
