package solanaswapgo

import "github.com/gagliardetto/solana-go"

var (
	TOKEN_PROGRAM_ID           = solana.TokenProgramID
	NATIVE_SOL_MINT_PROGRAM_ID = solana.SolMint
	SWAP_DEX_PROGRAM_ID        = solana.TokenSwapProgramID
	SPL_ASSOCIATED_TOKEN_ID    = solana.SPLAssociatedTokenAccountProgramID
)

var (
	OKX_PROGRAM_ID         = solana.MustPublicKeyFromBase58("6m2CDdhRgxpH4WjvdzxAYbGxwdGUz5MziiL5jek2kBma")
	JUPITER_PROGRAM_ID     = solana.MustPublicKeyFromBase58("JUP6LkbZbjS1jKKwapdHNy74zcZ3tLUZoi5QNyVTaV4")
	JUPITER_DCA_PROGRAM_ID = solana.MustPublicKeyFromBase58("DCAK36VfExkPdAkYUQg6ewgxyinvcEyPLyHjRbmveKFw")
	PUMP_FUN_PROGRAM_ID    = solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")
	PHOENIX_PROGRAM_ID     = solana.MustPublicKeyFromBase58("PhoeNiXZ8ByJGLkxNfZRnkUfjvmuYqLR89jjFHGqdXY")
	LIFINITY_V2_PROGRAM_ID = solana.MustPublicKeyFromBase58("2wT8Yq49kHgDzXuPxZSaeLaH1qbmGXtEyPy64bL7aD3c") // Lifinity Swap V2:
	FLUXBEAM_PROGRAM_ID    = solana.MustPublicKeyFromBase58("FLUXubRmkEi2q6K3Y9kBPg9248ggaZVsoSFhtJHSrm1X") // Fluxbeam

	// Trading Bots
	BANANA_GUN_PROGRAM_ID = solana.MustPublicKeyFromBase58("BANANAjs7FJiPQqJTGFzkZJndT9o7UmKiYYGaJz6frGu")
	MINTECH_PROGRAM_ID    = solana.MustPublicKeyFromBase58("minTcHYRLVPubRK8nt6sqe2ZpWrGDLQoNLipDJCGocY")
	BLOOM_PROGRAM_ID      = solana.MustPublicKeyFromBase58("b1oomGGqPKGD6errbyfbVMBuzSC8WtAAYo8MwNafWW1")
	MAESTRO_PROGRAM_ID    = solana.MustPublicKeyFromBase58("MaestroAAe9ge5HTc64VbBQZ6fP77pwvrhM8i1XWSAx")

	RAYDIUM_V4_PROGRAM_ID                     = solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")
	RAYDIUM_AMM_PROGRAM_ID                    = solana.MustPublicKeyFromBase58("routeUGWgWzqBWFcrCfv8tritsqukccJPu3q5GPP3xS")
	RAYDIUM_AMM_LIQUIDITY_POOL_PROGRAM_ID     = solana.MustPublicKeyFromBase58("5quBtoiQqxF9Jv6KYKctB59NT3gtJD2Y65kdnB1Uev3h")
	RAYDIUM_CPMM_PROGRAM_ID                   = solana.MustPublicKeyFromBase58("CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C")
	RAYDIUM_CONCENTRATED_LIQUIDITY_PROGRAM_ID = solana.MustPublicKeyFromBase58("CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK")
	METEORA_PROGRAM_ID                        = solana.MustPublicKeyFromBase58("LBUZKhRxPF3XUpBCjp4YzTKgLccjZhTSDM9YuVaPwxo")
	METEORA_POOLS_PROGRAM_ID                  = solana.MustPublicKeyFromBase58("Eo7WjKq67rjJQSZxS6z3YkapzY3eMj6Xy8X5EQVn5UaB")
	MOONSHOT_PROGRAM_ID                       = solana.MustPublicKeyFromBase58("MoonCVVNZFSYkqNXP6bxHLPL6QQJiMagDL3qcqUQTrG")
	ORCA_PROGRAM_ID                           = solana.MustPublicKeyFromBase58("whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc")
	ORCA_TOKEN_V2_PROGRAM_ID                  = solana.MustPublicKeyFromBase58("9W959DqEETiGZocYWCQPaJ6sBmUzgfxXfqGeTEdp3aQP")

	OPENBOOK_V2_PROGRAM_ID = solana.MustPublicKeyFromBase58("opnb2LAfJYbRMAHHvqjCwQxanZn7ReEHp1k81EohpZb")
	STABLE_SWAP_PROGRAM_ID = solana.MustPublicKeyFromBase58("SSwpkEEcbUqx4vtoEByFjSkhKdCT862DNVb52nZg1UZ")
	ALDRIN_AMM_PROGRAM_ID  = solana.MustPublicKeyFromBase58("CURVGoZn8zycx6FXwwevgBTB2gVvdbGTEpvMJDbgs2t4")
)

