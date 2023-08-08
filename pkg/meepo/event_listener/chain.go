package meepo_event_listener

type Chain []string

func (c Chain) Head() string {
	return c[0]
}

func (c Chain) Rest() Chain {
	return Chain(c[1:])
}

func (c Chain) IsNull() bool {
	return c.Head() == ""
}
