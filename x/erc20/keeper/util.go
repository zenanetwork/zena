package keeper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	types2 "github.com/zenanetwork/zena/x/erc20/types"
	"github.com/zenanetwork/zena/x/vm/types"

	"cosmossdk.io/errors"
)

// validateApprovalEventDoesNotExist returns an error if the given transactions logs include
// an unexpected `Approval` event
func validateApprovalEventDoesNotExist(logs []*types.Log) error {
	for _, log := range logs {
		if log.Topics[0] == logApprovalSigHash.Hex() {
			return errors.Wrapf(
				types2.ErrUnexpectedEvent, "unexpected Approval event",
			)
		}
	}

	return nil
}

// validateTransferEventExists returns an error if the given transactions logs DO NOT include
// an expected `Transfer` event from the expected address with matching from, to and value (H-01).
func validateTransferEventExists(
	logs []*types.Log,
	tokenAddress common.Address,
	expectedFrom common.Address,
	expectedTo common.Address,
	expectedAmount *big.Int,
) error {
	if len(logs) == 0 {
		return errors.Wrapf(
			types2.ErrExpectedEvent, "expected Transfer event",
		)
	}
	found := false
	for _, log := range logs {
		if log.Topics[0] == logTransferSigHash.Hex() {
			if log.Address != tokenAddress.Hex() {
				return errors.Wrapf(
					types2.ErrUnexpectedEvent, "Transfer event from unexpected address",
				)
			}
			if found {
				return errors.Wrapf(
					types2.ErrUnexpectedEvent, "duplicate Transfer event",
				)
			}

			// Validate indexed parameters: from (Topics[1]) and to (Topics[2])
			if len(log.Topics) < 3 {
				return errors.Wrapf(
					types2.ErrUnexpectedEvent, "Transfer event has insufficient topics: got %d, expected 3", len(log.Topics),
				)
			}

			from := common.HexToAddress(log.Topics[1])
			if from != expectedFrom {
				return errors.Wrapf(
					types2.ErrUnexpectedEvent, "Transfer from mismatch: expected %s, got %s", expectedFrom.Hex(), from.Hex(),
				)
			}

			to := common.HexToAddress(log.Topics[2])
			if to != expectedTo {
				return errors.Wrapf(
					types2.ErrUnexpectedEvent, "Transfer to mismatch: expected %s, got %s", expectedTo.Hex(), to.Hex(),
				)
			}

			// Validate non-indexed parameter: value (Data field, ABI-encoded uint256)
			if len(log.Data) >= 32 {
				value := new(big.Int).SetBytes(log.Data[:32])
				if value.Cmp(expectedAmount) != 0 {
					return errors.Wrapf(
						types2.ErrUnexpectedEvent, "Transfer value mismatch: expected %s, got %s", expectedAmount.String(), value.String(),
					)
				}
			}

			found = true
		}
	}

	if !found {
		return errors.Wrapf(
			types2.ErrExpectedEvent, "expected Transfer event",
		)
	}

	return nil
}
