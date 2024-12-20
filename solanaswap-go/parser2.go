package solanaswapgo

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func (p *Parser) NewTxParser() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("========================ERROR BREAK========================")
			fmt.Println(r)
			spew.Dump(p.SwapData, "========================ERROR BREAK========================")
		}
	}()
	for i, outerInstruction := range p.TxInfo.Message.Instructions {
		progID := p.AllAccountKeys[outerInstruction.ProgramIDIndex]
		switch progID {
		case solana.VoteProgramID, solana.SystemProgramID, solana.ComputeBudget:
		case TOKEN_PROGRAM_ID, SPL_ASSOCIATED_TOKEN_ID:
		case ALDRIN_AMM_PROGRAM_ID, STABLE_SWAP_PROGRAM_ID, SWAP_DEX_PROGRAM_ID, PHOENIX_PROGRAM_ID, LIFINITY_V2_PROGRAM_ID, FLUXBEAM_PROGRAM_ID, RAYDIUM_V4_PROGRAM_ID, RAYDIUM_CPMM_PROGRAM_ID, RAYDIUM_AMM_PROGRAM_ID, RAYDIUM_AMM_LIQUIDITY_POOL_PROGRAM_ID, RAYDIUM_CONCENTRATED_LIQUIDITY_PROGRAM_ID, ORCA_PROGRAM_ID, ORCA_TOKEN_V2_PROGRAM_ID, METEORA_PROGRAM_ID, METEORA_POOLS_PROGRAM_ID:
			p.programParseTo(p.processTransferSwapDexByProgID(i, progID), progID)
		case JUPITER_PROGRAM_ID:
			p.programParseTo(p.processJupiterSwaps(i), progID)
		case MOONSHOT_PROGRAM_ID:
			p.programParseTo(p.MoonshotInstruction(outerInstruction, progID, i), progID)
		case BANANA_GUN_PROGRAM_ID, MINTECH_PROGRAM_ID, BLOOM_PROGRAM_ID, MAESTRO_PROGRAM_ID:
			// Check inner instructions to determine which swap protocol is being used
			p.programParseTo(p.processTradingBotSwaps(i), progID)
		case OPENBOOK_V2_PROGRAM_ID:
			decodedBytes, err := base58.Decode(outerInstruction.Data.String())
			if err != nil {
				continue
			}

			swapDiscriminator := []byte{3, 44, 71, 3, 26, 199, 203, 85} //[3 44 71 3 26 199 203 85]
			if bytes.Equal(decodedBytes[:8], swapDiscriminator) {
				p.programParseTo(p.processTransferSwapDex(i, OPENBOOK), progID)
			}

		case solana.MustPublicKeyFromBase58("AP51WLiiqTdbZfgyRMs35PsZpdmLuPDdHYmrB23pEtMU"):
			p.programParseTo(p.processTransferSwapDex(i, RAYDIUM), progID)
		case OKX_PROGRAM_ID:
			p.programParseTo(p.InnerParseInstruction(outerInstruction, progID, i), progID)
		case PUMP_FUN_PROGRAM_ID:
			p.programParseTo(p.processPumpfunSwaps(i), progID)
		default:
			p.programParseTo(p.InnerParseInstruction(outerInstruction, progID, i), progID)
			// fmt.Println("UNKNOWN Program", progID, p.TxInfo.Signatures[0].String())
		}
	}
}

func (p *Parser) programParseTo(datas []SwapData, progID solana.PublicKey) {
	if len(datas) == 0 {
		return
	}
	p.SwapData = append(p.SwapData, datas...)
	p.parseDataToAction(datas, progID)
}

func (p *Parser) parseDataToAction(datas []SwapData, progID solana.PublicKey) {
	switch progID {
	case TOKEN_PROGRAM_ID:
		p.Actions = append(p.Actions, NewUnknownAction(progID, p.TxInfo.Signatures[0].String(), nil))
	case JUPITER_PROGRAM_ID:
		data := datas[0].Data.(*JupiterSwapEventData)
		last := lo.LastOrEmpty(datas).Data.(*JupiterSwapEventData)
		p.Actions = append(p.Actions, CommonSwapAction{
			BaseAction: BaseAction{
				ProgramID:       progID.String(),
				ProgramName:     string(ProgramName(data.Amm)),
				InstructionName: string(ProgramName(data.Amm)),
				Signature:       p.TxInfo.Signatures[0].String(),
			},
			Who:               p.AllAccountKeys[0].String(),
			FromToken:         data.InputMint.String(),
			FromTokenAmount:   data.InputAmount,
			FromTokenDecimals: p.SplDecimalsMap[data.InputMint.String()],
			ToToken:           last.OutputMint.String(),
			ToTokenAmount:     last.OutputAmount,
			ToTokenDecimals:   p.SplDecimalsMap[last.OutputMint.String()],
		})
	case METEORA_POOLS_PROGRAM_ID:
		p.parseMeteoraPoolsSwapData(progID, datas)
	case MOONSHOT_PROGRAM_ID:
		p.parseMoonshotSwapData(progID, datas)
	case solana.MustPublicKeyFromBase58("AP51WLiiqTdbZfgyRMs35PsZpdmLuPDdHYmrB23pEtMU"):
		p.parseOneTransferSwapData(progID, datas)
	default:
		switch datas[0].Type {
		case PUMP_FUN:
			p.parsePumpfunSwapData(progID, datas)
		default:
			p.parseGroupTransferSwapData(progID, datas)
		}
		// p.Actions = append(p.Actions, NewUnknownAction(progID, p.TxInfo.Signatures[0].String(), fmt.Errorf("unknown parser action, %s", progID.String())))
	}

}

