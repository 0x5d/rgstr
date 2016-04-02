package rkt

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/castillobg/rgstr/registries"
	"github.com/castillobg/rgstr/runtimes"
	"github.com/coreos/rkt/api/v1alpha"
	"google.golang.org/grpc"
)

var pods = make(map[string]*v1alpha.Pod)

func init() {
	runtimes.Register(new(Factory), "rkt")
}

// Factory implements runtimes.AdapterFactory
type Factory struct{}

// Adapter implements runtimes.RuntimeAdapter
type Adapter struct {
	Address  string
	Registry *registries.RegistryAdapter
}

// New builds a new runtimes.RuntimeAdapter.
func (*Factory) New(address string, registry *registries.RegistryAdapter) (runtimes.RuntimeAdapter, error) {
	rktAdapter := &Adapter{
		Address:  address,
		Registry: registry,
	}
	return rktAdapter, nil
}

// Listen triggers event listening.
func (adapter *Adapter) Listen(errs chan error) {
	conn, err := grpc.Dial(adapter.Address, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c := v1alpha.NewPublicAPIClient(conn)
	defer conn.Close()

	errs <- startPolling(c)
}

func startPolling(c v1alpha.PublicAPIClient) error {
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
					fmt.Println("Registered Pod: " + pod.Id)
					// TODO: register it.
				}
				continue
			}

			if pod.State == v1alpha.PodState_POD_STATE_EXITED {
				// TODO: Pod stopped. Deregister it.
				fmt.Println("Deregistered Pod: " + pod.Id)
				delete(pods, pod.Id)
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
