package acl

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Entity interface {
	ID() string
	Network() string
	Host() string
	Port() string
	Contains(Entity) bool
	String() string
}

type entity struct {
	id      string
	network string
	host    string
	port    string

	ipnet *net.IPNet
}

func (e *entity) ID() string      { return e.id }
func (e *entity) Network() string { return e.network }
func (e *entity) Host() string    { return e.host }
func (e *entity) Port() string    { return e.port }

func (e *entity) String() string {
	return fmt.Sprintf("<Entity/%s/%s:%s:%s>", e.ID(), e.Network(), e.Host(), e.Port())
}

func (e *entity) Contains(x Entity) bool {
	if e.ID() != "*" {
		if e.ID() != x.ID() {
			return false
		}
	}

	if e.Network() != "*" {
		if e.Network() != x.Network() {
			return false
		}
	}

	if e.Host() != "*" {
		if e.ipnet != nil {
			xip := net.ParseIP(x.Host())

			if !e.ipnet.Contains(xip) {
				return false
			}
		} else {
			if e.Host() != x.Host() {
				return false
			}
		}

	}

	if e.Port() != "*" {
		if e.Port() != x.Port() {
			return false
		}
	}

	return true
}

func NewEntity(id, network, host, port string) Entity {
	var ipnet *net.IPNet
	ss := strings.Split(host, "/")
	switch len(ss) {
	case 1:
	case 2:
		ip := net.ParseIP(ss[0])
		if ip != nil {
			ones, err := strconv.Atoi(ss[1])
			if err == nil {
				ipnet = &net.IPNet{
					IP:   ip,
					Mask: net.CIDRMask(ones, 32),
				}
			}
		}
	}

	return &entity{id, network, host, port, ipnet}
}

/*
 * <id>,<network>,<host>:<port>
 * examples:
 *   "a,tcp,192.168.1.1:8080" => id=a, network=tcp, host=192.168.1.1, port=8080
 *   "a,tcp,*:8080" => id=a, network=tcp, host=*, port=8080
 *   "*,*,127.0.0.1:8080" => id=a, network=*, host=127.0.0.1 port=8080
 *   "a,tcp,*" => id=a, network=tcp, host=*, port=*
 *   "a,*,*" => id=a, network=*, host=*, port=*

 *   "a,*" => id=a, network=*, host=*, port=*
 *   "a,tcp" => id=a, network=tcp, host=*, port=*

 *   "a" => id=a, network=*, host=*, port=*
 *   "*" => id=*, network=*, host=*, port=*
 */
func Parse(s string) (Entity, error) {
	if len(s) == 0 {
		return nil, ErrInvalidEntityFn(s)
	}

	ss := strings.Split(s, ",")
	switch len(ss) {
	case 1:
		return NewEntity(ss[0], ANY_SYM, ANY_SYM, ANY_SYM), nil
	case 2:
		return NewEntity(ss[0], ss[1], ANY_SYM, ANY_SYM), nil
	case 3:
		var host, port string
		var err error

		if ss[2] == ANY_SYM {
			host = ANY_SYM
			port = ANY_SYM
		} else {
			host, port, err = net.SplitHostPort(ss[2])
			if err != nil {
				return nil, err
			}
		}

		return NewEntity(ss[0], ss[1], host, port), nil
	default:
		return nil, ErrInvalidEntityFn(s)
	}
}
