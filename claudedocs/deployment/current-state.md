# Zena Mainnet Deployment — Current State

**작성일**: 2026-04-15
**목적**: `release/v0.6.0` 업그레이드 준비를 위해 현재 mainnet 운영 상태를 추적

이 문서는 "빈 슬롯"으로 시작합니다. Validator 운영자/팀 내부 기록으로부터 정보를 수집해서 채워 넣으세요. 모든 슬롯이 채워져야 v0.5.0 소급 태그 + v0.6.0 업그레이드 진행 가능.

---

## 1. 현재 운영 바이너리 식별 (🔴 가장 긴급)

### 목표
Mainnet validator들이 현재 돌리는 `zenad` 바이너리가 **어느 commit hash에서 빌드되었는지** 파악.

### 수집 방법

**방법 A**: Validator 운영자에게 다음 명령 실행 요청
```bash
zenad version --long 2>&1 | grep -E "^(version|commit):"
```
출력 예:
```
version: 0.4.0
commit: abc1234def5678...
```

**방법 B**: 팀 내부 빌드 기록 확인
- [ ] CI 파이프라인(GitHub Actions)의 workflow run 로그
- [ ] Docker Hub / Container Registry에 푸시된 이미지 태그
- [ ] 팀 서버의 `~/.bash_history` / `~/.zsh_history`에서 `make build` 시각
- [ ] Notion / Slack / Discord의 배포 공지

**방법 C**: 바이너리 sha256로 역추적
```bash
sha256sum /path/to/zenad-binary
# 이 해시와 매치되는 과거 CI 빌드 아티팩트를 찾기
```

### 수집 결과 (↓ 채워 넣기)

| Validator | 운영자 연락처 | `zenad version` | commit hash | 확인 일자 |
|---|---|---|---|---|
| validator-1 | | | | |
| validator-2 | | | | |
| validator-3 | | | | |

**모든 validator가 동일 commit에서 빌드되었는가?**
- [ ] Yes — 해당 commit: `__________`
- [ ] No — validator 간 버전 상이. 긴급 조사 필요 (소프트 포크 위험)

---

## 2. 체인 런타임 상태

### 수집 명령
```bash
# Tendermint/CometBFT 상태
curl -s <mainnet-rpc>/status | jq '.result.sync_info'

# Validator 목록
curl -s <mainnet-rpc>/validators | jq '.result.validators | length'

# ABCI 정보 (app version 확인)
curl -s <mainnet-rpc>/abci_info | jq
```

### 수집 결과 (↓ 채워 넣기)

| 항목 | 값 |
|---|---|
| Mainnet RPC endpoint | `https://____________` |
| Chain ID | (예: `zena_15150-1`) |
| 현재 block height | |
| 활성 validator 수 | |
| App version (abci_info) | |
| 최신 블록 시각 | |

---

## 3. x/upgrade 모듈 사용 가능 여부 (🟡 중요)

### 확인 방법
```bash
# 현재 바이너리로 이 명령이 응답하면 x/upgrade가 wired 상태
zenad query upgrade module-versions --node <mainnet-rpc>

# 과거 업그레이드가 적용된 적 있는지 확인
zenad query upgrade applied v0.4.0-to-v0.5.0 --node <mainnet-rpc>
```

### 수집 결과

- [ ] `module-versions` 응답 정상 → ✅ Governance upgrade 경로 사용 가능 (Phase 6)
- [ ] 응답 없음/에러 → ⚠️ Coordinated Hard Fork 필요 (Phase 6 Contingency)

**`v0.4.0-to-v0.5.0` 업그레이드가 적용된 적 있는가?**
- [ ] Yes — 체인에 `atest` 더미 denom metadata가 박혀있을 수 있음. v0.6.0 핸들러에 정리 로직 검토 필요.
- [ ] No — 체인 상태 깨끗. v0.6.0 업그레이드 시 추가 마이그레이션 불필요.

---

## 4. Validator 커뮤니케이션 채널

v0.6.0 업그레이드 공지 시 사용할 채널:

- [ ] Discord 서버: `____________`
- [ ] Telegram 그룹: `____________`
- [ ] 이메일 배포 리스트: `____________`
- [ ] 공식 공지 페이지: `____________`

**예상 공지 리드타임**: upgrade height 설정은 공지 후 최소 **7일** 권장 (cosmovisor 미사용 validator의 수동 대응 시간 확보).

---

## 5. 리스크 체크리스트

업그레이드 실행 전 반드시 확인:

- [ ] 모든 validator가 동일 commit으로 운영 중 (§1)
- [ ] `x/upgrade` 모듈 wired 확인 (§3)
- [ ] v0.5.0 소급 태그가 GitHub에 push됨
- [ ] v0.6.0 바이너리의 sha256 체크섬이 공개됨
- [ ] 테스트넷에서 v0.5.0 → v0.6.0 업그레이드 리허설 성공 (Plan Phase 5)
- [ ] 외부 보안 감사 완료 또는 최소 self-audit 리포트 완료
- [ ] Validator 사전 공지 7일 이상 완료
- [ ] 업그레이드 height 결정 (권장: 공지 시점 + 7일 + 버퍼)
- [ ] 롤백 계획 문서화 (업그레이드 실패 시 절차)

---

## 6. 다음 단계 (이 문서 완성 후)

1. **Plan Phase 1 진행**: §1의 commit에 `v0.5.0` 태그 부여 + GitHub Release
2. **Plan Phase 5 진행**: 테스트넷 리허설
3. **외부 감사 의뢰**: Oak Security / Zellic / Halborn 중 선정 (commit hash 지정 필요)
4. **Plan Phase 6 진행**: Mainnet governance upgrade (성공 조건 충족 후)

---

## Appendix — 정보 수집용 메시지 템플릿

Validator 운영자에게 보낼 메시지 초안:

```
안녕하세요, Zena 팀입니다.

v0.5.0 → v0.6.0 업그레이드 준비를 위해 현재 운영 중인 바이너리 정보를
수집하고 있습니다. 아래 명령을 실행한 결과를 회신해 주시면 감사하겠습니다.

실행 명령:
    zenad version --long

또는 간단히:
    zenad version --long 2>&1 | grep -E "^(version|commit|build_tags):"

보내주실 정보:
    - version: ...
    - commit: ...
    - build_tags: ...

수집 마감: YYYY-MM-DD
문의: <team-contact>

이 정보는 업그레이드 제안서 작성과 보안 감사 의뢰에 필수적입니다.
감사합니다.
```
