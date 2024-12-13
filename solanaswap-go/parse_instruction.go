package solanaswapgo

import (
	"github.com/gagliardetto/solana-go"
	"github.com/samber/lo"
)

func (p *Parser) OkxInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {
	var swaps []SwapData
	for _, innerInstructionSet := range p.Tx.Meta.InnerInstructions {
		if innerInstructionSet.Index == uint16(index) {
			for _, innerInstruction := range innerInstructionSet.Instructions {
				programId := p.AllAccountKeys[innerInstruction.ProgramIDIndex]
				switch programId {
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
					return p.processTransferSwapDex(index, RAYDIUM)
				case PUMP_FUN_PROGRAM_ID:
					return p.processPumpfunSwaps(index)
				case ORCA_PROGRAM_ID, ORCA_TOKEN_V2_PROGRAM_ID:
					return p.processTransferSwapDex(index, ORCA)
				case METEORA_PROGRAM_ID, METEORA_POOLS_PROGRAM_ID:
					return p.processTransferSwapDex(index, METEORA)
					return p.processTransferSwapDex(index, METEORA)
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
