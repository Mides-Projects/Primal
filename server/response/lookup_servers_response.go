package response

import "github.com/holypvp/primal/server"

type LookupServersResponse struct {
	Servers []server.ServerInfo `json:"values"`
}