func (p *Parser) parseMoonshotSwapData(progID solana.PublicKey, swapDatas []SwapData) {
	for _, data := range swapDatas {
		switch v := data.Data.(type) {
		case *MoonshotCreateTokenEvent:
			p.Actions = append(p.Actions, CommonDataAction{
				BaseAction: BaseAction{
					ProgramID:       progID.String(),
					ProgramName:     "Moonshot",
					InstructionName: "CreateToken",
					Signature:       p.TxInfo.Signatures[0].String(),
				},
				Data: v,
			})
		case *MoonshotTradeInstructionWithMint:
			if v.TradeType == TradeTypeBuy {
				p.Actions = append(p.Actions, CommonSwapAction{
					BaseAction: BaseAction{
						ProgramID:       progID.String(),
						ProgramName:     "Moonshot",
						InstructionName: "Swap",
						Signature:       p.TxInfo.Signatures[0].String(),
					},
					Who:               p.AllAccountKeys[0].String(),
					FromToken:         NATIVE_SOL_MINT_PROGRAM_ID.String(),
					FromTokenAmount:   v.CollateralAmount,
					FromTokenDecimals: 9,
					ToToken:           v.Mint.String(),
					ToTokenAmount:     v.TokenAmount,
					ToTokenDecimals:   p.SplDecimalsMap[v.Mint.String()],
				})
			} else if v.TradeType == TradeTypeSell {
				p.Actions = append(p.Actions, CommonSwapAction{
					BaseAction: BaseAction{
						ProgramID:       progID.String(),
						ProgramName:     "Moonshot",
						InstructionName: "Swap",
						Signature:       p.TxInfo.Signatures[0].String(),
					},
					Who:               p.AllAccountKeys[0].String(),
					FromToken:         v.Mint.String(),
					FromTokenAmount:   v.TokenAmount,
					FromTokenDecimals: p.SplDecimalsMap[v.Mint.String()],
					ToToken:           NATIVE_SOL_MINT_PROGRAM_ID.String(),
					ToTokenAmount:     v.CollateralAmount,
					ToTokenDecimals:   9,
				})
			}
		}
	}
}

func (p *Parser) parsePumpfunSwapData(progID solana.PublicKey, swapDatas []SwapData) {
	for _, data := range swapDatas {
		switch v := data.Data.(type) {
		case *PumpfunCreateEvent:
			p.Actions = append(p.Actions, CommonDataAction{
				BaseAction: BaseAction{
					ProgramID:       progID.String(),
					ProgramName:     "PumpFun",
					InstructionName: "Create",
					Signature:       p.TxInfo.Signatures[0].String(),
				},
				Data: v,
			})
		case *PumpfunCompleteEvent:
			p.Actions = append(p.Actions, CommonDataAction{
				BaseAction: BaseAction{
					ProgramID:       progID.String(),
					ProgramName:     "PumpFun",
					InstructionName: "Complete",
					Signature:       p.TxInfo.Signatures[0].String(),
				},
				Data: v,
			})
		case *PumpfunTradeEvent:
			action := CommonSwapAction{
				BaseAction: BaseAction{
					ProgramID:       progID.String(),
					ProgramName:     "PumpFun",
					InstructionName: "Swap",
					Signature:       p.TxInfo.Signatures[0].String(),
				},
				Who: p.AllAccountKeys[0].String(),
			}
			if v.IsBuy {
				action.FromToken = solana.SolMint.String()
				action.FromTokenAmount = v.SolAmount
				action.FromTokenDecimals = 9
				action.ToToken = v.Mint.String()
				action.ToTokenAmount = v.TokenAmount
				action.ToTokenDecimals = p.SplDecimalsMap[v.Mint.String()]
			} else {
				action.FromToken = v.Mint.String()
				action.FromTokenAmount = v.TokenAmount
				action.FromTokenDecimals = p.SplDecimalsMap[v.Mint.String()]
				action.ToToken = solana.SolMint.String()
				action.ToTokenAmount = v.SolAmount
				action.ToTokenDecimals = 9
			}
			p.Actions = append(p.Actions, action)
		}
	}

}

