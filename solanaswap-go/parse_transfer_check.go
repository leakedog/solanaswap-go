package solanaswapgo

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	"github.com/gagliardetto/solana-go"
)

type TransferCheck struct {
	Info struct {
		Authority   string `json:"authority"`
		Destination string `json:"destination"`
		Mint        string `json:"mint"`
		Source      string `json:"source"`
		TokenAmount struct {
			Amount         string  `json:"amount"`
			Decimals       uint8   `json:"decimals"`
			UIAmount       float64 `json:"uiAmount"`
			UIAmountString string  `json:"uiAmountString"`
		} `json:"tokenAmount"`
	} `json:"info"`
	Type string `json:"type"`
}

func (p *Parser) processTransferCheck(instr solana.CompiledInstruction) *TransferCheck {

	amount := binary.LittleEndian.Uint64(instr.Data[1:9])

	transferData := &TransferCheck{
		Type: "transferChecked",
	}

	transferData.Info.Source = p.AllAccountKeys[instr.Accounts[0]].String()
	transferData.Info.Destination = p.AllAccountKeys[instr.Accounts[2]].String()
	transferData.Info.Mint = p.AllAccountKeys[instr.Accounts[1]].String()
	transferData.Info.Authority = p.AllAccountKeys[instr.Accounts[3]].String()

	transferData.Info.TokenAmount.Amount = fmt.Sprintf("%d", amount)
	transferData.Info.TokenAmount.Decimals = p.SplDecimalsMap[transferData.Info.Mint]
	uiAmount := float64(amount) / math.Pow10(int(transferData.Info.TokenAmount.Decimals))
	transferData.Info.TokenAmount.UIAmount = uiAmount
	transferData.Info.TokenAmount.UIAmountString = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.9f", uiAmount), "0"), ".")

	return transferData
}
