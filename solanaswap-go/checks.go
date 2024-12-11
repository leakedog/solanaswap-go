package solanaswapgo

import (
	"bytes"
	"fmt"

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

var RaydiumAddLiquidityEventDiscriminator = [8]byte{133, 29, 89, 223, 69, 238, 176, 10}

func (p *Parser) isRaydiumAddLiquidityEventInstruction(inst solana.CompiledInstruction) bool {
	if !p.AllAccountKeys[inst.ProgramIDIndex].Equals(RAYDIUM_CONCENTRATED_LIQUIDITY_PROGRAM_ID) || len(inst.Data) < 8 {
		return false
	}

	decodedBytes, err := base58.Decode(inst.Data.String())
	if err != nil {
		return false
	}
	return bytes.Equal(decodedBytes[:8], RaydiumAddLiquidityEventDiscriminator[:])
}

var OrcaRemoveLiquidityEventDiscriminator = [8]byte{164, 152, 207, 99, 30, 186, 19, 182}

func (p *Parser) isOrcaRemoveLiquidityEventInstruction(inst solana.CompiledInstruction) bool {
	if !p.AllAccountKeys[inst.ProgramIDIndex].Equals(ORCA_PROGRAM_ID) || len(inst.Data) < 8 {
		return false
	}

	decodedBytes, err := base58.Decode(inst.Data.String())
	if err != nil {
		return false
	}

	return bytes.Equal(decodedBytes[:8], OrcaRemoveLiquidityEventDiscriminator[:])
}

var MeteoraRemoveLiquidityEventDiscriminator = [8]byte{26, 82, 102, 152, 240, 74, 105, 26}

func (p *Parser) isMeteoraRemoveLiquidityEventInstruction(inst solana.CompiledInstruction) bool {
	if !p.AllAccountKeys[inst.ProgramIDIndex].Equals(METEORA_PROGRAM_ID) || len(inst.Data) < 8 {
		return false
	}
	decodedBytes, err := base58.Decode(inst.Data.String())
	if err != nil {
		return false
	}
	fmt.Println(decodedBytes[:8])
	return bytes.Equal(decodedBytes[:8], MeteoraRemoveLiquidityEventDiscriminator[:])
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
