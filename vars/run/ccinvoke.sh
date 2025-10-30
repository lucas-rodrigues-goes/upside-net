#!/bin/bash
# Script to invoke chaincode
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=peer.hawkins.com:7051
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/hawkins.com/peers/peer.hawkins.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=hawkins-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/hawkins.com/users/Admin@hawkins.com/msp
export ORDERER_ADDRESS=orderer1.example.com:7050
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/example.com/orderers/orderer1.example.com/tls/ca.crt
peer chaincode invoke -o $ORDERER_ADDRESS --cafile $ORDERER_TLS_CA \
  --tls -C mychannel -n blendchaincode  \
  --peerAddresses peer.montauk.com:7051 \
  --tlsRootCertFiles /vars/discover/mychannel/montauk-com/tlscert \
  -c '{"Args":["initDimensionalEnergy","1","lab","55","2","Lucas"]}'
