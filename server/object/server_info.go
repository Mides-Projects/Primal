package object

import "github.com/holypvp/primal/server/model"

type ServerInfo struct {
	id     string
	port   int64
	groups []string

	playersCount   int
	maxSlots       int64
	heartbeat      int64
	players        []string
	bungeeCord     bool
	onlineMode     bool
	activeThreads  int
	daemonThreads  int
	motd           string
	ticksPerSecond float64
	directory      string
	fullTicks      float64
	initialTime    int64
	plugins        []string
}

// NewServerInfo creates a new server info.
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
		motd:           "",
		ticksPerSecond: 0,
		directory:      "",
		fullTicks:      0,
		initialTime:    0,
		plugins:        []string{},
	}
}

// Id returns the server's id.
func (i *ServerInfo) Id() string {
	return i.id
}

// Port returns the server's port.
func (i *ServerInfo) Port() int64 {
	return i.port
}

// SetPort sets the server's port.
func (i *ServerInfo) SetPort(port int64) {
	i.port = port
}

// Groups returns the server's groups.
func (i *ServerInfo) Groups() []string {
	return i.groups
}

// AddGroup adds a group to the server's groups.
func (i *ServerInfo) AddGroup(group string) {
	i.groups = append(i.groups, group)
}

// SetGroups sets the server's groups.
func (i *ServerInfo) SetGroups(groups []string) {
	i.groups = groups
}

// PlayersCount returns the server's player count.
func (i *ServerInfo) PlayersCount() int {
	return i.playersCount
}

// SetPlayersCount sets the server's player count.
func (i *ServerInfo) SetPlayersCount(count int) {
	i.playersCount = count
}

// MaxSlots returns the server's max slots.
func (i *ServerInfo) MaxSlots() int64 {
	return i.maxSlots
}

// SetMaxSlots sets the server's max slots.
func (i *ServerInfo) SetMaxSlots(slots int64) {
	i.maxSlots = slots
}

// Heartbeat returns the server's heartbeat.
func (i *ServerInfo) Heartbeat() int64 {
	return i.heartbeat
}

// SetHeartbeat sets the server's heartbeat.
func (i *ServerInfo) SetHeartbeat(heartbeat int64) {
	i.heartbeat = heartbeat
}

// Players returns the server's players.
func (i *ServerInfo) Players() []string {
	return i.players
}

// SetPlayers sets the server's players.
func (i *ServerInfo) SetPlayers(players []string) {
	i.players = players
}

// BungeeCord returns the server's bungee cord status.
func (i *ServerInfo) BungeeCord() bool {
	return i.bungeeCord
}

// SetBungeeCord sets the server's bungee cord status.
func (i *ServerInfo) SetBungeeCord(bungeeCord bool) {
	i.bungeeCord = bungeeCord
}

// OnlineMode returns the server's online mode status.
func (i *ServerInfo) OnlineMode() bool {
	return i.onlineMode
}

// SetOnlineMode sets the server's online mode status.
func (i *ServerInfo) SetOnlineMode(onlineMode bool) {
	i.onlineMode = onlineMode
}

// ActiveThreads returns the server's active threads.
func (i *ServerInfo) ActiveThreads() int {
	return i.activeThreads
}

// SetActiveThreads sets the server's active threads.
func (i *ServerInfo) SetActiveThreads(activeThreads int) {
	i.activeThreads = activeThreads
}

// DaemonThreads returns the server's daemon threads.
func (i *ServerInfo) DaemonThreads() int {
	return i.daemonThreads
}

// SetDaemonThreads sets the server's daemon threads.
func (i *ServerInfo) SetDaemonThreads(daemonThreads int) {
	i.daemonThreads = daemonThreads
}

// Motd returns the server's motd.
func (i *ServerInfo) Motd() string {
	return i.motd
}

// SetMotd sets the server's motd.
func (i *ServerInfo) SetMotd(motd string) {
	i.motd = motd
}

// TicksPerSecond returns the server's ticks per second.
func (i *ServerInfo) TicksPerSecond() float64 {
	return i.ticksPerSecond
}

// SetTicksPerSecond sets the server's ticks per second.
func (i *ServerInfo) SetTicksPerSecond(ticksPerSecond float64) {
	i.ticksPerSecond = ticksPerSecond
}

// Directory returns the server's directory.
func (i *ServerInfo) Directory() string {
	return i.directory
}

// SetDirectory sets the server's directory.
func (i *ServerInfo) SetDirectory(directory string) {
	i.directory = directory
}

// FullTicks returns the server's full ticks.
func (i *ServerInfo) FullTicks() float64 {
	return i.fullTicks
}

// SetFullTicks sets the server's full ticks.
func (i *ServerInfo) SetFullTicks(fullTicks float64) {
	i.fullTicks = fullTicks
}

// InitialTime returns the server's initial time.
func (i *ServerInfo) InitialTime() int64 {
	return i.initialTime
}

// SetInitialTime sets the server's initial time.
func (i *ServerInfo) SetInitialTime(initialTime int64) {
	i.initialTime = initialTime
}

// Plugins returns the server's plugins.
func (i *ServerInfo) Plugins() []string {
	return i.plugins
}

// SetPlugins sets the server's plugins.
func (i *ServerInfo) SetPlugins(plugins []string) {
	i.plugins = plugins
}

func (i *ServerInfo) ToModel() model.ServerInfoModel {
	return model.ServerInfoModel{
		Id:          i.Id(),
		Port:        i.Port(),
		Groups:      i.Groups(),
		MaxSlots:    i.MaxSlots(),
		Heartbeat:   i.Heartbeat(),
		BungeeCord:  i.BungeeCord(),
		OnlineMode:  i.OnlineMode(),
		InitialTime: i.InitialTime(),
	}
}
