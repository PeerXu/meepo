package meepo_test

import (
	"net"
	"testing"

	"github.com/PeerXu/meepo/pkg/meepo"
	"github.com/stretchr/testify/suite"
)

type AclTestSuite struct {
	suite.Suite
}

func (s *AclTestSuite) TestAddrContains() {
	for _, t := range []struct {
		a net.Addr
		b net.Addr
	}{
		{meepo.ACL_ANY_ADDR, meepo.NewAclAddr("socks5", "127.0.0.1:22")},
		{meepo.ACL_ANY_ADDR, meepo.NewAclAddr("tcp", "127.0.0.1:8080")},
		{meepo.NewAclAddr("tcp", "192.168.1.0/24:*"), meepo.NewAclAddr("tcp", "192.168.1.1:8080")},
	} {
		s.True(meepo.AddrContains(t.a, t.b), "%s:%s %s:%s", t.a.Network(), t.a.String(), t.b.Network(), t.b.String())
	}

	for _, t := range []struct {
		a net.Addr
		b net.Addr
	}{
		{meepo.NewAclAddr(meepo.ACL_ANY, "192.168.2.1:22"), meepo.NewAclAddr("tcp", "192.168.2.2:22")},
		{meepo.NewAclAddr(meepo.ACL_ANY, "192.168.2.1:22"), meepo.NewAclAddr("tcp", "192.168.2.1:44")},
		{meepo.NewAclAddr(meepo.ACL_ANY, "192.168.1.0/24:*"), meepo.NewAclAddr("tcp", "192.168.2.1:22")},
	} {
		s.False(meepo.AddrContains(t.a, t.b), "%s:%s %s:%s", t.a.Network(), t.a.String(), t.b.Network(), t.b.String())
	}
}

func (s *AclTestSuite) TestParseAclPolicy() {
	for _, t := range []struct {
		s string
		p meepo.AclPolicy
	}{
		{"*", meepo.NewAclPolicy(meepo.ACL_ANY_ENTITY, meepo.ACL_ANY_ENTITY)},
		{"192.168.1.1:22", meepo.NewAclPolicy(meepo.ACL_ANY_ENTITY, meepo.NewAclEntity(meepo.ACL_ANY, meepo.NewAclAddr(meepo.ACL_ANY, "192.168.1.1:22")))},
		{"10.0.0.0/24:*", meepo.NewAclPolicy(meepo.ACL_ANY_ENTITY, meepo.NewAclEntity(meepo.ACL_ANY, meepo.NewAclAddr(meepo.ACL_ANY, "10.0.0.0/24:*")))},
		{"tcp:10.0.0.1:22", meepo.NewAclPolicy(meepo.ACL_ANY_ENTITY, meepo.NewAclEntity(meepo.ACL_ANY, meepo.NewAclAddr("tcp", "10.0.0.1:22")))},
		{"a:*:10.0.0.1:22", meepo.NewAclPolicy(meepo.ACL_ANY_ENTITY, meepo.NewAclEntity("a", meepo.NewAclAddr(meepo.ACL_ANY, "10.0.0.1:22")))},
		{"socks5:*:*,127.0.0.1:80", meepo.NewAclPolicy(meepo.NewAclEntity(meepo.ACL_ANY, meepo.NewAclAddr("socks5", "*:*")), meepo.NewAclEntity(meepo.ACL_ANY, meepo.NewAclAddr(meepo.ACL_ANY, "127.0.0.1:80")))},
	} {
		p, err := meepo.ParseAclPolicy(t.s)
		s.Require().Nil(err)
		s.Equal(t.p.Source().ID(), p.Source().ID(), "%s %s", s, p)
		s.Equal(t.p.Source().Addr().Network(), p.Source().Addr().Network(), "%s %s", s, p)
		s.Equal(t.p.Source().Addr().String(), p.Source().Addr().String(), "%s %s", s, p)
		s.Equal(t.p.Destination().ID(), p.Destination().ID(), "%s %s", s, p)
		s.Equal(t.p.Destination().Addr().Network(), p.Destination().Addr().Network(), "%s %s", s, p)
		s.Equal(t.p.Destination().Addr().String(), p.Destination().Addr().String(), "%s %s", s, p)
	}
}

func (ts *AclTestSuite) TestAclAllowed() {
	for _, t := range []struct {
		ap []string
		bp []string
		ac []meepo.AclChallenge
		bc []meepo.AclChallenge
	}{
		{
			[]string{"*"},
			nil,
			[]meepo.AclChallenge{
				meepo.NewAclChallenge(meepo.ACL_ANY_ENTITY, meepo.NewAclEntity("a", meepo.NewAclAddr("tcp", "127.0.0.1:22"))),
			},
			nil,
		},
		{
			[]string{"127.0.0.1:22"},
			nil,
			[]meepo.AclChallenge{
				meepo.NewAclChallenge(meepo.ACL_ANY_ENTITY, meepo.NewAclEntity("a", meepo.NewAclAddr("tcp", "127.0.0.1:22"))),
			},
			[]meepo.AclChallenge{
				meepo.NewAclChallenge(meepo.ACL_ANY_ENTITY, meepo.NewAclEntity("a", meepo.NewAclAddr("tcp", "127.0.0.1:80"))),
			},
		},
		{
			[]string{"*"},
			[]string{"127.0.0.1:*"},
			[]meepo.AclChallenge{
				meepo.NewAclChallenge(meepo.ACL_ANY_ENTITY, meepo.NewAclEntity("a", meepo.NewAclAddr("tcp", "10.1.1.1:22"))),
			},
			[]meepo.AclChallenge{
				meepo.NewAclChallenge(meepo.ACL_ANY_ENTITY, meepo.NewAclEntity("a", meepo.NewAclAddr("tcp", "127.0.0.1:22"))),
				meepo.NewAclChallenge(meepo.ACL_ANY_ENTITY, meepo.NewAclEntity("a", meepo.NewAclAddr("tcp", "127.0.0.1:80"))),
			},
		},
		{
			[]string{"*:socks5:*:*,*"},
			nil,
			[]meepo.AclChallenge{
				meepo.NewAclChallenge(meepo.NewAclEntity("a", meepo.NewAclAddr("socks5", "0.0.0.0:0")), meepo.NewAclEntity("b", meepo.NewAclAddr("tcp", "127.0.0.1:22"))),
			},
			[]meepo.AclChallenge{
				meepo.NewAclChallenge(meepo.NewAclEntity("a", meepo.NewAclAddr("tcp", "127.0.0.1:2222")), meepo.NewAclEntity("b", meepo.NewAclAddr("tcp", "127.0.0.1:22"))),
			},
		},
	} {
		var allows []meepo.AclPolicy
		var blocks []meepo.AclPolicy
		for _, s := range t.ap {
			p, err := meepo.ParseAclPolicy(s)
			ts.Require().Nil(err)
			allows = append(allows, p)
		}
		for _, s := range t.bp {
			p, err := meepo.ParseAclPolicy(s)
			ts.Require().Nil(err)
			blocks = append(blocks, p)
		}
		acl := meepo.NewAcl(meepo.WithAllowPolicies(allows), meepo.WithBlockPolicies(blocks))
		for _, c := range t.ac {
			err := acl.Allowed(c)
			ts.Nil(err)
		}
		for _, c := range t.bc {
			err := acl.Allowed(c)
			ts.Equal(meepo.ErrAclNotAllowed, err)
		}
	}
}

func TestAclTestSuite(t *testing.T) {
	suite.Run(t, new(AclTestSuite))
}
