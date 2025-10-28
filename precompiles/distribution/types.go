package distribution

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	cmn "github.com/zenanetwork/zena/precompiles/common"
	"github.com/zenanetwork/zena/utils"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// EventSetWithdrawAddress defines the event data for the SetWithdrawAddress transaction.
type EventSetWithdrawAddress struct {
	Caller            common.Address
	WithdrawerAddress string
}

// EventWithdrawDelegatorReward defines the event data for the WithdrawDelegatorReward transaction.
type EventWithdrawDelegatorReward struct {
	DelegatorAddress common.Address
	ValidatorAddress common.Address
	Amount           *big.Int
}

// EventWithdrawValidatorRewards defines the event data for the WithdrawValidatorRewards transaction.
type EventWithdrawValidatorRewards struct {
	ValidatorAddress common.Hash
	Commission       *big.Int
}

// EventClaimRewards defines the event data for the ClaimRewards transaction.
type EventClaimRewards struct {
	DelegatorAddress common.Address
	Amount           *big.Int
}

// EventFundCommunityPool defines the event data for the FundCommunityPool transaction.
type EventFundCommunityPool struct {
	Depositor common.Address
	Denom     string
	Amount    *big.Int
}

// EventDepositValidatorRewardsPool defines the event data for the DepositValidatorRewardsPool transaction.
type EventDepositValidatorRewardsPool struct {
	Depositor        common.Address
	ValidatorAddress common.Address
	Denom            string
	Amount           *big.Int
}

// parseClaimRewardsArgs parses the arguments for the ClaimRewards method.
func parseClaimRewardsArgs(args []interface{}) (common.Address, uint32, error) {
	if len(args) != 2 {
		return common.Address{}, 0, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 2, len(args))
	}

	delegatorAddress, ok := args[0].(common.Address)
	if !ok || delegatorAddress == (common.Address{}) {
		return common.Address{}, 0, fmt.Errorf(cmn.ErrInvalidDelegator, args[0])
	}

	maxRetrieve, ok := args[1].(uint32)
	if !ok {
		return common.Address{}, 0, fmt.Errorf(cmn.ErrInvalidType, "maxRetrieve", uint32(0), args[1])
	}

	return delegatorAddress, maxRetrieve, nil
}

// NewMsgSetWithdrawAddress creates a new MsgSetWithdrawAddress instance.
func NewMsgSetWithdrawAddress(args []interface{}, addrCdc address.Codec) (*distributiontypes.MsgSetWithdrawAddress, common.Address, error) {
	if len(args) != 2 {
		return nil, common.Address{}, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 2, len(args))
	}

	delegatorAddress, ok := args[0].(common.Address)
	if !ok || delegatorAddress == (common.Address{}) {
		return nil, common.Address{}, fmt.Errorf(cmn.ErrInvalidDelegator, args[0])
	}

	withdrawerAddress, _ := args[1].(string)

	// If the withdrawer address is a hex address, convert it to a bech32 address.
	if common.IsHexAddress(withdrawerAddress) {
		var err error
		withdrawerAddress, err = sdk.Bech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), common.HexToAddress(withdrawerAddress).Bytes())
		if err != nil {
			return nil, common.Address{}, err
		}
	}

	delAddr, err := addrCdc.BytesToString(delegatorAddress.Bytes())
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("failed to decode delegator address: %w", err)
	}
	msg := &distributiontypes.MsgSetWithdrawAddress{
		DelegatorAddress: delAddr,
		WithdrawAddress:  withdrawerAddress,
	}

	return msg, delegatorAddress, nil
}

// NewMsgWithdrawDelegatorReward creates a new MsgWithdrawDelegatorReward instance.
func NewMsgWithdrawDelegatorReward(args []interface{}, addrCdc address.Codec) (*distributiontypes.MsgWithdrawDelegatorReward, common.Address, error) {
	if len(args) != 2 {
		return nil, common.Address{}, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 2, len(args))
	}

	delegatorAddress, ok := args[0].(common.Address)
	if !ok || delegatorAddress == (common.Address{}) {
		return nil, common.Address{}, fmt.Errorf(cmn.ErrInvalidDelegator, args[0])
	}

	validatorAddress, _ := args[1].(string)

	delAddr, err := addrCdc.BytesToString(delegatorAddress.Bytes())
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("failed to decode delegator address: %w", err)
	}
	msg := &distributiontypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: delAddr,
		ValidatorAddress: validatorAddress,
	}

	return msg, delegatorAddress, nil
}

