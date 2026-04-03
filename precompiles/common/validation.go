package common

import (
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
)

// ValidatePositiveAmount validates that amount is positive and within SDK bounds (H-02).
func ValidatePositiveAmount(amount *big.Int) error {
	if amount == nil || amount.Sign() <= 0 {
		return fmt.Errorf("amount must be positive, got %v", amount)
	}
	if amount.BitLen() > sdkmath.MaxBitLen {
		return fmt.Errorf("amount exceeds maximum bit length (%d)", sdkmath.MaxBitLen)
	}
	return nil
}

// ValidateNonNegativeAmount validates that amount is non-negative and within SDK bounds (H-02).
func ValidateNonNegativeAmount(amount *big.Int) error {
	if amount == nil || amount.Sign() < 0 {
		return fmt.Errorf("amount must be non-negative, got %v", amount)
	}
	if amount.BitLen() > sdkmath.MaxBitLen {
		return fmt.Errorf("amount exceeds maximum bit length (%d)", sdkmath.MaxBitLen)
	}
	return nil
}

// ValidateNonZeroAddress validates that address is not the zero address (H-02).
func ValidateNonZeroAddress(addr common.Address) error {
	if addr == (common.Address{}) {
		return fmt.Errorf("zero address is not allowed")
	}
	return nil
}
