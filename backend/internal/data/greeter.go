package data

import (
	"context"
	"net/http"

	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"backend/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

type greeterRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *greeterRepo) Save(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	return g, nil
}

func (r *greeterRepo) Update(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	return g, nil
}

func (r *greeterRepo) FindByID(context.Context, int64) (*biz.Greeter, error) {
	return nil, nil
}

func (r *greeterRepo) ListByHello(context.Context, string) ([]*biz.Greeter, error) {
	return nil, nil
}

func (r *greeterRepo) ListAll(context.Context) ([]*biz.Greeter, error) {
	return nil, nil
}

func (r *greeterRepo) Query(ctx context.Context,span trace.Span, query *biz.QueryRequest) (*biz.QueryReply, error) {
	type User struct {
		bun.BaseModel `bun:"table:users"`

		ID uint64 `bun:",pk,autoincrement"`
		Username string `bun:"username"`
		Password string `bun:"password"`
		Age      uint32    `bun:"age"`
		Gender   biz.Gender `bun:"gender"`
	}

	user :=User{
		Username: query.Name,
	}
	if _, err := r.data.db.NewSelect().
		Model(&user).
		Where("username = ?", query.Name).
		Exec(ctx);
	err != nil {
		// Rely on Zap to record the error and set status code.
		otelzap.L().Ctx(ctx).Error("query failed", zap.Error(err))
		return nil, err
	}

	// 使用属性记录上下文信息
	if span.IsRecording() {
		span.SetAttributes(
			attribute.Int64("user.id", int64(user.ID)),
			attribute.String("user.username", user.Username),
		)
	}

	return &biz.QueryReply{
		Data: biz.QueryReplyData{
			ID:       user.ID,
			Age:      user.Age,
			Gender:   user.Gender,
			Username: user.Username,
		},
		Code:    http.StatusOK,
		Message: "Query successful",
	}, nil
}
