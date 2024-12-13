package solanaswapgo

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

func CommonParseInnerData[T any](result *Parser, instructionIndex int, progID solana.PublicKey, discriminator []byte) ([]*T, error) {
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

	return CommonParseData[T](result, instructions, progID, discriminator)
}

func CommonParseData[T any](result *Parser, instructions []solana.CompiledInstruction, progID solana.PublicKey, discriminator []byte) ([]*T, error) {
	actions := make([]*T, 0)
	for _, instr := range instructions {
		programId := result.AllAccountKeys[instr.ProgramIDIndex]

		switch programId {
		case progID:
			decodedBytes, err := base58.Decode(instr.Data.String())
			if err != nil {
				log.Printf("error decoding instruction data: %s\n", err)
				return nil, err
			}

			if !bytes.Equal(decodedBytes[:8], discriminator) {
				continue
			}

			if len(decodedBytes) < 16 {
				return nil, fmt.Errorf("instruction data is error")
			}

			action := new(T)
			decoder := bin.NewBorshDecoder(decodedBytes[16:])

			// var data SwapData
			err = decoder.Decode(action)
			if err != nil {
				log.Printf("Borsh decode data error: %s\n", err)
				return nil, err
			}
			actions = append(actions, action)
		}

		if len(actions) > 0 {
			return actions, nil
		}
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
