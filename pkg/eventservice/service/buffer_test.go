package service

import (
	"context"
	"event/pkg/eventservice/service/data"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBufferMethods(t *testing.T) {
	buf := newBuffer(10)
	require.NotNil(t, buf)

	require.Equal(t, 0, len(buf.data))
	require.Equal(t, 10, cap(buf.data))

	empty := buf.isEmpty()
	require.True(t, empty)

	event := data.DataEventModel{ParamInt: 1}
	buf.append([]data.DataEventModel{event})

	empty = buf.isEmpty()
	require.False(t, empty)
	require.Equal(t, 1, len(buf.data))

	events := make([]data.DataEventModel, 0, 10)
	for i := 0; i < 10; i++ {
		events = append(events, event)
	}

	buf.append(events)
	require.Equal(t, 11, len(buf.data))
	require.Equal(t, 22, cap(buf.data))

	fromBuf := buf.get()
	require.Equal(t, 11, len(fromBuf))

	buf.flush()
	require.Equal(t, 0, len(buf.data))
	require.Equal(t, 22, cap(buf.data))
}

func TestLogicLopp(t *testing.T) {
	events := []data.DataEventModel{
		{
			ClientTime:      time.Now(),
			ServerTime:      time.Now(),
			DeviceId:        "12121212121212121212",
			Session:         "dclvvf",
			ParamStr:        "123",
			Ip:              1234,
			Sequence:        1,
			ParamInt:        1,
			DeviceOs:        1,
			DeviceOsVersion: 1,
			Event:           1,
		},
	}

	ctx := context.Background()
	ticker := time.NewTicker(time.Duration(1) * time.Second)
	s := &Service{
		buffer: *newBuffer(10),
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.buffer.append(events)
			if !s.buffer.isEmpty() {
				_ = s.buffer.get()
				s.buffer.flush()
			}
		}
	}
}

func TestBufTickerCase(t *testing.T) {
	events := []data.DataEventModel{
		{
			ClientTime:      time.Now(),
			ServerTime:      time.Now(),
			DeviceId:        "12121212121212121212",
			Session:         "dclvvf",
			ParamStr:        "123",
			Ip:              1234,
			Sequence:        1,
			ParamInt:        1,
			DeviceOs:        1,
			DeviceOsVersion: 1,
			Event:           1,
		},
	}

	s := &Service{
		buffer: *newBuffer(10),
	}

	s.buffer.append(events)
	if !s.buffer.isEmpty() {
		_ = s.buffer.get()
		s.buffer.flush()
	}
	s.buffer.append(events)
}
