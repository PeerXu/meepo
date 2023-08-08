package meepo_event_listener

type Set map[string]struct{}

func NewSet() Set {
	return Set(make(map[string]struct{}))
}

func (s Set) Has(val string) bool {
	_, ok := s[val]
	return ok
}

func (s Set) Add(val string) {
	s[val] = struct{}{}
}

func (s Set) Remove(val string) {
	delete(s, val)
}
