#!/bin/bash
# Script to install chaincode onto a peer node
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer.montauk.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/montauk.com/peers/peer.montauk.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=montauk-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/montauk.com/users/Admin@montauk.com/msp
cd /opt/gopath/src/github.com/chaincode/simple
if [ -f 'go/go.mod' ] && [ ! -d 'go/vendor' ]; then
  cd go
  export GO111MODULE=on
  go mod vendor
  export GO111MODULE=off
  cd -
fi
peer chaincode install --tls -l golang -v 1.0 -n simple \
  -p github.com/chaincode/simple/go
