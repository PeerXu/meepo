package acl

var (
	ANY_SYM = "*"
	ANY     = NewEntity(ANY_SYM, ANY_SYM, ANY_SYM, ANY_SYM)
)

type Acl interface {
	Permit(Entity) error
}

func FromString(s string) (Acl, error) {
	rs, err := ParseRules(s)
	if err != nil {
		return nil, err
	}
	return FromRules(rs)
}

func FromRules(rs []Rule) (Acl, error) {
	c := &chain{}
	for _, r := range rs {
		var ss []string
		var es []Entity
		var fn func([]Entity) Acl
		var emitNow func(error) bool
		if r.Allow != "" {
			ss = append(ss, r.Allow)
			fn = func(rs []Entity) Acl { return &allows{rs} }
			emitNow = func(err error) bool { return err == nil }
		} else if r.Block != "" {
			ss = append(ss, r.Block)
			fn = func(rs []Entity) Acl { return &blocks{rs} }
			emitNow = func(err error) bool { return err != nil }
		} else if len(r.Allows) > 0 {
			ss = r.Allows
			fn = func(rs []Entity) Acl { return &allows{rs} }
			emitNow = func(err error) bool { return err == nil }
		} else if len(r.Block) > 0 {
			ss = r.Blocks
			fn = func(rs []Entity) Acl { return &blocks{rs} }
			emitNow = func(err error) bool { return err != nil }
		}

		if len(ss) == 0 {
			continue
		}

		for _, s := range ss {
			e, err := Parse(s)
			if err != nil {
				return nil, err
			}
			es = append(es, e)
		}
		c.nodes = append(c.nodes, struct {
			acl     Acl
			emitNow func(error) bool
		}{
			acl:     fn(es),
			emitNow: emitNow,
		})
	}

	return c, nil
}
