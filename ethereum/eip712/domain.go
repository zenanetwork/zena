package eip712

import (
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// createEIP712Domain creates the typed data domain for the given chainID.
func createEIP712Domain(chainID uint64) apitypes.TypedDataDomain {
	domain := apitypes.TypedDataDomain{
		Name:              "Zenanet Web3",
		Version:           "1.0.0",
		ChainId:           math.NewHexOrDecimal256(int64(chainID)), // #nosec G115
		VerifyingContract: "zena",
		Salt:              "0",
	}

	return domain
}
