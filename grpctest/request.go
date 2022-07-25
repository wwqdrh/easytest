package grpctest

import (
	"context"
	"math"
	"time"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"errors"
	"fmt"
	"io"

	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc/metadata"
)

type RunConfig struct {
	Host       string
	Proto      string
	ImportPath []string
	Handle     string
	Expect     []string

	defaultCallOptions []grpc.CallOption

	// security settings
	creds credentials.TransportCredentials

	insecure  bool
	authority string

	timeout       time.Duration
	dialTimeout   time.Duration
	keepaliveTime time.Duration

	streamInterval     time.Duration
	streamCallDuration time.Duration
	streamCallCount    uint
}

// Requester is used for doing the requests
type Requester struct {
	conn *grpc.ClientConn
	stub grpcdynamic.Stub

	mtd *desc.MethodDescriptor

	config *RunConfig
}

func NewRequester(c *RunConfig) (*Requester, error) {
	mtd, err := GetMethodDescFromProto(c.Handle, c.Proto, c.ImportPath)
	if err != nil {
		return nil, err
	}

	req := &Requester{
		config: c,
		mtd:    mtd,
	}
	// 构建conn
	conn, err := req.newClientConn()
	if err != nil {
		return nil, err
	}
	req.conn = conn
	req.stub = grpcdynamic.NewStub(conn)

	return req, nil
}

func (b *Requester) newClientConn() (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	if b.config.insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(b.config.creds))
	}

	if b.config.authority != "" {
		opts = append(opts, grpc.WithAuthority(b.config.authority))
	}

	if len(b.config.defaultCallOptions) > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(b.config.defaultCallOptions...))
	} else {
		// increase max receive and send message sizes
		opts = append(opts,
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(math.MaxInt32),
				grpc.MaxCallSendMsgSize(math.MaxInt32),
			))

	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, b.config.dialTimeout)
	defer cancel()
	if b.config.keepaliveTime > 0 {
		opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    b.config.keepaliveTime,
			Timeout: b.config.keepaliveTime,
		}))
	}

	return grpc.DialContext(ctx, b.config.Host, opts...)
}

// 关闭grpc connection
func (b *Requester) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()
	shutdownCh := connectionOnState(ctx, b.conn, connectivity.Shutdown)

	err := b.conn.Close()

	<-shutdownCh

	return err
}

func checkState(conn *grpc.ClientConn, states ...connectivity.State) bool {
	currentState := conn.GetState()
	for _, s := range states {
		if currentState == s {
			return true
		}
	}

	return false
}

func connectionOnState(ctx context.Context, conn *grpc.ClientConn, states ...connectivity.State) <-chan bool {
	stateCh := make(chan bool)
	go func() {
		defer close(stateCh)
		if checkState(conn, states...) {
			stateCh <- true
			return
		}

		for {
			change := conn.WaitForStateChange(ctx, conn.GetState())
			if !change {
				stateCh <- checkState(conn, states...)
				return
			}

			if checkState(conn, states...) {
				stateCh <- true
				return
			}
		}
	}()

	return stateCh
}

////////////////////
// grpc call worker
////////////////////

// ErrEndStream is a signal from message providers that worker should close the stream
// It should not be used for erronous states
var ErrEndStream = errors.New("ending stream")

// ErrLastMessage is a signal from message providers that the returned payload is the last one of the stream
// This is optional but encouraged for optimized performance
// Message payload returned along with this error must be valid and may not be nil
var ErrLastMessage = errors.New("last message")

type StreamRecvMsgInterceptFunc func(*dynamic.Message, error) error

type StreamMessageProviderFunc func(*CallData) (*dynamic.Message, error)

// Worker is used for doing a single stream of requests in parallel
type Worker struct {
	stub grpcdynamic.Stub
	mtd  *desc.MethodDescriptor

	config   *RunConfig
	workerID string
	active   bool
	stopCh   chan bool

	msgProvider StreamMessageProviderFunc

	streamRecv StreamRecvMsgInterceptFunc
}

func (w *Worker) runWorker() error {
	return w.makeRequest(nil, nil)
}

// Stop stops the worker. It has to be started with Run() again.
func (w *Worker) Stop() {
	if !w.active {
		return
	}

	w.active = false
	w.stopCh <- true
}

