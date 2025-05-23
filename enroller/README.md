# Install Vault with Helm

## Add HashiCorp Helm repository
```shell
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update
```

## Create namespace
```shell
kubectl create namespace vault

```

## Install Vault in standalone mode
```shell
helm install vault hashicorp/vault \
    --namespace vault \
    --set "server.standalone.enabled=true" \
    --set "server.ha.enabled=false" \
    --set "server.dataStorage.enabled=true" \
    --set "server.dataStorage.size=1Gi"
```

## Initialize and Unseal Vault
```shell
# Wait for pod to be ready
kubectl -n vault wait pod/vault-0 --for=condition=Ready --timeout=120s

# Initialize Vault
kubectl -n vault exec vault-0 -- vault operator init \
  -key-shares=1 \
  -key-threshold=1 \
  -format=json > vault-keys.json

# Extract unseal key and root token
UNSEAL_KEY=$(jq -r ".unseal_keys_b64[]" vault-keys.json)
ROOT_TOKEN=$(jq -r ".root_token" vault-keys.json)

# Unseal Vault
kubectl -n vault exec vault-0 -- vault operator unseal $UNSEAL_KEY
```
## Configure Vault for Fabric
```shell
# Enable KV v2 secrets engine
kubectl -n vault exec vault-0 -- /bin/sh -c \
  "VAULT_TOKEN=$ROOT_TOKEN vault secrets enable -path=fabric/msp kv-v2"
```

## Vault Policy
```shell
kubectl -n vault exec vault-0 -- /bin/sh
export VAULT_TOKEN=<token>

vault policy write fabric-msp - <<EOF
path "fabric/msp/data/users/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

path "fabric/msp/metadata/users/*" {
  capabilities = ["list", "delete", "sudo"]
}
EOF
```

## In. cluster access
```shell
VAULT_ADDR="http://vault.vault.svc.cluster.local:8200"  # For in-cluster access
VAULT_TOKEN="<token>"
```

## Enrollment service all curl request:

1. Register User
```shell
curl -X POST http://enrollerapi.localho.st/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user124",
    "user_id": "org1-124"
  }'
```
2. Enroll User
```shell
curl -X POST http://enrollerapi.localho.st/enroll \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user124",
    "password": "f96a805591c0bee941abfe67670359e1"
  }'
```
3. Register & Enroll (Combined)
```shell
curl -X POST http://enrollerapi.localho.st/register-enroll \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user124",
    "user_id": "org1-124"
  }'
```
4. Get User MSP
```shell
curl -X GET http://enrollerapi.localho.st/msp/user124
```

5. Revoke User
```shell
curl -X DELETE http://enrollerapi.localho.st/revoke/user124
```
