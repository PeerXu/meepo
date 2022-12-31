package acl

type allows struct {
	rules []Entity
}

func (a *allows) Permit(x Entity) error {
	for _, r := range a.rules {
		if r.Contains(x) {
			return nil
		}
	}

	return ErrNotPermitted
}
