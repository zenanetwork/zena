package contracts

import (
	_ "embed"

	contractutils "github.com/zenanetwork/zena/contracts/utils"
	evmtypes "github.com/zenanetwork/zena/x/vm/types"
)

var (
	// WATOMJSON are the compiled bytes of the WATOMContract
	//
	//go:embed solidity/WATOM.json
	WATOMJSON []byte

	// WATOMContract is the compiled watom contract
	WATOMContract evmtypes.CompiledContract
)

func init() {
	var err error
	if WATOMContract, err = contractutils.ConvertHardhatBytesToCompiledContract(
		WATOMJSON,
	); err != nil {
		panic(err)
	}
}
