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
	switch {
	case p.isTransfer(instruction):
		return p.processRaydSwaps(index)
	default:
		// fmt.Println("Jupiter Unknown discriminator", discriminator)
		return []SwapData{
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
		}
	}
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
