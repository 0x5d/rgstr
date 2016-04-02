package consul

import (
	"errors"

	"github.com/castillobg/rgstr/registries"
	"github.com/hashicorp/consul/api"
)

func init() {
	registries.Register(new(Factory), "consul")
}

// Factory implements registries.AdapterFactory.
type Factory struct{}

// Adapter represents a Consul RegistryAdapter.
type Adapter struct {
	client *api.Client
}

// New builds a Consul RegistryAdapter
func (*Factory) New(address string) (registries.RegistryAdapter, error) {
	config := api.DefaultConfig()
	config.Address = address
	consulClient, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	consulAdapter := &Adapter{client: consulClient}
	return consulAdapter, nil
}

// Register registers a new service on Consul.
func (adapter *Adapter) Register(service *registries.Service) error {
	return errors.New("Not yet implemented.")
}

// Deregister deregisters a service from Consul.
func (adapter *Adapter) Deregister(service *registries.Service) error {
	return errors.New("Not yet implemented.")
}
