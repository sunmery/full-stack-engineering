package biz

import (
	"context"
	"go.opentelemetry.io/otel/trace"

	v1 "backend/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Gender 枚举在 Go 中的表示
type Gender uint32

const (
	GenderMale   Gender = 1
	GenderFemale Gender = 2
)

func (g Gender) String() string {
	switch g {
	case GenderMale:
		return "male"
	case GenderFemale:
		return "female"
	default:
		return "unknown"
	}
}

// Greeter is a Greeter model.
type Greeter struct {
	Hello string
}

type QueryRequest struct {
	Name string
}

type QueryReplyData struct {
	ID uint64 `json:"id"`
	Age uint32
	Gender Gender
	Username string
}
type QueryReply struct {
	Data QueryReplyData
	Code uint32
	Message string
}

// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	Save(context.Context, *Greeter) (*Greeter, error)
	Update(context.Context, *Greeter) (*Greeter, error)
	FindByID(context.Context, int64) (*Greeter, error)
	ListByHello(context.Context, string) ([]*Greeter, error)
	ListAll(context.Context) ([]*Greeter, error)
	Query(ctx context.Context, span trace.Span, q *QueryRequest) (*QueryReply, error)
}

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo GreeterRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo, logger log.Logger) *GreeterUsecase {
	return &GreeterUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *GreeterUsecase) CreateGreeter(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.WithContext(ctx).Infof("CreateGreeter: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}

// Query 查询名字
func (uc *GreeterUsecase) Query(ctx context.Context, span trace.Span, q *QueryRequest) (*QueryReply, error){
	uc.log.WithContext(ctx).Infof("Query: %v", q.Name)
	return uc.repo.Query(ctx, span, q)
}
