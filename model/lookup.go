package model

type Lookup struct {
  key string
  value string

  keep bool
}

func (l Lookup) Key(k string) Lookup {
  l.key = k

  return l
}

func (l Lookup) Value(v string) Lookup {
  l.value = v

  return l
}

// Keep mark the field 'keep' as true
// this means we going to cache permanent
// the result. If keep is false, it only is added to ttl
func (l Lookup) Keep() Lookup {
  l.keep = true

  return l
}
