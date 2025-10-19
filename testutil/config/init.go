//go:build test

package config

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// init은 모든 테스트 패키지에서 자동으로 실행되어 Bech32 prefix를 설정합니다.
// 이 함수는 testutil/config 패키지를 import하는 모든 테스트 파일에서
// 테스트 실행 전에 자동으로 호출됩니다.
//
// 주의: 이 파일은 일반 .go 파일이므로 프로덕션 빌드에서도 포함됩니다.
// 하지만 init() 함수는 테스트에서만 사용되는 설정을 수행하므로
// 프로덕션 환경에서는 영향을 미치지 않습니다.
func init() {
	cfg := sdk.GetConfig()

	// 중요: EvmAppOptions가 SDK config를 reset할 수 있으므로
	// Bech32 prefix 설정을 EvmAppOptions 이후에 수행합니다
	_ = EvmAppOptions(EVMChainID)

	// Bech32 prefix 설정 (zenanet) - EvmAppOptions 이후에 설정
	SetBech32Prefixes(cfg)

	// SDK config를 seal하여 이후 변경을 방지
	cfg.Seal()
}
