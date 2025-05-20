#!/bin/bash

function launch_rest() {
    export ORG=org1
    export WORKSHOP_INGRESS_DOMAIN=localho.st
    export WORKSHOP_NAMESPACE=test-network
    export WORKSHOP_CRYPTO=/root/digital-identity/network/temp
    export CHANNEL_MSP_DIR=${WORKSHOP_CRYPTO}/channel-msp

    local tls_pem=$CHANNEL_MSP_DIR/peerOrganizations/$ORG/msp/tlscacerts/tlsca-signcert.pem

    echo "creating rest secret gateway-tls-cert"
    kubectl create secret generic gateway-tls-cert --from-file=tlsCertPath="$tls_pem" -n $WORKSHOP_NAMESPACE
    #build docker image and push to local registry
    echo "building restapi docker image"
    docker build -t localhost:5000/rest-api .
    echo "pushing restapi docker image to local registry"
    docker push localhost:5000/rest-api
    #deploy rest api image to k8s
    echo "deploying rest api to k8s"
    kubectl -n $WORKSHOP_NAMESPACE apply -f ./deployment.yaml
}

function delete_rest() {
  export WORKSHOP_NAMESPACE=test-network
  echo "deleting secret gateway-tls-cert"
  kubectl delete secret gateway-tls-cert -n $WORKSHOP_NAMESPACE
  echo "deleting rest deploy"
  kubectl -n $WORKSHOP_NAMESPACE delete -f ./deployment.yaml
}

delete_rest
launch_rest
