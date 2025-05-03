## Example rest curl command:

## port-forward the gateway rest deployment with 8080 port then do the followings:

### create identity
```curl
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"id":"org1-102","phone":"01777188559","nationalID":"101010101","firstName":"sagor"}' \
  localhost:8080/create
```  
### update identity
```curl
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"id":"org1-102","phone":"01777188559","nationalID":"101010101","firstName":"sagor","lastName":"azad","gender":"male"}' \
  localhost:8080/update 
``` 
### delete identity
```curl
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"id":"org1-102"}' \
  localhost:8080/delete   
```
### get identity
```curl
curl localhost:8080/get/org1-102  
```