// NewMsgWithdrawValidatorCommission creates a new MsgWithdrawValidatorCommission message.
func NewMsgWithdrawValidatorCommission(args []interface{}) (*distributiontypes.MsgWithdrawValidatorCommission, common.Address, error) {
	if len(args) != 1 {
		return nil, common.Address{}, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 1, len(args))
	}

	validatorAddress, _ := args[0].(string)

	msg := &distributiontypes.MsgWithdrawValidatorCommission{
		ValidatorAddress: validatorAddress,
	}

	validatorHexAddr, err := utils.HexAddressFromBech32String(msg.ValidatorAddress)
	if err != nil {
		return nil, common.Address{}, err
	}

	return msg, validatorHexAddr, nil
}

// NewMsgFundCommunityPool creates a new NewMsgFundCommunityPool message.
func NewMsgFundCommunityPool(args []interface{}, addrCdc address.Codec) (*distributiontypes.MsgFundCommunityPool, common.Address, error) {
	emptyAddr := common.Address{}
	if len(args) != 2 {
		return nil, emptyAddr, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 2, len(args))
	}

	depositorAddress, ok := args[0].(common.Address)
	if !ok || depositorAddress == emptyAddr {
		return nil, emptyAddr, fmt.Errorf(cmn.ErrInvalidHexAddress, args[0])
	}

	coins, err := cmn.ToCoins(args[1])
	if err != nil {
		return nil, emptyAddr, fmt.Errorf(ErrInvalidAmount, "amount arg")
	}

	amt, err := cmn.NewSdkCoinsFromCoins(coins)
	if err != nil {
		return nil, emptyAddr, fmt.Errorf(ErrInvalidAmount, "amount arg")
	}

	depAddr, err := addrCdc.BytesToString(depositorAddress.Bytes())
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("failed to decode depositor address: %w", err)
	}
	msg := &distributiontypes.MsgFundCommunityPool{
		Depositor: depAddr,
		Amount:    amt,
	}

	return msg, depositorAddress, nil
}

// NewMsgDepositValidatorRewardsPool creates a new MsgDepositValidatorRewardsPool message.
func NewMsgDepositValidatorRewardsPool(args []interface{}, addrCdc address.Codec) (*distributiontypes.MsgDepositValidatorRewardsPool, common.Address, error) {
	if len(args) != 3 {
		return nil, common.Address{}, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 3, len(args))
	}

	depositorAddress, ok := args[0].(common.Address)
	if !ok || depositorAddress == (common.Address{}) {
		return nil, common.Address{}, fmt.Errorf(cmn.ErrInvalidHexAddress, args[0])
	}

	validatorAddress, _ := args[1].(string)

	coins, err := cmn.ToCoins(args[2])
	if err != nil {
		return nil, common.Address{}, fmt.Errorf(cmn.ErrInvalidAmount, args[2])
	}

	amount, err := cmn.NewSdkCoinsFromCoins(coins)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf(cmn.ErrInvalidAmount, err.Error())
	}

	depAddr, err := addrCdc.BytesToString(depositorAddress.Bytes())
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("failed to decode depositor address: %w", err)
	}

	msg := &distributiontypes.MsgDepositValidatorRewardsPool{
		Depositor:        depAddr,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
	}

	return msg, depositorAddress, nil
}

// NewValidatorDistributionInfoRequest creates a new QueryValidatorDistributionInfoRequest  instance and does sanity
// checks on the provided arguments.
func NewValidatorDistributionInfoRequest(args []interface{}) (*distributiontypes.QueryValidatorDistributionInfoRequest, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 1, len(args))
	}

	validatorAddress, _ := args[0].(string)

	return &distributiontypes.QueryValidatorDistributionInfoRequest{
		ValidatorAddress: validatorAddress,
	}, nil
}

// NewValidatorOutstandingRewardsRequest creates a new QueryValidatorOutstandingRewardsRequest  instance and does sanity
// checks on the provided arguments.
func NewValidatorOutstandingRewardsRequest(args []interface{}) (*distributiontypes.QueryValidatorOutstandingRewardsRequest, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 1, len(args))
	}

	validatorAddress, _ := args[0].(string)

	return &distributiontypes.QueryValidatorOutstandingRewardsRequest{
		ValidatorAddress: validatorAddress,
	}, nil
}

// NewValidatorCommissionRequest creates a new QueryValidatorCommissionRequest  instance and does sanity
// checks on the provided arguments.
func NewValidatorCommissionRequest(args []interface{}) (*distributiontypes.QueryValidatorCommissionRequest, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 1, len(args))
	}

	validatorAddress, _ := args[0].(string)

	return &distributiontypes.QueryValidatorCommissionRequest{
		ValidatorAddress: validatorAddress,
	}, nil
}

