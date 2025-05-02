#!/bin/bash

function launch_rest() {
    export WORKSHOP_INGRESS_DOMAIN=localho.st
    export WORKSHOP_NAMESPACE=test-network

    local peer_pem=$CHANNEL_MSP_DIR/peerOrganizations/org1/msp/tlscacerts/tlsca-signcert.pem
    local ca_pem=$ENROLLMENT_DIR/org1/users/$USERNAME/msp/signcerts/cert.pem
    local keyPath=$ENROLLMENT_DIR/org1/users/$USERNAME/msp/keystore/key.pem

    #configure secrets
    kubectl -n $WORKSHOP_NAMESPACE delete secret my-secret || true
    kubectl create secret generic my-secret --from-file=keyPath="$keyPath"  --from-file=certPath="$ca_pem" --from-file=tlsCertPath="$peer_pem" -n $WORKSHOP_NAMESPACE
    #build docker image and push to local registry
    log "building restapi docker image"
    docker build -t localhost:5000/rest-api .
    log "pushing restapi docker image to local registry"
    docker push localhost:5000/rest-api
    #deploy rest api image to k8s
    log "deploying rest api to k8s"
    kubectl -n $WORKSHOP_NAMESPACE apply -f ./deployment.yaml
}

launch_rest
