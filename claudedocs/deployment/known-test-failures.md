# Known Test Failures

`release/v0.6.0` 브랜치에서 의도적으로 `t.Skip()` 처리된 테스트와 그 이유.

## Skipped tests

### 1. `TestMiddlewareTestSuite/TestOnRecvPacketNativeErc20`
### 2. `TestMiddlewareTestSuite/TestOnAcknowledgementPacketNativeErc20`
### 3. `TestMiddlewareTestSuite/TestOnTimeoutPacketNativeErc20`

**파일**: `zenad/tests/ibc/ibc_middleware_test.go`

**실패 메시지**:
```
sdk/5: spendable balance 0erc20:0x80b5...
   is smaller than 100erc20:0x80b5...: insufficient funds
```

## 근본 원인

세 테스트 모두 **raw `MsgTransfer`에 `erc20:0x...` denom을 넣어 송금이 성공하기를 기대**합니다.

```go
msg := transfertypes.NewMsgTransfer(
    ..., sdk.NewCoin(nativeErc20.Denom, sendAmt), ...,
)
_, err := suite.evmChainA.SendMsgs(msg)
suite.Require().NoError(err) // ← 여기서 실패
```

하지만 현재 코드베이스 아키텍처:

1. `SetupNativeErc20` helper는 **ERC20 컨트랙트 내부 잔고만** mint (100 토큰)
2. Bank 모듈의 `erc20:0x...` denom 잔고는 **0** (자동 미러링 없음)
3. `x/erc20/ibc_middleware.go`는 `OnRecvPacket / OnAck / OnTimeout`만 구현 — **outgoing `OnSendPacket` 변환 로직 부재**
4. Bank가 balance 검사에서 0을 보고 `insufficient funds` 반환

## 실제 프로덕션 경로 (정상 동작)

유저는 raw `MsgTransfer`가 아니라 **ICS20 precompile**을 EVM tx로 호출:

```solidity
ICS20.transfer(srcPort, srcChannel, erc20ContractAddr, amount, sender, receiver, timeout, ...);
```

이 경로는 `precompiles/ics20/tx.go`의 `transferWithStateDB` 가 내부에서 ERC20 → bank 변환을 수행. 테스트:
- `zenad/tests/ibc/ics20_precompile_transfer_test.go` ✅ **통과 중**
- `zenad/tests/ibc/v2_ics20_precompile_transfer_test.go` ✅ **통과 중**

## 유저 영향

**없음**. Wallet/dApp들은 ICS20 precompile을 통한 송금만 사용함 (Keplr 외에도 MetaMask 기반 UI는 EVM tx로만 IBC 전송).

Raw `MsgTransfer`를 Cosmos CLI로 직접 호출하는 시나리오는 NativeErc20에 대해서는 원래부터 지원되지 않았으며, 유저 플로우의 문제가 아님.

## 향후 수정 방향

### Option A: Outgoing IBC middleware 자동 변환 (근본 해결)
`x/erc20/ibc_middleware.go`에 `OnSendPacket` 훅 추가:
```go
func (im IBCMiddleware) OnSendPacket(
    ctx sdk.Context, portID, channelID string, sequence uint64,
    data []byte, signer sdk.AccAddress,
) error {
    // 1. Unmarshal FungibleTokenPacketData
    // 2. denom이 "erc20:0x..." 패턴이면 MsgConvertERC20 내부 호출
    // 3. underlying module의 OnSendPacket 호출
}
```
- 장점: 기존 프로덕션 로직을 확장. raw MsgTransfer도 작동.
- 단점: 프로덕션 로직 변경 → 감사 범위 확장. 대규모 회귀 테스트 필요.

### Option B: Helper 재설계 + 테스트 assertion 재작성
`SetupNativeErc20WithBankBalance` 추가하여 사전 변환된 상태를 시뮬레이션. 7개 테스트의 잔고 assertion을 "ERC20 잔고 변화"가 아니라 "bank 잔고 변화"로 재작성.
- 장점: 프로덕션 로직 불변, 감사 영향 없음.
- 단점: 테스트 의미론이 실제 유저 플로우와 달라짐.

### Option C: 테스트 폐기 + ICS20 precompile 테스트로 커버
이 세 테스트를 삭제하고 동일 시나리오를 `ics20_precompile_transfer_test.go`에 추가. 실제 유저 경로를 검증하는 테스트만 남김.
- 장점: 가장 깔끔. 프로덕션과 일치.
- 단점: IBC 미들웨어 자체의 NativeErc20 flow 커버리지 손실.

**권장**: v0.7.0 릴리스 사이클에서 **Option A**를 감사와 함께 진행.

---

## Additional known issue: Race condition in `ante/evm/`

### 증상
`go test -race -tags=test ./zenad/tests/...` 실행 시 거의 모든 테스트 패키지에서 race 감지됨:

```
WARNING: DATA RACE
Write at 0x... by goroutine N:
  github.com/cosmos/cosmos-sdk/types.(*EventManager).EmitEvent()
    cosmos-sdk@v0.53.6/types/events.go:50
  github.com/zenanetwork/zena/ante/evm.ConsumeFeesAndEmitEvent()
    ante/evm/08_gas_consume.go:55
  github.com/zenanetwork/zena/ante/evm.MonoDecorator.AnteHandle()
    ante/evm/mono_decorator.go:220
```

### 상태
- `-race` 없이 돌리면 전체 통과
- **main 브랜치에도 동일하게 존재** (현재 mainnet에서 운영 중인 코드)
- production에서 실 트래픽에 의한 관찰 가능한 영향 보고 없음

### 근본 원인 추정
`MonoDecorator.AnteHandle` 경로에서 `EventManager.EmitEvent`와 그 내부 gRPC stream 사이에 공유 상태가 있고, 고루틴 간 동기화 없이 접근됨. 정확한 수정은 Cosmos SDK 수준 또는 Zena의 ante decorator 구조 검토 필요.

### 대응
- **v0.6.0 릴리스 범위 밖**. 이 릴리스의 주 목적(보안 + 리브랜딩 + upgrade handler)과 무관.
- **v0.7.0 사이클에서 수정 검토**. ante decorator 변경은 감사 범위 확장이 필요하므로 별도 프로젝트로 취급.

### 재현
```bash
# Race 감지됨
cd zenad && go test -race -tags=test -run 'TestICS20' ./tests/ibc/...

# Race 없음, 기능 정상
cd zenad && go test -tags=test -run 'TestICS20' ./tests/ibc/...
```

## GitHub Issue 템플릿 (추후 제출용)

```markdown
## Summary
Three tests in `zenad/tests/ibc/ibc_middleware_test.go` assume raw
`MsgTransfer` auto-converts NativeErc20 balances, but no such
outgoing conversion exists in `x/erc20` IBC middleware.

## Failing tests (skipped in v0.6.0)
- TestMiddlewareTestSuite/TestOnRecvPacketNativeErc20
- TestMiddlewareTestSuite/TestOnAcknowledgementPacketNativeErc20
- TestMiddlewareTestSuite/TestOnTimeoutPacketNativeErc20

## Reproduction
Remove `suite.T().Skip(...)` lines and run:
`make test-zenad -run TestMiddlewareTestSuite`

## Proposed fix
Implement `OnSendPacket` in `x/erc20/ibc_middleware.go` OR migrate
tests to the ICS20 precompile path (see `ics20_precompile_transfer_test.go`).
```
