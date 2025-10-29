#!/bin/bash
# Script to discover endorsers and channel config
cd /vars

export PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/hawkins.com/users/Admin@hawkins.com/tls/ca.crt
export ADMINPRIVATEKEY=/vars/keyfiles/peerOrganizations/hawkins.com/users/Admin@hawkins.com/msp/keystore/2d6386957ed9569a191b1939f796a4f544c3e34ee09a4d87ebcd7e9c5ead0f81_sk
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
