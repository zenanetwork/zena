package types

const (
	P256PrecompileAddress   = "0x0000000000000000000000000000000000000100"
	Bech32PrecompileAddress = "0x0000000000000000000000000000000000000400"
)

const (
	StakingPrecompileAddress      = "0x0000000000000000000000000000000000000800"
	DistributionPrecompileAddress = "0x0000000000000000000000000000000000000801"
	ICS20PrecompileAddress        = "0x0000000000000000000000000000000000000802"
	BankPrecompileAddress         = "0x0000000000000000000000000000000000000804"
	GovPrecompileAddress          = "0x0000000000000000000000000000000000000805"
	SlashingPrecompileAddress     = "0x0000000000000000000000000000000000000806"
)

const (
	// VestingPrecompileAddress is the address reserved by upstream Cosmos EVM for the
	// vesting precompile. Zenanet does not ship an implementation for it.
	VestingPrecompileAddress = "0x0000000000000000000000000000000000000803"
	// ICS02PrecompileAddress is the address used by newer upstream Cosmos EVM releases for
	// the ICS-02 client precompile. Zenanet does not ship an implementation for it.
	ICS02PrecompileAddress = "0x0000000000000000000000000000000000000807"
)

// AvailableStaticPrecompiles defines the full list of all available EVM extension addresses.
//
// NOTE: To be explicit, this list does not include the dynamically registered EVM extensions
// like the ERC-20 extensions.
var AvailableStaticPrecompiles = []string{
	P256PrecompileAddress,
	Bech32PrecompileAddress,
	StakingPrecompileAddress,
	DistributionPrecompileAddress,
	ICS20PrecompileAddress,
	BankPrecompileAddress,
	GovPrecompileAddress,
	SlashingPrecompileAddress,
}

// UnimplementedStaticPrecompiles lists precompile addresses that have no implementation in
// this binary but were enabled in active_static_precompiles by earlier genesis scripts.
// They must never be active: calling an active address without an implementation panics
// in Keeper.GetStaticPrecompileInstance.
var UnimplementedStaticPrecompiles = []string{
	VestingPrecompileAddress,
	ICS02PrecompileAddress,
}
