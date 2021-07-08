package meepo

import (
	"fmt"
	"net"
	"strings"

	"github.com/PeerXu/meepo/pkg/ofn"
	"github.com/stretchr/objx"
)

var (
	ACL_ANY        string = "*"
	ACL_ANY_ADDR          = &aclAddr{n: "*", s: "*"}
	ACL_ANY_ENTITY        = NewAclEntity(ACL_ANY, ACL_ANY_ADDR)
)

type aclAddr struct {
	n string
	s string
}

func (a *aclAddr) Network() string {
	return a.n
}

func (a *aclAddr) String() string {
	return a.s
}

func NewAclAddr(n, s string) net.Addr {
	return &aclAddr{n: n, s: s}
}

func AddrContains(a, b net.Addr) bool {
	if a.Network() != ACL_ANY && a.Network() != b.Network() {
		return false
	}

	if a.String() == ACL_ANY {
		return true
	}

	aHost, aPort, err := net.SplitHostPort(a.String())
	if err != nil {
		return false
	}

	bHost, bPort, err := net.SplitHostPort(b.String())
	if err != nil {
		return false
	}

	if aPort != ACL_ANY && aPort != bPort {
		return false
	}

	if aHost != ACL_ANY && aHost != bHost {
		if _, aIPNet, err := net.ParseCIDR(aHost); err != nil {
			if aHost != bHost {
				return false
			}
		} else {
			if !aIPNet.Contains(net.ParseIP(bHost)) {
				return false
			}
		}
	}

	return true
}

type AclEntity interface {
	ID() string
	Addr() net.Addr
	Contains(AclEntity) bool
	String() string
}

type aclEntity struct {
	id   string
	addr net.Addr
}

func NewAclEntity(id string, addr net.Addr) AclEntity {
	return &aclEntity{
		id:   id,
		addr: addr,
	}
}

func (e *aclEntity) ID() string {
	return e.id
}

func (e *aclEntity) Addr() net.Addr {
	return e.addr
}

func (a *aclEntity) Contains(b AclEntity) bool {
	if a.ID() != ACL_ANY && a.ID() != b.ID() {
		return false
	}

	if !AddrContains(a.Addr(), b.Addr()) {
		return false
	}

	return true
}

func (e *aclEntity) String() string {
	return fmt.Sprintf("<AclEntity %s:%s:%s>", e.ID(), e.Addr().Network(), e.Addr().String())
}

type AclChallenge interface {
	Source() AclEntity
	Destination() AclEntity
	String() string
}

type aclChallenge struct {
	source      AclEntity
	destination AclEntity
}

func (c *aclChallenge) Source() AclEntity {
	return c.source
}

func (c *aclChallenge) Destination() AclEntity {
	return c.destination
}

func (c *aclChallenge) String() string {
	return fmt.Sprintf("<AclChallenge source:%s destination:%s>", c.Source(), c.Destination())
}

func NewAclChallenge(source, destination AclEntity) AclChallenge {
	return &aclChallenge{
		source:      source,
		destination: destination,
	}
}

type AclPolicy interface {
	Source() AclEntity
	Destination() AclEntity
	Contains(AclChallenge) bool
	String() string
}

// ParseAclPolicy:
//   Policy: [[[<src-id>:]<src-addr-network>:]<src-addr-host>:<src-addr-port>,][[<dst-id>:]<dst-addr-network>:]<dst-addr-host>:<dst-addr-port>
//   Rules:
//     src-addr-network and dst-addr-network only support tcp or socks5 now.
//   Examples:
//     * => *:*:*:*,*:*:*:*
//     192.168.1.1:22 => *:*:*:*,*:*:192.168.1.1:22
//     10.0.0.0/24:* => *:*:*:*,*:*:10.0.0.0/24:*
//     tcp:10.0.0.1:22 => *:*:*:*,*:tcp:10.0.0.1:22
//     a:*:10.0.0.1:22 => *:*:*:*,a:*:10.0.0.1:22
//     socks5:*:*,127.0.0.1:80 => *:socks5:*:*,*:*:127.0.0.1:80
func ParseAclPolicy(s string) (p AclPolicy, err error) {
	var src, dst AclEntity

	ss := strings.SplitN(s, ",", 2)
	switch len(ss) {
	case 0:
		return nil, ErrInvalidAclPolicyString
	case 1:
		src = ACL_ANY_ENTITY
		if dst, err = parseAclEntity(ss[0]); err != nil {
			return nil, err
		}
	case 2:
		if src, err = parseAclEntity(ss[0]); err != nil {
			return nil, err
		}
		if dst, err = parseAclEntity(ss[1]); err != nil {
			return nil, err
		}
	default:
		return nil, ErrInvalidAclPolicyString
	}

	p = NewAclPolicy(src, dst)
	return
}

