package main

import (
	"context"
	"flag"
	"fmt"
	"sol-example/solgateway"

	"github.com/gagliardetto/solana-go"
	cb "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var url = flag.String("url", "solana-fra.blocksmith.org:16060", "url")
var token = flag.String("token", "", "token")
var privkey = flag.String("privkey", "", "privkey")
var amount = flag.Uint64("amount", 1e6, "amount")
var to = flag.String("to", "", "to")

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

	var dest solana.PublicKey
	if *to == "" {
		dest = privkey.PublicKey()
	} else {
		dest, err = solana.PublicKeyFromBase58(*to)
		if err != nil {
			panic(err)
		}
	}
	tx := transfer(privkey.PublicKey(), *amount, dest)

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.String() == privkey.PublicKey().String() {
			return &privkey
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	txData := tx.MustToBase64()
	sendResp, err := g.SendTransaction(ctx, &solgateway.SendTransactionRequest{
		Transaction:   txData,
		SkipPreFlight: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("txSig: %s\n", sendResp.Signature)

}

func transfer(pubkey solana.PublicKey, amount uint64, to solana.PublicKey) *solana.Transaction {
	client := rpc.New(rpc.MainNetBeta.RPC)
	cu := cb.NewSetComputeUnitLimitInstruction(500).Build()
	xfer := system.NewTransferInstruction(amount, pubkey, to).Build()
	instructions := []solana.Instruction{cu, xfer}
	recent, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentConfirmed)
	if err != nil {
		panic(err)
	}

	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(pubkey),
	)
	if err != nil {
		panic(err)
	}
	return tx
}
