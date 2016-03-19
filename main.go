package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/coreos/rkt/api/v1alpha"
)

var delay *int
var pods = make(map[string]*v1alpha.Pod)

func main() {
	address := flag.String("a", "localhost:15441", "The address where rkt's API service is listening.")
	delay = flag.Int("d", 100, "The polling interval (in milliseconds).")
	flag.Parse()
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c := v1alpha.NewPublicAPIClient(conn)
	defer conn.Close()

	done := make(chan bool)
	go startPolling(c, done)
	<-done
}

func startPolling(c v1alpha.PublicAPIClient, done chan bool) {
	for {
		res, err := getRunningPods(c)
		if err != nil {
			panic(err)
		}

		for _, pod := range res.Pods {
			_, ok := pods[pod.Id]
			if !ok {
				// Check if it's running.
				if pod.State == v1alpha.PodState_POD_STATE_RUNNING {
					// Map the pod.
					pods[pod.Id] = pod
					// TODO: register it.
				}
				continue
			}

			if pod.State == v1alpha.PodState_POD_STATE_EXITED {
				// TODO: Pod stopped. Deregister it.
				delete(pods, pod.Id)
			}
		}
		time.Sleep(time.Duration(*delay) * time.Millisecond)
	}
	done <- true
}

func getRunningPods(c v1alpha.PublicAPIClient) (*v1alpha.ListPodsResponse, error) {
	req := &v1alpha.ListPodsRequest{
		// Specify the request: Fetch and print only running pods and their details.
		Detail: true,
		Filters: []*v1alpha.PodFilter{
			{
				States: []v1alpha.PodState{v1alpha.PodState_POD_STATE_RUNNING},
			},
		},
	}
	return c.ListPods(context.Background(), req)
}
