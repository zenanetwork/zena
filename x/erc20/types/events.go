package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// erc20 events
const (
	EventTypeConvertERC20           = "convert_erc20"
	EventTypeRegisterERC20          = "register_erc20"
	EventTypeToggleTokenConversion  = "toggle_token_conversion" // #nosec
	EventTypeRegisterERC20Extension = "register_erc20_extension"

	EventTypeFailedConvertERC20 = "failed_convert_erc20"

	AttributeCoinSourceChannel = "source_channel"
	AttributeKeyCosmosCoin     = "cosmos_coin"
	AttributeKeyERC20Token     = "erc20_token" // #nosec
	AttributeKeyReceiver       = "receiver"
)

// LogTransfer Event type for Transfer(address from, address to, uint256 value)
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}
