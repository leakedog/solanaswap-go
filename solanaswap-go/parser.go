package solanaswapgo

import (
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/sirupsen/logrus"
)

type Parser struct {
	Tx              *rpc.GetTransactionResult
	TxInfo          *solana.Transaction
	AllAccountKeys  solana.PublicKeySlice
	SplTokenInfoMap map[string]TokenInfo // map[authority]TokenInfo
	SplDecimalsMap  map[string]uint8     // map[mint]decimals
	Log             *logrus.Logger
	Actions         []Action
	SwapData        []SwapData
}

func NewTransactionParser(tx *rpc.GetTransactionResult) (*Parser, error) {

	txInfo, err := tx.Transaction.GetTransaction()
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	allAccountKeys := append(txInfo.Message.AccountKeys, tx.Meta.LoadedAddresses.Writable...)
	allAccountKeys = append(allAccountKeys, tx.Meta.LoadedAddresses.ReadOnly...)

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	parser := &Parser{
		Tx:             tx,
		TxInfo:         txInfo,
		AllAccountKeys: allAccountKeys,
		Log:            log,
	}

	if err := parser.extractSPLTokenInfo(); err != nil {
		return nil, fmt.Errorf("failed to extract SPL Token Addresses: %w", err)
	}

	if err := parser.extractSPLDecimals(); err != nil {
		return nil, fmt.Errorf("failed to extract SPL decimals: %w", err)
	}

	return parser, nil
}

func NewBlockTransactionParser(tx rpc.TransactionWithMeta) (*Parser, error) {

	txInfo, err := tx.GetTransaction()
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	allAccountKeys := append(txInfo.Message.AccountKeys, tx.Meta.LoadedAddresses.Writable...)
	allAccountKeys = append(allAccountKeys, tx.Meta.LoadedAddresses.ReadOnly...)

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	parser := &Parser{
		Tx: &rpc.GetTransactionResult{
			Slot:      tx.Slot,
			BlockTime: tx.BlockTime,
			// Transaction: tx.Transaction,
			Meta:    tx.Meta,
			Version: tx.Version,
		},
		TxInfo:         txInfo,
		AllAccountKeys: allAccountKeys,
		Log:            log,
	}

	if err := parser.extractSPLTokenInfo(); err != nil {
		return nil, fmt.Errorf("failed to extract SPL Token Addresses: %w", err)
	}

	if err := parser.extractSPLDecimals(); err != nil {
		return nil, fmt.Errorf("failed to extract SPL decimals: %w", err)
	}

	return parser, nil
}

type SwapData struct {
	Type   SwapType
	Action string
	Data   interface{}
}

func (p *Parser) ParseTransaction() ([]SwapData, error) {
	var parsedSwaps []SwapData

	skip := false
	for i, outerInstruction := range p.TxInfo.Message.Instructions {
		progID := p.AllAccountKeys[outerInstruction.ProgramIDIndex]
		switch {
		case progID.Equals(JUPITER_PROGRAM_ID):
			skip = true
			parsedSwaps = append(parsedSwaps, p.processJupiterSwaps(i)...)
		case progID.Equals(MOONSHOT_PROGRAM_ID):
			skip = true
			parsedSwaps = append(parsedSwaps, p.processMoonshotSwaps()...)
		case progID.Equals(BANANA_GUN_PROGRAM_ID) ||
			progID.Equals(MINTECH_PROGRAM_ID) ||
			progID.Equals(BLOOM_PROGRAM_ID) ||
			progID.Equals(MAESTRO_PROGRAM_ID):
			// Check inner instructions to determine which swap protocol is being used
			if innerSwaps := p.processTradingBotSwaps(i); len(innerSwaps) > 0 {
				parsedSwaps = append(parsedSwaps, innerSwaps...)
			}
		}
	}
	if skip {
		return parsedSwaps, nil
	}

	for i, outerInstruction := range p.TxInfo.Message.Instructions {
		progID := p.AllAccountKeys[outerInstruction.ProgramIDIndex]
		switch {
		case progID.Equals(RAYDIUM_V4_PROGRAM_ID) ||
			progID.Equals(RAYDIUM_CPMM_PROGRAM_ID) ||
			progID.Equals(RAYDIUM_AMM_PROGRAM_ID) ||
			progID.Equals(RAYDIUM_CONCENTRATED_LIQUIDITY_PROGRAM_ID) ||
			progID.Equals(solana.MustPublicKeyFromBase58("AP51WLiiqTdbZfgyRMs35PsZpdmLuPDdHYmrB23pEtMU")):
			parsedSwaps = append(parsedSwaps, p.processTransferSwapDex(i, RAYDIUM)...)
		case progID.Equals(ORCA_PROGRAM_ID):
			parsedSwaps = append(parsedSwaps, p.processTransferSwapDex(i, ORCA)...)
		case progID.Equals(METEORA_PROGRAM_ID) || progID.Equals(METEORA_POOLS_PROGRAM_ID):
			parsedSwaps = append(parsedSwaps, p.processTransferSwapDex(i, METEORA)...)
		case progID.Equals(PUMP_FUN_PROGRAM_ID) ||
			progID.Equals(solana.MustPublicKeyFromBase58("BSfD6SHZigAfDWSjzD5Q41jw8LmKwtmjskPH9XW1mrRW")): // PumpFun
			parsedSwaps = append(parsedSwaps, p.processPumpfunSwaps(i)...)
		}
	}

	return parsedSwaps, nil
}

