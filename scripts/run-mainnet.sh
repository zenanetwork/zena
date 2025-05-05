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
  echo ""
  echo "============================================================"
  echo "메인넷 제네시스 파일을 다운로드해야 합니다."
  echo ""
  echo "다음 단계를 따르세요:"
  echo "1. 공식 웹사이트나 GitHub 저장소에서 메인넷 제네시스 파일을 다운로드하세요."
  echo "   (예: https://github.com/zena-network/mainnet/genesis.json)"
  echo ""
  echo "2. 다운로드한 파일을 다음 위치에 저장하세요:"
  echo "   $HOMEDIR/config/genesis.json"
  echo ""
  echo "3. 또는 다음 명령어로 직접 다운로드할 수 있습니다:"
  echo "   curl -s https://raw.githubusercontent.com/zena-network/mainnet/main/genesis.json > $HOMEDIR/config/genesis.json"
  echo ""
  echo "다운로드 후 엔터 키를 누르세요..."
  echo "============================================================"
  echo ""
  read -r
  
  # 제네시스 파일 존재 확인
  if [ ! -f "$GENESIS" ]; then
    echo "❌ 제네시스 파일이 존재하지 않습니다: $GENESIS"
    echo "제네시스 파일을 다운로드하고 다시 시도하세요."
    exit 1
  fi
  
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
  echo "시드 노드와 피어 설정 (필수)"
  echo "메인넷에 성공적으로 연결하려면 신뢰할 수 있는 시드 노드와 피어를 설정해야 합니다."
  echo "config.toml 파일에서 다음 섹션을 찾아 업데이트하세요:"
  echo ""
  echo "[p2p]"
  echo "seeds = \"시드노드ID@IP:26656,시드노드ID@IP:26656,...\""
  echo "persistent_peers = \"피어노드ID@IP:26656,피어노드ID@IP:26656,...\""
  echo ""
  echo "공식 문서나 커뮤니티에서 최신 시드 노드와 피어 정보를 확인하세요."
  echo "설정을 완료한 후 엔터 키를 눌러 계속 진행하세요..."
  read -r
fi

