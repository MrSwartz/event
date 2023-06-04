package service

import (
	"event/internal/config"
	"event/pkg/eventservice/service/data"
	"sync"

	"github.com/sirupsen/logrus"
)

type buffer struct {
	data        []data.DataEventModel
	lock        sync.Mutex
	Size        int
	RetriesLeft int
}

func newBuffer(cnf config.Buffer) *buffer {
	if cnf.Size == 0 && cnf.LoopTimeout == 0 {
		//буфер не создан и данные будут идти напрямую в бд
		return &buffer{
			Size: cnf.Size,
		}
	}
	return &buffer{
		data:        make([]data.DataEventModel, 0, cnf.Size),
		Size:        cnf.Size,
		RetriesLeft: cnf.RetriesLeft,
		lock:        sync.Mutex{},
	}
}

func (b *buffer) append(data []data.DataEventModel) {
	b.lock.Lock()
	b.data = append(b.data, data...)
	b.lock.Unlock()
	logrus.Infof("appended %d objects to buffer", len(data))
}

func (b *buffer) extractAndFlush() []data.DataEventModel {
	logrus.Infof("get objects from buffer")
	b.lock.Lock()
	tmpBuf := b.data
	b.data = b.data[:0]
	b.lock.Unlock()
	logrus.Infof("extracted %d objects", len(tmpBuf))
	logrus.Infof("buffer flushed")
	return tmpBuf
}

func (b *buffer) isEmpty() bool {
	b.lock.Lock()
	length := len(b.data)
	b.lock.Unlock()
	return length == 0
}
