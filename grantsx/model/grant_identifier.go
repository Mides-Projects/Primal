package model

import "errors"

type GrantIdentifier struct {
	key   string
	value string
}

// Key returns the key of the grant.
func (g GrantIdentifier) Key() string {
	return g.key
}

// Value returns the value of the grant.
func (g GrantIdentifier) Value() string {
	return g.value
}

func (g GrantIdentifier) Marshal() map[string]interface{} {
	return map[string]interface{}{
		"key":   g.key,
		"value": g.value,
	}
}

// UnmarshalIdentifier unmarshals an value from a map.
func UnmarshalIdentifier(body map[string]interface{}) (GrantIdentifier, error) {
	k, ok := body["key"].(string)
	if !ok {
		return GrantIdentifier{}, errors.New("key is not a string")
	}

	i, ok := body["value"].(string)
	if !ok {
		return GrantIdentifier{}, errors.New("value is not a string")
	}

	return GrantIdentifier{
		key:   k,
		value: i,
	}, nil
}
