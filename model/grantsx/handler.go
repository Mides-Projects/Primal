package grantsx

type Handler interface {
	// HandleLookup is a function that handles the lookup of a grant
	HandleLookup(filter string, src string, state string, value string) error
}
