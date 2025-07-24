# Zena EVM Chain Setup Guide

이 가이드는 Zena EVM 체인을 완전히 설정하고 운영하기 위한 bash 스크립트(`setup_zena_evm.sh`)의 사용법을 설명합니다.

## 🚀 주요 기능

### ✅ 1. EVM 기능 완전 활성화
- 모든 precompile 활성화 (10개 전체)
- P256, Bech32, Staking, Distribution, ICS20, Vesting, Bank, Gov, Slashing, Evidence
- EVM 체인 ID 설정 (기본값: 4221)
- 이더리움 호환 트랜잭션 지원

### ✅ 2. 이더리움 주소 표시
- 모든 키에 대해 Cosmos 주소와 Ethereum 주소 동시 표시
- `addresses` 명령어로 언제든지 확인 가능
- eth_secp256k1 키 알고리즘 사용

### ✅ 3. 백그라운드 실행
- `start-bg` 명령어로 백그라운드 실행
- PID 파일 관리 및 프로세스 모니터링
- 로그 파일 자동 생성 및 관리
- 안전한 종료 및 재시작 기능

### ✅ 4. 전체 기능 활성화
- JSON-RPC API 전체 활성화 (eth, txpool, personal, net, debug, web3, miner)
- WebSocket 지원
- REST API 및 gRPC 활성화
- 메트릭스 수집 활성화
- 트랜잭션 인덱서 활성화

### ✅ 5. 실제 운영 환경 설정
- 파일 키링 사용 (보안 강화)
- 적절한 가스 설정 및 베이스 수수료
- 프로덕션 환경에 맞는 포트 설정
- 로그 및 모니터링 완전 구성

## 📋 사전 요구사항

### 필수 소프트웨어
```bash
# jq 설치
sudo apt-get install jq  # Ubuntu/Debian
brew install jq          # macOS

# curl 설치 (대부분 기본 설치됨)
sudo apt-get install curl  # Ubuntu/Debian
brew install curl          # macOS

# Go 1.21+ 설치
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### 프로젝트 빌드
```bash
# 프로젝트 디렉토리에서
make install
```

## ⚠️ 실제 운영 환경 주의사항

### 1. 키 보안 관리
- 스크립트는 빠른 설정을 위해 미리 정의된 mnemonic을 사용합니다
- **🔐 실제 운영 환경에서는 반드시 새로운 키를 생성하세요:**
```bash
zenad keys add validator --keyring-backend file
zenad keys add user1 --keyring-backend file
zenad keys add user2 --keyring-backend file
```

### 2. 네트워크 설정
- 기본 체인 ID: 4221 (Zena EVM 메인넷)
- 기본 denomination: azena (atomic unit)
- 파일 키링 사용으로 보안 강화
- 최소 가스 가격: 0.001azena

### 3. 브랜딩 설정
- 모든 denomination이 "zena" 브랜딩으로 설정됨
- DENOM: azena (atomic unit)
- DISPLAY_DENOM: zena (display unit)
- 주소 prefix: zenanet

### 4. 보안 설정
- 파일 키링 사용 (--keyring file)
- 최소 가스 가격 설정
- 연결 제한 설정 (max-open-connections: 2000)
- 트랜잭션 제한 설정 (max-txs-per-conn: 200)

## 🔧 설치 및 기본 설정

### 1. 스크립트 다운로드 및 권한 설정
```bash
# 실행 권한 부여
chmod +x setup_zena_evm.sh

# 사용법 확인
./setup_zena_evm.sh --help
```

### 2. 체인 초기화
```bash
# 기본 설정으로 초기화
./setup_zena_evm.sh init

# 사용자 정의 설정으로 초기화
./setup_zena_evm.sh --chain-id 9001 --moniker "my-node" init
```

### 3. 노드 실행

#### 포그라운드 실행 (개발용)
```bash
./setup_zena_evm.sh start
```

#### 백그라운드 실행 (운영용)
```bash
./setup_zena_evm.sh start-bg
```

## 📖 명령어 상세 사용법

### 기본 명령어
```bash
# 체인 초기화
./setup_zena_evm.sh init

# 포그라운드에서 노드 시작
./setup_zena_evm.sh start

# 백그라운드에서 노드 시작
./setup_zena_evm.sh start-bg

# 노드 중지
./setup_zena_evm.sh stop

# 노드 재시작
./setup_zena_evm.sh restart

# 노드 상태 확인
./setup_zena_evm.sh status

# 이더리움 주소 확인
./setup_zena_evm.sh addresses

# 로그 확인
./setup_zena_evm.sh logs

# 모든 데이터 삭제 (주의!)
./setup_zena_evm.sh clean

# 데이터 백업
./setup_zena_evm.sh backup

# 데이터 복원
./setup_zena_evm.sh restore

# 버전 정보
./setup_zena_evm.sh version
```

### 설정 옵션
```bash
# 체인 ID 설정
./setup_zena_evm.sh --chain-id 9001 init

# 노드 모니커 설정
./setup_zena_evm.sh --moniker "my-validator" init