func (w *Worker) makeRequest(reqMD *metadata.MD, inputs []*dynamic.Message) error {
	ctd := newCallData(w.mtd, w.workerID, 0)

	// if w.config.enableCompression {
	// 	reqMD.Append("grpc-accept-encoding", gzip.Name)
	// }

	ctx := context.Background()
	var cancel context.CancelFunc

	if w.config.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, w.config.timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	// include the metadata
	if reqMD != nil {
		ctx = metadata.NewOutgoingContext(ctx, *reqMD)
	}

	var msgProvider StreamMessageProviderFunc
	if w.msgProvider != nil {
		msgProvider = w.msgProvider
	} else if w.mtd.IsClientStreaming() {
		msgProvider = func(cd *CallData) (*dynamic.Message, error) { return nil, nil }
	}

	if len(inputs) == 0 && msgProvider == nil {
		return fmt.Errorf("no data provided for request")
	}

	// RPC errors are handled via stats handler
	if w.mtd.IsClientStreaming() && w.mtd.IsServerStreaming() {
		_ = w.makeBidiRequest(&ctx, ctd, msgProvider)
	} else if w.mtd.IsClientStreaming() {
		_ = w.makeClientStreamingRequest(&ctx, ctd, msgProvider)
	} else if w.mtd.IsServerStreaming() {
		_ = w.makeServerStreamingRequest(&ctx, inputs[0])
	} else {
		_ = w.makeUnaryRequest(&ctx, reqMD, inputs[0])
	}

	return nil
}

func (w *Worker) makeUnaryRequest(ctx *context.Context, reqMD *metadata.MD, input *dynamic.Message) error {
	var resErr error
	var callOptions = []grpc.CallOption{}

	_, resErr = w.stub.InvokeRpc(*ctx, w.mtd, input, callOptions...)

	return resErr
}

func (w *Worker) makeClientStreamingRequest(ctx *context.Context, ctd *CallData, messageProvider StreamMessageProviderFunc) error {
	var err error
	var str *grpcdynamic.ClientStream
	// var callOptions = []grpc.CallOption{}

	closeStream := func() {
		if _, err := str.CloseAndReceive(); err != nil {
			fmt.Println(err.Error())
		}
	}

	performSend := func(payload *dynamic.Message) (bool, error) {
		err := str.SendMsg(payload)

		if err == io.EOF {
			return true, nil
		}

		return false, err
	}

	doneCh := make(chan struct{})
	cancel := make(chan struct{}, 1)
	if w.config.streamCallDuration > 0 {
		go func() {
			sct := time.NewTimer(w.config.streamCallDuration)
			select {
			case <-sct.C:
				cancel <- struct{}{}
				return
			case <-doneCh:
				if !sct.Stop() {
					<-sct.C
				}
				return
			}
		}()
	}

	done := false
	counter := uint(0)
	end := false
	for !done && len(cancel) == 0 {
		// default message provider checks counter
		// but we also need to keep our own counts
		// in case of custom client providers

		var payload *dynamic.Message
		payload, err = messageProvider(ctd)

		isLast := false
		if errors.Is(err, ErrLastMessage) {
			isLast = true
			err = nil
		}

		if err != nil {
			if errors.Is(err, ErrEndStream) {
				err = nil
			}
			break
		}

		end, err = performSend(payload)
		if end || err != nil || isLast || len(cancel) > 0 {
			break
		}

		counter++

		if w.config.streamCallCount > 0 && counter >= w.config.streamCallCount {
			break
		}

		if w.config.streamInterval > 0 {
			wait := time.NewTimer(w.config.streamInterval)
			select {
			case <-wait.C:
				break
			case <-cancel:
				if !wait.Stop() {
					<-wait.C
				}
				done = true
				break
			}
		}
	}

	for len(cancel) > 0 {
		<-cancel
	}

	closeStream()

	close(doneCh)
	close(cancel)

	return err
}

