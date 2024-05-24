package service

import (
	v1 "backend/api/helloworld/v1"
	"backend/internal/biz"
	"backend/internal/helper/log"
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"time"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServiceServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	tr := otel.Tracer("examples-hello-sayHello")
	spCtx, span := tr.Start(ctx, "sayHello")
	span.SetAttributes(attribute.String("name", "sayHello"))
	type _LogStruct struct {
		CurrentTime time.Time `json:"currentTime"`
		PassWho     string    `json:"passWho"`
		Name        string    `json:"name"`
	}
	logTest := _LogStruct{
		CurrentTime: time.Time{},
		PassWho:     "jzin",
		Name:        "sayHello",
	}
	log.InfofC(spCtx, "this is sayHello logs")
	b, _ := json.Marshal(logTest)
	log.InfofC(spCtx, string(b))
	span.SetAttributes(attribute.Key("测试key").String(string(b)))
	time.Sleep(time.Second)

	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	span.End()
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}

func (s *GreeterService) Query(ctx context.Context, q *v1.QueryRequest) (*v1.QueryReply, error)  {
	// 链路追踪, 包含本函数的信息, 方便定位
	tr := otel.Tracer("examples-hello-query")
	spCtx, span := tr.Start(ctx, "query")
	defer span.End()

	fmt.Println("q.Name",q.Name)

	// 传递ctx
	result, err := s.uc.Query(spCtx,span, &biz.QueryRequest{
		Name: q.Name,
	})
	if err != nil {
		return nil, err
	}

	return &v1.QueryReply{
		Data:    &v1.Data{
			Username:   result.Data.Username,
			Gender: v1.Data_Gender(result.Data.Gender),
		},
		Code: result.Code,
		Message: result.Message,
	},nil
}
