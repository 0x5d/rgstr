package rkt

import (
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/appc/spec/schema"
	"github.com/castillobg/rgstr/registries"
	"github.com/castillobg/rgstr/runtimes"
	"github.com/coreos/rkt/api/v1alpha"
	"google.golang.org/grpc"
)

// TODO FIXME: Set a mutex on this map.
var pods = make(map[string]*v1alpha.Pod)

func init() {
	runtimes.Register(new(Factory), "rkt")
}

// Factory implements runtimes.AdapterFactory
type Factory struct{}

// Adapter implements runtimes.RuntimeAdapter
type Adapter struct {
	Address  string
	Registry registries.RegistryAdapter
}

// New builds a new runtimes.RuntimeAdapter.
func (*Factory) New(address string, registry registries.RegistryAdapter) (runtimes.RuntimeAdapter, error) {
	rktAdapter := &Adapter{
		Address:  address,
		Registry: registry,
	}
	return rktAdapter, nil
}

// Listen triggers event listening.
func (adapter *Adapter) Listen(errs chan error) {
	conn, err := grpc.Dial(
		adapter.Address,
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Duration(10)*time.Second),
	)
	defer conn.Close()
	if err != nil {
		errs <- err
	}
	c := v1alpha.NewPublicAPIClient(conn)

	errs <- startPolling(adapter, c)
}

func startPolling(adapter *Adapter, c v1alpha.PublicAPIClient) error {
	for {
		res, err := getPods(c)
		if err != nil {
			return err
		}

		for _, pod := range res.Pods {
			_, ok := pods[pod.Id]
			if !ok {
				// Check if it's running.
				if pod.State == v1alpha.PodState_POD_STATE_RUNNING {
					// Map the pod.
					pods[pod.Id] = pod
					services, err := getPodServices(pod)
					if err != nil {
						return err
					}
					for _, service := range services {
						err = adapter.Registry.Register(service)
						if err != nil {
							return err
						}
					}
				}
				continue
			}

			if pod.State == v1alpha.PodState_POD_STATE_EXITED {
				services, err := getPodServices(pods[pod.Id])
				delete(pods, pod.Id)
				if err != nil {
					return err
				}
				// Deregister every pod service from the registry.
				for _, service := range services {
					err = adapter.Registry.Deregister(service)
					if err != nil {
						return err
					}
				}
			}
		}
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}

func getPods(c v1alpha.PublicAPIClient) (*v1alpha.ListPodsResponse, error) {
	req := &v1alpha.ListPodsRequest{
		// Specify the request: Fetch and print only running pods and their details.
		Detail: true,
		Filters: []*v1alpha.PodFilter{
			{
				States: []v1alpha.PodState{
					v1alpha.PodState_POD_STATE_ABORTED_PREPARE,
					v1alpha.PodState_POD_STATE_DELETING,
					v1alpha.PodState_POD_STATE_EMBRYO,
					v1alpha.PodState_POD_STATE_EXITED,
					v1alpha.PodState_POD_STATE_GARBAGE,
					v1alpha.PodState_POD_STATE_PREPARED,
					v1alpha.PodState_POD_STATE_PREPARING,
					v1alpha.PodState_POD_STATE_RUNNING,
					v1alpha.PodState_POD_STATE_UNDEFINED,
				},
			},
		},
	}
	return c.ListPods(context.Background(), req)
}

func getPodManifest(pod *v1alpha.Pod) (*schema.PodManifest, error) {
	// The pod manifest is a JSON string. We have to unmarshal it.
	var manifest = new(schema.PodManifest)
	err := manifest.UnmarshalJSON(pod.Manifest)
	if err != nil {
		return nil, err
	}
	return manifest, nil
}

func getPodServices(pod *v1alpha.Pod) ([]*registries.Service, error) {
	manifest, err := getPodManifest(pod)
	if err != nil {
		return nil, err
	}
	// Register one service for each port on each app.
	services := []*registries.Service{}
	for _, app := range manifest.Apps {
		for _, port := range app.App.Ports {
			for _, network := range pod.Networks {
				// FIXME: What happens if there's no IPv4 addresses?
				if len(network.Ipv4) != 0 {
					service := &registries.Service{
						ID:   network.Ipv4 + ":" + strconv.Itoa(int(port.Port)),
						Name: string(app.Name) + "-" + string(port.Name),
						Port: port.Port,
						IP:   network.Ipv4,
					}
					services = append(services, service)
					break
				}
			}
		}
	}
	return services, nil
}
