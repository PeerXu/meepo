package packet_test

import (
	"testing"

	"github.com/PeerXu/meepo/pkg/meepo/packet"
	"github.com/stretchr/testify/suite"
)

type HeaderTestSuite struct {
	suite.Suite
}

func (s *HeaderTestSuite) TestMarshalAndUnmarshal() {
	h := packet.NewHeader(1, "a", "b", packet.Request, "test")
	h = h.SetSignature([]byte("abc"))
	buf, err := packet.MarshalHeader(h)
	s.Require().Nil(err)

	h1, err := packet.UnmarshalHeader(buf)
	s.Require().Nil(err)

	s.Equal(h.Session(), h1.Session())
	s.Equal(h.Source(), h1.Source())
	s.Equal(h.Destination(), h1.Destination())
	s.Equal(h.Type(), h1.Type())
	s.Equal(h.Method(), h1.Method())
	s.Equal(h.Signature(), h1.Signature())
}

func (s *HeaderTestSuite) TestInvertHeader() {
	h := packet.NewHeader(1, "a", "b", packet.Request, "test")
	h = h.SetSignature([]byte("abc"))
	ih := packet.InvertHeader(h)

	s.Equal(h.Session(), ih.Session())
	s.Equal(h.Destination(), ih.Source())
	s.Equal(h.Source(), ih.Destination())
	s.Equal(packet.Response, ih.Type())
	s.Equal(h.Method(), ih.Method())
	s.Nil(ih.Signature())
}

func (s *HeaderTestSuite) TestSetAndUnsetSignature() {
	h := packet.NewHeader(1, "a", "b", packet.Request, "test")
	hs := h.SetSignature([]byte("ttt"))
	hus := h.UnsetSignature()
	s.Nil(h.Signature())
	s.Equal([]byte("ttt"), hs.Signature())
	hb, err := packet.MarshalHeader(h)
	s.Require().Nil(err)
	husb, err := packet.MarshalHeader(hus)
	s.Require().Nil(err)
	s.Equal(hb, husb)
}

func TestHeaderTestSuite(t *testing.T) {
	suite.Run(t, new(HeaderTestSuite))
}
