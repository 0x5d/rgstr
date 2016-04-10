package consul

import (
	"fmt"

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
	fmt.Println("Registering service", service.ID, "on Consul.")
	consulService := &api.AgentServiceRegistration{
		ID:      service.ID,
		Address: service.IP,
		Port:    int(service.Port),
		Name:    service.Name,
	}
	return adapter.client.Agent().ServiceRegister(consulService)
}

// Deregister deregisters a service from Consul.
func (adapter *Adapter) Deregister(service *registries.Service) error {
	fmt.Println("Deregistering service", service.ID, "from Consul.")
	return adapter.client.Agent().ServiceDeregister(service.ID)
}