# WZENA 컨트랙트 배포 여부 확인
if grep -q "0x0000000000000000000000000000000000000000" zenad/token_pair.go; then
  echo "⚠️ WZENA 컨트랙트 주소가 아직 설정되지 않았습니다."
  echo "컨트랙트를 배포하시겠습니까? (y/n)"
  read -r deploy_choice
  
  if [[ "$deploy_choice" == "y" || "$deploy_choice" == "Y" ]]; then
    echo "WZENA 컨트랙트 배포를 진행합니다..."
    
    # 기존 계정 목록 표시 및 선택
    echo "사용 가능한 계정 목록:"
    KEYS_OUTPUT=$(./build/zenad keys list --keyring-backend "$KEYRING" --home "$HOMEDIR" --output json)
    KEY_COUNT=$(echo "$KEYS_OUTPUT" | jq '. | length')
    
    if [ "$KEY_COUNT" -eq 0 ]; then
      echo "사용 가능한 계정이 없습니다. 계정을 생성하겠습니까? (y/n)"
      read -r create_key
      
      if [[ "$create_key" == "y" || "$create_key" == "Y" ]]; then
        echo "계정 이름을 입력하세요:"
        read -r key_name
        
        # 새 계정 생성
        ./build/zenad keys add "$key_name" --keyring-backend "$KEYRING" --algo "$KEYALGO" --home "$HOMEDIR"
        DEPLOY_KEY="$key_name"
      else
        echo "계정이 없어 컨트랙트 배포를 중단합니다."
        exit 1
      fi
    else
      # 계정 목록 표시
      echo "$KEYS_OUTPUT" | jq -r '.[] | "\(.index). \(.name) (\(.address))"'
      
      # 사용자에게 계정 선택 요청
      echo "컨트랙트 배포에 사용할 계정 번호를 입력하세요 (1-$KEY_COUNT):"
      read -r key_index
      
      # 선택한 계정의 이름 가져오기
      DEPLOY_KEY=$(echo "$KEYS_OUTPUT" | jq -r ".[$((key_index-1))].name")
      
      if [ -z "$DEPLOY_KEY" ]; then
        echo "잘못된 계정 번호입니다. 컨트랙트 배포를 중단합니다."
        exit 1
      fi
      
      echo "선택한 계정: $DEPLOY_KEY"
    fi
    
    # 선택한 계정 주소 가져오기
    DEPLOYER_ADDRESS=$(./build/zenad keys show "$DEPLOY_KEY" -a --keyring-backend "$KEYRING" --home "$HOMEDIR")
    echo "컨트랙트 배포에 사용할 주소: $DEPLOYER_ADDRESS"
    
    # 계정 잔액 확인
    ACCOUNT_BALANCE=$(./build/zenad query bank balances "$DEPLOYER_ADDRESS" --output json 2>/dev/null || echo '{"balances":[]}')
    AZENA_BALANCE=$(echo "$ACCOUNT_BALANCE" | jq -r '.balances[] | select(.denom=="azena") | .amount // "0"')
    
    if [ -z "$AZENA_BALANCE" ] || [ "$AZENA_BALANCE" == "null" ]; then
      AZENA_BALANCE="0"
    fi
    
    if [ "$AZENA_BALANCE" == "0" ]; then
      echo "⚠️ 선택한 계정에 azena 토큰이 없습니다. 컨트랙트 배포에 필요한 토큰을 충전하세요."
      echo "계속 진행하려면 충전 후 엔터 키를 누르세요..."
      read -r
    else
      echo "계정 잔액: $AZENA_BALANCE azena"
    fi
    
    # WZENA9.json 파일이 존재하는지 확인
    WZENA_JSON="precompiles/werc20/testdata/WEVMOS9.json"
    if [ ! -f "$WZENA_JSON" ]; then
      echo "❌ WZENA 컨트랙트 JSON 파일을 찾을 수 없습니다: $WZENA_JSON"
      echo "컨트랙트 배포를 위해 필요한 파일입니다."
      echo "계속 진행하려면 엔터 키를 누르세요..."
      read -r
    else
      echo "WZENA 컨트랙트 JSON 파일을 찾았습니다."
      
      # 컨트랙트 데이터 추출
      CONTRACT_BYTECODE=$(jq -r '.bytecode' "$WZENA_JSON")
      
      # 임시 컨트랙트 배포 JSON 파일 생성
      DEPLOY_JSON="$HOMEDIR/deploy_wzena.json"
      cat > "$DEPLOY_JSON" << EOF
{
  "from": "$DEPLOYER_ADDRESS",
  "data": "$CONTRACT_BYTECODE",
  "gasLimit": "3000000",
  "value": "0"
}
EOF
      
      echo "컨트랙트 배포 트랜잭션을 전송합니다..."
      DEPLOY_RESULT=$(./build/zenad tx vm raw-call --data @"$DEPLOY_JSON" --from "$DEPLOY_KEY" --keyring-backend "$KEYRING" --chain-id "$CHAINID" --home "$HOMEDIR" --gas auto --gas-adjustment 1.5 -y --output json)
      
      if [ $? -eq 0 ]; then
        # 트랜잭션 해시 추출
        TX_HASH=$(echo "$DEPLOY_RESULT" | jq -r '.txhash')
        echo "컨트랙트 배포 트랜잭션 해시: $TX_HASH"
        
        # 트랜잭션 결과를 기다림
        echo "트랜잭션이 처리될 때까지 기다립니다..."
        sleep 5  # 트랜잭션이 처리될 때까지 대기
        
        # 트랜잭션 조회
        TX_RESULT=$(./build/zenad q tx "$TX_HASH" --output json)
        if [ $? -eq 0 ]; then
          TX_SUCCESS=$(echo "$TX_RESULT" | jq -r '.code == 0')
          
          if [ "$TX_SUCCESS" == "true" ]; then
            # 컨트랙트 주소 추출
            LOGS=$(echo "$TX_RESULT" | jq -r '.logs')
            CONTRACT_ADDRESS=$(echo "$LOGS" | jq -r '.[0].events[] | select(.type=="ethereum_tx") | .attributes[] | select(.key=="ethereumTxHash" or .key=="contractAddress") | select(.key=="contractAddress") | .value')
            
            if [ -n "$CONTRACT_ADDRESS" ]; then
              echo "✅ WZENA 컨트랙트 배포 성공! 주소: $CONTRACT_ADDRESS"
              
              # 토큰 페어 등록 안내
              echo ""
              echo "============================================================"
              echo "WZENA 컨트랙트가 성공적으로 배포되었습니다!"
              echo "컨트랙트 주소: $CONTRACT_ADDRESS"
              echo ""
              echo "이 주소를 zenad/token_pair.go 파일의 WZENAContractMainnet 상수에 업데이트하거나"
              echo "거버넌스 제안서를 통해 토큰 페어를 등록하세요."
              echo ""
              echo "토큰 정보:"
              echo "이름: Wrapped ZENA"
              echo "심볼: WZENA"
              echo "소수점: 18"
              echo "============================================================"
              echo ""
            else
              echo "❓ 컨트랙트 주소를 찾을 수 없습니다. 트랜잭션 해시로 조회해보세요: $TX_HASH"
            fi
          else
            echo "❌ 트랜잭션이 실패했습니다. 자세한 내용은 트랜잭션 결과를 확인하세요."
            echo "$TX_RESULT" | jq '.raw_log'
          fi
        else
          echo "❌ 트랜잭션 조회에 실패했습니다. 해시를 확인하세요: $TX_HASH"
        fi
      else
        echo "❌ 컨트랙트 배포 트랜잭션 전송에 실패했습니다."
        echo "오류: $DEPLOY_RESULT"
      fi
      
      # 임시 파일 정리
      rm -f "$DEPLOY_JSON"
    fi
  else
    echo "경고: WZENA 컨트랙트 주소가 아직 설정되지 않았습니다."
    echo "메인넷에서 컨트랙트 배포 후 zenad/token_pair.go 파일의 WZENAContractMainnet 값을 업데이트하거나"
    echo "거버넌스 제안을 통해 토큰 페어를 등록할 수 있습니다."
  fi
else
  echo "✅ WZENA 컨트랙트가 이미 설정되어 있습니다."
fi

# 노드 실행
echo "zena 메인넷 노드를 시작합니다..."
./build/zenad start "$TRACE" \
  --log_level="$LOGLEVEL" \
  --minimum-gas-prices="$MIN_GAS_PRICES" \
  --home="$HOMEDIR" \
  --json-rpc.api eth,txpool,personal,net,debug,web3 \
  --chain-id="$CHAINID" 