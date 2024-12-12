package solanaswapgo

import (
	"github.com/gagliardetto/solana-go"
	"github.com/samber/lo"
)

type Action interface {
	GetProgramID() string
	GetProgramName() string
	GetInstructionName() string
	GetSignature() string
}

type BaseAction struct {
	ProgramID       string `json:"programId"`
	ProgramName     string `json:"programName"`
	InstructionName string `json:"instructionName"`
	Signature       string `json:"signature"`
}

func (a BaseAction) GetProgramID() string {
	return a.ProgramID
}

func (a BaseAction) GetProgramName() string {
	return a.ProgramName
}

func (a BaseAction) GetInstructionName() string {
	return a.InstructionName
}

func (a BaseAction) GetSignature() string {
	return a.Signature
}

type UnknownAction struct {
	BaseAction
	Error error `json:"error"`
}

func NewUnknownAction(progID solana.PublicKey, signature string, err error) UnknownAction {
	return UnknownAction{
		BaseAction: BaseAction{
			ProgramID:       progID.String(),
			ProgramName:     lo.Ternary(lo.HasKey(ProgramName, progID), ProgramName[progID].String(), "Unknown"),
			InstructionName: "Unknown",
			Signature:       signature,
		},
		Error: err,
	}
}

type CommonDataAction struct {
	BaseAction
	Data interface{}
}

func NewCommonDataAction(progID solana.PublicKey, signature string, data interface{}, instructionName ...string) CommonDataAction {
	var instructionNameStr string
	if len(instructionName) > 0 {
		instructionNameStr = instructionName[0]
	}
	return CommonDataAction{
		BaseAction{
			ProgramID:       progID.String(),
			ProgramName:     lo.Ternary(lo.HasKey(ProgramName, progID), ProgramName[progID].String(), "Unknown"),
			InstructionName: instructionNameStr,
			Signature:       signature,
		},
		data,
	}
}

type CommonSwapAction struct {
	BaseAction
	Who               string `json:"who"`
	FromToken         string `json:"fromToken"`
	FromTokenAmount   uint64 `json:"fromTokenAmount"`
	FromTokenDecimals uint8  `json:"fromTokenDecimals"`
	ToToken           string `json:"toToken"`
	ToTokenAmount     uint64 `json:"toTokenAmount"`
	ToTokenDecimals   uint8  `json:"toTokenDecimals"`
}

type CommonAddLiquidityAction struct {
	BaseAction
	Who            string `json:"who"`
	Token1         string `json:"token1"`
	Token1Amount   uint64 `json:"token1Amount"`
	Token1Decimals uint8  `json:"token1Decimals"`
	Token2         string `json:"token2"`
	Token2Amount   uint64 `json:"token2Amount"`
	Token2Decimals uint8  `json:"token2Decimals"`
}

type CommonRemoveLiquidityAction struct {
	BaseAction
	Who            string `json:"who"`
	Token1         string `json:"token1"`
	Token1Amount   uint64 `json:"token1Amount"`
	Token1Decimals uint8  `json:"token1Decimals"`
	Token2         string `json:"token2"`
	Token2Amount   uint64 `json:"token2Amount"`
	Token2Decimals uint8  `json:"token2Decimals"`
}
