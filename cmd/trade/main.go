package main

import (
	"context"
	"flag"
	"fmt"
	"sol-example/solgateway"

	"github.com/gagliardetto/solana-go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var url = flag.String("url", "solana-fra.blocksmith.org:16060", "url")
var token = flag.String("token", "", "token")
var privkey = flag.String("privkey", "", "privkey")
var tipAccount = flag.String("tipAccount", "jito", "tipAccount: jito|bloxroute|nextblock|temporal|(address)")
var tipAmount = flag.Uint64("tipAmount", 200000, "tipAmount")

func main() {
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

	privkey, err := solana.PrivateKeyFromBase58(*privkey)
	if err != nil {
		panic(err)
	}

	computeLimit := uint32(120_000)
	fee := &solgateway.Fee{
		ComputeLimit: &computeLimit,
	}
	if *tipAccount != "" && *tipAmount != 0 {
		fee.TipAccount = tipAccount
		fee.TipAmount = tipAmount
	}
	swapResp, err := g.RaydiumSwap(ctx, &solgateway.RaydiumSwapRequest{
		User:                 privkey.PublicKey().String(),
		Amm:                  "9Tb2ohu5P16BpBarqd3N27WnkF51Ukfs8Z1GzzLDxVZW",
		PoolCoinTokenAccount: "HQD2eNuCRbDCFfaPjFt6ZttM8EiD4BPg1MXy2ALVwobg",
		PoolPcTokenAccount:   "FhMHm2TVY9ULmZJvjjRx849Jn2ZRirVhJbTncHULTSvH",
		TokenIn:              solana.WrappedSol.String(),
		TokenOut:             "CzLSujWBLFsSjncfkh59rUFqvafWcY5tzedWJSuypump",
		AmountIn:             "1000000",
		AmountOut:            "1",
		CheckAta:             true,
		Fee:                  fee,
	})
	if err != nil {
		panic(err)
	}

	tx, err := solana.TransactionFromBase64(swapResp.Transaction)
	if err != nil {
		panic(err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.String() == privkey.PublicKey().String() {
			return &privkey
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(tx)

	txData := tx.MustToBase64()
	sendResp, err := g.SendTransaction(ctx, &solgateway.SendTransactionRequest{
		Transaction:   txData,
		SkipPreFlight: true,
		OpenPlatform:  *tipAccount,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("txSig: %s\n", sendResp.Signature)

}
