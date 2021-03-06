package grpctester

import (
	"context"
	"encoding/json"
	"time"

	"github.com/douyu/juno/internal/pkg/packages/xtest"
	"github.com/douyu/juno/internal/pkg/service/grpctest/grpcinvoker"
	"github.com/jhump/protoreflect/desc"
)

type (
	Response     map[string]interface{}
	RequestInput map[string]interface{}
	Metadata     map[string]string

	RequestPayload struct {
		PackageName string
		ServiceName string
		MethodName  string
		Input       RequestInput
		MetaData    Metadata
		ProtoFile   string
		Host        string
		Timeout     time.Duration
		TestScript  string

		MethodDescriptor *desc.MethodDescriptor
	}
)

type (
	GrpcTester struct {
		tester *xtest.XTest
	}
)

func New() *GrpcTester {
	return &GrpcTester{
		tester: xtest.New(
			xtest.WithInterpreter(xtest.InterpreterTypeJS),
			xtest.WithGlobalStore(true),
		),
	}
}

func (g *GrpcTester) registerFunctions(payload *RequestPayload) {
	_ = g.tester.Interpreter().RegisterFunc("setInput", func(input RequestInput) {
		payload.Input = input
	})

	_ = g.tester.Interpreter().RegisterFunc("getInput", func() RequestInput {
		return payload.Input
	})

	_ = g.tester.Interpreter().RegisterFunc("getMetadata", func() Metadata {
		return payload.MetaData
	})

	_ = g.tester.Interpreter().RegisterFunc("setMetadata", func(m Metadata) {
		payload.MetaData = m
	})

	_ = g.tester.Interpreter().RegisterFunc("setHost", func(host string) {
		payload.Host = host
	})
}

func (g *GrpcTester) Run(c context.Context, payload RequestPayload) xtest.TestResult {
	g.registerFunctions(&payload)

	testScript := xtest.TestScript{Source: payload.TestScript}
	result := g.tester.Run(c, testScript, func() (data xtest.Response, err error) {
		return g.send(payload)
	})

	return result
}

func (g *GrpcTester) send(payload RequestPayload) (data xtest.Response, err error) {
	inputBytes, err := json.Marshal(payload.Input)
	if err != nil {
		return
	}

	md, err := json.Marshal(payload.MetaData)
	if err != nil {
		return
	}

	resp, err := grpcinvoker.MakeRequest(grpcinvoker.ReqProtoConfig{
		PackageName:      payload.PackageName,
		ServiceName:      payload.ServiceName,
		MethodName:       payload.MethodName,
		InputParams:      string(inputBytes),
		MetaData:         string(md),
		MethodDescriptor: payload.MethodDescriptor,
		Host:             payload.Host,
		Timeout:          payload.Timeout,
	})
	if err != nil {
		return
	}

	jsonBytes, err := resp.MarshalJSON()
	if err != nil {
		return
	}

	response := make(Response)
	err = json.Unmarshal(jsonBytes, &response)
	if err != nil {
		return
	}

	return response, err
}
