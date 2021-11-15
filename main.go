package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/serbrech/broadcast/target"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"k8s.io/klog"
)

func main() {
	mode, ok := os.LookupEnv("MODE")
	if ok && mode == "target"{
		runTarget()
	} else {
		runBroadcaster()
	}
}

func runBroadcaster(){
	ctx := context.Background()
	for {
		cname, srvs, err := net.LookupSRV("grpc", "tcp", "target-svc.default.svc.cluster.local")
		if err != nil {
			klog.Error(err)
		}

		klog.Infof("cname: %s \n", cname)
		for _, srv := range srvs {
			klog.Infof("\t- %v:%v:%d:%d", srv.Target, srv.Port, srv.Priority, srv.Weight)
			callTarget(ctx, fmt.Sprintf("%s:%d", srv.Target, srv.Port))
		}

		time.Sleep(2*time.Second)
	}
}

func callTarget(ctx context.Context, dialTarget string) {
	conn, err := grpc.DialContext(ctx, dialTarget, grpc.WithInsecure(), grpc.WithConnectParams(grpc.ConnectParams{
		Backoff:           backoff.Config{
			BaseDelay:  100*time.Millisecond,
			Multiplier: 2,
			MaxDelay:   10*time.Second,
		},
	}))
	if err!=nil{
		klog.Errorf("\tfailed to connect to %s: %s", dialTarget, err)
		return
	}
	defer conn.Close()
	greeter := target.NewGreeterClient(conn)
	reply, err := greeter.SayHello(ctx, &target.HelloRequest{Name: dialTarget})
	if err != nil {
		klog.Errorf("\tfailed to talk to %s: %s", dialTarget, err)
		return
	}
	klog.Infof("\tSuccess: %s", reply.Message)
}

func runTarget() {
	klog.Info("Starting target server")
	address := ":9000"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		klog.Fatalf("failed to listen tcp %v", err)
	}
	srv := grpc.NewServer()
	target.RegisterGreeterServer(srv, target.Server{})
	klog.Infof("Listening at %s", address)
	if err := srv.Serve(lis); err!=nil {
		klog.Fatal(err)
	}
}