# 홈 디렉토리 설정
./setup_zena_evm.sh --home /custom/path init

# 키링 백엔드 설정
./setup_zena_evm.sh --keyring file init

# 로그 레벨 설정
./setup_zena_evm.sh --log-level debug init

# 바이너리 설치 건너뛰기
./setup_zena_evm.sh --no-install init
```

## 🌐 네트워크 엔드포인트

체인이 실행되면 다음 엔드포인트들을 사용할 수 있습니다:

### JSON-RPC (이더리움 호환)
```bash
# 기본 JSON-RPC 엔드포인트
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
  http://localhost:8545

# 계정 잔액 조회
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x...", "latest"],"id":1}' \
  http://localhost:8545
```

### WebSocket
```bash
# WebSocket 연결
wscat -c ws://localhost:8546
```

### REST API
```bash
# 블록 높이 조회
curl http://localhost:1317/cosmos/base/tendermint/v1beta1/blocks/latest

# 계정 정보 조회
curl http://localhost:1317/cosmos/auth/v1beta1/accounts/{address}
```

### gRPC
```bash
# grpcurl 사용 예제
grpcurl -plaintext localhost:9090 cosmos.base.tendermint.v1beta1.Service/GetLatestBlock
```

## 💡 이더리움 주소 사용법

### 1. 키 생성 및 주소 확인
```bash
# 체인 초기화 후 기본 키들이 생성됨
./setup_zena_evm.sh addresses

# 출력 예시:
# validator:
#   Cosmos:   cosmos1abc...
#   Ethereum: 0x123...
# user1:
#   Cosmos:   cosmos1def...
#   Ethereum: 0x456...
```

### 2. 새 키 생성
```bash
# 새 키 생성
zenad keys add newkey --keyring-backend file --algo eth_secp256k1 --home ~/.zenad

# 이더리움 주소 확인
zenad debug addr $(zenad keys show newkey -a --keyring-backend file --home ~/.zenad)
```

### 3. MetaMask 연결
```bash
# MetaMask 네트워크 설정
Network Name: Zena EVM
RPC URL: http://localhost:8545
Chain ID: 4221
Currency Symbol: ZENA
Block Explorer: (선택사항)
```

## 🛠️ 고급 설정

### 환경 변수 설정
```bash
# 체인 ID 설정
export CHAIN_ID=9001

# 모니커 설정
export MONIKER="my-validator"

# 홈 디렉토리 설정
export HOMEDIR="/custom/path"

# 키링 백엔드 설정
export KEYRING="file"

# 로그 레벨 설정
export LOGLEVEL="debug"

# 베이스 수수료 설정
export BASEFEE="20000000000000000"

# 포트 설정
export RPC_PORT=26657
export JSONRPC_PORT=8545
export WS_PORT=8546
```

### 프로덕션 환경 설정
```bash
# 프로덕션 환경에서는 다음 설정을 권장합니다:
export KEYRING="file"           # 파일 키링 사용
export LOGLEVEL="info"          # 적절한 로그 레벨
export BASEFEE="10000000000000000"  # 적절한 베이스 수수료

# 방화벽 설정 (예시)
sudo ufw allow 26657/tcp  # RPC
sudo ufw allow 8545/tcp   # JSON-RPC
sudo ufw allow 8546/tcp   # WebSocket
```

## 🔍 모니터링 및 디버깅

### 로그 확인
```bash
# 실시간 로그 모니터링
./setup_zena_evm.sh logs

# 또는 직접 로그 파일 확인
tail -f ~/.zenad/logs/zenad.log
```

### 노드 상태 확인
```bash
# 전체 상태 확인
./setup_zena_evm.sh status

# 블록 동기화 확인
curl -s http://localhost:26657/status | jq '.result.sync_info'

# 피어 연결 확인
curl -s http://localhost:26657/net_info | jq '.result.peers'
```

### 메트릭스 확인
```bash
# Prometheus 메트릭스 확인
curl http://localhost:26660/metrics
```

## 🚨 문제 해결

### 일반적인 문제들

#### 1. 포트 충돌
```bash
# 포트 사용 중 확인
lsof -i :8545
lsof -i :26657

# 다른 포트 사용
export JSONRPC_PORT=8555
export RPC_PORT=26667
```

#### 2. 키링 문제
```bash
# 키링 백엔드 확인
zenad keys list --keyring-backend file --home ~/.zenad

# 키링 문제 시 재생성
rm -rf ~/.zenad/keyring-file
./setup_zena_evm.sh init
```

#### 3. 제네시스 검증 실패
```bash
# 제네시스 파일 재생성
./setup_zena_evm.sh clean
./setup_zena_evm.sh init
```

#### 4. 메모리 부족
```bash
# 프루닝 설정 조정 (app.toml)
pruning = "custom"
pruning-keep-recent = "100"
pruning-interval = "10"
```

### 로그 분석
```bash
# 오류 로그 필터링
grep -i error ~/.zenad/logs/zenad.log

