package config

type Meepo struct {
	Daemon    bool        `yaml:"daemon"`
	Mode      string      `yaml:"mode"`
	Pprof     string      `yaml:"pprof"`
	Poof      Poof        `yaml:"poof"`
	Identity  Identity    `yaml:"identity,omitempty"`
	API       API         `yaml:"api,omitempty"`
	Socks5    Socks5      `yaml:"socks5,omitempty"`
	Acl       string      `yaml:"acl,omitempty"`
	Log       Log         `yaml:"log,omitempty"`
	Trackerd  *Trackerd   `yaml:"trackerd,omitempty"`
	Trackerds []*Trackerd `yaml:"trackerds,omitempty"`
	Tracker   *Tracker    `yaml:"tracker,omitempty"`
	Trackers  []*Tracker  `yaml:"trackers,omitempty"`
	Webrtc    Webrtc      `yaml:"webrtc,omitempty"`
	Smux      Smux        `yaml:"smux,omitempty"`
	Kcp       Kcp         `yaml:"kcp,omitempty"`
}
