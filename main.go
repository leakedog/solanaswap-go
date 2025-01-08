package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	solanaswapgo "github.com/leakedog/solanaswap-go/solanaswap-go"
)

/*
Example Transactions:
- Orca: 2kAW5GAhPZjM3NoSrhJVHdEpwjmq9neWtckWnjopCfsmCGB27e3v2ZyMM79FdsL4VWGEtYSFi1sF1Zhs7bqdoaVT
- Pumpfun: 4Cod1cNGv6RboJ7rSB79yeVCR4Lfd25rFgLY3eiPJfTJjTGyYP1r2i1upAYZHQsWDqUbGd1bhTRm1bpSQcpWMnEz
- Banana Gun: oXUd22GQ1d45a6XNzfdpHAX6NfFEfFa9o2Awn2oimY89Rms3PmXL1uBJx3CnTYjULJw6uim174b3PLBFkaAxKzK
- Jupiter: DBctXdTTtvn7Rr4ikeJFCBz4AtHmJRyjHGQFpE59LuY3Shb7UcRJThAXC7TGRXXskXuu9LEm9RqtU6mWxe5cjPF
- Jupiter DCA: 4mxr44yo5Qi7Rabwbknkh8MNUEWAMKmzFQEmqUVdx5JpHEEuh59TrqiMCjZ7mgZMozRK1zW8me34w8Myi8Qi1tWP
- Meteora DLMM: 125MRda3h1pwGZpPRwSRdesTPiETaKvy4gdiizyc3SWAik4cECqKGw2gggwyA1sb2uekQVkupA2X9S4vKjbstxx3
- Rayd V4: 5kaAWK5X9DdMmsWm6skaUXLd6prFisuYJavd9B62A941nRGcrmwvncg3tRtUfn7TcMLsrrmjCChdEjK3sjxS6YG9
- Rayd Routing: 51nj5GtAmDC23QkeyfCNfTJ6Pdgwx7eq4BARfq1sMmeEaPeLsx9stFA3Dzt9MeLV5xFujBgvghLGcayC3ZevaQYi
- Rayd CPMM: afUCiFQ6amxuxx2AAwsghLt7Q9GYqHfZiF4u3AHhAzs8p1ThzmrtSUFMbcdJy8UnQNTa35Fb1YqxR6F9JMZynYp
- Rayd Concentrated Liquidity SwapV2: 2durZHGFkK4vjpWFGc5GWh5miDs8ke8nWkuee8AUYJA8F9qqT2Um76Q5jGsbK3w2MMgqwZKbnENTLWZoi3d6o2Ds
- Rayd Concentrated Liquidity Swap: 4MSVpVBwxnYTQSF3bSrAB99a3pVr6P6bgoCRDsrBbDMA77WeQqoBDDDXqEh8WpnUy5U4GeotdCG9xyExjNTjYE1u // Not parse bot multiple Amms, do parse one
- Maestro: mWaH4FELcPj4zeY4Cgk5gxUirQDM7yE54VgMEVaqiUDQjStyzwNrxLx4FMEaKEHQoYsgCRhc1YdmBvhGDRVgRrq
- Meteora Pools Program: 4uuw76SPksFw6PvxLFkG9jRyReV1F4EyPYNc3DdSECip8tM22ewqGWJUaRZ1SJEZpuLJz1qPTEPb2es8Zuegng9Z
- Moonshot: AhiFQX1Z3VYbkKQH64ryPDRwxUv8oEPzQVjSvT7zY58UYDm4Yvkkt2Ee9VtSXtF6fJz8fXmb5j3xYVDF17Gr9CG (Buy)
- Moonshot: 2XYu86VrUXiwNNj8WvngcXGytrCsSrpay69Rt3XBz9YZvCQcZJLjvDfh9UWETFtFW47vi4xG2CkiarRJwSe6VekE (Sell)
- Multiple AMMs: 46Jp5EEUrmdCVcE3jeewqUmsMHhqiWWtj243UZNDFZ3mmma6h2DF4AkgPE9ToRYVLVrfKQCJphrvxbNk68Lub9vw
- Okx: 4JrQfHYrAsLzJNQPpiEjc53jLfAqhV6bigX6RxHW5esUecX6rND5vA2QyEpScjfvm9hGAN5wjhbN7xqWwDLvv8ew
- Phoenix: 5hUxS3e6hidTUsnHcrnK7ZvENgTyZr86rRyxAZZzSAexXdT3aLDu7Akvq4MB94bkd7WhMibzZRZoiDY8nKFGAnxn
*/

