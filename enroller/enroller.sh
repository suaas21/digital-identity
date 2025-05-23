#!/bin/bash

function launch_enroller() {
  export WORKSHOP_CRYPTO=/root/digital-identity/network/temp
  export ORG=org1
  export VAULT_TOKEN=token
  export WORKSHOP_NAMESPACE=test-network

  # Create separate secrets for each MSP component
  kubectl create secret -n $WORKSHOP_NAMESPACE generic rcaadmin-cacerts \
    --from-file=cacert.pem=$WORKSHOP_CRYPTO/enrollments/$ORG/users/rcaadmin/msp/cacerts/test-network-org1-ca-ca-localho-st-443.pem

  kubectl create secret -n $WORKSHOP_NAMESPACE generic rcaadmin-keystore \
    --from-file=key.pem=$WORKSHOP_CRYPTO/enrollments/$ORG/users/rcaadmin/msp/keystore/6c28bc5f13c33f014cc69d7f2cd45f2c4471390d686e9021d5427a2f9725da23_sk

  kubectl create secret -n $WORKSHOP_NAMESPACE generic rcaadmin-signcerts \
    --from-file=cert.pem=$WORKSHOP_CRYPTO/enrollments/$ORG/users/rcaadmin/msp/signcerts/cert.pem

  # Create secret for TLS certificates
  kubectl create secret generic tls-certs \
    --from-file=tls-cert.pem=$WORKSHOP_CRYPTO/cas/$ORG-ca/tls-cert.pem -n $WORKSHOP_NAMESPACE

  # Create secret for Vault token
  kubectl create secret -n $WORKSHOP_NAMESPACE generic vault-secret \
    --from-literal=token=$VAULT_TOKEN

  #build docker image and push to local registry
  echo "building restapi docker image"
  docker build -t localhost:5000/enroller-api .
  echo "pushing restapi docker image to local registry"
  docker push localhost:5000/enroller-api
  #deploy rest api image to k8s
  echo "deploying rest api to k8s"
  kubectl -n $WORKSHOP_NAMESPACE apply -f ./deployment.yaml

}

function delete_enroller() {
  export WORKSHOP_NAMESPACE=test-network
  echo "deleting all secret"
  kubectl delete secret  rcaadmin-signcerts -n $WORKSHOP_NAMESPACE
  kubectl delete secret  rcaadmin-keystore -n $WORKSHOP_NAMESPACE
  kubectl delete secret  rcaadmin-cacerts -n $WORKSHOP_NAMESPACE
  kubectl delete secret tls-certs -n $WORKSHOP_NAMESPACE
  kubectl delete secret vault-secret -n $WORKSHOP_NAMESPACE
  echo "deleting rest deploy"
  kubectl -n $WORKSHOP_NAMESPACE delete -f ./deployment.yaml
}

delete_enroller
launch_enroller