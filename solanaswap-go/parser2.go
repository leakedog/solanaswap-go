package solanaswapgo

import (
	"fmt"
	"reflect"
	"strconv"
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
		case strings.Contains(progID.String(), "Vote111111111111111111111111111111111111111"):
		case strings.Contains(progID.String(), "ComputeBudget111111111111111111111111111111"):
		case strings.Contains(progID.String(), "11111111111111111111111111111111"):
		case progID.Equals(TOKEN_PROGRAM_ID):
			// p.programParseTo(p.processTransferSwapDex(i, TOKEN), progID)
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
		case progID.Equals(LIFINITY_V2_PROGRAM_ID):
			p.programParseTo(p.processTransferSwapDex(i, LIFINITY), progID)
		case progID.Equals(PHOENIX_PROGRAM_ID):
			p.programParseTo(p.processTransferSwapDex(i, PHOENIX), progID)
		case progID.Equals(OPENBOOK_V2_PROGRAM_ID):
			p.programParseTo(p.processTransferSwapDex(i, OPENBOOK), progID)
		case progID.Equals(RAYDIUM_V4_PROGRAM_ID) ||
			progID.Equals(RAYDIUM_CPMM_PROGRAM_ID) ||
			progID.Equals(RAYDIUM_AMM_PROGRAM_ID) ||
			progID.Equals(RAYDIUM_AMM_LIQUIDITY_POOL_PROGRAM_ID) ||
			progID.Equals(RAYDIUM_CONCENTRATED_LIQUIDITY_PROGRAM_ID) ||
			progID.Equals(solana.MustPublicKeyFromBase58("AP51WLiiqTdbZfgyRMs35PsZpdmLuPDdHYmrB23pEtMU")):
			p.programParseTo(p.processTransferSwapDex(i, RAYDIUM), progID)
		case progID.Equals(OKX_PROGRAM_ID):
			p.programParseTo(p.OkxInstruction(outerInstruction, progID, i), progID)
			// p.programParseTo(p.processOkxSwaps(i), progID)
		case progID.Equals(ORCA_PROGRAM_ID) || progID.Equals(ORCA_TOKEN_V2_PROGRAM_ID):
			p.programParseTo(p.processTransferSwapDex(i, ORCA), progID)
		case progID.Equals(METEORA_PROGRAM_ID) || progID.Equals(METEORA_POOLS_PROGRAM_ID):
			p.programParseTo(p.processTransferSwapDex(i, METEORA), progID)
		case progID.Equals(PUMP_FUN_PROGRAM_ID) ||
			progID.Equals(solana.MustPublicKeyFromBase58("BSfD6SHZigAfDWSjzD5Q41jw8LmKwtmjskPH9XW1mrRW")): // PumpFun
			p.programParseTo(p.processPumpfunSwaps(i), progID)
		default:
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
				ProgramName:     "Jupiter",
				InstructionName: "Jupiter",
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
		ProgramName:     string(ProgramName[progID]),
		InstructionName: in.Type.String(),
		Signature:       p.TxInfo.Signatures[0].String(),
	}
	if in.Action == "add_liquidity" {
		baseAction.InstructionName = "AddLiquidity"
		switch in := in.Data.(type) {
		case *TransferData:
			out := out.Data.(*TransferData)
			action = CommonAddLiquidityAction{
				BaseAction:     baseAction,
				Who:            who,
				Token1:         in.Mint,
				Token1Amount:   cast.ToUint64(in.Info.Amount),
				Token1Decimals: p.SplDecimalsMap[in.Mint],
				Token2:         out.Mint,
				Token2Amount:   cast.ToUint64(out.Info.Amount),
				Token2Decimals: p.SplDecimalsMap[out.Mint],
			}

		case *TransferCheck:
			out := out.Data.(*TransferCheck)
			action = CommonAddLiquidityAction{
				BaseAction:     baseAction,
				Who:            who,
				Token1:         in.Info.Mint,
				Token1Amount:   cast.ToUint64(in.Info.TokenAmount.Amount),
				Token1Decimals: p.SplDecimalsMap[in.Info.Mint],
				Token2:         out.Info.Mint,
				Token2Amount:   cast.ToUint64(out.Info.TokenAmount.Amount),
				Token2Decimals: p.SplDecimalsMap[out.Info.Mint],
			}
		}

	} else if in.Action == "remove_liquidity" {
		baseAction.InstructionName = "RemoveLiquidity"
		switch in := in.Data.(type) {
		case *TransferData:
			out := out.Data.(*TransferData)
			action = CommonRemoveLiquidityAction{
				BaseAction:     baseAction,
				Who:            who,
				Token1:         in.Mint,
				Token1Amount:   cast.ToUint64(in.Info.Amount),
				Token1Decimals: p.SplDecimalsMap[in.Mint],
				Token2:         out.Mint,
				Token2Amount:   cast.ToUint64(out.Info.Amount),
				Token2Decimals: p.SplDecimalsMap[out.Mint],
			}

		case *TransferCheck:
			out := out.Data.(*TransferCheck)
			action = CommonRemoveLiquidityAction{
				BaseAction:     baseAction,
				Who:            who,
				Token1:         in.Info.Mint,
				Token1Amount:   cast.ToUint64(in.Info.TokenAmount.Amount),
				Token1Decimals: p.SplDecimalsMap[in.Info.Mint],
				Token2:         out.Info.Mint,
				Token2Amount:   cast.ToUint64(out.Info.TokenAmount.Amount),
				Token2Decimals: p.SplDecimalsMap[out.Info.Mint],
			}
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
		switch in := in.Data.(type) {
		case *TransferData:
			out := out.Data.(*TransferData)
			action = CommonSwapAction{
				BaseAction:        baseAction,
				Who:               who,
				FromToken:         in.Mint,
				FromTokenAmount:   in.Info.Amount,
				FromTokenDecimals: p.SplDecimalsMap[in.Mint],
				ToToken:           out.Mint,
				ToTokenAmount:     out.Info.Amount,
				ToTokenDecimals:   p.SplDecimalsMap[out.Mint],
			}

		case *TransferCheck:
			out := out.Data.(*TransferCheck)
			action = CommonSwapAction{
				BaseAction:        baseAction,
				Who:               who,
				FromToken:         in.Info.Mint,
				FromTokenAmount:   cast.ToUint64(in.Info.TokenAmount.Amount),
				FromTokenDecimals: p.SplDecimalsMap[in.Info.Mint],
				ToToken:           out.Info.Mint,
				ToTokenAmount:     cast.ToUint64(out.Info.TokenAmount.Amount),
				ToTokenDecimals:   p.SplDecimalsMap[out.Info.Mint],
			}
		}
	}
	p.Actions = append(p.Actions, action)
}

