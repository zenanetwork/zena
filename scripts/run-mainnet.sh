#!/bin/bash

# zena 메인넷 시작 스크립트
# 사용법: ./scripts/run-mainnet.sh [--reset-data] [--init-validator]

set -e

# 메인넷 설정 (zenad/constants.go에서 정의된 값과 코스모스 체인 명명 규칙 준수)
CHAINID="zena_1-1"    # 메인넷 체인 ID 형식: {이름}_{버전}-{네트워크번호}
MONIKER=${MONIKER:-"zena-mainnet-node"}
KEYRING=${KEYRING:-"file"}  # 실제 메인넷은 file 타입 사용 (test는 개발용)
KEYALGO="eth_secp256k1"
LOGLEVEL=${LOGLEVEL:-"info"}
HOMEDIR=${HOMEDIR:-"${HOME}/.zenad"}
MIN_GAS_PRICES=${MIN_GAS_PRICES:-"0.0001azena"}
PRUNING=${PRUNING:-"custom"}
PRUNING_KEEP_RECENT=${PRUNING_KEEP_RECENT:-"100"}
PRUNING_INTERVAL=${PRUNING_INTERVAL:-"10"}
TRACE=${TRACE:-""}

# 설정 파일 경로
CONFIG=$HOMEDIR/config/config.toml
APP_TOML=$HOMEDIR/config/app.toml
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json

# 의존성 검사
command -v jq >/dev/null 2>&1 || {
  echo >&2 "jq가 설치되어 있지 않습니다. 설치해주세요: https://stedolan.github.io/jq/download/"
  exit 1
}

# 인자 처리
RESET_DATA=false
INIT_VALIDATOR=false

while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --reset-data)
      RESET_DATA=true
      shift
      ;;
    --init-validator)
      INIT_VALIDATOR=true
      shift
      ;;
    *)
      echo "알 수 없는 옵션: $key"
      echo "사용법: ./scripts/run-mainnet.sh [--reset-data] [--init-validator]"
      exit 1
      ;;
  esac
done

# 바이너리 확인 및 빌드
if [ ! -f "./build/zenad" ]; then
  echo "zenad 바이너리를 찾을 수 없습니다. 빌드합니다..."
  make build
fi

# 초기화
if [ "$RESET_DATA" = true ] || [ ! -d "$HOMEDIR" ]; then
  echo "노드 데이터를 초기화합니다..."
  
  # 이전 데이터가 있으면 백업
  if [ -d "$HOMEDIR" ]; then
    BACKUP_DIR="$HOMEDIR.backup.$(date +%Y%m%d%H%M%S)"
    echo "기존 데이터를 백업합니다: $BACKUP_DIR"
    mv "$HOMEDIR" "$BACKUP_DIR"
  fi
  
  # 메인넷 초기화
  ./build/zenad init "$MONIKER" --chain-id "$CHAINID" --home "$HOMEDIR"
  
  # 클라이언트 설정
  ./build/zenad config set client chain-id "$CHAINID" --home "$HOMEDIR"
  ./build/zenad config set client keyring-backend "$KEYRING" --home "$HOMEDIR"
  
  # 메인넷 제네시스 파일 다운로드 안내
  echo "메인넷 제네시스 파일을 다운로드하세요..."
  echo "제네시스 파일을 $HOMEDIR/config/genesis.json 위치에 저장하세요."
  echo "다운로드 후 엔터 키를 누르세요..."
  read -r
  
  # 제네시스 파일 검증
  ./build/zenad genesis validate-genesis --home "$HOMEDIR"
  
  # 기본 설정 변경
  if [[ "$OSTYPE" == "darwin"* ]]; then
    # 프루닝 설정
    sed -i '' "s/pruning = \"default\"/pruning = \"$PRUNING\"/g" "$APP_TOML"
    sed -i '' "s/pruning-keep-recent = \"0\"/pruning-keep-recent = \"$PRUNING_KEEP_RECENT\"/g" "$APP_TOML"
    sed -i '' "s/pruning-interval = \"0\"/pruning-interval = \"$PRUNING_INTERVAL\"/g" "$APP_TOML"
    
    # API 활성화
    sed -i '' 's/enable = false/enable = true/g' "$APP_TOML"
  else
    # 프루닝 설정
    sed -i "s/pruning = \"default\"/pruning = \"$PRUNING\"/g" "$APP_TOML"
    sed -i "s/pruning-keep-recent = \"0\"/pruning-keep-recent = \"$PRUNING_KEEP_RECENT\"/g" "$APP_TOML"
    sed -i "s/pruning-interval = \"0\"/pruning-interval = \"$PRUNING_INTERVAL\"/g" "$APP_TOML"
    
    # API 활성화
    sed -i 's/enable = false/enable = true/g' "$APP_TOML"
  fi
  
  # 밸리데이터 초기화 (선택적)
  if [ "$INIT_VALIDATOR" = true ]; then
    echo "밸리데이터를 초기화합니다. 키를 생성하거나 복구하세요."
    echo "1. 새 키 생성"
    echo "2. 기존 키 복구"
    echo "선택하세요 (1 또는 2):"
    read -r choice
    
    VAL_KEY="validator"
    
    if [ "$choice" = "1" ]; then
      # 새 키 생성
      ./build/zenad keys add "$VAL_KEY" --keyring-backend "$KEYRING" --algo "$KEYALGO" --home "$HOMEDIR"
    elif [ "$choice" = "2" ]; then
      # 기존 키 복구
      echo "니모닉 시드를 입력하세요:"
      read -r mnemonic
      echo "$mnemonic" | ./build/zenad keys add "$VAL_KEY" --recover --keyring-backend "$KEYRING" --algo "$KEYALGO" --home "$HOMEDIR"
    else
      echo "잘못된 선택입니다."
      exit 1
    fi
    
    echo "밸리데이터 설정이 완료되었습니다."
    echo "메인넷이 시작된 후 staking 트랜잭션을 통해 밸리데이터를 생성하세요."
  fi
  
  # 시드 노드 및 피어 설정
  echo "시드 노드와 피어 설정 (선택사항)"
  echo "config.toml 파일을 직접 수정하거나 필요한 경우 엔터 키를 눌러 계속 진행하세요..."
  read -r
fi

# WZENA 컨트랙트 배포 여부 확인
if grep -q "0x0000000000000000000000000000000000000000" zenad/token_pair.go; then
  echo "⚠️ 경고: WZENA 컨트랙트 주소가 아직 설정되지 않았습니다."
  echo "메인넷에서 컨트랙트 배포 후 zenad/token_pair.go 파일의 WZENAContractMainnet 값을 업데이트하거나"
  echo "거버넌스 제안을 통해 토큰 페어를 등록할 수 있습니다."
fi

# 노드 실행
echo "zena 메인넷 노드를 시작합니다..."
./build/zenad start "$TRACE" \
  --log_level="$LOGLEVEL" \
  --minimum-gas-prices="$MIN_GAS_PRICES" \
  --home="$HOMEDIR" \
  --json-rpc.api eth,txpool,personal,net,debug,web3 \
  --chain-id="$CHAINID" 