// NewValidatorSlashesRequest creates a new QueryValidatorSlashesRequest  instance and does sanity
// checks on the provided arguments.
func NewValidatorSlashesRequest(method *abi.Method, args []interface{}) (*distributiontypes.QueryValidatorSlashesRequest, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 4, len(args))
	}

	if _, ok := args[1].(uint64); !ok {
		return nil, fmt.Errorf(cmn.ErrInvalidType, "startingHeight", uint64(0), args[1])
	}
	if _, ok := args[2].(uint64); !ok {
		return nil, fmt.Errorf(cmn.ErrInvalidType, "endingHeight", uint64(0), args[2])
	}

	var input ValidatorSlashesInput
	if err := method.Inputs.Copy(&input, args); err != nil {
		return nil, fmt.Errorf("error while unpacking args to ValidatorSlashesInput struct: %s", err)
	}

	return &distributiontypes.QueryValidatorSlashesRequest{
		ValidatorAddress: input.ValidatorAddress,
		StartingHeight:   input.StartingHeight,
		EndingHeight:     input.EndingHeight,
		Pagination:       &input.PageRequest,
	}, nil
}

// NewDelegationRewardsRequest creates a new QueryDelegationRewardsRequest  instance and does sanity
// checks on the provided arguments.
func NewDelegationRewardsRequest(args []interface{}, addrCdc address.Codec) (*distributiontypes.QueryDelegationRewardsRequest, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 2, len(args))
	}

	delegatorAddress, ok := args[0].(common.Address)
	if !ok || delegatorAddress == (common.Address{}) {
		return nil, fmt.Errorf(cmn.ErrInvalidDelegator, args[0])
	}

	validatorAddress, _ := args[1].(string)

	delAddr, err := addrCdc.BytesToString(delegatorAddress.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to decode delegator address: %w", err)
	}
	return &distributiontypes.QueryDelegationRewardsRequest{
		DelegatorAddress: delAddr,
		ValidatorAddress: validatorAddress,
	}, nil
}

// NewDelegationTotalRewardsRequest creates a new QueryDelegationTotalRewardsRequest  instance and does sanity
// checks on the provided arguments.
func NewDelegationTotalRewardsRequest(args []interface{}, addrCdc address.Codec) (*distributiontypes.QueryDelegationTotalRewardsRequest, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 1, len(args))
	}

	delegatorAddress, ok := args[0].(common.Address)
	if !ok || delegatorAddress == (common.Address{}) {
		return nil, fmt.Errorf(cmn.ErrInvalidDelegator, args[0])
	}

	delAddr, err := addrCdc.BytesToString(delegatorAddress.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to decode delegator address: %w", err)
	}
	return &distributiontypes.QueryDelegationTotalRewardsRequest{
		DelegatorAddress: delAddr,
	}, nil
}

// NewDelegatorValidatorsRequest creates a new QueryDelegatorValidatorsRequest  instance and does sanity
// checks on the provided arguments.
func NewDelegatorValidatorsRequest(args []interface{}, addrCdc address.Codec) (*distributiontypes.QueryDelegatorValidatorsRequest, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 1, len(args))
	}

	delegatorAddress, ok := args[0].(common.Address)
	if !ok || delegatorAddress == (common.Address{}) {
		return nil, fmt.Errorf(cmn.ErrInvalidDelegator, args[0])
	}

	delAddr, err := addrCdc.BytesToString(delegatorAddress.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to decode delegator address: %w", err)
	}
	return &distributiontypes.QueryDelegatorValidatorsRequest{
		DelegatorAddress: delAddr,
	}, nil
}

// NewDelegatorWithdrawAddressRequest creates a new QueryDelegatorWithdrawAddressRequest  instance and does sanity
// checks on the provided arguments.
func NewDelegatorWithdrawAddressRequest(args []interface{}, addrCdc address.Codec) (*distributiontypes.QueryDelegatorWithdrawAddressRequest, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 1, len(args))
	}

	delegatorAddress, ok := args[0].(common.Address)
	if !ok || delegatorAddress == (common.Address{}) {
		return nil, fmt.Errorf(cmn.ErrInvalidDelegator, args[0])
	}

	delAddr, err := addrCdc.BytesToString(delegatorAddress.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to decode delegator address: %w", err)
	}
	return &distributiontypes.QueryDelegatorWithdrawAddressRequest{
		DelegatorAddress: delAddr,
	}, nil
}

// NewCommunityPoolRequest creates a new QueryCommunityPoolRequest instance and does sanity
// checks on the provided arguments.
func NewCommunityPoolRequest(args []interface{}) (*distributiontypes.QueryCommunityPoolRequest, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 0, len(args))
	}

	return &distributiontypes.QueryCommunityPoolRequest{}, nil
}

// ValidatorDistributionInfo is a struct to represent the key information from
// a ValidatorDistributionInfoResponse.
type ValidatorDistributionInfo struct {
	OperatorAddress string        `abi:"operatorAddress"`
	SelfBondRewards []cmn.DecCoin `abi:"selfBondRewards"`
	Commission      []cmn.DecCoin `abi:"commission"`
}