// METEORA_POOLS
func (p *Parser) parseMeteoraPoolsSwapData(progID solana.PublicKey, swapDatas []SwapData) {
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
		case *TransferCheck:
			swapData := swapData.Data.(*TransferCheck)
			if i == 0 {
				tokenInAmount, _ := strconv.ParseInt(swapData.Info.TokenAmount.Amount, 10, 64)
				action.FromToken = swapData.Info.Mint
				action.FromTokenAmount = uint64(tokenInAmount)
				action.FromTokenDecimals = swapData.Info.TokenAmount.Decimals
			} else {
				tokenOutAmount, _ := strconv.ParseFloat(swapData.Info.TokenAmount.Amount, 64)
				action.ToToken = swapData.Info.Mint
				action.ToTokenAmount = uint64(tokenOutAmount)
				action.ToTokenDecimals = swapData.Info.TokenAmount.Decimals
			}
		case *TransferData: // Meteora Pools
			swapData := swapData.Data.(*TransferData)
			if i == 0 {
				action.FromToken = swapData.Mint
				action.FromTokenAmount = swapData.Info.Amount
				action.FromTokenDecimals = swapData.Decimals
			} else {
				if swapData.Info.Authority == p.AllAccountKeys[0].String() && swapData.Mint == action.FromToken {
					action.FromTokenAmount += swapData.Info.Amount
				}
				action.ToToken = swapData.Mint
				action.ToTokenAmount = swapData.Info.Amount
				action.ToTokenDecimals = swapData.Decimals
			}
		}
	}
	p.Actions = append(p.Actions, action)
}
