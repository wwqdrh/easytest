package grpctest

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/wwqdrh/easytest/grpctest/testdata/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// 连接已经启动的grpc服务获取服务信息

type GrpcDescribeTestSuite struct {
	suite.Suite

	secure        bool
	TestPort      string
	TestLocalhost string
	server        *grpc.Server
}

func TestGrpcDescribe(t *testing.T) {
	suite.Run(t, &GrpcDescribeTestSuite{secure: false})
}

func (suite *GrpcDescribeTestSuite) SetupSuite() {
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

}

func (suite *GrpcDescribeTestSuite) TearDownSuite() {
	if suite.server != nil {
		suite.server.Stop()
	}
}

func (suite *GrpcDescribeTestSuite) TestProtodesc_GetMethodDescFromReflect() {
	suite.T().Run("test known call", func(t *testing.T) {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, suite.TestLocalhost, opts...)
		assert.NoError(t, err)

		md := make(metadata.MD)

		refCtx := metadata.NewOutgoingContext(ctx, md)

		refClient := grpcreflect.NewClient(refCtx, reflectpb.NewServerReflectionClient(conn))

		mtd, err := GetMethodDescFromReflect("helloworld.Greeter.SayHello", refClient)
		assert.NoError(t, err)
		assert.NotNil(t, mtd)
		assert.Equal(t, "SayHello", mtd.GetName())
	})

	suite.T().Run("test known call with /", func(t *testing.T) {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, suite.TestLocalhost, opts...)
		assert.NoError(t, err)

		md := make(metadata.MD)

		refCtx := metadata.NewOutgoingContext(ctx, md)

		refClient := grpcreflect.NewClient(refCtx, reflectpb.NewServerReflectionClient(conn))

		mtd, err := GetMethodDescFromReflect("helloworld.Greeter/SayHello", refClient)
		assert.NoError(t, err)
		assert.NotNil(t, mtd)
		assert.Equal(t, "SayHello", mtd.GetName())
	})

	suite.T().Run("test unknown known call", func(t *testing.T) {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, suite.TestLocalhost, opts...)
		assert.NoError(t, err)

		md := make(metadata.MD)

		refCtx := metadata.NewOutgoingContext(ctx, md)

		refClient := grpcreflect.NewClient(refCtx, reflectpb.NewServerReflectionClient(conn))

		mtd, err := GetMethodDescFromReflect("helloworld.Greeter/SayHelloAsdf", refClient)
		assert.Error(t, err)
		assert.Nil(t, mtd)
	})
}
