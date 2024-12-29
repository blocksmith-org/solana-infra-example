package main

import (
	"context"
	"flag"
	"fmt"
	"sol-example/solgateway"

	"github.com/greyireland/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var url = flag.String("url", "solana-fra.blocksmith.org:16060", "url")
var token = flag.String("token", "", "token")
var topic = flag.String("topic", "raydium-swap", "topic")

const (
	topicPumpCreate    = "pump-create"
	topicPumpSwap      = "pump-swap"
	topicPumpWithdraw  = "pump-withdraw"
	topicRaydiumCreate = "raydium-create"
	topicRaydiumSwap   = "raydium-swap"
)

func main() {
	log.Root().SetHandler(log.StdoutHandler)
	flag.Parse()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(*url, opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ctx := metadata.NewOutgoingContext(context.TODO(), metadata.Pairs("authorization", *token))

	g := solgateway.NewSolgatewayClient(conn)
	switch *topic {
	case topicPumpCreate:
		subscribePumpCreate(ctx, g)
	case topicPumpSwap:
		subscribePumpSwap(ctx, g)
	case topicPumpWithdraw:
		subscribePumpWithdraw(ctx, g)
	case topicRaydiumCreate:
		subscribeRaydiumCreate(ctx, g)
	case topicRaydiumSwap:
		subscribeRadiumSwap(ctx, g)
	}
}

func subscribePumpCreate(ctx context.Context, gateway solgateway.SolgatewayClient) {
	log.Info("start subscribe pump create")
	resp, err := gateway.SubscribePumpFunCreateStream(ctx, &solgateway.SubscribePumpFunCreateStreamRequest{})
	if err != nil {
		log.Warn("Subscribe pump create error", "err", err)
	}
	log.Info("subscribe pump create success")
	for {
		data, err := resp.Recv()
		if err != nil {
			log.Warn("recv error", "err", err)
		}
		fmt.Println(data)
	}

}

func subscribePumpSwap(ctx context.Context, gateway solgateway.SolgatewayClient) {
	log.Info("start subscribe pump swap")
	resp, err := gateway.SubscribePumpFunSwapStream(ctx, &solgateway.SubscribePumpFunSwapStreamRequest{})
	if err != nil {
		log.Warn("Subscribe pump swap error", "err", err)
	}
	log.Info("subscribe pump swap success")
	for {
		data, err := resp.Recv()
		if err != nil {
			log.Warn("recv error", "err", err)
		}
		fmt.Println(data)
	}

}

func subscribePumpWithdraw(ctx context.Context, gateway solgateway.SolgatewayClient) {
	log.Info("start subscribe pump withdraw")
	resp, err := gateway.SubscribePumpFunWithdrawStream(ctx, &solgateway.SubscribePumpFunWithdrawStreamRequest{})
	if err != nil {
		log.Warn("Subscribe pump withdraw error", "err", err)
	}
	log.Info("subscribe pump withdraw success")
	for {
		data, err := resp.Recv()
		if err != nil {
			log.Warn("recv error", "err", err)
		}

		fmt.Println(data)
	}
}

func subscribeRaydiumCreate(ctx context.Context, gateway solgateway.SolgatewayClient) {
	log.Info("start subscribe raydium create")
	resp, err := gateway.SubscribeRaydiumCreateStream(ctx, &solgateway.SubscribeRaydiumCreateStreamRequest{})
	if err != nil {
		log.Warn("Subscribe raydium create error", "err", err)
	}
	log.Info("subscribe raydium create success")
	for {
		data, err := resp.Recv()
		if err != nil {
			log.Warn("recv error", "err", err)
		}
		fmt.Println(data)
	}

}

func subscribeRadiumSwap(ctx context.Context, gateway solgateway.SolgatewayClient) {
	log.Info("start subscribe raydium swap")
	resp, err := gateway.SubscribeRaydiumSwapStream(ctx, &solgateway.SubscribeRaydiumSwapStreamRequest{})
	if err != nil {
		log.Warn("Subscribe raydium swap error", "err", err)
	}
	log.Info("subscribe raydium swap success")
	for {
		data, err := resp.Recv()
		if err != nil {
			log.Warn("recv error", "err", err)
		}
		fmt.Println(data)
	}
}
