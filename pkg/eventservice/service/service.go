package service

import (
	"context"
	"event/internal/config"
	"event/pkg/eventservice/service/data"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	osVersion = make(map[string]data.DeviceOSVersion, 0)
	osType    = make(map[string]data.DeviceOS, 0)
	eventType = make(map[string]data.EventType, 0)
)

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

type Event interface {
	Insert(ctx context.Context, events []data.DataEventModel) error
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
			event := data.DataEventModel{ClientTime: time.Now().UTC(), ServerTime: time.Now().UTC(), DeviceId: "0287D9AA-4ADF-4B37-A60F-3E9E645C821E", Session: "dfb", ParamStr: "fgg", Ip: 1234567890, Sequence: 1, ParamInt: 1234, DeviceOs: 0, DeviceOsVersion: 1, Event: 10}
			buf := []data.DataEventModel{}
			for i := 0; i < 1000; i++ {
				buf = append(buf, event)
			}
			s.buffer.append(buf)
			if !s.buffer.isEmpty() {
				logrus.Info("start inserting data in loop")

				buf := s.buffer.get()
				s.buffer.flush()

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

	if len(ees) < 1000 {
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

	srvc := &Service{
		event:  conn,
		buffer: *newBuffer(cnf.Buffer.Size),
	}

	initMappers(cnf.Mappers)
	fmt.Println(osVersion, osType, eventType)

	go srvc.startLoop(ctx, cnf.Buffer.LoopTimeout)
	return srvc, nil
}