// ValidatorDistributionInfoOutput is a wrapper for ValidatorDistributionInfo to return in the response.
type ValidatorDistributionInfoOutput struct {
	DistributionInfo ValidatorDistributionInfo `abi:"distributionInfo"`
}

// FromResponse converts a response to a ValidatorDistributionInfo.
func (o *ValidatorDistributionInfoOutput) FromResponse(res *distributiontypes.QueryValidatorDistributionInfoResponse) ValidatorDistributionInfoOutput {
	return ValidatorDistributionInfoOutput{
		DistributionInfo: ValidatorDistributionInfo{
			OperatorAddress: res.OperatorAddress,
			SelfBondRewards: cmn.NewDecCoinsResponse(res.SelfBondRewards),
			Commission:      cmn.NewDecCoinsResponse(res.Commission),
		},
	}
}

// ValidatorSlashEvent is a struct to represent the key information from
// a ValidatorSlashEvent response.
type ValidatorSlashEvent struct {
	ValidatorPeriod uint64  `abi:"validatorPeriod"`
	Fraction        cmn.Dec `abi:"fraction"`
}

// ValidatorSlashesInput is a struct to represent the key information
// to perform a ValidatorSlashes query.
type ValidatorSlashesInput struct {
	ValidatorAddress string
	StartingHeight   uint64
	EndingHeight     uint64
	PageRequest      query.PageRequest
}

// ValidatorSlashesOutput is a struct to represent the key information from
// a ValidatorSlashes response.
type ValidatorSlashesOutput struct {
	Slashes      []ValidatorSlashEvent
	PageResponse query.PageResponse
}

// FromResponse populates the ValidatorSlashesOutput from a QueryValidatorSlashesResponse.
func (vs *ValidatorSlashesOutput) FromResponse(res *distributiontypes.QueryValidatorSlashesResponse) *ValidatorSlashesOutput {
	vs.Slashes = make([]ValidatorSlashEvent, len(res.Slashes))
	for i, s := range res.Slashes {
		vs.Slashes[i] = ValidatorSlashEvent{
			ValidatorPeriod: s.ValidatorPeriod,
			Fraction: cmn.Dec{
				Value:     s.Fraction.BigInt(),
				Precision: math.LegacyPrecision,
			},
		}
	}

	if res.Pagination != nil {
		vs.PageResponse.Total = res.Pagination.Total
		vs.PageResponse.NextKey = res.Pagination.NextKey
	}

	return vs
}

// Pack packs a given slice of abi arguments into a byte array.
func (vs *ValidatorSlashesOutput) Pack(args abi.Arguments) ([]byte, error) {
	return args.Pack(vs.Slashes, vs.PageResponse)
}

// DelegationDelegatorReward is a struct to represent the key information from
// a query for the rewards of a delegation to a given validator.
type DelegationDelegatorReward struct {
	ValidatorAddress string
	Reward           []cmn.DecCoin
}

// DelegationTotalRewardsOutput is a struct to represent the key information from
// a DelegationTotalRewards response.
type DelegationTotalRewardsOutput struct {
	Rewards []DelegationDelegatorReward
	Total   []cmn.DecCoin
}

// FromResponse populates the DelegationTotalRewardsOutput from a QueryDelegationTotalRewardsResponse.
func (dtr *DelegationTotalRewardsOutput) FromResponse(res *distributiontypes.QueryDelegationTotalRewardsResponse) *DelegationTotalRewardsOutput {
	dtr.Rewards = make([]DelegationDelegatorReward, len(res.Rewards))
	for i, r := range res.Rewards {
		dtr.Rewards[i] = DelegationDelegatorReward{
			ValidatorAddress: r.ValidatorAddress,
			Reward:           cmn.NewDecCoinsResponse(r.Reward),
		}
	}
	dtr.Total = cmn.NewDecCoinsResponse(res.Total)
	return dtr
}

// Pack packs a given slice of abi arguments into a byte array.
func (dtr *DelegationTotalRewardsOutput) Pack(args abi.Arguments) ([]byte, error) {
	return args.Pack(dtr.Rewards, dtr.Total)
}

// CommunityPoolOutput is a struct to represent the key information from
// a CommunityPool response.
type CommunityPoolOutput struct {
	Pool []cmn.DecCoin
}

// FromResponse populates the CommunityPoolOutput from a QueryCommunityPoolResponse.
func (cp *CommunityPoolOutput) FromResponse(res *distributiontypes.QueryCommunityPoolResponse) *CommunityPoolOutput {
	cp.Pool = cmn.NewDecCoinsResponse(res.Pool)
	return cp
}

// Pack packs a given slice of abi arguments into a byte array.
func (cp *CommunityPoolOutput) Pack(args abi.Arguments) ([]byte, error) {
	return args.Pack(cp.Pool)
}
