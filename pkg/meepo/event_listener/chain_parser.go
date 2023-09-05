package meepo_event_listener

import "strings"

type ChainParser struct {
	delimiter string
}

var (
	DefaultChainParser = &ChainParser{delimiter: DELIMITER}
)

func (p *ChainParser) Parse(name string) Chain {
	if !strings.HasSuffix(name, p.delimiter) {
		name = name + p.delimiter
	}

	return strings.Split(name, p.delimiter)
}
