package request

type ServerUpBody struct {
	Port int64 `json:"port"`

	Directory string `json:"directory"`
	Motd      string `json:"motd"`

	BungeeCord bool `json:"bungee-cord"`
	OnlineMode bool `json:"online-mode"`

	MaxSlots int64    `json:"max-slots"`
	Plugins  []string `json:"plugins"`
}
