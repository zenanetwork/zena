package contracts

import (
	_ "embed"

	contractutils "github.com/cosmos/evm/contracts/utils"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

var (
	// WZENAJSON are the compiled bytes of the WZENAContract
	//
	//go:embed solidity/WZENA.json
	WZENAJSON []byte

	// WATOMContract is the compiled watom contract
	WZENAContract evmtypes.CompiledContract
)

func init() {
	var err error
	if WZENAContract, err = contractutils.ConvertHardhatBytesToCompiledContract(
		WZENAJSON,
	); err != nil {
		panic(err)
	}
}
