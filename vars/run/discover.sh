#!/bin/bash
# Script to discover endorsers and channel config
cd /vars

export PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/hawkins.com/users/Admin@hawkins.com/tls/ca.crt
export ADMINPRIVATEKEY=/vars/keyfiles/peerOrganizations/hawkins.com/users/Admin@hawkins.com/msp/keystore/92ced0134637e4597d7adef3a8fcf28af3bf98772814ea9f7ee7d631f14dede3_sk
export ADMINCERT=/vars/keyfiles/peerOrganizations/hawkins.com/users/Admin@hawkins.com/msp/signcerts/Admin@hawkins.com-cert.pem

discover endorsers --peerTLSCA $PEER_TLS_ROOTCERT_FILE \
  --userKey $ADMINPRIVATEKEY \
  --userCert $ADMINCERT \
  --MSP hawkins-com --channel mychannel \
  --server peer.hawkins.com:7051 \
  --chaincode simple | jq '.[0]' | \
  jq 'del(.. | .Identity?)' | jq 'del(.. | .LedgerHeight?)' \
  > /vars/discover/mychannel_simple_endorsers.json

discover config --peerTLSCA $PEER_TLS_ROOTCERT_FILE \
  --userKey $ADMINPRIVATEKEY \
  --userCert $ADMINCERT \
  --MSP hawkins-com --channel mychannel \
  --server peer.hawkins.com:7051 > /vars/discover/mychannel_config.json
