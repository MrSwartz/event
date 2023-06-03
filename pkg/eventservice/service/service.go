package service

import (
	"context"
	"event/internal/config"
	"event/pkg/eventservice/service/data"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	osVersion = make(map[string]data.DeviceOSVersion, 0)
	osType    = make(map[string]data.DeviceOS, 0)
	eventType = make(map[string]data.EventType, 0)
)

var _ Event = (*Service)(nil)

type Event interface {
	Insert(ctx context.Context, events []ServiceEventModel) error
	Ping(ctx context.Context) error
	CloseData() error
}

func (s *Service) startLoop(ctx context.Context, timeout int) {
	t := time.NewTicker(time.Duration(timeout) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if !s.buffer.isEmpty() {
				logrus.Info("start inserting data in loop")

				buf := s.buffer.extractAndFlush()

				if err := s.event.Insert(ctx, buf); err != nil {
					logrus.Errorf("can't insert data in loop, error: %v", err)
				}

				logrus.Info("end inserting data in loop")
			} else {
				logrus.Info("buffer is empty, nothing to insert")
			}
		}
	}
}

type Service struct {
	event  data.Events
	buffer buffer
}

func (s *Service) Insert(ctx context.Context, events []ServiceEventModel) error {
	ees := make([]data.DataEventModel, 0, len(events))
	for _, event := range events {
		if e, ok := event.toDataModel(); ok {
			ees = append(ees, *e)
		}
	}

	if len(ees) < s.buffer.maxEventsToBuffer {
		s.buffer.append(ees)
		logrus.Infof("data appended to buffer: %d", len(ees))
		return nil
	}

	logrus.Infof("data insert directly to db")
	return s.event.Insert(ctx, ees)
}

func (s *Service) Ping(ctx context.Context) error {
	return s.event.Ping(ctx)
}

func (s *Service) CloseData() error {
	return s.event.CloseData()
}

func NewService(ctx context.Context, cnf config.Config) (*Service, error) {
	conn, err := data.NewClickHouseDB(ctx, cnf.DataBase)
	if err != nil {
		return nil, err
	}

	if cnf.Buffer.MaxEventsToBuffer == 0 {
		cnf.Buffer.MaxEventsToBuffer = 1000
	}

	if cnf.Buffer.LoopTimeout == 0 {
		cnf.Buffer.LoopTimeout = 10
	}

	srvc := &Service{
		event: conn,

		// буффер будет пересоздаваться
		// но после того, как отработает GC, эта память будет перевыделяться эффективнее
		buffer: newBuffer(cnf.Buffer),
	}

	initMappers(cnf.Mappers)

	go srvc.startLoop(ctx, cnf.Buffer.LoopTimeout)
	return srvc, nil
}

func initMappers(cnf config.Mappers) {
	for k, v := range cnf.OsVersion {
		osVersion[k] = data.DeviceOSVersion(v)
	}
	logrus.Infof("os version mapper initialized")

	for k, v := range cnf.DeviceOs {
		osType[k] = data.DeviceOS(v)
	}
	logrus.Infof("device os mapper initialized")

	for k, v := range cnf.Events {
		eventType[k] = data.EventType(v)
	}
	logrus.Infof("events mapper initialized")
}
