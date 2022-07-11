package grpctest

import (
	"errors"

	"github.com/jhump/protoreflect/dynamic"
)

type GrpcContext struct {
}

type HandleOption struct {
	Url         string
	Proto       string
	ImportPaths []string
	Call        string
	Expect      []string
	Event       []string
}

func NewGrpcContext() *GrpcContext {
	return &GrpcContext{}
}

func (*GrpcContext) Do(option *HandleOption) error {
	if option.Proto == "" {
		return errors.New("未指定proto文件")
	}
	mtd, err := GetMethodDescFromProto(option.Call, option.Proto, option.ImportPaths)
	if err != nil {
		return err
	}

	md := mtd.GetInputType()
	payloadMessage := dynamic.NewMessage(md)
	if payloadMessage == nil {
		return errors.New("No input type of method: " + mtd.GetName())
	}

	return nil
}
