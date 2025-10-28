package utils

import (
	"fmt"

	"github.com/zenanetwork/zena/testutil/integration/evm/factory"
	"github.com/zenanetwork/zena/testutil/integration/evm/network"
	erc20types "github.com/zenanetwork/zena/x/erc20/types"
	feemarkettypes "github.com/zenanetwork/zena/x/feemarket/types"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

type UpdateParamsInput struct {
	Tf      factory.TxFactory
	Network network.Network
	Pk      cryptotypes.PrivKey
	Params  interface{}
}

var authority = authtypes.NewModuleAddress(govtypes.ModuleName).String()

// UpdateEvmParams helper function to update the EVM module parameters
// It submits an update params proposal, votes for it, and waits till it passes
func UpdateEvmParams(input UpdateParamsInput) error {
	return updateModuleParams[evmtypes.Params](input, evmtypes.ModuleName)
}

// UpdateGovParams helper function to update the governance module parameters
// It submits an update params proposal, votes for it, and waits till it passes
func UpdateGovParams(input UpdateParamsInput) error {
	return updateModuleParams[govv1types.Params](input, govtypes.ModuleName)
}

// UpdateFeeMarketParams helper function to update the feemarket module parameters
// It submits an update params proposal, votes for it, and waits till it passes
func UpdateFeeMarketParams(input UpdateParamsInput) error {
	return updateModuleParams[feemarkettypes.Params](input, feemarkettypes.ModuleName)
}

// UpdateERC20Params helper function to update the erc20 module parameters
// It submits an update params proposal, votes for it, and waits till it passes
func UpdateERC20Params(input UpdateParamsInput) error {
	return updateModuleParams[erc20types.Params](input, erc20types.ModuleName)
}

// updateModuleParams helper function to update module parameters
// It submits an update params proposal, votes for it, and waits till it passes
func updateModuleParams[T interface{}](input UpdateParamsInput, moduleName string) error {
	newParams, ok := input.Params.(T)
	if !ok {
		return fmt.Errorf("invalid params type %T for module %s", input.Params, moduleName)
	}

	proposalMsg := createProposalMsg(newParams, moduleName)

	title := fmt.Sprintf("Update %s params", moduleName)
	proposalID, err := SubmitProposal(input.Tf, input.Network, input.Pk, title, proposalMsg)
	if err != nil {
		return err
	}

	return ApproveProposal(input.Tf, input.Network, input.Pk, proposalID)
}

// createProposalMsg creates the module-specific update params message
func createProposalMsg(params interface{}, name string) sdk.Msg {
	switch name {
	case evmtypes.ModuleName:
		return &evmtypes.MsgUpdateParams{Authority: authority, Params: params.(evmtypes.Params)}
	case govtypes.ModuleName:
		return &govv1types.MsgUpdateParams{Authority: authority, Params: params.(govv1types.Params)}
	case feemarkettypes.ModuleName:
		return &feemarkettypes.MsgUpdateParams{Authority: authority, Params: params.(feemarkettypes.Params)}
	case erc20types.ModuleName:
		return &erc20types.MsgUpdateParams{Authority: authority, Params: params.(erc20types.Params)}
	default:
		return nil
	}
}
