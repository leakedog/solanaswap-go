package solanaswapgo

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

func (p *Parser) InnerParseInstruction(instruction solana.CompiledInstruction, progID solana.PublicKey, index int) []SwapData {
	var swaps []SwapData
	for _, innerInstructionSet := range p.Tx.Meta.InnerInstructions {
		if innerInstructionSet.Index == uint16(index) {
			for _, innerInstruction := range innerInstructionSet.Instructions {
				programId := p.AllAccountKeys[innerInstruction.ProgramIDIndex]
				switch programId {
				case SPL_ASSOCIATED_TOKEN_ID, TOKEN_PROGRAM_ID:
				case ALDRIN_AMM_PROGRAM_ID, STABLE_SWAP_PROGRAM_ID, STABLE_SWAP_V2_PROGRAM_ID, SWAP_DEX_PROGRAM_ID, OPENBOOK_V2_PROGRAM_ID, PHOENIX_PROGRAM_ID,
					LIFINITY_V2_PROGRAM_ID, FLUXBEAM_PROGRAM_ID, RAYDIUM_V4_PROGRAM_ID, RAYDIUM_CPMM_PROGRAM_ID, RAYDIUM_AMM_PROGRAM_ID,
					RAYDIUM_AMM_LIQUIDITY_POOL_PROGRAM_ID, RAYDIUM_CONCENTRATED_LIQUIDITY_PROGRAM_ID, ORCA_PROGRAM_ID, ORCA_TOKEN_V2_PROGRAM_ID,
					METEORA_PROGRAM_ID, METEORA_POOLS_PROGRAM_ID:

					datas := p.processTransferSwapDexByProgID(index, programId)
					who := p.AllAccountKeys[0].String()
					in := SwapData{}
					out := SwapData{}
					for _, v := range datas {
						item := v.Data.(*TransferSwapData)
						if item.Authority == who && in.Data == nil {
							in = v
							continue
						}

						tokenAccount, _, err := solana.FindAssociatedTokenAddress(solana.MustPublicKeyFromBase58(who), solana.MustPublicKeyFromBase58(item.Mint))
						if err != nil {
							continue
						}
						if tokenAccount.String() == item.Destination {
							out = v
						}
					}

					if in.Data != nil && out.Data != nil {
						in.Type = out.Type
						return append(swaps, in, out)
					}

					return datas
				case MOONSHOT_PROGRAM_ID:
					return p.processMoonshotSwaps()
				case solana.MustPublicKeyFromBase58("1MooN32fuBBgApc8ujknKJw5sef3BVwPGgz3pto1BAh"):
					return p.processTransferSwapDex(index, "1MooN")
				case solana.MustPublicKeyFromBase58("AP51WLiiqTdbZfgyRMs35PsZpdmLuPDdHYmrB23pEtMU"):
					return p.processTransferSwapDex(index, RAYDIUM)
				case PUMP_FUN_PROGRAM_ID:
					return p.processPumpfunSwaps(index)
				default:
					if progID == OKX_PROGRAM_ID {
						fmt.Println("OKX", instruction.Data.String(), programId)
						swaps = append(swaps, []SwapData{
							{
								Type:   "OKX",
								Action: UNKNOWN.String(),
								Data: UnknownAction{
									BaseAction: BaseAction{
										ProgramID:       programId.String(),
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
	}

	// return swaps
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
