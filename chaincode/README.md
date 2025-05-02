# Chaincode

## Fabric k8s Builder


## Ready?

- Workshop environment variables:
```shell

export WORKSHOP_PATH=/root/fabric-samples/full-stack-asset-transfer-guide
export FABRIC_CFG_PATH=${WORKSHOP_PATH}/config  
export PATH=${WORKSHOP_PATH}/bin:$PATH

```

- Set the location for the network's TLS certificates, channel MSP, and user enrollments:
```shell

export WORKSHOP_CRYPTO=$WORKSHOP_PATH/infrastructure/sample-network/temp
```

## Kubernetes IN Docker (KIND)

- Set the cluster ingress domain and target k8s namespace.  The `localho.st` domain is a public DNS wildcard resolver
  mapping `*.localho.st` to 127.0.0.1.
```shell

export WORKSHOP_INGRESS_DOMAIN=localho.st
export WORKSHOP_NAMESPACE=test-network

```


## Set the peer client environment

```shell

export ORG1_PEER1_ADDRESS=${WORKSHOP_NAMESPACE}-org1-peer1-peer.${WORKSHOP_INGRESS_DOMAIN}:443
export ORG1_PEER2_ADDRESS=${WORKSHOP_NAMESPACE}-org1-peer2-peer.${WORKSHOP_INGRESS_DOMAIN}:443

# org1-peer1: 
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_ADDRESS=${ORG1_PEER1_ADDRESS}
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_MSPCONFIGPATH=${WORKSHOP_CRYPTO}/enrollments/org1/users/org1admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${WORKSHOP_CRYPTO}/channel-msp/peerOrganizations/org1/msp/tlscacerts/tlsca-signcert.pem
export CORE_PEER_CLIENT_CONNTIMEOUT=15s
export CORE_PEER_DELIVERYCLIENT_CONNTIMEOUT=15s
export ORDERER_ENDPOINT=${WORKSHOP_NAMESPACE}-org0-orderersnode1-orderer.${WORKSHOP_INGRESS_DOMAIN}:443
export ORDERER_TLS_CERT=${WORKSHOP_CRYPTO}/channel-msp/ordererOrganizations/org0/orderers/org0-orderersnode1/tls/signcerts/tls-cert.pem

```

## Docker Engine Configuration

**NOTE: SKIP THIS STEP IF USING `localho.st` AS THE INGRESS DOMAIN**

Configure the docker engine with the insecure container registry `${WORKSHOP_INGRESS_DOMAIN}:5000`

For example:  (Docker -> Preferences -> Docker Engine)
```json
{
  "insecure-registries": [
    "192-168-205-6.nip.io:5000"
  ]
}
```

- apply and restart

## Chaincode Revision

```shell

CHANNEL_NAME=mychannel
VERSION=v0.0.1
SEQUENCE=1

```

## Build the Chaincode Docker Image

```shell

CHAINCODE_NAME=identity
CHAINCODE_PACKAGE=${CHAINCODE_NAME}.tgz
CONTAINER_REGISTRY=$WORKSHOP_INGRESS_DOMAIN:5000
CHAINCODE_IMAGE=$CONTAINER_REGISTRY/$CHAINCODE_NAME

# Build the chaincode image
docker build -t $CHAINCODE_IMAGE .

# Push the image to the insecure container registry
docker push $CHAINCODE_IMAGE

```


## Prepare a k8s Chaincode Package

```shell

IMAGE_DIGEST=$(docker inspect --format='{{index .RepoDigests 0}}' $CHAINCODE_IMAGE | cut -d'@' -f2)

./pkgcc.sh -l $CHAINCODE_NAME -n localhost:5000/$CHAINCODE_NAME -d $IMAGE_DIGEST

```

## Install the Chaincode

```shell

# Install the chaincode package on both peers in the org 
CORE_PEER_ADDRESS=${ORG1_PEER1_ADDRESS} peer lifecycle chaincode install $CHAINCODE_PACKAGE
CORE_PEER_ADDRESS=${ORG1_PEER2_ADDRESS} peer lifecycle chaincode install $CHAINCODE_PACKAGE

export PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid $CHAINCODE_PACKAGE) && echo $PACKAGE_ID

# Approve the contract for org1 
peer lifecycle \
	chaincode       approveformyorg \
	--channelID     ${CHANNEL_NAME} \
	--name          ${CHAINCODE_NAME} \
	--version       ${VERSION} \
	--package-id    ${PACKAGE_ID} \
	--sequence      ${SEQUENCE} \
	--orderer       ${ORDERER_ENDPOINT} \
	--tls --cafile  ${ORDERER_TLS_CERT} \
	--connTimeout   15s

# Commit the contract on the channel
peer lifecycle \
	chaincode       commit \
	--channelID     ${CHANNEL_NAME} \
	--name          ${CHAINCODE_NAME} \
	--version       ${VERSION} \
	--sequence      ${SEQUENCE} \
	--orderer       ${ORDERER_ENDPOINT} \
	--tls --cafile  ${ORDERER_TLS_CERT} \
	--connTimeout   15s

```

```shell

peer chaincode query -n $CHAINCODE_NAME -C mychannel -c '{"Args":["org.hyperledger.fabric:GetMetadata"]}' | jq

```