func main() {
	rpcClient := rpc.New(rpc.MainNetBeta.RPC)
	// rpcClient := rpc.New("https://solana-rpc.publicnode.com")
	fetchSignatureTx(rpcClient)
	// fetchBlockTx(rpcClient)
}

func fetchSignatureTx(rpcClient *rpc.Client) {
	txSig := solana.MustSignatureFromBase58("27yQ87jYN7hdkhaQoWn7jyM5Q1Kx649Yx8zt65eZCcnpJ3iseNEaBTh438bFL3meSheBHZVwU2oYuELPcViyixAz")
	var maxTxVersion uint64 = 0
	tx, err := rpcClient.GetTransaction(
		context.TODO(),
		txSig,
		&rpc.GetTransactionOpts{
			Commitment:                     rpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: &maxTxVersion,
		},
	)
	if err != nil {
		log.Fatalf("error getting tx: %s", err)
	}

	parser, err := solanaswapgo.NewParser(tx)
	if err != nil {
		log.Fatalf("error creating orca parser: %s", err)
	}

	parser.NewTxParser()
	// spew.Dump(parser.Actions)
	// spew.Dump("========================BREAK========================", parser.SwapData)

	// swapData, err := parser.ProcessSwapData(parser.SwapData)
	// if err != nil {
	// 	log.Fatalf("Error processing swap data: %s", err)
	// }

	// fmt.Println("Swap Data:", swapData)
	// Marshal the parsed result to JSON
	parsedJson, err := json.Marshal(parser)
	if err != nil {
		fmt.Println("Error marshalling JSON:", parsedJson)
		return
	}

	jsonData, err := json.MarshalIndent(parser, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Print the indented JSON string
	fmt.Println(string(jsonData))
}

// func fetchBlockTx(rpcClient *rpc.Client) {
// 	blockNumber := uint64(306035317)
// 	maxSupportedTransactionVersion := uint64(0)
// 	out, err := rpcClient.GetBlockWithOpts(context.TODO(), blockNumber, &rpc.GetBlockOpts{
// 		MaxSupportedTransactionVersion: &maxSupportedTransactionVersion,
// 	})
// 	if err != nil {
// 		if strings.HasSuffix(err.(*jsonrpc.RPCError).Message, "was skipped, or missing due to ledger jump to recent snapshot") {
// 			log.Println("skipped block: ", blockNumber)
// 			return
// 		}
// 		panic(err)
// 	}

// 	for _, tx := range out.Transactions {
// 		tx.BlockTime = out.BlockTime
// 		tx.Slot = blockNumber // out.ParentSlot+1 It's wrong. The previous block is not continuous and there will be bad blocks in the middle.

// 		parser, err := solanaswapgo.NewParser(&tx)
// 		if err != nil {
// 			log.Fatalf("error creating orca parser: %s", err)
// 		}

// 		parser.NewTxParser()
// 		spew.Dump(parser.Actions)
// 		spew.Dump("========================BREAK========================", parser.SwapData)

// 		// transactionData, err := parser.ParseTransaction()
// 		// if err != nil {
// 		// 	log.Fatalf("error parsing transaction: %s", err)
// 		// }

// 		// // marshalledData, _ := json.MarshalIndent(transactionData, "", "  ")
// 		// // spew.Dump(transactionData)
// 		// // return
// 		// swapData, err := parser.ProcessSwapData(transactionData)
// 		// if err != nil {
// 		// 	log.Fatalf("error processing swap data: %s", err)
// 		// }

// 		// marshalledSwapData, _ := json.MarshalIndent(swapData, "", "  ")
// 		// if "3gtHe9aUoBiiyZBNuiWpAjqMz6LCK7V6LXijHJNbL1DtRdFyBsTrXgMV4sDcU7PEHeyYhyPxXt6Ynt1mkE6wsRVz" == swapData.Signatures[0].String() {
// 		// 	fmt.Println(string(marshalledSwapData))
// 		// 	return
// 		// }
// 	}
// }