func ParseAclPolicies(ss []string) (ps []AclPolicy, err error) {
	for _, s := range ss {
		p, err := ParseAclPolicy(s)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func parseAclEntity(s string) (e AclEntity, err error) {
	id, network, host, port := ACL_ANY, ACL_ANY, ACL_ANY, ACL_ANY

	ss := strings.SplitN(s, ":", 4)
	switch len(ss) {
	case 0:
		return nil, ErrInvalidAclPolicyString
	case 1:
		if ss[0] != ACL_ANY {
			return nil, ErrInvalidAclPolicyString
		}
		e = ACL_ANY_ENTITY
	case 2:
		host = ss[0]
		port = ss[1]
	case 3:
		network = ss[0]
		host = ss[1]
		port = ss[2]
	case 4:
		id = ss[0]
		network = ss[1]
		host = ss[2]
		port = ss[3]
	default:
		return nil, ErrInvalidAclPolicyString
	}

	if e == nil {
		e = NewAclEntity(id, NewAclAddr(network, net.JoinHostPort(host, port)))
	}

	return
}

type aclPolicy struct {
	source      AclEntity
	destination AclEntity
}

func (p *aclPolicy) Source() AclEntity {
	return p.source
}

func (p *aclPolicy) Destination() AclEntity {
	return p.destination
}

func (p *aclPolicy) Contains(c AclChallenge) bool {
	return p.Source().Contains(c.Source()) && p.Destination().Contains(c.Destination())
}

func (p *aclPolicy) String() string {
	return fmt.Sprintf("<AclPolicy source:%s destination:%s>", p.Source(), p.Destination())
}

func NewAclPolicy(source, destination AclEntity) AclPolicy {
	return &aclPolicy{
		source:      source,
		destination: destination,
	}
}

type Acl interface {
	Allowed(AclChallenge) error
}

type chainList struct {
	allows []*allowList
	blocks []*blockList
}

func (l *chainList) Allowed(c AclChallenge) error {
	for _, b := range l.blocks {
		if err := b.Allowd(c); err != nil {
			return err
		}
	}

	for _, a := range l.allows {
		if err := a.Allowed(c); err == nil {
			return nil
		}
	}

	return ErrAclNotAllowed
}

type allowList struct {
	policies []AclPolicy
}

func (l *allowList) Allowed(c AclChallenge) error {
	for _, p := range l.policies {
		if p.Contains(c) {
			return nil
		}
	}

	return ErrAclNotAllowed
}

type blockList struct {
	policies []AclPolicy
}

func (l *blockList) Allowd(c AclChallenge) error {
	for _, p := range l.policies {
		if p.Contains(c) {
			return ErrAclNotAllowed
		}
	}

	return nil
}

type NewAclOption = ofn.OFN

func WithAllowPolicies(ps []AclPolicy) NewAclOption {
	return func(o objx.Map) {
		o["allowPolicies"] = append(o["allowPolicies"].([]AclPolicy), ps...)
	}
}

func WithBlockPolicies(ps []AclPolicy) NewAclOption {
	return func(o objx.Map) {
		o["blockPolicies"] = append(o["blockPolicies"].([]AclPolicy), ps...)
	}
}

func defaultNewAclOption() objx.Map {
	return objx.New(map[string]interface{}{
		"allowPolicies": []AclPolicy{},
		"blockPolicies": []AclPolicy{},
	})
}

func NewAcl(opts ...NewAclOption) Acl {
	o := defaultNewAclOption()
	for _, opt := range opts {
		opt(o)
	}

	allows := o["allowPolicies"].([]AclPolicy)
	blocks := o["blockPolicies"].([]AclPolicy)

	if len(allows) == 0 {
		allows, _ = ParseAclPolicies([]string{"*"})
	}

	return &chainList{
		allows: []*allowList{{allows}},
		blocks: []*blockList{{blocks}},
	}
}