# 경고 로그 필터링
grep -i warning ~/.zenad/logs/zenad.log

# 특정 모듈 로그 필터링
grep -i "module=evm" ~/.zenad/logs/zenad.log
```

## 🔒 보안 고려사항

### 1. 키링 보안
```bash
# 프로덕션에서는 항상 file 키링 사용
export KEYRING="file"

# 키링 디렉토리 권한 설정
chmod 700 ~/.zenad/keyring-file
```

### 2. 네트워크 보안
```bash
# 방화벽 설정
sudo ufw enable
sudo ufw allow 22/tcp      # SSH
sudo ufw allow 26657/tcp   # RPC (필요시)
sudo ufw allow 8545/tcp    # JSON-RPC (필요시)
```

### 3. 로그 보안
```bash
# 로그 파일 권한 설정
chmod 600 ~/.zenad/logs/zenad.log

# 로그 로테이션 설정
sudo logrotate -f /etc/logrotate.conf
```

## 📚 추가 자료

### 공식 문서
- [Cosmos EVM 공식 문서](https://evm.cosmos.network/)
- [Cosmos SDK 문서](https://docs.cosmos.network/)
- [Ethereum JSON-RPC API](https://ethereum.github.io/execution-apis/api-documentation/)

### 커뮤니티 및 지원
- [GitHub 이슈](https://github.com/cosmos/cosmos-sdk/issues)
- [Discord 커뮤니티](https://discord.gg/cosmos)
- [공식 포럼](https://forum.cosmos.network/)

### 개발 도구
- [Hardhat](https://hardhat.org/) - 스마트 계약 개발
- [Remix](https://remix.ethereum.org/) - 온라인 IDE
- [MetaMask](https://metamask.io/) - 브라우저 지갑

## 🎯 성능 최적화

### 하드웨어 권장사항
```bash
# 최소 요구사항
CPU: 4 코어
RAM: 8GB
Storage: 100GB SSD

# 권장 사항
CPU: 8 코어
RAM: 16GB
Storage: 500GB NVMe SSD
```

### 설정 최적화
```bash
# config.toml 최적화
create_empty_blocks = false
skip_timeout_commit = false
timeout_commit = "1s"

# app.toml 최적화
pruning = "custom"
pruning-keep-recent = "100"
pruning-keep-every = "0"
pruning-interval = "10"
```

이 가이드를 따라 진행하면 Zena EVM 체인을 완전히 설정하고 운영할 수 있습니다. 추가 질문이나 문제가 있으시면 언제든지 문의해 주세요! 

---

## 🚀 메인넷 환경 배포 가이드

### 1. 프로덕션 환경 설정
```bash
# 프로덕션 환경 초기화
./setup_zena_evm.sh --keyring file --moniker "zena-mainnet-node" init

# 새로운 키 생성 (보안 강화)
zenad keys add validator --keyring-backend file
zenad keys add operator --keyring-backend file

# 백업 생성
./setup_zena_evm.sh backup
```

### 2. 시스템 서비스 등록 (선택사항)
```bash
# systemd 서비스 파일 생성
sudo tee /etc/systemd/system/zenad.service > /dev/null <<EOF
[Unit]
Description=Zena EVM Node
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$HOME/workspace/blockchain/cosmos-sdk-evm-v2/zena
ExecStart=$HOME/workspace/blockchain/cosmos-sdk-evm-v2/zena/setup_zena_evm.sh --keyring file start
Restart=always
RestartSec=3
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF

# 서비스 활성화
sudo systemctl daemon-reload
sudo systemctl enable zenad
sudo systemctl start zenad
```

### 3. 모니터링 설정
```bash
# 프로메테우스 메트릭 확인
curl http://localhost:26660/metrics

# 상태 모니터링
./setup_zena_evm.sh status

# 로그 모니터링
tail -f ~/.zenad/logs/zenad.log
```

---

## 📋 메인넷 배포 체크리스트

### 배포 전 확인사항
- [ ] 새로운 키 생성 완료
- [ ] 백업 생성 완료
- [ ] 보안 설정 확인 (file keyring 사용)
- [ ] 네트워크 포트 설정 확인
- [ ] 모니터링 설정 완료
- [ ] 시스템 리소스 충분성 확인
- [ ] 방화벽 설정 확인

### 정기 유지보수
- [ ] 주기적 백업 생성
- [ ] 로그 정리
- [ ] 시스템 업데이트
- [ ] 성능 모니터링
- [ ] 보안 패치 적용

---

**⚠️ 메인넷 환경 주의사항:**
- 🔐 실제 운영 환경에서는 반드시 새로운 키를 생성하세요
- 💾 백업을 정기적으로 생성하고 안전한 곳에 보관하세요
- �� 보안을 위해 파일 키링을 사용하세요
- 🌐 네트워크 보안 설정을 적절히 구성하세요
- 📊 모니터링을 통해 노드 상태를 지속적으로 확인하세요

이 스크립트를 사용하여 Zena EVM 체인을 성공적으로 설정하고 운영할 수 있습니다!