var (
	programName = map[solana.PublicKey]SwapType{
		TOKEN_PROGRAM_ID:                          TOKEN,
		OKX_PROGRAM_ID:                            OKX,
		JUPITER_PROGRAM_ID:                        JUPITER,
		JUPITER_DCA_PROGRAM_ID:                    JUPITER,
		PUMP_FUN_PROGRAM_ID:                       PUMP_FUN,
		PHOENIX_PROGRAM_ID:                        PHOENIX,
		BANANA_GUN_PROGRAM_ID:                     BANANA_GUN,
		MINTECH_PROGRAM_ID:                        MINTECH,
		BLOOM_PROGRAM_ID:                          BLOOM,
		MAESTRO_PROGRAM_ID:                        MAESTRO,
		RAYDIUM_V4_PROGRAM_ID:                     RAYDIUM,
		RAYDIUM_AMM_PROGRAM_ID:                    RAYDIUM,
		RAYDIUM_CPMM_PROGRAM_ID:                   RAYDIUM,
		RAYDIUM_CONCENTRATED_LIQUIDITY_PROGRAM_ID: RAYDIUM,
		METEORA_PROGRAM_ID:                        METEORA,
		METEORA_POOLS_PROGRAM_ID:                  METEORA,
		MOONSHOT_PROGRAM_ID:                       MOONSHOT,
		ORCA_PROGRAM_ID:                           ORCA,
		ORCA_TOKEN_V2_PROGRAM_ID:                  ORCA,
		FLUXBEAM_PROGRAM_ID:                       FLUXBEAM,
		OPENBOOK_V2_PROGRAM_ID:                    OPENBOOK,
		SWAP_DEX_PROGRAM_ID:                       SWAP_DEX,
		LIFINITY_V2_PROGRAM_ID:                    LIFINITY,
		STABLE_SWAP_PROGRAM_ID:                    STABLE_SWAP,
		ALDRIN_AMM_PROGRAM_ID:                     ALDRIN_AMM,
	}
)

type SwapType string

const (
	TOKEN       SwapType = "Token Transfer"
	PUMP_FUN    SwapType = "PumpFun"
	JUPITER     SwapType = "Jupiter"
	OKX         SwapType = "OKX"
	RAYDIUM     SwapType = "Raydium"
	SWAP_DEX    SwapType = "Swap Dex"
	ORCA        SwapType = "Orca"
	PHOENIX     SwapType = "Phoenix"
	LIFINITY    SwapType = "LifinityV2"
	BANANA_GUN  SwapType = "Banana Gun"
	MINTECH     SwapType = "Mintech"
	BLOOM       SwapType = "Bloom"
	MAESTRO     SwapType = "Maestro"
	METEORA     SwapType = "Meteora"
	MOONSHOT    SwapType = "Moonshot"
	FLUXBEAM    SwapType = "Fluxbeam"
	OPENBOOK    SwapType = "Openbook"
	STABLE_SWAP SwapType = "StableSwap"
	ALDRIN_AMM  SwapType = "Aldrin"
	UNKNOWN     SwapType = "Unknown"
)

func ProgramName(programId solana.PublicKey) SwapType {
	if name, ok := programName[programId]; ok {
		return name
	}
	return UNKNOWN
}

func (s SwapType) String() string {
	return string(s)
}
