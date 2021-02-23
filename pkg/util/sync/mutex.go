package sync

import (
	gsync "sync"

	"github.com/sirupsen/logrus"
)

type Mutex struct {
	name   string
	logger logrus.FieldLogger
	m      gsync.Mutex
}

func (m *Mutex) getLogger() logrus.FieldLogger {
	return m.logger.WithFields(logrus.Fields{
		"#instance": "Mutex",
		"name":      m.name,
	})
}

func (m *Mutex) Lock() {
	logger := m.getLogger()
	logger.Debugf("lock")
	m.m.Lock()
}

func (m *Mutex) Unlock() {
	logger := m.getLogger()
	m.m.Unlock()
	logger.Debugf("unlock")
}

func NewMutex(name string, logger logrus.FieldLogger) *Mutex {
	return &Mutex{
		name:   name,
		logger: logger,
	}
}