func (w *Worker) makeServerStreamingRequest(ctx *context.Context, input *dynamic.Message) error {
	var callOptions = []grpc.CallOption{}

	callCtx, callCancel := context.WithCancel(*ctx)
	defer callCancel()

	str, err := w.stub.InvokeRpcServerStream(callCtx, w.mtd, input, callOptions...)

	if err != nil {
		return err
	}

	doneCh := make(chan struct{})
	cancel := make(chan struct{}, 1)
	if w.config.streamCallDuration > 0 {
		go func() {
			sct := time.NewTimer(w.config.streamCallDuration)
			select {
			case <-sct.C:
				cancel <- struct{}{}
				return
			case <-doneCh:
				if !sct.Stop() {
					<-sct.C
				}
				return
			}
		}()
	}

	interceptCanceled := false
	counter := uint(0)
	for err == nil {
		// we should check before receiving a message too
		if w.config.streamCallDuration > 0 && len(cancel) > 0 {
			<-cancel
			callCancel()
			break
		}

		var res proto.Message
		res, err = str.RecvMsg()

		// with any of the cancellation operations we can't just bail
		// we have to drain the messages until the server gets the cancel and ends their side of the stream

		if w.streamRecv != nil {
			if converted, ok := res.(*dynamic.Message); ok {
				err = w.streamRecv(converted, err)
				if errors.Is(err, ErrEndStream) && !interceptCanceled {
					interceptCanceled = true
					err = nil

					callCancel()
				}
			}
		}

		if err != nil {
			if err == io.EOF {
				err = nil
			}

			break
		}

		counter++

		if w.config.streamCallCount > 0 && counter >= w.config.streamCallCount {
			callCancel()
		}

		if w.config.streamCallDuration > 0 && len(cancel) > 0 {
			<-cancel
			callCancel()
		}
	}

	close(doneCh)
	close(cancel)

	return err
}

func (w *Worker) makeBidiRequest(ctx *context.Context,
	ctd *CallData, messageProvider StreamMessageProviderFunc) error {

	var callOptions = []grpc.CallOption{}

	str, err := w.stub.InvokeRpcBidiStream(*ctx, w.mtd, callOptions...)

	if err != nil {
		return err
	}

	counter := uint(0)
	indexCounter := 0
	recvDone := make(chan bool)
	sendDone := make(chan bool)

	closeStream := func() {
		if err := str.CloseSend(); err != nil {
			fmt.Println(err.Error())
		}
	}

	doneCh := make(chan struct{})
	cancel := make(chan struct{}, 1)
	if w.config.streamCallDuration > 0 {
		go func() {
			sct := time.NewTimer(w.config.streamCallDuration)
			select {
			case <-sct.C:
				cancel <- struct{}{}
				return
			case <-doneCh:
				if !sct.Stop() {
					<-sct.C
				}
				return
			}
		}()
	}

	var recvErr error

	go func() {
		interceptCanceled := false

		for recvErr == nil {
			var res proto.Message
			res, recvErr = str.RecvMsg()
			if w.streamRecv != nil {
				if converted, ok := res.(*dynamic.Message); ok {
					iErr := w.streamRecv(converted, recvErr)
					if errors.Is(iErr, ErrEndStream) && !interceptCanceled {
						interceptCanceled = true
						if len(cancel) == 0 {
							cancel <- struct{}{}
						}
						recvErr = nil
					}
				}
			}

			if recvErr != nil {
				close(recvDone)
				break
			}
		}
	}()

	go func() {
		done := false

		for err == nil && !done {

			// check at start before send too
			if len(cancel) > 0 {
				<-cancel
				closeStream()
				break
			}

			// default message provider checks counter
			// but we also need to keep our own counts
			// in case of custom client providers

			var payload *dynamic.Message
			payload, err = messageProvider(ctd)

			isLast := false
			if errors.Is(err, ErrLastMessage) {
				isLast = true
				err = nil
			}

			if err != nil {
				if errors.Is(err, ErrEndStream) {
					err = nil
				}

				closeStream()
				break
			}

			err = str.SendMsg(payload)
			if err != nil {
				if err == io.EOF {
					err = nil
				}

				break
			}

			if isLast {
				closeStream()
				break
			}

			counter++
			indexCounter++

			if w.config.streamCallCount > 0 && counter >= w.config.streamCallCount {
				closeStream()
				break
			}

			if len(cancel) > 0 {
				<-cancel
				closeStream()
				break
			}

			if w.config.streamInterval > 0 {
				wait := time.NewTimer(w.config.streamInterval)
				select {
				case <-wait.C:
					break
				case <-cancel:
					if !wait.Stop() {
						<-wait.C
					}
					closeStream()
					done = true
					break
				}
			}
		}

		close(sendDone)
	}()

	_, _ = <-recvDone, <-sendDone

	for len(cancel) > 0 {
		<-cancel
	}

	close(doneCh)
	close(cancel)

	if err == nil && recvErr != nil {
		err = recvErr
	}

	return err
}