type SwapInfo struct {
	Signers    []solana.PublicKey
	Signatures []solana.Signature
	AMMs       []string
	Slot       uint64
	Timestamp  time.Time

	TokenInMint     solana.PublicKey
	TokenInAmount   uint64
	TokenInDecimals uint8

	TokenOutMint     solana.PublicKey
	TokenOutAmount   uint64
	TokenOutDecimals uint8
}

func (p *Parser) ProcessSwapData(swapDatas []SwapData) (*SwapInfo, error) {

	// txInfo, err := p.tx.Transaction.GetTransaction()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get transaction: %w", err)
	// }
	txInfo := p.TxInfo

	swapInfo := &SwapInfo{
		Signers:    txInfo.Message.Signers(),
		Signatures: txInfo.Signatures,
		Timestamp:  p.Tx.BlockTime.Time(),
		Slot:       p.Tx.Slot,
	}

	for i, swapData := range swapDatas {
		switch swapData.Type {
		case JUPITER:
			intermediateInfo, err := parseJupiterEvents(swapDatas)
			if err != nil {
				return nil, fmt.Errorf("failed to parse Jupiter events: %w", err)
			}
			jupiterSwapInfo, err := p.convertToSwapInfo(intermediateInfo)
			if err != nil {
				return nil, fmt.Errorf("failed to convert to swap info: %w", err)
			}
			jupiterSwapInfo.Signatures = swapInfo.Signatures
			return jupiterSwapInfo, nil
		case PUMP_FUN:
			if swapData.Data.(*PumpfunTradeEvent).IsBuy {
				swapInfo.TokenInMint = NATIVE_SOL_MINT_PROGRAM_ID // TokenIn info is always SOL for Pumpfun
				swapInfo.TokenInAmount = swapData.Data.(*PumpfunTradeEvent).SolAmount
				swapInfo.TokenInDecimals = 9
				swapInfo.TokenOutMint = swapData.Data.(*PumpfunTradeEvent).Mint
				swapInfo.TokenOutAmount = swapData.Data.(*PumpfunTradeEvent).TokenAmount
				swapInfo.TokenOutDecimals = p.SplDecimalsMap[swapInfo.TokenOutMint.String()]
			} else {
				swapInfo.TokenInMint = swapData.Data.(*PumpfunTradeEvent).Mint
				swapInfo.TokenInAmount = swapData.Data.(*PumpfunTradeEvent).TokenAmount
				swapInfo.TokenInDecimals = p.SplDecimalsMap[swapInfo.TokenInMint.String()]
				swapInfo.TokenOutMint = NATIVE_SOL_MINT_PROGRAM_ID // TokenOut info is always SOL for Pumpfun
				swapInfo.TokenOutAmount = swapData.Data.(*PumpfunTradeEvent).SolAmount
				swapInfo.TokenOutDecimals = 9
			}
			swapInfo.AMMs = append(swapInfo.AMMs, string(swapData.Type))
			swapInfo.Timestamp = time.Unix(int64(swapData.Data.(*PumpfunTradeEvent).Timestamp), 0)
			return swapInfo, nil // Pumpfun only has one swap event
		case METEORA:
			switch swapData.Data.(type) {
			case *TransferSwapData: // Meteora Pools
				swapData := swapData.Data.(*TransferSwapData)
				if i == 0 {
					swapInfo.TokenInMint = solana.MustPublicKeyFromBase58(swapData.Mint)
					swapInfo.TokenInAmount = swapData.Amount
					swapInfo.TokenInDecimals = swapData.Decimals
				} else {
					if swapData.Authority == swapInfo.Signers[0].String() && swapData.Mint == swapInfo.TokenInMint.String() {
						swapInfo.TokenInAmount += swapData.Amount
					}
					swapInfo.TokenOutMint = solana.MustPublicKeyFromBase58(swapData.Mint)
					swapInfo.TokenOutAmount = swapData.Amount
					swapInfo.TokenOutDecimals = swapData.Decimals
				}
			}
		case RAYDIUM, ORCA:
			switch swapData.Data.(type) {
			case *TransferSwapData: // Raydium V4 and Orca
				swapData := swapData.Data.(*TransferSwapData)
				if i == 0 {
					swapInfo.TokenInMint = solana.MustPublicKeyFromBase58(swapData.Mint)
					swapInfo.TokenInAmount = swapData.Amount
					swapInfo.TokenInDecimals = swapData.Decimals
				} else {
					swapInfo.TokenOutMint = solana.MustPublicKeyFromBase58(swapData.Mint)
					swapInfo.TokenOutAmount = swapData.Amount
					swapInfo.TokenOutDecimals = swapData.Decimals
				}
			}
		case MOONSHOT:
			swapData := swapData.Data.(*MoonshotTradeInstructionWithMint)
			switch swapData.TradeType {
			case TradeTypeBuy: // BUY
				swapInfo.TokenInMint = NATIVE_SOL_MINT_PROGRAM_ID
				swapInfo.TokenInAmount = swapData.CollateralAmount
				swapInfo.TokenInDecimals = 9
				swapInfo.TokenOutMint = swapData.Mint
				swapInfo.TokenOutAmount = swapData.TokenAmount
				swapInfo.TokenOutDecimals = 9
			case TradeTypeSell: // SELL
				swapInfo.TokenInMint = swapData.Mint
				swapInfo.TokenInAmount = swapData.TokenAmount
				swapInfo.TokenInDecimals = 9
				swapInfo.TokenOutMint = NATIVE_SOL_MINT_PROGRAM_ID
				swapInfo.TokenOutAmount = swapData.CollateralAmount
				swapInfo.TokenOutDecimals = 9
			default:
				return nil, fmt.Errorf("invalid trade type: %d", swapData.TradeType)
			}
		}
		swapInfo.AMMs = append(swapInfo.AMMs, string(swapData.Type))
	}
	return swapInfo, nil
}

