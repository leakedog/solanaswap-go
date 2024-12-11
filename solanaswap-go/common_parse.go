package solanaswapgo

import (
	"crypto/sha256"
	"fmt"
	"log"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

func CommonParseSwap[T any](result *Parser, instructionIndex int, programID solana.PublicKey, innerProgramIds ...solana.PublicKey) ([]*T, error) {
	// var instructionIndex uint16
	// for idx, instr := range result.TxInfo.Message.Instructions {
	// 	if result.AllAccountKeys[instr.ProgramIDIndex] == programID && instr.Data.String() == instruction.Data.String() {
	// 		instructionIndex = uint16(idx)
	// 		break
	// 	}
	// }

	var instructions []solana.CompiledInstruction
	for _, innerInstruction := range result.Tx.Meta.InnerInstructions {
		if innerInstruction.Index == uint16(instructionIndex) {
			instructions = innerInstruction.Instructions
			break
		}
	}

	innerProgramId := programID
	if len(innerProgramIds) > 0 {
		innerProgramId = innerProgramIds[0]
	}

	actions := make([]*T, 0)
	for _, instr := range instructions {
		programId := result.AllAccountKeys[instr.ProgramIDIndex]
		switch programId {
		case innerProgramId:
			action := new(T)
			decodedBytes, err := base58.Decode(instr.Data.String())
			if err != nil {
				log.Printf("error decoding instruction data: %s", err)
				continue
			}
			if len(decodedBytes) < 16 {
				continue
			}
			decoder := bin.NewBorshDecoder(decodedBytes[16:])
			// var data SwapData
			err = decoder.Decode(action)
			if err != nil {
				log.Printf("error decoding instruction data: %s", err)
				continue
			}
			actions = append(actions, action)
		default:
			// fmt.Printf("commonParse: Program:%s, InnerProgram:%s unknown inner program:%s", programID, innerProgramId, programId)
		}
	}

	if len(actions) > 0 {
		return actions, nil
	}

	return nil, fmt.Errorf("unknown instruction")
}

func CommonParseDecimals(result *Parser, mint solana.PublicKey) (uint8, error) {
	if mint == solana.SolMint {
		return 9, nil
	}

	decimals := result.SplDecimalsMap[mint.String()]
	if decimals > 0 {
		return decimals, nil
	}

	return 0, fmt.Errorf("unknown mint")
}

func CalculateDiscriminator(instructionName string) [8]byte {
	hash := sha256.Sum256([]byte(instructionName))
	var discriminator [8]byte
	copy(discriminator[:], hash[:8])
	return discriminator
}
