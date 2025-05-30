package common

const (
	// ErrNotRunInEvm is raised when a function is not called inside the EVM.
	ErrNotRunInEvm = "not run in EVM"
	// ErrRequesterIsNotMsgSender is raised when the requester address is not the same as the msg.sender.
	ErrRequesterIsNotMsgSender = "msg.sender address %s does not match the requester address %s"
	// ErrInvalidABI is raised when the ABI cannot be parsed.
	ErrInvalidABI = "invalid ABI: %w"
	// ErrInvalidAmount is raised when the amount cannot be cast to a big.Int.
	ErrInvalidAmount = "invalid amount: %v"
	// ErrInvalidHexAddress is raised when the hex address is not valid.
	ErrInvalidHexAddress = "invalid hex address address: %s"
	// ErrInvalidDelegator is raised when the delegator address is not valid.
	ErrInvalidDelegator = "invalid delegator address: %s"
	// ErrInvalidValidator is raised when the validator address is not valid.
	ErrInvalidValidator = "invalid validator address: %s"
	// ErrInvalidDenom is raised when the denom is not valid.
	ErrInvalidDenom = "invalid denom: %s"
	// ErrInvalidMsgType is raised when the transaction type is not valid for the given precompile.
	ErrInvalidMsgType = "invalid %s transaction type: %s"
	// ErrInvalidNumberOfArgs is raised when the number of arguments is not what is expected.
	ErrInvalidNumberOfArgs = "invalid number of arguments; expected %d; got: %d"
	// ErrUnknownMethod is raised when the method is not known.
	ErrUnknownMethod = "unknown method: %s"
	// ErrIntegerOverflow is raised when an integer overflow occurs.
	ErrIntegerOverflow = "integer overflow when increasing allowance"
	// ErrNegativeAmount is raised when an amount is negative.
	ErrNegativeAmount = "negative amount when decreasing allowance"
	// ErrInvalidType is raised when the provided type is different than the expected.
	ErrInvalidType = "invalid type for %s: expected %T, received %T"
	// ErrInvalidDescription is raised when the input description cannot be cast to stakingtypes.Description{}.
	ErrInvalidDescription = "invalid description: %v"
	// ErrInvalidCommission is raised when the input commission cannot be cast to stakingtypes.CommissionRates{}.
	ErrInvalidCommission = "invalid commission: %v"
)
