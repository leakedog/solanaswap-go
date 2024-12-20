package solanaswapgo

import (
	"bytes"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

// isTransfer checks if the instruction is a token transfer (Raydium, Orca)
func (p *Parser) isTransfer(instr solana.CompiledInstruction) bool {
	progID := p.AllAccountKeys[instr.ProgramIDIndex]

	if !progID.Equals(solana.TokenProgramID) {
		return false
	}

	if len(instr.Accounts) < 3 || len(instr.Data) < 9 {
		return false
	}

	if instr.Data[0] != 3 {
		return false
	}

	for i := 0; i < 3; i++ {
		if int(instr.Accounts[i]) >= len(p.AllAccountKeys) {
			return false
		}
	}

	return true
}

// isTransferCheck checks if the instruction is a token transfer check (Meteora)
func (p *Parser) isTransferCheck(instr solana.CompiledInstruction) bool {
	progID := p.AllAccountKeys[instr.ProgramIDIndex]

	if !progID.Equals(solana.TokenProgramID) && !progID.Equals(solana.Token2022ProgramID) {
		return false
	}

	if len(instr.Accounts) < 4 || len(instr.Data) < 9 {
		return false
	}

	if instr.Data[0] != 12 {
		return false
	}

	for i := 0; i < 4; i++ {
		if int(instr.Accounts[i]) >= len(p.AllAccountKeys) {
			return false
		}
	}

	return true
}

var RaydiumSwapEventDiscriminator = CalculateDiscriminator("global:swap")

func (p *Parser) isRaydiumSwapEventInstruction(inst solana.CompiledInstruction) bool {
	if !p.AllAccountKeys[inst.ProgramIDIndex].Equals(RAYDIUM_V4_PROGRAM_ID) || len(inst.Data) < 8 {
		return false
	}

	decodedBytes, err := base58.Decode(inst.Data.String())
	if err != nil {
		return false
	}
	return bytes.Equal(decodedBytes[:8], RaydiumSwapEventDiscriminator[:])
}

type LiquidityEventType string

const (
	AddLiquidity    LiquidityEventType = "add_liquidity"
	RemoveLiquidity LiquidityEventType = "remove_liquidity"
	NoLiquidity     LiquidityEventType = "unknown"
)

func (s LiquidityEventType) String() string {
	return string(s)
}

var RaydiumAddLiquidityEventDiscriminator = [8]byte{133, 29, 89, 223, 69, 238, 176, 10}
var OrcaRemoveLiquidityEventDiscriminator = [8]byte{164, 152, 207, 99, 30, 186, 19, 182}
var OrcaRemoveLiquidityEventDiscriminator2 = [8]byte{160, 38, 208, 111, 104, 91, 44, 1}
var MeteoraRemoveLiquidityEventDiscriminator = [8]byte{26, 82, 102, 152, 240, 74, 105, 26}
var MeteoraAddLiquidityEventDiscriminator = [8]byte{7, 3, 150, 127, 148, 40, 61, 200}

func (p *Parser) isLiquidityEventInstruction(inst solana.CompiledInstruction) LiquidityEventType {
	if len(inst.Data) < 8 {
		return NoLiquidity
	}

	decodedBytes, err := base58.Decode(inst.Data.String())
	if err != nil {
		return NoLiquidity
	}

	discriminator := *(*[8]byte)(decodedBytes[:8])
	switch discriminator {
	case RaydiumAddLiquidityEventDiscriminator, MeteoraAddLiquidityEventDiscriminator:
		return AddLiquidity
	case OrcaRemoveLiquidityEventDiscriminator, OrcaRemoveLiquidityEventDiscriminator2, MeteoraRemoveLiquidityEventDiscriminator:
		return RemoveLiquidity
	default:
		// fmt.Println(discriminator, CalculateDiscriminator("global:add_liquidity_by_strategy"))
		return NoLiquidity
	}
}

func (p *Parser) isPumpFunTradeEventInstruction(inst solana.CompiledInstruction) bool {
	if !p.AllAccountKeys[inst.ProgramIDIndex].Equals(PUMP_FUN_PROGRAM_ID) || len(inst.Data) < 16 {
		return false
	}
	decodedBytes, err := base58.Decode(inst.Data.String())
	if err != nil {
		return false
	}
	return bytes.Equal(decodedBytes[:16], PumpfunTradeEventDiscriminator[:])
}

func (p *Parser) isPumpFunCreateEventInstruction(inst solana.CompiledInstruction) bool {
	if !p.AllAccountKeys[inst.ProgramIDIndex].Equals(PUMP_FUN_PROGRAM_ID) || len(inst.Data) < 16 {
		return false
	}
	decodedBytes, err := base58.Decode(inst.Data.String())
	if err != nil {
		return false
	}
	return bytes.Equal(decodedBytes[:16], PumpfunCreateEventDiscriminator[:])
}

func (p *Parser) isPumpFunCompleteEventInstruction(inst solana.CompiledInstruction) bool {
	if !p.AllAccountKeys[inst.ProgramIDIndex].Equals(PUMP_FUN_PROGRAM_ID) || len(inst.Data) < 16 {
		return false
	}
	decodedBytes, err := base58.Decode(inst.Data.String())

	if err != nil {
		return false
	}
	return bytes.Equal(decodedBytes[8:16], PumpfunCompleteEventDiscriminator[:])
}

func (p *Parser) isJupiterRouteEventInstruction(inst solana.CompiledInstruction) bool {
	if !p.AllAccountKeys[inst.ProgramIDIndex].Equals(JUPITER_PROGRAM_ID) || len(inst.Data) < 16 {
		return false
	}
	decodedBytes, err := base58.Decode(inst.Data.String())
	if err != nil {
		return false
	}
	return bytes.Equal(decodedBytes[:16], JupiterRouteEventDiscriminator[:])
}