func (p *Parser) processTradingBotSwaps(instructionIndex int) []SwapData {
	var swaps []SwapData

	// get inner instructions for this index
	innerInstructions := p.getInnerInstructions(instructionIndex)
	if len(innerInstructions) == 0 {
		return swaps
	}

	// track which protocols we've processed to avoid duplicates
	processedProtocols := make(map[string]bool)

	// check program IDs of inner instructions to determine swap type
	for _, inner := range innerInstructions {
		progID := p.AllAccountKeys[inner.ProgramIDIndex]

		switch {
		case (progID.Equals(RAYDIUM_V4_PROGRAM_ID) ||
			progID.Equals(RAYDIUM_CPMM_PROGRAM_ID) ||
			progID.Equals(RAYDIUM_AMM_PROGRAM_ID) ||
			progID.Equals(RAYDIUM_CONCENTRATED_LIQUIDITY_PROGRAM_ID)) && !processedProtocols["raydium"]:
			processedProtocols["raydium"] = true
			if raydSwaps := p.processTransferSwapDex(instructionIndex, RAYDIUM); len(raydSwaps) > 0 {
				swaps = append(swaps, raydSwaps...)
			}

		case progID.Equals(ORCA_PROGRAM_ID) && !processedProtocols["orca"]:
			processedProtocols["orca"] = true
			if orcaSwaps := p.processTransferSwapDex(instructionIndex, ORCA); len(orcaSwaps) > 0 {
				swaps = append(swaps, orcaSwaps...)
			}

		case (progID.Equals(METEORA_PROGRAM_ID) ||
			progID.Equals(METEORA_POOLS_PROGRAM_ID)) && !processedProtocols["meteora"]:
			processedProtocols["meteora"] = true
			if meteoraSwaps := p.processTransferSwapDex(instructionIndex, METEORA); len(meteoraSwaps) > 0 {
				swaps = append(swaps, meteoraSwaps...)
			}

		case (progID.Equals(PUMP_FUN_PROGRAM_ID) ||
			progID.Equals(solana.MustPublicKeyFromBase58("BSfD6SHZigAfDWSjzD5Q41jw8LmKwtmjskPH9XW1mrRW"))) && !processedProtocols["pumpfun"]:
			processedProtocols["pumpfun"] = true
			if pumpfunSwaps := p.processPumpfunSwaps(instructionIndex); len(pumpfunSwaps) > 0 {
				swaps = append(swaps, pumpfunSwaps...)
			}
		}
	}

	return swaps
}

// helper function to get inner instructions for a given instruction index
func (p *Parser) getInnerInstructions(index int) []solana.CompiledInstruction {
	if p.Tx.Meta == nil || p.Tx.Meta.InnerInstructions == nil {
		return nil
	}

	for _, inner := range p.Tx.Meta.InnerInstructions {
		if inner.Index == uint16(index) {
			return inner.Instructions
		}
	}

	return nil
}
