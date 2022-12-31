package acl

import "gopkg.in/yaml.v3"

type Rule struct {
	Allow  string   `yaml:"allow,omitempty"`
	Block  string   `yaml:"block,omitempty"`
	Allows []string `yaml:"allows,omitempty"`
	Blocks []string `yaml:"blocks,omitempty"`
}

// - blocks:
//   - "*,*,10.1.1.0/24:22"
// - allow: "*,tcp,127.0.0.1:80"
// - allows:
//   - "a,tcp,127.0.0.1:22"
// - block: "*"
func ParseRules(s string) (rs []Rule, err error) {
	if err = yaml.Unmarshal([]byte(s), &rs); err != nil {
		return nil, err
	}
	return rs, nil
}
