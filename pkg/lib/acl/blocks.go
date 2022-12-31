package acl

type blocks struct {
	rlues []Entity
}

func (b *blocks) Permit(x Entity) error {
	for _, r := range b.rlues {
		if r.Contains(x) {
			return ErrNotPermitted
		}
	}

	return nil
}
