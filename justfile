#
# Copyright contributors to the Hyperledgendary Full Stack Asset Transfer project
#
# SPDX-License-Identifier: Apache-2.0
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at:
#
# 	  http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Main justfile to run all the development scripts
# To install 'just' see https://github.com/casey/just#installation


###############################################################################
# COMMON TARGETS                                                              #
###############################################################################


# Ensure all properties are exported as shell env-vars
set export

# set the current directory, and the location of the test dats
CWDIR := justfile_directory()

_default:
  @just -f {{justfile()}} --list

# Run the check script to validate tool versions installed
check:
  ${CWDIR}/check.sh

cluster_name    := env_var_or_default("WORKSHOP_CLUSTER_NAME",       "kind")
cluster_runtime := env_var_or_default("WORKSHOP_CLUSTER_RUNTIME",    "kind")
ingress_domain  := env_var_or_default("WORKSHOP_INGRESS_DOMAIN",     "localho.st")
storage_class   := env_var_or_default("WORKSHOP_STORAGE_CLASS",      "standard")
chaincode_name  := env_var_or_default("WORKSHOP_CHAINCODE_NAME",     "identity")
internal_repo_endpoint  := env_var_or_default("WORKSHOP_INTERNAL_REPO",     "localhost:5000")
external_repo_endpoint  := env_var_or_default("WORKSHOP_EXTERNAL_REPO",     "localhost:5000")
cluster_type    := env_var_or_default("WORKSHOP_CLUSTER_TYPE",       "k8s")


# Start a local KIND cluster with nginx, localhost:5000 registry, and *.localho.st alias in kube DNS
kind: unkind
    #!/bin/bash
    set -e -o pipefail

    k8s-cluster/kind_with_nginx.sh {{cluster_name}}
    ls -lart ~/.kube/config
    chmod o+r ~/.kube/config

    # check connectivity to local k8s
    kubectl cluster-info &>/dev/null

# Shut down the KIND cluster
unkind:
    #!/bin/bash
    kind delete cluster --name {{cluster_name}}

    if docker inspect kind-registry &>/dev/null; then
        echo "Stopping container registry"
        docker kill kind-registry
        docker rm kind-registry
    fi

# Bring up the nginx ingress controller on the target k8s cluster
nginx:
    #!/bin/bash
    kubectl apply -k https://github.com/hyperledger-labs/fabric-operator.git/config/ingress/{{ cluster_runtime }}

    sleep 20

    kubectl wait --namespace ingress-nginx \
      --for=condition=ready pod \
      --selector=app.kubernetes.io/component=controller \
      --timeout=3m

# Just start the operator
operator: operator-crds
    network/network operator

# Just start the console
console: operator
    network/network console

# Just install the operator CRDs
operator-crds: check-kube
    kubectl apply -k https://github.com/hyperledger-labs/fabric-operator.git/config/crd

###############################################################################
# CLOUD NATIVE TARGETS                                                        #
###############################################################################

# Deploy the operator sample network and create a channel
cloud-network: cloud-network-down check-kube
    network/network up

# Tear down the operator sample network
cloud-network-down:
    network/network down

# Create channel 'mychannel'
cloud-channel:
    network/network channel create

# Check that the cloud setup has been performed
check-setup: check

# Check that the k8s API controller is ready
check-kube: check-setup
    checks/check-kube.sh

# Check that the sample network and channel have been deployed
check-network: check-kube
    checks/check-network.sh

# Check that the smart contract has been deployed
check-chaincode: check-network
    checks/check-chaincode.sh

