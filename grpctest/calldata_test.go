package grpctest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallData_New(t *testing.T) {
	md, err := GetMethodDescFromProto("helloworld.Greeter/SayHello", "./testdata/greeter.proto", []string{})
	require.NoError(t, err)
	require.NotNil(t, md)

	ctd := newCallData(md, "worker_id_123", 100)

	assert.NotNil(t, ctd)
	assert.Equal(t, "worker_id_123", ctd.WorkerID)
	assert.Equal(t, int64(100), ctd.RequestNumber)
	assert.Equal(t, "helloworld.Greeter.SayHello", ctd.FullyQualifiedName)
	assert.Equal(t, "SayHello", ctd.MethodName)
	assert.Equal(t, "Greeter", ctd.ServiceName)
	assert.Equal(t, "HelloRequest", ctd.InputName)
	assert.Equal(t, "HelloReply", ctd.OutputName)
	assert.Equal(t, false, ctd.IsClientStreaming)
	assert.Equal(t, false, ctd.IsServerStreaming)
	assert.NotEmpty(t, ctd.Timestamp)
	assert.NotZero(t, ctd.TimestampUnix)
	assert.NotZero(t, ctd.TimestampUnixMilli)
	assert.NotZero(t, ctd.TimestampUnixNano)
	assert.Equal(t, ctd.TimestampUnix, ctd.TimestampUnixMilli/1000)
	assert.NotEmpty(t, ctd.UUID)
	assert.Equal(t, 36, len(ctd.UUID))
}
