package acl

type chain struct {
	nodes []struct {
		acl     Acl
		emitNow func(error) bool
	}
}

func (c *chain) Permit(x Entity) error {
	for _, n := range c.nodes {
		err := n.acl.Permit(x)
		if n.emitNow(err) {
			return err
		}
	}

	// DROP ANY CHALLENGE if not match before
	return ErrNotPermitted
}
