package transport_webrtc

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/snappy"
	"github.com/xtaci/kcp-go/v5"
	"github.com/xtaci/smux"
	"golang.org/x/crypto/pbkdf2"

	"github.com/PeerXu/meepo/pkg/lib/dialer"
	dialer_interface "github.com/PeerXu/meepo/pkg/lib/dialer/interface"
	"github.com/PeerXu/meepo/pkg/lib/logging"
)

func (t *WebrtcTransport) nextChannelID() uint16 {
	return uint16(atomic.AddUint32(&t.currentChannelID, 2) & 0xffff)
}

func (t *WebrtcTransport) readyError() error {
	return t.readyErrVal.Load()
}

type LockableChannel struct {
	Ch   chan Message
	Once sync.Once
}

func (lch *LockableChannel) Close() {
	lch.Once.Do(func() {
		close(lch.Ch)
		lch.Ch = nil
	})
}

func parseResponseSession(s string) string {
	sessU64, _ := strconv.ParseUint(s, 16, 32)
	return fmt.Sprintf("%08x", sessU64-1)
}

func (t *WebrtcTransport) parseMuxStreamLabel(stm *smux.Stream) string {
	return fmt.Sprintf("mux#%d", stm.ID())
}

func (t *WebrtcTransport) parseKcpStreamLabel(stm *smux.Stream) string {
	return fmt.Sprintf("kcp#%d", stm.ID())
}

type packetConn struct {
	rwc        dialer_interface.Conn
	remoteAddr net.Addr
	localAddr  net.Addr
}

func (c *packetConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	n, err = c.rwc.Read(p)
	addr = c.remoteAddr
	return
}

func (c *packetConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	n, err = c.rwc.Write(p)
	return
}

func (c *packetConn) Close() error                       { return c.rwc.Close() }
func (c *packetConn) LocalAddr() net.Addr                { return c.localAddr }
func (c *packetConn) SetDeadline(t time.Time) error      { return nil }
func (c *packetConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *packetConn) SetWriteDeadline(t time.Time) error { return nil }

func (t *WebrtcTransport) upgradeKcpConn(conn io.ReadWriteCloser) (io.ReadWriteCloser, error) {
	logger := t.GetLogger().WithFields(logging.Fields{
		"#method":     "upgradeKcpConn",
		"preset":      t.kcpPreset,
		"crypt":       t.kcpCrypt,
		"mtu":         t.kcpMtu,
		"sndwnd":      t.kcpSndwnd,
		"rcvwnd":      t.kcpRcvwnd,
		"dataShard":   t.kcpDataShard,
		"parityShard": t.kcpParityShard,
	})
	pc := &packetConn{
		rwc:        conn,
		localAddr:  dialer.NewAddr("webrtc", fmt.Sprintf("%016x", t.randSrc.Int63())),
		remoteAddr: dialer.NewAddr("webrtc", fmt.Sprintf("%016x", t.randSrc.Int63())),
	}

	pass := pbkdf2.Key([]byte(t.kcpKey), []byte(SALT), 4096, 32, sha1.New)
	var block kcp.BlockCrypt
	switch t.kcpCrypt {
	case "null":
		block = nil
	case "sm4":
		block, _ = kcp.NewSM4BlockCrypt(pass[:16])
	case "tea":
		block, _ = kcp.NewTEABlockCrypt(pass[:16])
	case "xor":
		block, _ = kcp.NewSimpleXORBlockCrypt(pass)
	case "none":
		block, _ = kcp.NewNoneBlockCrypt(pass)
	case "aes-128":
		block, _ = kcp.NewAESBlockCrypt(pass[:16])
	case "aes-192":
		block, _ = kcp.NewAESBlockCrypt(pass[:24])
	case "blowfish":
		block, _ = kcp.NewBlowfishBlockCrypt(pass)
	case "twofish":
		block, _ = kcp.NewTwofishBlockCrypt(pass)
	case "cast5":
		block, _ = kcp.NewCast5BlockCrypt(pass[:16])
	case "3des":
		block, _ = kcp.NewTripleDESBlockCrypt(pass[:24])
	case "xtea":
		block, _ = kcp.NewXTEABlockCrypt(pass[:16])
	case "salsa20":
		block, _ = kcp.NewSalsa20BlockCrypt(pass)
	default:
		block, _ = kcp.NewAESBlockCrypt(pass)
	}

	upgradedConn, err := kcp.NewConn3(0, dialer.NewAddr("kcp", ""), block, t.kcpDataShard, t.kcpParityShard, pc)
	if err != nil {
		logger.WithError(err).Debugf("failed to upgrade kcp conn")
		return nil, err
	}

	upgradedConn.SetStreamMode(true)
	upgradedConn.SetWriteDelay(false)
	upgradedConn.SetMtu(t.kcpMtu)
	upgradedConn.SetWindowSize(t.kcpSndwnd, t.kcpRcvwnd)
	upgradedConn.SetACKNoDelay(true)

	switch t.kcpPreset {
	case "normal":
		upgradedConn.SetNoDelay(0, 40, 2, 1)
	case "fast":
		upgradedConn.SetNoDelay(0, 30, 2, 1)
	case "fast2":
		upgradedConn.SetNoDelay(1, 20, 2, 1)
	case "fast3":
		upgradedConn.SetNoDelay(1, 10, 2, 1)
	}

	logger.Tracef("upgrade kcp conn")

	return upgradedConn, nil
}

func (t *WebrtcTransport) getSmuxConfig() *smux.Config {
	cfg := smux.DefaultConfig()
	cfg.Version = t.muxVer
	cfg.MaxReceiveBuffer = t.muxBuf
	cfg.MaxStreamBuffer = t.muxStreamBuf
	cfg.KeepAliveDisabled = true
	return cfg
}

type compStream struct {
	io.Reader
	io.Writer
	io.Closer
}

func NewCompStream(conn io.ReadWriteCloser) io.ReadWriteCloser {
	return &compStream{
		Reader: snappy.NewReader(conn),
		Writer: snappy.NewWriter(conn), // nolint:staticcheck
		Closer: conn,
	}
}
