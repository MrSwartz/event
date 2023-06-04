package service

import (
	"context"
	"event/internal/config"
	"event/pkg/eventservice/service/data"
	"fmt"
	"os"
	"runtime/debug"
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
			if s.buffer.RetriesLeft != 0 {
				// даю пару попыток справиться с пересылкой сохранённых значений в бд
				// если разрешить бесконечно заполнять буфер, то случится Exponential Backoff
				// и исчерпав всю доступную память сервис упадёт сам с потерей данных
				s.buffer.RetriesLeft--

				if s.buffer.isEmpty() {
					// если код зайдёт в эту ветку это будет очень странно
					// так как происходит попытка записать(передать в бд) пустой массив
					// который по странной причине не успевает выполниться за установленное время
					panic("critical error in buffer\n" + string(debug.Stack()))
				} else {
					buf := s.buffer.extractAndFlush()

					if err := s.event.Insert(ctx, buf); err != nil {
						logrus.Errorf("can't insert data in loop, error: %v", err)
					}
				}
			} else {
				// тут важно знать конфигурацию сервиса и ос, если разрешён swapping, то дисковая память
				// может быть исчерпана и сохранить буфер на диске не получится
				file, err := os.Create("dump_" + time.Now().UTC().String() + ".txt")
				if err != nil {
					logrus.Error("can't save buffer dump")
					os.Exit(1)
				}
				// такое себе решение, но для малых значений RetriesLeft подойдёт
				// по условию сервис дожен выдерживать 200 RPS со средним запросом в 30 событий,
				// каждое из которых занимает примерно 120 байт, таймаут выставлен в 10 секунд
				// 200(RPS) * 10(ticker) * 30(events) * 120(byte) * 3(RetriesLeft) = 20.59 MB
				fmt.Fprintf(file, "%v", s.buffer.extractAndFlush())
				file.Close()
			}
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
	buffer *buffer
}

func (s *Service) Insert(ctx context.Context, events []ServiceEventModel) error {
	ees := make([]data.DataEventModel, 0, len(events))
	for _, event := range events {
		if e, ok := event.toDataModel(); ok {
			ees = append(ees, *e)
		}
	}

	if len(ees) < s.buffer.Size {
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
		event: conn,

		// буффер будет пересоздаваться
		// но после того, как отработает GC, эта память будет перевыделяться эффективнее
		buffer: newBuffer(cnf.Buffer),
	}

	initMappers(cnf.Mappers)

	if cnf.Buffer.LoopTimeout != 0 && cnf.Buffer.Size != 0 {
		logrus.Info("buffer initialized")
		go srvc.startLoop(ctx, cnf.Buffer.LoopTimeout)
	}

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
