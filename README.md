# Hyperledger fabric network 
- Deploy kind cluster, local image registry, and nginx ingress
- Deploy Hyperledger fabric crds, operator
- Deploy simple network 
- Create 2 organization's channel
- Build chain as a service for user digital identity
- Enroll user msp using enroller service
- Deploy gateway service for rest API



## Required Tools for Kubernetes Deployment

- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [jq](https://stedolan.github.io/jq/)
- [just](https://github.com/casey/just#installation) to run all the comamnds here directly
- [kind](https://kind.sigs.k8s.io/) if you want to create a cluster locally, see below for other options
- [k9s](https://k9scli.io) (recommended, but not essential)


## Install Fabric peer CLI and set environment variables

```shell
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
./install-fabric.sh binary
export WORKSHOP_PATH=$(pwd)
export PATH=${WORKSHOP_PATH}/bin:$PATH
export FABRIC_CFG_PATH=${WORKSHOP_PATH}/config
```

## Set Ingress and namespace environment variables

```shell
export WORKSHOP_INGRESS_DOMAIN=localho.st
export WORKSHOP_NAMESPACE=test-network
```

## just magic

```shell
just
```