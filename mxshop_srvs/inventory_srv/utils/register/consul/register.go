package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

type RegistryClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceId string) error
}

type Registry struct {
	Host string
	Port int
}

func NewRegistryClient(host string,port int) RegistryClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}
//服务注册
func (r Registry)Register(address string, port int, name string, tags []string, id string) error {
	//注册服务健康检查
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client,err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name:    name,
		ID:      id,
		Port:    port,
		Tags:    tags,
		Address: address,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d",address, port),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "20s",
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (r Registry)DeRegister(serviceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client,err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	err =  client.Agent().ServiceDeregister(serviceId)
	return err
}