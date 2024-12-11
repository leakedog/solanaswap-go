package solanaswapgo

import (
	"encoding/binary"

	"github.com/gagliardetto/solana-go"
)

type TransferInfo struct {
	Amount      uint64 `json:"amount"`
	Authority   string `json:"authority"`
	Destination string `json:"destination"`
	Source      string `json:"source"`
}

type TransferData struct {
	Info     TransferInfo `json:"info"`
	Type     string       `json:"type"`
	Mint     string       `json:"mint"`
	Decimals uint8        `json:"decimals"`
}

type TokenInfo struct {
	Mint     string
	Decimals uint8
}

func (p *Parser) processRaydSwaps(instructionIndex int) []SwapData {
	var swaps []SwapData

	for _, innerInstructionSet := range p.Tx.Meta.InnerInstructions {
		if innerInstructionSet.Index == uint16(instructionIndex) {
			for _, innerInstruction := range innerInstructionSet.Instructions {
				switch {
				case p.isTransfer(innerInstruction):
					transfer := p.processTransfer(innerInstruction)
					if transfer != nil {
						swapData := SwapData{Type: RAYDIUM, Data: transfer}
						if p.isRaydiumSwapEventInstruction(p.TxInfo.Message.Instructions[instructionIndex]) {
							swapData.Action = "swap"
						}
						swaps = append(swaps, swapData)
					}

				case p.isTransferCheck(innerInstruction):
					transfer := p.processTransferCheck(innerInstruction)
					if transfer != nil {
						swapData := SwapData{Type: RAYDIUM, Data: transfer}
						if p.isRaydiumAddLiquidityEventInstruction(p.TxInfo.Message.Instructions[instructionIndex]) {
							swapData.Action = "add_liquidity"
						}
						swaps = append(swaps, swapData)
					}
				}
			}
		}
	}
	return swaps
}

func (p *Parser) processOrcaSwaps(instructionIndex int) []SwapData {
	var swaps []SwapData
	for _, innerInstructionSet := range p.Tx.Meta.InnerInstructions {
		if innerInstructionSet.Index == uint16(instructionIndex) {
			for _, innerInstruction := range innerInstructionSet.Instructions {
				if p.isTransfer(innerInstruction) {
					transfer := p.processTransfer(innerInstruction)
					if transfer != nil {
						swaps = append(swaps, SwapData{Type: ORCA, Data: transfer})
					}
				}
			}
		}
	}
	return swaps
}

func (p *Parser) processOkxSwaps(instructionIndex int) []SwapData {
	var swaps []SwapData

	for _, innerInstructionSet := range p.Tx.Meta.InnerInstructions {
		if innerInstructionSet.Index == uint16(instructionIndex) {
			for _, innerInstruction := range innerInstructionSet.Instructions {
				switch {
				case p.isTransfer(innerInstruction):
					transfer := p.processTransfer(innerInstruction)
					if transfer != nil {
						swapData := SwapData{Type: OKX, Data: transfer}
						swaps = append(swaps, swapData)
					}

				case p.isTransferCheck(innerInstruction):
					transfer := p.processTransferCheck(innerInstruction)
					if transfer != nil {
						swapData := SwapData{Type: OKX, Data: transfer}
						swaps = append(swaps, swapData)
					}
				}
			}
		}
	}
	return swaps
}

func (p *Parser) processTransfer(instr solana.CompiledInstruction) *TransferData {

	amount := binary.LittleEndian.Uint64(instr.Data[1:9])

	transferData := &TransferData{
		Info: TransferInfo{
			Amount:      amount,
			Source:      p.AllAccountKeys[instr.Accounts[0]].String(),
			Destination: p.AllAccountKeys[instr.Accounts[1]].String(),
			Authority:   p.AllAccountKeys[instr.Accounts[2]].String(),
		},
		Type:     "transfer",
		Mint:     p.SplTokenInfoMap[p.AllAccountKeys[instr.Accounts[1]].String()].Mint,
		Decimals: p.SplTokenInfoMap[p.AllAccountKeys[instr.Accounts[1]].String()].Decimals,
	}

	if transferData.Mint == "" {
		transferData.Mint = "Unknown"
	}

	return transferData
}

func (p *Parser) extractSPLTokenInfo() error {
	splTokenAddresses := make(map[string]TokenInfo)

	for _, accountInfo := range p.Tx.Meta.PostTokenBalances {
		if !accountInfo.Mint.IsZero() {
			accountKey := p.AllAccountKeys[accountInfo.AccountIndex].String()
			splTokenAddresses[accountKey] = TokenInfo{
				Mint:     accountInfo.Mint.String(),
				Decimals: accountInfo.UiTokenAmount.Decimals,
			}
		}
	}

	processInstruction := func(instr solana.CompiledInstruction) {
		if !p.AllAccountKeys[instr.ProgramIDIndex].Equals(solana.TokenProgramID) {
			return
		}

		if len(instr.Data) == 0 || (instr.Data[0] != 3 && instr.Data[0] != 12) {
			return
		}

		if len(instr.Accounts) < 3 {
			return
		}

		source := p.AllAccountKeys[instr.Accounts[0]].String()
		destination := p.AllAccountKeys[instr.Accounts[1]].String()

		if _, exists := splTokenAddresses[source]; !exists {
			splTokenAddresses[source] = TokenInfo{Mint: "", Decimals: 0}
		}
		if _, exists := splTokenAddresses[destination]; !exists {
			splTokenAddresses[destination] = TokenInfo{Mint: "", Decimals: 0}
		}
	}

	for _, instr := range p.TxInfo.Message.Instructions {
		processInstruction(instr)
	}
	for _, innerSet := range p.Tx.Meta.InnerInstructions {
		for _, instr := range innerSet.Instructions {
			processInstruction(instr)
		}
	}

	for account, info := range splTokenAddresses {
		if info.Mint == "" {
			splTokenAddresses[account] = TokenInfo{
				Mint:     NATIVE_SOL_MINT_PROGRAM_ID.String(),
				Decimals: 9, // Native SOL has 9 decimal places
			}
		}
	}

	p.SplTokenInfoMap = splTokenAddresses

	return nil
}
