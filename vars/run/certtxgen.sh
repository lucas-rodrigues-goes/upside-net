#!/bin/bash
cd $FABRIC_CFG_PATH
# cryptogen generate --config crypto-config.yaml --output keyfiles
configtxgen -profile OrdererGenesis -outputBlock genesis.block -channelID systemchannel

configtxgen -printOrg hawkins-com > JoinRequest_hawkins-com.json
configtxgen -printOrg montauk-com > JoinRequest_montauk-com.json
configtxgen -printOrg oakridge-com > JoinRequest_oakridge-com.json
