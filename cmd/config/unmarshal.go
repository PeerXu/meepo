package config

import (
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

func joinNamespaceName(namespace, name string) string {
	return namespace + "." + name
}

func WrapKeyYaml(key, text string) string {
	var sb strings.Builder
	ss := strings.Split(text, "\n")
	sb.WriteString(key + ":\n")
	for _, s := range ss {
		if s != "" {
			sb.WriteString("  " + s + "\n")
		}
	}
	return sb.String()
}

func UnmarshalConfig(namespace, name, text string) (interface{}, error) {
	return unmarshalConfig(namespace, name, yaml.NewDecoder(strings.NewReader(text)).Decode)
}

func unmarshalConfig(namespace, name string, unmarshal func(interface{}) error) (interface{}, error) {
	fn, ok := unmarshalConfigFuncs.Load(joinNamespaceName(namespace, name))
	if !ok {
		return nil, UnsupportedError{Namespace: namespace, Name: name}
	}

	return fn.(UnmarshalConfigFunc)(unmarshal)
}

type UnmarshalConfigFunc = func(unmarshal func(interface{}) error) (interface{}, error)

var unmarshalConfigFuncs sync.Map

func RegisterUnmarshalConfigFunc(namespace, name string, fn UnmarshalConfigFunc) {
	unmarshalConfigFuncs.Store(joinNamespaceName(namespace, name), fn)
}
