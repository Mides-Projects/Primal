package server

type ServerInfo struct {
	id     string
	port   int64
	groups []string

	playersCount   int
	maxSlots       int
	heartbeat      int64
	players        []string
	bungeeCord     bool
	onlineMode     bool
	activeThreads  int
	daemonThreads  int
	motd           *string
	ticksPerSecond float64
	directory      string
	fullTicks      int64
	initialTime    int64
	plugins        []string
}

func NewServerInfo(id string, port int64) *ServerInfo {
	return &ServerInfo{
		id:     id,
		port:   port,
		groups: []string{},

		playersCount:   0,
		maxSlots:       0,
		heartbeat:      0,
		players:        []string{},
		bungeeCord:     false,
		onlineMode:     false,
		activeThreads:  0,
		daemonThreads:  0,
		motd:           nil,
		ticksPerSecond: 0,
		directory:      "",
		fullTicks:      0,
		initialTime:    0,
		plugins:        []string{},
	}
}

func (i *ServerInfo) Id() string {
	return i.id
}

func (i *ServerInfo) Port() int64 {
	return i.port
}

func (i *ServerInfo) SetPort(port int64) {
	i.port = port
}

func (i *ServerInfo) Groups() []string {
	return i.groups
}

func (i *ServerInfo) AddGroup(group string) {
	i.groups = append(i.groups, group)
}

func (i *ServerInfo) SetGroups(groups []string) {
	i.groups = groups
}

func (i *ServerInfo) PlayersCount() int {
	return i.playersCount
}

func (i *ServerInfo) SetPlayersCount(count int) {
	i.playersCount = count
}

func (i *ServerInfo) MaxSlots() int {
	return i.maxSlots
}

func (i *ServerInfo) SetMaxSlots(slots int) {
	i.maxSlots = slots
}

func (i *ServerInfo) Heartbeat() int64 {
	return i.heartbeat
}

func (i *ServerInfo) SetHeartbeat(heartbeat int64) {
	i.heartbeat = heartbeat
}

func (i *ServerInfo) Players() []string {
	return i.players
}

func (i *ServerInfo) SetPlayers(players []string) {
	i.players = players
}

func (i *ServerInfo) BungeeCord() bool {
	return i.bungeeCord
}

func (i *ServerInfo) SetBungeeCord(bungeeCord bool) {
	i.bungeeCord = bungeeCord
}

func (i *ServerInfo) OnlineMode() bool {
	return i.onlineMode
}

func (i *ServerInfo) SetOnlineMode(onlineMode bool) {
	i.onlineMode = onlineMode
}

func (i *ServerInfo) ActiveThreads() int {
	return i.activeThreads
}

func (i *ServerInfo) SetActiveThreads(activeThreads int) {
	i.activeThreads = activeThreads
}

func (i *ServerInfo) DaemonThreads() int {
	return i.daemonThreads
}

func (i *ServerInfo) SetDaemonThreads(daemonThreads int) {
	i.daemonThreads = daemonThreads
}

func (i *ServerInfo) Motd() *string {
	return i.motd
}

func (i *ServerInfo) SetMotd(motd *string) {
	i.motd = motd
}

func (i *ServerInfo) TicksPerSecond() float64 {
	return i.ticksPerSecond
}

func (i *ServerInfo) SetTicksPerSecond(ticksPerSecond float64) {
	i.ticksPerSecond = ticksPerSecond
}

func (i *ServerInfo) Directory() string {
	return i.directory
}

func (i *ServerInfo) SetDirectory(directory string) {
	i.directory = directory
}

func (i *ServerInfo) FullTicks() int64 {
	return i.fullTicks
}

func (i *ServerInfo) SetFullTicks(fullTicks int64) {
	i.fullTicks = fullTicks
}

func (i *ServerInfo) InitialTime() int64 {
	return i.initialTime
}

func (i *ServerInfo) SetInitialTime(initialTime int64) {
	i.initialTime = initialTime
}

func (i *ServerInfo) Plugins() []string {
	return i.plugins
}

func (i *ServerInfo) SetPlugins(plugins []string) {
	i.plugins = plugins
}