func (p *Parser) parseGroupTransferSwapData(progID solana.PublicKey, swapDatas []SwapData) {
	if len(swapDatas) == 0 {
		return
	}

	var resultGroup [][2]SwapData
	if len(swapDatas) == 1 {
		p.formatTransferData(swapDatas[0], swapDatas[0], progID, "OnlyTransfer")
		return
	}

	for i := 0; i < len(swapDatas)-1; i += 2 {
		resultGroup = append(resultGroup, [2]SwapData{swapDatas[i], swapDatas[i+1]})
	}

	for _, v := range resultGroup {
		in := v[0]
		out := v[1]
		if reflect.TypeOf(in.Data) == reflect.TypeOf(out.Data) {
			p.formatTransferData(in, out, progID)
		}
	}
}

func (p *Parser) parseOneTransferSwapData(progID solana.PublicKey, swapDatas []SwapData) {
	in := swapDatas[0]
	out := lo.LastOrEmpty(swapDatas)
	p.formatTransferData(in, out, progID, "Unknown Group Swap")
}

func (p *Parser) formatTransferData(in, out SwapData, progID solana.PublicKey, instructionName ...string) {
	who := p.AllAccountKeys[0].String()
	var action Action
	baseAction := BaseAction{
		ProgramID:       progID.String(),
		ProgramName:     in.Type.String(),
		InstructionName: in.Type.String(),
		Signature:       p.TxInfo.Signatures[0].String(),
	}
	if in.Action == "add_liquidity" {
		baseAction.InstructionName = "AddLiquidity"
		in := in.Data.(*TransferSwapData)
		out := out.Data.(*TransferSwapData)
		action = CommonAddLiquidityAction{
			BaseAction:     baseAction,
			Who:            who,
			Token1:         in.Mint,
			Token1Amount:   cast.ToUint64(in.Amount),
			Token1Decimals: p.SplDecimalsMap[in.Mint],
			Token2:         out.Mint,
			Token2Amount:   cast.ToUint64(out.Amount),
			Token2Decimals: p.SplDecimalsMap[out.Mint],
		}

	} else if in.Action == "remove_liquidity" {
		baseAction.InstructionName = "RemoveLiquidity"
		in := in.Data.(*TransferSwapData)
		out := out.Data.(*TransferSwapData)
		action = CommonRemoveLiquidityAction{
			BaseAction:     baseAction,
			Who:            who,
			Token1:         in.Mint,
			Token1Amount:   cast.ToUint64(in.Amount),
			Token1Decimals: p.SplDecimalsMap[in.Mint],
			Token2:         out.Mint,
			Token2Amount:   cast.ToUint64(out.Amount),
			Token2Decimals: p.SplDecimalsMap[out.Mint],
		}

	} else if in == out {
		if len(instructionName) > 0 {
			baseAction.InstructionName = instructionName[0]
		}
		p.Actions = append(p.Actions, NewCommonDataAction(progID, p.TxInfo.Signatures[0].String(), in.Data, baseAction.InstructionName))
		return
	} else {
		if len(instructionName) > 0 {
			baseAction.InstructionName = instructionName[0]
		}
		in := in.Data.(*TransferSwapData)
		out := out.Data.(*TransferSwapData)
		action = CommonSwapAction{
			BaseAction:        baseAction,
			Who:               who,
			FromToken:         in.Mint,
			FromTokenAmount:   in.Amount,
			FromTokenDecimals: p.SplDecimalsMap[in.Mint],
			ToToken:           out.Mint,
			ToTokenAmount:     out.Amount,
			ToTokenDecimals:   p.SplDecimalsMap[out.Mint],
		}
	}
	p.Actions = append(p.Actions, action)
}

// METEORA_POOLS
func (p *Parser) parseMeteoraPoolsSwapData(progID solana.PublicKey, swapDatas []SwapData) {
	if len(swapDatas)%2 < 2 {
		return
	}

	action := CommonSwapAction{
		BaseAction: BaseAction{
			ProgramID:       progID.String(),
			ProgramName:     "Meteora",
			InstructionName: "Meteora",
			Signature:       p.TxInfo.Signatures[0].String(),
		},
		Who: p.AllAccountKeys[0].String(),
	}
	for i, swapData := range swapDatas {
		switch swapData.Data.(type) {
		case *TransferSwapData: // Meteora Pools
			swapData := swapData.Data.(*TransferSwapData)
			if i == 0 {
				action.FromToken = swapData.Mint
				action.FromTokenAmount = swapData.Amount
				action.FromTokenDecimals = swapData.Decimals
			} else {
				if swapData.Authority == p.AllAccountKeys[0].String() && swapData.Mint == action.FromToken {
					action.FromTokenAmount += swapData.Amount
				}
				action.ToToken = swapData.Mint
				action.ToTokenAmount = swapData.Amount
				action.ToTokenDecimals = swapData.Decimals
			}
		}
	}
	p.Actions = append(p.Actions, action)
}
