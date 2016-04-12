package consul

import (
	"fmt"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

func TestToGenericService(t *testing.T) {
	name := "some-service"
	ip := "127.0.0.1"
	port := 5000
	id := fmt.Sprintf("%s:%d", ip, port)
	consulService := &api.AgentService{
		Service: name,
		Address: ip,
		Port:    port,
		ID:      id,
	}
	service := toGenericService(consulService)

	assert.Equal(t, service.Name, consulService.Service)
	assert.Equal(t, service.IP, consulService.Address)
	assert.Equal(t, service.Port, uint(consulService.Port))
	assert.Equal(t, service.ID, consulService.ID)
}
