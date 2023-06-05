package service

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/MrSwartz/event/internal/config"
	"github.com/MrSwartz/event/pkg/eventservice/service/data"

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

			if _, ok := <-s.buffer.closer; ok {
				logrus.Info("loop stopped")
				s.buffer.isFreed <- struct{}{}
				return
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

			// чтобы не городить кучу конструкций и каналов сделал завершение так
			if _, ok := <-s.buffer.closer; ok {
				logrus.Info("loop stopped")
				s.buffer.isFreed <- struct{}{}
				return
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
	// всё это сделано не очень красиво, но работает. сначала посылаем сигнал закрытия в буфер
	// в цикле по тикеру раз в N секунд отсылаются данные, а сразу после отправки данных
	// проверяем был ли получен сигнал на "закрытие" буфера. в этот момент данных в буфере уже нету
	// и буфер отсылает сигнал о том, что он пустой обратно в эту же функцию и можно спокойно
	// закрывать соединение с бд
	s.buffer.closer <- struct{}{}
	close(s.buffer.closer)

	if _, ok := <-s.buffer.isFreed; ok {
		close(s.buffer.isFreed)
		return s.event.CloseData()
	}

	// сюда код зайти не должен
	panic("data closing error")
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

	initMappers()

	if cnf.Buffer.LoopTimeout != 0 && cnf.Buffer.Size != 0 {
		logrus.Info("buffer initialized")
		go srvc.startLoop(ctx, cnf.Buffer.LoopTimeout)
	}

	return srvc, nil
}

func initMappers() {
	// тут я сотворил грех и захардкодил, что неочень хорошо,
	// но нужно срочно заканчивать сервис
	osType = map[string]data.DeviceOS{
		"unsupported": 1,
		"ios":         2,
		"android":     3,
		"windows":     4,
		"linux":       5,
		"unix":        6,
		"macos":       7,
	}

	osVersion = map[string]data.DeviceOSVersion{
		"13.5.1": 11351,
		"13.5.2": 11352,
		"13.5.3": 11353,
		"4.4.4":  2444,
		"5.0.1":  2501,
		"10.0.1": 21001,
	}

	eventType = map[string]data.EventType{
		"app_start": 1,
		"onPause":   2,
		"onRotate":  3,
		"onCreate":  4,
		"onDestroy": 5,
		"onClose":   6,
		"panic":     7,
	}
}
