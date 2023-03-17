package acl

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	Allow  string   `yaml:"allow,omitempty"`
	Block  string   `yaml:"block,omitempty"`
	Allows []string `yaml:"allows,omitempty"`
	Blocks []string `yaml:"blocks,omitempty"`
}

// - blocks:
//   - "*,*,10.1.1.0/24:22"
//
// - allow: "*,tcp,127.0.0.1:80"
// - allows:
//   - "a,tcp,127.0.0.1:22"
//
// - block: "*"
//
// OR
//
// #allow=*,*,10.1.1.0/24:22;#block=*

func ParseRules(s string) (rs []Rule, err error) {
	trimed := strings.Trim(s, " ")
	if strings.HasPrefix(trimed, "#") {
		return parseOneLineRules(trimed)
	}

	if err = yaml.Unmarshal([]byte(s), &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func parseOneLineRules(s string) (rs []Rule, err error) {
	rulesStrSlice := strings.Split(s, ";")
	for _, ruleStr := range rulesStrSlice {
		ruleStr = strings.Trim(ruleStr, " ")

		if len(ruleStr) == 0 {
			continue
		}

		if !strings.HasPrefix(ruleStr, "#") {
			return nil, ErrInvalidRuleFn(ruleStr)
		}
		splitRule := strings.SplitN(ruleStr, "=", 2)
		if len(splitRule) != 2 {
			return nil, ErrInvalidRuleFn(ruleStr)
		}
		op := splitRule[0]
		ent := splitRule[1]
		switch op {
		case "#allow":
			rs = append(rs, Rule{Allow: ent})
		case "#block":
			rs = append(rs, Rule{Block: ent})
		default:
			return nil, ErrInvalidRuleFn(ruleStr)
		}
	}
	return
}
