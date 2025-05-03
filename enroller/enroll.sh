#!/bin/bash

function enrollUser() {
    export MSP_ID=Org1MSP
    export ORG=org1
    export USERNAME=org1user1
    export PASSWORD=org1user1
    export USER_ID=org1-102

    export WORKSHOP_CRYPTO=/root/digital-identity/network/temp
    export USER_MSP_DIR=$WORKSHOP_CRYPTO/enrollments/${ORG}/users/${USERNAME}/msp

    export WORKSHOP_INGRESS_DOMAIN=localho.st
    export WORKSHOP_NAMESPACE=test-network

    export ENROLLMENT_DIR=${WORKSHOP_CRYPTO}/enrollments
    export CHANNEL_MSP_DIR=${WORKSHOP_CRYPTO}/channel-msp

    log "Register and enroll a new user at the org CA"
    log "registering $USERNAME"
    /root/fabric-samples/bin/fabric-ca-client  register \
      --id.name       $USERNAME \
      --id.secret     $PASSWORD \
      --id.type       client \
      --id.affiliation $ORG \
      --id.attrs      "identity.id=$USER_ID:ecert" \
      --url           https://$WORKSHOP_NAMESPACE-$ORG-ca-ca.$WORKSHOP_INGRESS_DOMAIN \
      --tls.certfiles $WORKSHOP_CRYPTO/cas/$ORG-ca/tls-cert.pem \
      --mspdir        $WORKSHOP_CRYPTO/enrollments/$ORG/users/rcaadmin/msp

    /root/fabric-samples/bin/fabric-ca-client enroll \
      --url           https://$USERNAME:$PASSWORD@$WORKSHOP_NAMESPACE-$ORG-ca-ca.$WORKSHOP_INGRESS_DOMAIN \
      --tls.certfiles $WORKSHOP_CRYPTO/cas/$ORG-ca/tls-cert.pem \
      --mspdir        $WORKSHOP_CRYPTO/enrollments/$ORG/users/$USERNAME/msp

    mv $USER_MSP_DIR/keystore/*_sk $USER_MSP_DIR/keystore/key.pem
}

enrollUser