//go:build test
// +build test

package config

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TestAddresses는 테스트에서 사용할 표준 Bech32 주소들을 제공합니다.
// 이 주소들은 올바른 zenanet prefix와 체크섬을 가집니다.
//
// 주소는 0xd33bFF38Bfc79df581815BED0c779FA99BaFAbf5 hex 주소에서 생성됩니다.
var (
	// TestAccAddr1은 표준 계정 주소입니다 (zenanet1...)
	TestAccAddr1 = mustBech32ify(Bech32PrefixAccAddr, hexToBytes("d33bFF38Bfc79df581815BED0c779FA99BaFAbf5"))

	// TestValAddr1은 validator operator 주소입니다 (zenanetvaloper1...)
	TestValAddr1 = mustBech32ify(Bech32PrefixValAddr, hexToBytes("d33bFF38Bfc79df581815BED0c779FA99BaFAbf5"))

	// TestConsAddr1은 consensus node 주소입니다 (zenanetvalcons1...)
	TestConsAddr1 = mustBech32ify(Bech32PrefixConsAddr, hexToBytes("d33bFF38Bfc79df581815BED0c779FA99BaFAbf5"))

	// TestHexAddr는 원본 hex 주소입니다
	TestHexAddr = "0xd33bFF38Bfc79df581815BED0c779FA99BaFAbf5"
)

// mustBech32ify는 주어진 prefix와 bytes를 Bech32 주소로 변환합니다.
// 변환 실패 시 panic합니다.
func mustBech32ify(prefix string, data []byte) string {
	addr, err := sdk.Bech32ifyAddressBytes(prefix, data)
	if err != nil {
		panic(err)
	}
	return addr
}

// hexToBytes는 hex string (0x 접두사 없이)을 bytes로 변환합니다.
func hexToBytes(hexStr string) []byte {
	// 0x 접두사 제거
	if len(hexStr) >= 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	return bytes
}

// GenerateTestAddress는 주어진 hex 주소에서 모든 타입의 Bech32 주소를 생성합니다.
func GenerateTestAddress(hexAddr string) (accAddr, valAddr, consAddr string) {
	addrBytes := hexToBytes(hexAddr)

	accAddr = mustBech32ify(Bech32PrefixAccAddr, addrBytes)
	valAddr = mustBech32ify(Bech32PrefixValAddr, addrBytes)
	consAddr = mustBech32ify(Bech32PrefixConsAddr, addrBytes)

	return
}
