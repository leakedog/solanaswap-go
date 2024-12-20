package solanaswapgo

import (
	"fmt"

	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

var (
	PumpfunTradeEventDiscriminator    = [16]byte{228, 69, 165, 46, 81, 203, 154, 29, 189, 219, 127, 211, 78, 230, 97, 238}
	PumpfunCreateEventDiscriminator   = [16]byte{228, 69, 165, 46, 81, 203, 154, 29, 27, 114, 169, 77, 222, 235, 99, 118}
	PumpfunCompleteEventDiscriminator = [8]byte{95, 114, 97, 156, 212, 46, 152, 8}
)

type PumpfunTradeEvent struct {
	Mint                 solana.PublicKey
	SolAmount            uint64
	TokenAmount          uint64
	IsBuy                bool
	User                 solana.PublicKey
	Timestamp            int64
	VirtualSolReserves   uint64
	VirtualTokenReserves uint64
}

type PumpfunCreateEvent struct {
	Name         string
	Symbol       string
	Uri          string
	Mint         solana.PublicKey
	BondingCurve solana.PublicKey
	User         solana.PublicKey
}

type PumpfunCompleteEvent struct {
	User         solana.PublicKey
	Mint         solana.PublicKey
	BondingCurve solana.PublicKey
	Timestamp    int64
}

func (p *Parser) processPumpfunSwaps(instructionIndex int) []SwapData {
	var swaps []SwapData
	for _, innerInstructionSet := range p.Tx.Meta.InnerInstructions {
		if innerInstructionSet.Index == uint16(instructionIndex) {
			for _, innerInstruction := range innerInstructionSet.Instructions {
				switch {
				case p.isPumpFunTradeEventInstruction(innerInstruction):
					trade, err := p.parsePumpfunTradeEventInstruction(innerInstruction)
					if err != nil {
						p.Log.Errorf("error processing Pumpfun trade event: %s", err)
					}
					if trade != nil {
						swaps = append(swaps, SwapData{Type: PUMP_FUN, Data: trade})
					}
				case p.isPumpFunCreateEventInstruction(innerInstruction):
					create, err := p.parsePumpfunCreateEventInstruction(innerInstruction)
					if err != nil {
						p.Log.Errorf("error processing Pumpfun create event: %s", err)
					}
					if create != nil {
						swaps = append(swaps, SwapData{Type: PUMP_FUN, Data: create})
					}
				case p.isPumpFunCompleteEventInstruction(innerInstruction):
					complete, err := p.parsePumpfunCompleteEventInstruction(innerInstruction)
					if err != nil {
						p.Log.Errorf("error processing Pumpfun complete event: %s", err)
					}
					if complete != nil {
						swaps = append(swaps, SwapData{Type: PUMP_FUN, Data: complete})
					}
				}

			}
		}
	}
	return swaps
}

func (p *Parser) parsePumpfunTradeEventInstruction(instruction solana.CompiledInstruction) (*PumpfunTradeEvent, error) {

	decodedBytes, err := base58.Decode(instruction.Data.String())
	if err != nil {
		return nil, fmt.Errorf("error decoding instruction data: %s", err)
	}
	decoder := ag_binary.NewBorshDecoder(decodedBytes[16:])

	return handlePumpfunTradeEvent(decoder)
}

func handlePumpfunTradeEvent(decoder *ag_binary.Decoder) (*PumpfunTradeEvent, error) {
	var trade PumpfunTradeEvent
	if err := decoder.Decode(&trade); err != nil {
		return nil, fmt.Errorf("error unmarshaling TradeEvent: %s", err)
	}

	return &trade, nil
}

func (p *Parser) parsePumpfunCreateEventInstruction(instruction solana.CompiledInstruction) (*PumpfunCreateEvent, error) {
	decodedBytes, err := base58.Decode(instruction.Data.String())
	if err != nil {
		return nil, fmt.Errorf("error decoding instruction data: %s", err)
	}
	decoder := ag_binary.NewBorshDecoder(decodedBytes[16:])
	return handlePumpfunCreateEvent(decoder)
}

func handlePumpfunCreateEvent(decoder *ag_binary.Decoder) (*PumpfunCreateEvent, error) {
	var create PumpfunCreateEvent
	if err := decoder.Decode(&create); err != nil {
		return nil, fmt.Errorf("error unmarshaling CreateEvent: %s", err)
	}

	return &create, nil
}

func (p *Parser) parsePumpfunCompleteEventInstruction(instruction solana.CompiledInstruction) (*PumpfunCompleteEvent, error) {
	decodedBytes, err := base58.Decode(instruction.Data.String())
	if err != nil {
		return nil, fmt.Errorf("error decoding instruction data: %s", err)
	}
	decoder := ag_binary.NewBorshDecoder(decodedBytes[16:])
	return handlePumpfunCompleteEvent(decoder)
}

func handlePumpfunCompleteEvent(decoder *ag_binary.Decoder) (*PumpfunCompleteEvent, error) {
	var complete PumpfunCompleteEvent
	if err := decoder.Decode(&complete); err != nil {
		return nil, fmt.Errorf("error unmarshaling CompleteEvent: %s", err)
	}

	return &complete, nil
}
