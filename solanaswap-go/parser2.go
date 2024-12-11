package solanaswapgo

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
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
		switch {
		case strings.Contains(progID.String(), "ComputeBudget"):
			// p.programParseTo(p.ComputeBudgetInstruction(outerInstruction, i), progID)
		case strings.Contains(progID.String(), "11111111111111111111111111111111"):
			// p.programParseTo(p.SystemProgramInstruction(outerInstruction, i), progID)
		case progID.Equals(JUPITER_PROGRAM_ID):
			p.programParseTo(p.processJupiterSwaps(i), progID)
		case progID.Equals(MOONSHOT_PROGRAM_ID):
			p.programParseTo(p.processMoonshotSwaps(), progID)
		case progID.Equals(BANANA_GUN_PROGRAM_ID) ||
			progID.Equals(MINTECH_PROGRAM_ID) ||
			progID.Equals(BLOOM_PROGRAM_ID) ||
			progID.Equals(MAESTRO_PROGRAM_ID):
			// Check inner instructions to determine which swap protocol is being used
			p.programParseTo(p.processTradingBotSwaps(i), progID)

		case progID.Equals(RAYDIUM_V4_PROGRAM_ID) ||
			progID.Equals(RAYDIUM_CPMM_PROGRAM_ID) ||
			progID.Equals(RAYDIUM_AMM_PROGRAM_ID) ||
			progID.Equals(RAYDIUM_CONCENTRATED_LIQUIDITY_PROGRAM_ID) ||
			progID.Equals(solana.MustPublicKeyFromBase58("AP51WLiiqTdbZfgyRMs35PsZpdmLuPDdHYmrB23pEtMU")):
			p.programParseTo(p.processRaydSwaps(i), progID)
		case progID.Equals(OKX_PROGRAM_ID):
			// p.programParseTo(p.OkxInstruction(outerInstruction, progID, i), progID)
			p.programParseTo(p.processOkxSwaps(i), progID)
		case progID.Equals(ORCA_PROGRAM_ID):
			p.programParseTo(p.processOrcaSwaps(i), progID)
		case progID.Equals(METEORA_PROGRAM_ID) || progID.Equals(METEORA_POOLS_PROGRAM_ID):
			p.programParseTo(p.processMeteoraSwaps(i), progID)
		case progID.Equals(PUMP_FUN_PROGRAM_ID) ||
			progID.Equals(solana.MustPublicKeyFromBase58("BSfD6SHZigAfDWSjzD5Q41jw8LmKwtmjskPH9XW1mrRW")): // PumpFun
			p.programParseTo(p.processPumpfunSwaps(i), progID)
		default:
			p.parseDataToAction([]SwapData{
				{
					Type: UNKNOWN,
					Data: "UNKNOWN Program " + progID.String(),
				},
			}, progID)
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
	// spew.Dump("parseDataToAction", progID.String(), datas, "====")

	switch progID {
	case TOKEN_PROGRAM_ID:
		p.Actions = append(p.Actions, NewUnknownAction(progID, p.TxInfo.Signatures[0].String(), nil))
	case JUPITER_PROGRAM_ID:
		data := datas[0].Data.(*JupiterSwapEventData)
		last := lo.LastOrEmpty(datas).Data.(*JupiterSwapEventData)
		p.Actions = append(p.Actions, CommonSwapAction{
			BaseAction: BaseAction{
				ProgramID:       progID.String(),
				ProgramName:     string(datas[0].Type),
				InstructionName: "Swap",
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

	case PUMP_FUN_PROGRAM_ID:
		switch v := datas[0].Data.(type) {
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
	case solana.MustPublicKeyFromBase58("AP51WLiiqTdbZfgyRMs35PsZpdmLuPDdHYmrB23pEtMU"):
		p.parseOneTransferSwapData(progID, datas)
	case METEORA_PROGRAM_ID, RAYDIUM_AMM_PROGRAM_ID, RAYDIUM_CPMM_PROGRAM_ID, RAYDIUM_CONCENTRATED_LIQUIDITY_PROGRAM_ID:
		p.parseGroupTransferSwapData(progID, datas)
	case RAYDIUM_V4_PROGRAM_ID, OKX_PROGRAM_ID, ORCA_PROGRAM_ID, BANANA_GUN_PROGRAM_ID, MAESTRO_PROGRAM_ID,
		METEORA_POOLS_PROGRAM_ID:
		p.parseGroupTransferSwapData(progID, datas)
	case MOONSHOT_PROGRAM_ID:
		data := datas[0].Data.(*MoonshotTradeInstructionWithMint)
		switch data.TradeType {
		case TradeTypeBuy: // BUY
			p.Actions = append(p.Actions, CommonSwapAction{
				BaseAction: BaseAction{
					ProgramID:       progID.String(),
					ProgramName:     "Moonshot",
					InstructionName: "Swap",
					Signature:       p.TxInfo.Signatures[0].String(),
				},
				Who:               p.AllAccountKeys[0].String(),
				FromToken:         NATIVE_SOL_MINT_PROGRAM_ID.String(),
				FromTokenAmount:   data.CollateralAmount,
				FromTokenDecimals: 9,
				ToToken:           data.Mint.String(),
				ToTokenAmount:     data.TokenAmount,
				ToTokenDecimals:   p.SplDecimalsMap[data.Mint.String()],
			})

		case TradeTypeSell: // SELL
			p.Actions = append(p.Actions, CommonSwapAction{
				BaseAction: BaseAction{
					ProgramID:       progID.String(),
					ProgramName:     "Moonshot",
					InstructionName: "Swap",
					Signature:       p.TxInfo.Signatures[0].String(),
				},
				Who:               p.AllAccountKeys[0].String(),
				FromToken:         data.Mint.String(),
				FromTokenAmount:   data.TokenAmount,
				FromTokenDecimals: p.SplDecimalsMap[data.Mint.String()],
				ToToken:           NATIVE_SOL_MINT_PROGRAM_ID.String(),
				ToTokenAmount:     data.CollateralAmount,
				ToTokenDecimals:   9,
			})
		default:
			p.Actions = append(p.Actions, NewCommonDataAction(progID, p.TxInfo.Signatures[0].String(), data))
		}
	default:
		p.Actions = append(p.Actions, NewUnknownAction(progID, p.TxInfo.Signatures[0].String(), fmt.Errorf("unknown parser action, %s", progID.String())))
	}

}

func (p *Parser) parseGroupTransferSwapData(progID solana.PublicKey, swapDatas []SwapData) {
	var resultGroup [][2]SwapData
	who := p.AllAccountKeys[0].String()
	if len(swapDatas) < 2 {
		p.Actions = append(p.Actions, NewCommonDataAction(progID, p.TxInfo.Signatures[0].String(), swapDatas[0].Data))
		return
	}

	if len(swapDatas) == 3 {
		resultGroup = append(resultGroup, [2]SwapData{swapDatas[1], swapDatas[2]})
	} else {
		for i := 0; i < len(swapDatas)-1; i += 2 {
			resultGroup = append(resultGroup, [2]SwapData{swapDatas[i], swapDatas[i+1]})
		}
	}

	for _, v := range resultGroup {
		in := v[0]
		out := v[1]
		if reflect.TypeOf(in.Data) == reflect.TypeOf(out.Data) {
			if in.Type == "add_liquidity" {
				switch in := in.Data.(type) {
				case *TransferData:
					out := out.Data.(*TransferData)
					p.Actions = append(p.Actions, CommonAddLiquidityAction{
						BaseAction: BaseAction{
							ProgramID:       progID.String(),
							ProgramName:     string(ProgramName[progID]),
							InstructionName: "AddLiquidity",
							Signature:       p.TxInfo.Signatures[0].String(),
						},
						Who:            who,
						Token1:         in.Mint,
						Token1Amount:   cast.ToUint64(in.Info.Amount),
						Token1Decimals: p.SplDecimalsMap[in.Mint],
						Token2:         out.Mint,
						Token2Amount:   cast.ToUint64(out.Info.Amount),
						Token2Decimals: p.SplDecimalsMap[out.Mint],
					})

				case *TransferCheck:
					out := out.Data.(*TransferCheck)
					p.Actions = append(p.Actions, CommonAddLiquidityAction{
						BaseAction: BaseAction{
							ProgramID:       progID.String(),
							ProgramName:     string(ProgramName[progID]),
							InstructionName: "AddLiquidity",
							Signature:       p.TxInfo.Signatures[0].String(),
						},
						Who:            who,
						Token1:         in.Info.Mint,
						Token1Amount:   cast.ToUint64(in.Info.TokenAmount.Amount),
						Token1Decimals: p.SplDecimalsMap[in.Info.Mint],
						Token2:         out.Info.Mint,
						Token2Amount:   cast.ToUint64(out.Info.TokenAmount.Amount),
						Token2Decimals: p.SplDecimalsMap[out.Info.Mint],
					})
				}
			} else {
				switch in := in.Data.(type) {
				case *TransferData:
					out := out.Data.(*TransferData)
					p.Actions = append(p.Actions, CommonSwapAction{
						BaseAction: BaseAction{
							ProgramID:       progID.String(),
							ProgramName:     string(ProgramName[progID]),
							InstructionName: "Swap",
							Signature:       p.TxInfo.Signatures[0].String(),
						},
						Who:               who,
						FromToken:         in.Mint,
						FromTokenAmount:   in.Info.Amount,
						FromTokenDecimals: p.SplDecimalsMap[in.Mint],
						ToToken:           out.Mint,
						ToTokenAmount:     out.Info.Amount,
						ToTokenDecimals:   p.SplDecimalsMap[out.Mint],
					})

				case *TransferCheck:
					out := out.Data.(*TransferCheck)
					p.Actions = append(p.Actions, CommonSwapAction{
						BaseAction: BaseAction{
							ProgramID:       progID.String(),
							ProgramName:     string(ProgramName[progID]),
							InstructionName: "Swap",
							Signature:       p.TxInfo.Signatures[0].String(),
						},
						Who:               who,
						FromToken:         in.Info.Mint,
						FromTokenAmount:   cast.ToUint64(in.Info.TokenAmount.Amount),
						FromTokenDecimals: p.SplDecimalsMap[in.Info.Mint],
						ToToken:           out.Info.Mint,
						ToTokenAmount:     cast.ToUint64(out.Info.TokenAmount.Amount),
						ToTokenDecimals:   p.SplDecimalsMap[out.Info.Mint],
					})
				}
			}

		}
	}

}

func (p *Parser) parseOneTransferSwapData(progID solana.PublicKey, swapDatas []SwapData) {
	who := p.AllAccountKeys[0].String()
	in := swapDatas[0]
	out := lo.LastOrEmpty(swapDatas)
	switch in := in.Data.(type) {
	case *TransferData:
		out := out.Data.(*TransferData)
		p.Actions = append(p.Actions, CommonSwapAction{
			BaseAction: BaseAction{
				ProgramID:       progID.String(),
				ProgramName:     string(ProgramName[progID]),
				InstructionName: "Unknown Group Swap",
				Signature:       p.TxInfo.Signatures[0].String(),
			},
			Who:               who,
			FromToken:         in.Mint,
			FromTokenAmount:   in.Info.Amount,
			FromTokenDecimals: p.SplDecimalsMap[in.Mint],
			ToToken:           out.Mint,
			ToTokenAmount:     out.Info.Amount,
			ToTokenDecimals:   p.SplDecimalsMap[out.Mint],
		})
	case *TransferCheck:
		out := out.Data.(*TransferCheck)
		p.Actions = append(p.Actions, CommonSwapAction{
			BaseAction: BaseAction{
				ProgramID:       progID.String(),
				ProgramName:     string(ProgramName[progID]),
				InstructionName: "Unknown Group Swap",
				Signature:       p.TxInfo.Signatures[0].String(),
			},
			Who:               who,
			FromToken:         in.Info.Mint,
			FromTokenAmount:   cast.ToUint64(in.Info.TokenAmount.Amount),
			FromTokenDecimals: p.SplDecimalsMap[in.Info.Mint],
			ToToken:           out.Info.Mint,
			ToTokenAmount:     cast.ToUint64(out.Info.TokenAmount.Amount),
			ToTokenDecimals:   p.SplDecimalsMap[out.Info.Mint],
		})
	}

}
