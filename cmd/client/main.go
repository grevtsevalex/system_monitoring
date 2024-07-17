package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	serverpb "github.com/grevtsevalex/system_monitoring/internal/server/pb"
)

var serverPort int

func init() {
	flag.IntVar(&serverPort, "port", 0, "grpc server port")
}

func main() {
	flag.Parse()

	if serverPort == 0 {
		log.Println("No server port in args")
		os.Exit(1)
	}

	conn, err := grpc.NewClient(fmt.Sprintf(":%d", serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := serverpb.NewMonitoringServiceClient(conn)
	ctx := metadata.NewOutgoingContext(context.Background(), nil)
	req := &serverpb.GetMetricsRequest{RangeSec: 10, PeriodSec: 3}
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	data, err := client.GetMetrics(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	for {
		sn, err := data.Recv()
		if err != nil {
			log.Fatal(err)
			break
		}

		laMetricString := fmt.Sprintf("LoadAverages: perMin: %.2f    per5Min: %.2f,  per15Min: %.2f", sn.La.PerMinute, sn.La.Per5Minutes, sn.La.Per15Minutes)
		cpuMetricString := fmt.Sprintf("CPU:          usr:    %.2f   sys:     %.2f,  idle:     %.2f", sn.Cpu.Usr, sn.Cpu.Sys, sn.Cpu.Idle)

		output := laMetricString + "\n" + cpuMetricString
		fmt.Println(output)
		fmt.Println()
	}
}
