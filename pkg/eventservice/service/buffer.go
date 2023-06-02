package service

import (
	"event/pkg/eventservice/service/data"
	"sync"

	"github.com/sirupsen/logrus"
)

type buffer struct {
	data []data.DataEventModel
	lock sync.Mutex
}

func newBuffer(size int) *buffer {
	return &buffer{
		data: make([]data.DataEventModel, 0, size),
		lock: sync.Mutex{},
	}
}

func (b *buffer) append(data []data.DataEventModel) {
	b.lock.Lock()
	b.data = append(b.data, data...)
	b.lock.Unlock()
	logrus.Infof("appended %d objects to buffer", len(data))
}

func (b *buffer) get() []data.DataEventModel {
	logrus.Infof("get objects from buffer")
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.data
}

func (b *buffer) flush() {
	b.lock.Lock()
	b.data = b.data[:0]
	b.lock.Unlock()
	logrus.Info("buffer flushed")
}

func (b *buffer) isEmpty() bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return len(b.data) == 0
}
