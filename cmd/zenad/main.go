package main

import (
	"fmt"
	"os"

	"github.com/zenanetwork/zena/cmd/zenad/cmd"
	evmdconfig "github.com/zenanetwork/zena/cmd/zenad/config"
	examplechain "github.com/zenanetwork/zena/zenad"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {
	setupSDKConfig()

	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "zenad", examplechain.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}

func setupSDKConfig() {
	config := sdk.GetConfig()
	evmdconfig.SetBech32Prefixes(config)
	config.Seal()
}
