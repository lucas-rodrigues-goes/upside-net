#!/bin/bash
# Script to instantiate chaincode
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer.hawkins.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/hawkins.com/peers/peer.hawkins.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=hawkins-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/hawkins.com/users/Admin@hawkins.com/msp
export ORDERER_ADDRESS=orderer1.example.com:7050
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/example.com/orderers/orderer1.example.com/tls/ca.crt
ccdone=$(peer chaincode list -C mychannel --instantiated|grep "Name: simple,")
if [[ -z "$ccdone" ]]; then
  peer chaincode instantiate -o $ORDERER_ADDRESS --cafile $ORDERER_TLS_CA --tls \
  -C mychannel -n simple -v 1.0 \
  -c '{"Args":["init","a","200","b","300"]}'
else
  peer chaincode upgrade -o $ORDERER_ADDRESS --cafile $ORDERER_TLS_CA --tls \
  -C mychannel -n simple -v 1.0 \
  -c '{"Args":["init","a","200","b","300"]}'
fi
