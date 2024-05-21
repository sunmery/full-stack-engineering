package server

import (
	"backend/internal/conf"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
)

func NewRegister(conf *conf.Registry) (registry.Registrar, error) {
	fmt.Println("conf:", conf.Consul)
	c := &api.Config{
		// Address: "192.168.2.181:8500",
		// Scheme:  "http",
		Address:    conf.Consul.Address,
		Scheme:     conf.Consul.Schema,
		TLSConfig:  api.TLSConfig{},
	}
	cli, err := api.NewClient(c)
	if err != nil {
		return nil, err
	}
	r := consul.New(cli,consul.WithHealthCheck(conf.Consul.healthCheck))
	return r, nil
}
