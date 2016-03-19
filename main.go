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
	delay = flag.Int("d", 100, "The polling interval (in milliseconds).")
	flag.Parse()
	conn, err := grpc.Dial("localhost:15441", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c := v1alpha.NewPublicAPIClient(conn)
	defer conn.Close()

	if err = initPods(c); err != nil {
		panic(err)
	}
	done := make(chan bool)
	go startPolling(c, done)
	<-done
}

func initPods(c v1alpha.PublicAPIClient) error {
	res, err := getRunningAndExitedPods(c)
	if err != nil {
		return err
	}
	for _, p := range res.Pods {
		if p.State == v1alpha.PodState_POD_STATE_RUNNING {
			// TODO: register service.
		}
		pods[p.Id] = p
	}
	return nil
}

func startPolling(c v1alpha.PublicAPIClient, done chan bool) {
	for {
		res, err := getRunningAndExitedPods(c)
		if err != nil {
			panic(err)
		}

		for _, pod := range res.Pods {
			p, ok := pods[pod.Id]
			if !ok {
				// Pod wasn't mapped. Map it.
				pods[pod.Id] = pod
				// Check if it's running.
				if pod.State == v1alpha.PodState_POD_STATE_RUNNING {
					// TODO: register it.
				}
				continue
			}
			if pod.State == p.State {
				// No changes.
				continue
			}

			if pod.State == v1alpha.PodState_POD_STATE_RUNNING {
				pods[pod.Id] = pod
				// TODO: Pod started running. Register it.
			} else {
				// TODO: Deregister it.
				pods[pod.Id] = pod
			}
		}
		time.Sleep(time.Duration(*delay) * time.Millisecond)
	}
	done <- true
}

func getRunningAndExitedPods(c v1alpha.PublicAPIClient) (*v1alpha.ListPodsResponse, error) {
	req := &v1alpha.ListPodsRequest{
		// Specify the request: Fetch and print only running pods and their details.
		Detail: true,
		Filters: []*v1alpha.PodFilter{
			{
				States: []v1alpha.PodState{
					v1alpha.PodState_POD_STATE_RUNNING,
					v1alpha.PodState_POD_STATE_EXITED,
				},
			},
		},
	}
	return c.ListPods(context.Background(), req)
}
