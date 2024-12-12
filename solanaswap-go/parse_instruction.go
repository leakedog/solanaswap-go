package solanaswapgo

import (
	"github.com/gagliardetto/solana-go"
	"github.com/samber/lo"
)

func (p *Parser) SystemProgramInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {
	return nil
}

func (p *Parser) TokenProgramInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {
	return nil
}

func (p *Parser) ComputeBudgetInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {
	return nil
}

func (p *Parser) PumpfunInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {
	return nil
}

func (p *Parser) JupiterDCAInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {

	return nil
}

func (p *Parser) RaydiumInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {
	return nil
}

func (p *Parser) OkxInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {
	var swaps []SwapData
	for _, innerInstructionSet := range p.Tx.Meta.InnerInstructions {
		if innerInstructionSet.Index == uint16(index) {
			for _, innerInstruction := range innerInstructionSet.Instructions {
				programId := p.AllAccountKeys[innerInstruction.ProgramIDIndex]
				switch programId {
				case MOONSHOT_PROGRAM_ID:
					return p.processMoonshotSwaps()
				case RAYDIUM_V4_PROGRAM_ID,
					RAYDIUM_CPMM_PROGRAM_ID,
					RAYDIUM_AMM_PROGRAM_ID,
					RAYDIUM_CONCENTRATED_LIQUIDITY_PROGRAM_ID,
					solana.MustPublicKeyFromBase58("AP51WLiiqTdbZfgyRMs35PsZpdmLuPDdHYmrB23pEtMU"):
					return p.processRaydSwaps(index)
				case PUMP_FUN_PROGRAM_ID:
					return p.processPumpfunSwaps(index)
				case ORCA_PROGRAM_ID:
					return p.processOrcaSwaps(index)
				case METEORA_PROGRAM_ID, METEORA_POOLS_PROGRAM_ID:
					return p.processMeteoraSwaps(index)
				default:
					swaps = append(swaps, []SwapData{
						{
							Type:   OKX,
							Action: "Unknown",
							Data: UnknownAction{
								BaseAction: BaseAction{
									ProgramID:       progID.String(),
									ProgramName:     lo.Ternary(lo.HasKey(ProgramName, progID), ProgramName[progID].String(), "Unknown"),
									InstructionName: "Unknown",
								},
							},
						},
					}...)
				}
			}
		}
	}
	return swaps
}

func (p *Parser) JupiterInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {
	return nil
}

func (p *Parser) OrcaInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {
	return nil
}

func (p *Parser) MeteoraInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {
	return nil
}

func (p *Parser) MoonshotInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {
	return nil
}

func (p *Parser) TradingBotInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {
	return nil
}
