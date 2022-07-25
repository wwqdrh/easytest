package grpctest

import (
	"time"

	"github.com/google/uuid"
	"github.com/jhump/protoreflect/desc"
)

// CallData represents contextualized data available for templating
type CallData struct {
	WorkerID           string // unique worker ID
	RequestNumber      int64  // unique incremented request number for each request
	FullyQualifiedName string // fully-qualified name of the method call
	MethodName         string // shorter call method name
	ServiceName        string // the service name
	InputName          string // name of the input message type
	OutputName         string // name of the output message type
	IsClientStreaming  bool   // whether this call is client streaming
	IsServerStreaming  bool   // whether this call is server streaming
	Timestamp          string // timestamp of the call in RFC3339 format
	TimestampUnix      int64  // timestamp of the call as unix time in seconds
	TimestampUnixMilli int64  // timestamp of the call as unix time in milliseconds
	TimestampUnixNano  int64  // timestamp of the call as unix time in nanoseconds
	UUID               string // generated UUIDv4 for each call
}

// newCallData returns new CallData
func newCallData(
	mtd *desc.MethodDescriptor,
	workerID string, reqNum int64) *CallData {

	now := time.Now()
	newUUID, _ := uuid.NewRandom()

	return &CallData{
		WorkerID:           workerID,
		RequestNumber:      reqNum,
		FullyQualifiedName: mtd.GetFullyQualifiedName(),
		MethodName:         mtd.GetName(),
		ServiceName:        mtd.GetService().GetName(),
		InputName:          mtd.GetInputType().GetName(),
		OutputName:         mtd.GetOutputType().GetName(),
		IsClientStreaming:  mtd.IsClientStreaming(),
		IsServerStreaming:  mtd.IsServerStreaming(),
		Timestamp:          now.Format(time.RFC3339),
		TimestampUnix:      now.Unix(),
		TimestampUnixMilli: now.UnixNano() / 1000000,
		TimestampUnixNano:  now.UnixNano(),
		UUID:               newUUID.String(),
	}
}

// Regenerate generates a new instance of call data from this parent instance
// The dynamic data like timestamps and UUIDs are re-filled
func (td *CallData) Regenerate() *CallData {
	now := time.Now()
	newUUID, _ := uuid.NewRandom()

	return &CallData{
		WorkerID:           td.WorkerID,
		RequestNumber:      td.RequestNumber,
		FullyQualifiedName: td.FullyQualifiedName,
		MethodName:         td.MethodName,
		ServiceName:        td.ServiceName,
		InputName:          td.InputName,
		OutputName:         td.OutputName,
		IsClientStreaming:  td.IsClientStreaming,
		IsServerStreaming:  td.IsServerStreaming,
		Timestamp:          now.Format(time.RFC3339),
		TimestampUnix:      now.Unix(),
		TimestampUnixMilli: now.UnixNano() / 1000000,
		TimestampUnixNano:  now.UnixNano(),
		UUID:               newUUID.String(),
	}
}
