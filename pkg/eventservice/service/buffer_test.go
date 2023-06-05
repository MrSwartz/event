package service

import (
	"testing"

	"github.com/MrSwartz/event/internal/config"
	"github.com/MrSwartz/event/pkg/eventservice/service/data"

	"github.com/stretchr/testify/require"
)

func TestBufferMethods(t *testing.T) {
	buf := newBuffer(config.Buffer{Size: 10})
	require.NotNil(t, buf.data)

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

	fromBuf := buf.extractAndFlush()
	require.Equal(t, 11, len(fromBuf))

	require.Equal(t, 0, len(buf.data))
	require.Equal(t, 22, cap(buf.data))
}
