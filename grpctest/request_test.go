package grpctest

import (
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/wwqdrh/easytest/grpctest/testdata/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type GrpcRequestTestSuite struct {
	suite.Suite

	secure        bool
	TestPort      string
	TestLocalhost string
	server        *grpc.Server
}

func TestGrpcRequest(t *testing.T) {
	suite.Run(t, &GrpcRequestTestSuite{secure: false})
}

func (suite *GrpcRequestTestSuite) SetupSuite() {
	lis, err := net.Listen("tcp", ":0")
	require.Nil(suite.T(), err)

	var opts []grpc.ServerOption

	if suite.secure {
		creds, err := credentials.NewServerTLSFromFile("./testdata/localhost.crt", "./testdata/localhost.key")
		require.Nil(suite.T(), err)
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	stats := helloworld.NewHWStats()

	opts = append(opts, grpc.StatsHandler(stats))

	s := grpc.NewServer(opts...)
	suite.server = s

	gs := helloworld.NewGreeter()
	helloworld.RegisterGreeterServer(s, gs)
	reflection.Register(s)

	gs.Stats = stats

	suite.TestPort = strconv.Itoa(lis.Addr().(*net.TCPAddr).Port)
	suite.TestLocalhost = "localhost:" + suite.TestPort

	go func() {
		_ = s.Serve(lis)
	}()
	time.Sleep(5 * time.Second) // wait server start
}

func (suite *GrpcRequestTestSuite) TearDownSuite() {
	if suite.server != nil {
		suite.server.Stop()
	}
}

func (s *GrpcRequestTestSuite) TestWorkerCall() {
	r, err := NewRequester(&RunConfig{
		Host:        s.TestLocalhost,
		Proto:       "./testdata/greeter.proto",
		Handle:      "helloworld.Greeter.SayHello",
		insecure:    true,
		dialTimeout: 1 * time.Second,
		ImportPath:  []string{},
	})
	require.Nil(s.T(), err)

	w := &Worker{
		stub:   r.stub,
		mtd:    r.mtd,
		config: r.config,
	}
	fmt.Println(w.runWorker())
}
