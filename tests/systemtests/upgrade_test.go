//go:build system_test

package systemtests

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	systest "cosmossdk.io/systemtests"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	upgradeHeight int64 = 22
	upgradeName         = "v0.6.0-stabilization" // must match UpgradeName in zenad/upgrades.go
	nativeDenom         = "aznnt"
	bech32Prefix        = "zenanet"
)

// legacyActivePrecompiles is the active_static_precompiles list of the
// mainnet genesis, which enabled the unimplemented 0x0803 and 0x0807.
const legacyActivePrecompiles = `["0x0000000000000000000000000000000000000100",
"0x0000000000000000000000000000000000000400","0x0000000000000000000000000000000000000800",
"0x0000000000000000000000000000000000000801","0x0000000000000000000000000000000000000802",
"0x0000000000000000000000000000000000000803","0x0000000000000000000000000000000000000804",
"0x0000000000000000000000000000000000000805","0x0000000000000000000000000000000000000806",
"0x0000000000000000000000000000000000000807"]`

// legacyMainnetGenesis turns the genesis written by the legacy binary's
// testnet command into one shaped like the Zenanet mainnet genesis: the
// legacy command funds accounts with "atest" but only writes aznnt metadata,
// so every atest amount becomes aznnt, and the unimplemented precompiles are
// enabled as on mainnet.
func legacyMainnetGenesis(t *testing.T) systest.GenesisMutator {
	t.Helper()
	return func(genesis []byte) []byte {
		genesis = bytes.ReplaceAll(genesis, []byte(`"atest"`), []byte(fmt.Sprintf("%q", nativeDenom)))
		state, err := sjson.SetRawBytes(genesis, "app_state.evm.params.active_static_precompiles", []byte(legacyActivePrecompiles))
		require.NoError(t, err)
		return state
	}
}

func TestChainUpgrade(t *testing.T) {
	// Scenario:
	// start a legacy chain with some state
	// when a chain upgrade proposal is executed
	// then the chain upgrades successfully
	systest.Sut.StopChain()

	currentBranchBinary := systest.Sut.ExecBinary()
	currentInitializer := systest.Sut.TestnetInitializer()

	legacyBinary := systest.WorkDir + "/binaries/legacy/zenad"
	systest.Sut.SetExecBinary(legacyBinary)
	systest.Sut.SetTestnetInitializer(systest.InitializerWithBinary(legacyBinary, systest.Sut))
	systest.Sut.SetupChain()

	votingPeriod := 5 * time.Second // enough time to vote
	systest.Sut.ModifyGenesisJSON(t, systest.SetGovVotingPeriod(t, votingPeriod), legacyMainnetGenesis(t))

	systest.Sut.StartChain(t, fmt.Sprintf("--halt-height=%d", upgradeHeight+1), "--chain-id=local-4221", "--minimum-gas-prices=0.00"+nativeDenom)

	cli := systest.NewCLIWrapper(t, systest.Sut, systest.Verbose)
	govAddr, err := sdk.Bech32ifyAddressBytes(bech32Prefix, address.Module("gov"))
	require.NoError(t, err)
	// submit upgrade proposal
	proposal := fmt.Sprintf(`
{
 "messages": [
  {
   "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
   "authority": %q,
   "plan": {
    "name": %q,
    "height": "%d"
   }
  }
 ],
 "metadata": "ipfs://CID",
 "deposit": "100000000stake",
 "title": "my upgrade",
 "summary": "testing"
}`, govAddr, upgradeName, upgradeHeight)
	rsp := cli.SubmitGovProposal(proposal, "--fees=10000000000000000000"+nativeDenom, "--from=node0")
	systest.RequireTxSuccess(t, rsp)
	raw := cli.CustomQuery("q", "gov", "proposals", "--depositor", cli.GetKeyAddr("node0"))
	proposals := gjson.Get(raw, "proposals.#.id").Array()
	require.NotEmpty(t, proposals, raw)
	proposalID := proposals[len(proposals)-1].String()

	for i := range systest.Sut.NodesCount() {
		go func(i int) { // do parallel
			systest.Sut.Logf("Voting: validator %d\n", i)
			rsp := cli.Run("tx", "gov", "vote", proposalID, "yes", "--fees=10000000000000000000"+nativeDenom, "--from", cli.GetKeyAddr(fmt.Sprintf("node%d", i)))
			systest.RequireTxSuccess(t, rsp)
		}(i)
	}

	systest.Sut.AwaitBlockHeight(t, upgradeHeight-1, 3*time.Minute)
	t.Logf("current_height: %d\n", systest.Sut.CurrentHeight())
	raw = cli.CustomQuery("q", "gov", "proposal", proposalID)
	proposalStatus := gjson.Get(raw, "proposal.status").String()
	require.Equal(t, "PROPOSAL_STATUS_PASSED", proposalStatus, raw)

	t.Log("waiting for upgrade info")
	systest.Sut.AwaitUpgradeInfo(t)
	systest.Sut.StopChain()

	t.Log("Upgrade height was reached. Upgrading chain")
	systest.Sut.SetExecBinary(currentBranchBinary)
	systest.Sut.SetTestnetInitializer(currentInitializer)
	systest.Sut.StartChain(t, "--chain-id=local-4221")

	require.Equal(t, upgradeHeight+1, systest.Sut.CurrentHeight())

	// smoke test to make sure the chain still functions.
	cli = systest.NewCLIWrapper(t, systest.Sut, systest.Verbose)
	to := cli.GetKeyAddr("node1")
	from := cli.GetKeyAddr("node0")
	got := cli.Run("tx", "bank", "send", from, to, "1"+nativeDenom, "--from=node0", "--fees=10000000000000000000"+nativeDenom, "--chain-id=local-4221")
	systest.RequireTxSuccess(t, got)

	// the upgrade handler removed the unimplemented static precompiles
	raw = cli.CustomQuery("q", "evm", "params")
	active := gjson.Get(raw, "params.active_static_precompiles").String()
	require.NotContains(t, active, "0x0000000000000000000000000000000000000803", raw)
	require.NotContains(t, active, "0x0000000000000000000000000000000000000807", raw)
	require.Contains(t, active, "0x0000000000000000000000000000000000000802", raw)

	// the native denom metadata is untouched
	raw = cli.CustomQuery("q", "bank", "denom-metadata", nativeDenom)
	require.Equal(t, "ZNNT", gjson.Get(raw, "metadata.symbol").String(), raw)
}
