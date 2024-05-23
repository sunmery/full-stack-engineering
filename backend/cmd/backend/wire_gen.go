// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"backend/internal/biz"
	"backend/internal/conf"
	"backend/internal/data"
	"backend/internal/server"
	"backend/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, auth *conf.Auth, registry *conf.Registry, trace *conf.Trace, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	registrar, err := server.NewRegister(registry)
	if err != nil {
		return nil, nil, err
	}
	db := data.NewDB(confData)
	client := data.NewCache(confData)
	dataData, cleanup, err := data.NewData(confData, logger, db, client)
	if err != nil {
		return nil, nil, err
	}
	greeterRepo := data.NewGreeterRepo(dataData, logger)
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)
	greeterService := service.NewGreeterService(greeterUsecase)
	grpcServer := server.NewGRPCServer(confServer, greeterService, logger)
	httpServer := server.NewHTTPServer(auth, confServer, trace, greeterService, logger)
	app := newApp(logger, registrar, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
