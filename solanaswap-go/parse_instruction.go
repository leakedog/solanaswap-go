package solanaswapgo

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

func (p *Parser) OkxInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {
	var swaps []SwapData
	for _, innerInstructionSet := range p.Tx.Meta.InnerInstructions {
		if innerInstructionSet.Index == uint16(index) {
			for _, innerInstruction := range innerInstructionSet.Instructions {
				programId := p.AllAccountKeys[innerInstruction.ProgramIDIndex]
				switch programId {
				case SWAP_DEX_PROGRAM_ID:
					return p.processTransferSwapDex(index, SWAP_DEX)
				case OPENBOOK_V2_PROGRAM_ID:
					return p.processTransferSwapDex(index, OPENBOOK)
				case PHOENIX_PROGRAM_ID:
					return p.processTransferSwapDex(index, PHOENIX)
				case LIFINITY_V2_PROGRAM_ID:
					return p.processTransferSwapDex(index, LIFINITY)
				case FLUXBEAM_PROGRAM_ID:
					return p.processTransferSwapDex(index, FLUXBEAM)
				case MOONSHOT_PROGRAM_ID:
					return p.processMoonshotSwaps()
				case RAYDIUM_V4_PROGRAM_ID,
					RAYDIUM_CPMM_PROGRAM_ID,
					RAYDIUM_AMM_PROGRAM_ID,
					RAYDIUM_AMM_LIQUIDITY_POOL_PROGRAM_ID,
					RAYDIUM_CONCENTRATED_LIQUIDITY_PROGRAM_ID,
					solana.MustPublicKeyFromBase58("AP51WLiiqTdbZfgyRMs35PsZpdmLuPDdHYmrB23pEtMU"):
					return p.processTransferSwapDex(index, RAYDIUM)
				case PUMP_FUN_PROGRAM_ID:
					return p.processPumpfunSwaps(index)
				case ORCA_PROGRAM_ID, ORCA_TOKEN_V2_PROGRAM_ID:
					return p.processTransferSwapDex(index, ORCA)
				case METEORA_PROGRAM_ID, METEORA_POOLS_PROGRAM_ID:
					return p.processTransferSwapDex(index, METEORA)
				default:
					swaps = append(swaps, []SwapData{
						{
							Type:   OKX,
							Action: "Unknown",
							Data: UnknownAction{
								BaseAction: BaseAction{
									ProgramID:       progID.String(),
									ProgramName:     ProgramName(progID).String(),
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

func (p *Parser) MoonshotInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {
	var swaps []SwapData

	data := instruction.Data
	decode, err := base58.Decode(data.String())
	if err != nil {
		return nil
	}
	discriminator := *(*[8]byte)(decode[:8])

	//CalculateDiscriminator("global:sell")
	switch discriminator {
	case MOONSHOT_SELL_INSTRUCTION, MOONSHOT_BUY_INSTRUCTION:
		return p.processMoonshotSwaps()
	case MOONSHOT_CREATE_TOKEN:
		return p.processMoonshotCreateToken(progID, instruction)
	default:
		fmt.Println("MoonshotInstruction", discriminator, CalculateDiscriminator(":token_mint"))
	}

	return swaps
}
