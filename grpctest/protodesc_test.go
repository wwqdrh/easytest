package grpctest

import (
	"context"
	"testing"

	"net"
	"strconv"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/wwqdrh/easytest/grpctest/testdata/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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

func (suite *GrpcDescribeTestSuite) SetupTest() {
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

// func (suite *GrpcDescribeTestSuite) TestByCollection() {
// 	collections, err := NewCollections("./testdata/grpc_collection.json", nil)
// 	require.Nil(suite.T(), err)

// 	// mock server
// 	var opts []grpc.DialOption
// 	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	ctx := context.Background()
// 	conn, err := grpc.DialContext(ctx, suite.TestLocalhost, opts...)
// 	assert.NoError(suite.T(), err)

// 	md := make(metadata.MD)
// 	refCtx := metadata.NewOutgoingContext(ctx, md)
// 	refClient := grpcreflect.NewClient(refCtx, reflectpb.NewServerReflectionClient(conn))

// 	for _, collect := range collections {
// 		mtd, err := GetMethodDescFromReflect(collect.Call, refClient)
// 		assert.NoError(suite.T(), err)
// 		assert.NotNil(suite.T(), mtd)
// 		assert.Equal(suite.T(), "SayHello", mtd.GetName())
// 	}
// }

func TestProtodesc_GetMethodDescFromProto(t *testing.T) {
	t.Run("invalid path", func(t *testing.T) {
		md, err := GetMethodDescFromProto("pkg.Call", "invalid.proto", []string{})
		assert.Error(t, err)
		assert.Nil(t, md)
	})

	t.Run("invalid call symbol", func(t *testing.T) {
		md, err := GetMethodDescFromProto("pkg.Call", "./testdata/greeter.proto", []string{})
		assert.Error(t, err)
		assert.Nil(t, md)
	})

	t.Run("invalid package", func(t *testing.T) {
		md, err := GetMethodDescFromProto("helloworld.pkg.SayHello", "./testdata/greeter.proto", []string{})
		assert.Error(t, err)
		assert.Nil(t, md)
	})

	t.Run("invalid method", func(t *testing.T) {
		md, err := GetMethodDescFromProto("helloworld.Greeter.Foo", "./testdata/greeter.proto", []string{})
		assert.Error(t, err)
		assert.Nil(t, md)
	})

	t.Run("valid symbol", func(t *testing.T) {
		md, err := GetMethodDescFromProto("helloworld.Greeter.SayHello", "./testdata/greeter.proto", []string{})
		assert.NoError(t, err)
		assert.NotNil(t, md)
	})

	t.Run("valid symbol slashes", func(t *testing.T) {
		md, err := GetMethodDescFromProto("helloworld.Greeter/SayHello", "./testdata/greeter.proto", []string{})
		assert.NoError(t, err)
		assert.NotNil(t, md)
	})

	t.Run("proto3 optional support", func(t *testing.T) {
		md, err := GetMethodDescFromProto("helloworld.OptionalGreeter/SayHello", "./testdata/optional.proto", []string{})
		assert.NoError(t, err)
		assert.NotNil(t, md)
	})
}

func TestProtodesc_GetMethodDescFromProtoSet(t *testing.T) {
	t.Run("invalid path", func(t *testing.T) {
		md, err := GetMethodDescFromProtoSet("pkg.Call", "invalid.protoset")
		assert.Error(t, err)
		assert.Nil(t, md)
	})

	t.Run("invalid call symbol", func(t *testing.T) {
		md, err := GetMethodDescFromProtoSet("pkg.Call", "./testdata/bundle.protoset")
		assert.Error(t, err)
		assert.Nil(t, md)
	})

	t.Run("invalid package", func(t *testing.T) {
		md, err := GetMethodDescFromProtoSet("helloworld.pkg.SayHello", "./testdata/bundle.protoset")
		assert.Error(t, err)
		assert.Nil(t, md)
	})

	t.Run("invalid method", func(t *testing.T) {
		md, err := GetMethodDescFromProtoSet("helloworld.Greeter.Foo", "./testdata/bundle.protoset")
		assert.Error(t, err)
		assert.Nil(t, md)
	})

	t.Run("valid symbol", func(t *testing.T) {
		md, err := GetMethodDescFromProtoSet("helloworld.Greeter.SayHello", "./testdata/bundle.protoset")
		assert.NoError(t, err)
		assert.NotNil(t, md)
	})

	t.Run("valid symbol proto 2", func(t *testing.T) {
		md, err := GetMethodDescFromProtoSet("cap.Capper.Cap", "./testdata/bundle.protoset")
		assert.NoError(t, err)
		assert.NotNil(t, md)
	})

	t.Run("valid symbol slashes", func(t *testing.T) {
		md, err := GetMethodDescFromProtoSet("helloworld.Greeter/SayHello", "./testdata/bundle.protoset")
		assert.NoError(t, err)
		assert.NotNil(t, md)
	})
}

func TestParseServiceMethod(t *testing.T) {
	testParseServiceMethodSuccess(t, "package.Service.Method", "package.Service", "Method")
	testParseServiceMethodSuccess(t, ".package.Service.Method", "package.Service", "Method")
	testParseServiceMethodSuccess(t, "package.Service/Method", "package.Service", "Method")
	testParseServiceMethodSuccess(t, ".package.Service/Method", "package.Service", "Method")
	testParseServiceMethodSuccess(t, "Service.Method", "Service", "Method")
	testParseServiceMethodSuccess(t, ".Service.Method", "Service", "Method")
	testParseServiceMethodSuccess(t, "Service/Method", "Service", "Method")
	testParseServiceMethodSuccess(t, ".Service/Method", "Service", "Method")
	testParseServiceMethodError(t, "")
	testParseServiceMethodError(t, ".")
	testParseServiceMethodError(t, "package/Service/Method")
}

func testParseServiceMethodSuccess(t *testing.T, svcAndMethod string, expectedService string, expectedMethod string) {
	service, method, err := parseServiceMethod(svcAndMethod)
	assert.NoError(t, err)
	assert.Equal(t, expectedService, service)
	assert.Equal(t, expectedMethod, method)
}

func testParseServiceMethodError(t *testing.T, svcAndMethod string) {
	_, _, err := parseServiceMethod(svcAndMethod)
	assert.Error(t, err)
}
