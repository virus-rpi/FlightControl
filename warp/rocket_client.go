package warp

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/cskr/pubsub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"time"
)

var Client *rocketClient

type rocketClient struct {
	C   WaterRocketServiceClient
	Ctx *context.Context
}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no non-loopback IP address found")
}

func InitRocketClient(App fyne.App, ps *pubsub.PubSub) {
	rocketAddress := App.Preferences().StringWithFallback("RocketAddress", "Not set")
	log.Println("Rocket address: " + rocketAddress)

	if rocketAddress == "Not set" {
		return
	}

	conn, err := grpc.NewClient(rocketAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to rocket: %v", err)
	}
	defer conn.Close()
	c := NewWaterRocketServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ip, err := getLocalIP()
	_, err = c.SetControlServiceAddress(ctx, &SetControlServiceAddressRequest{Address: ip + ":50051"})
	if err != nil {
		log.Printf("failed to set control service address: %v", err)
	}
	go NewControlServiceServer(ps)

	Client = &rocketClient{C: c, Ctx: &ctx}
}

func RefreshRocketClient(App fyne.App) {
	rocketAddress := App.Preferences().StringWithFallback("RocketAddress", "Not set")
	log.Println("Rocket address: " + rocketAddress)

	if rocketAddress == "Not set" {
		return
	}

	conn, err := grpc.NewClient(rocketAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to rocket: %v", err)
	}
	defer conn.Close()
	c := NewWaterRocketServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ip, err := getLocalIP()
	_, err = c.SetControlServiceAddress(ctx, &SetControlServiceAddressRequest{Address: ip + ":50051"})
	if err != nil {
		log.Printf("failed to set control service address: %v", err)
	}

	Client.C = c
	Client.Ctx = &ctx
}
