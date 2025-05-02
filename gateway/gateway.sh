#!/bin/bash

function launch_rest() {
    export ORG=org1
    export USERNAME=org1user

    export WORKSHOP_INGRESS_DOMAIN=localho.st
    export WORKSHOP_NAMESPACE=test-network

    export WORKSHOP_CRYPTO=/root/fabric-samples/full-stack-asset-transfer-guide/infrastructure/sample-network/temp
    export ENROLLMENT_DIR=${WORKSHOP_CRYPTO}/enrollments
    export CHANNEL_MSP_DIR=${WORKSHOP_CRYPTO}/channel-msp

    local peer_pem=$CHANNEL_MSP_DIR/peerOrganizations/$ORG/msp/tlscacerts/tlsca-signcert.pem
    local ca_pem=$ENROLLMENT_DIR/$ORG/users/$USERNAME/msp/signcerts/cert.pem
    local keyPath=$ENROLLMENT_DIR/$ORG/users/$USERNAME/msp/keystore/key.pem

    #configure secrets
    kubectl -n $WORKSHOP_NAMESPACE delete secret my-secret || true
    kubectl create secret generic my-secret --from-file=keyPath="$keyPath"  --from-file=certPath="$ca_pem" --from-file=tlsCertPath="$peer_pem" -n $WORKSHOP_NAMESPACE
    #build docker image and push to local registry
    echo "building restapi docker image"
    docker build -t localhost:5000/rest-api .
    echo "pushing restapi docker image to local registry"
    docker push localhost:5000/rest-api
    #deploy rest api image to k8s
    echo "deploying rest api to k8s"
    kubectl -n $WORKSHOP_NAMESPACE apply -f ./deployment.yaml
}

launch_rest
