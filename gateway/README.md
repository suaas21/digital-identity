## Example rest curl command:

## port-forward the gateway rest deployment with 8080 port then do the followings:

### create identity
```curl
curl -X POST http://localhost:8080/create \
  -H "Content-Type: application/json" \
  -H "X-User-Cert: $(base64 -w 0 /root/digital-identity/network/temp/enrollments/org1/users/org1user111/msp/signcerts/cert.pem)" \
  -H "X-User-Key: $(base64 -w 0 /root/digital-identity/network/temp/enrollments/org1/users/org1user111/msp/keystore/key.pem)" \
  -H "X-User-MSPID: Org1MSP" \
  -d '{
    "id": "org1-111",
    "firstName": "sagor",
    "lastName": "azad",
    "phone": "01777188559",
    "nationalID": "101010101"
  }' 
```  
### update identity
```curl
curl -X POST http://localhost:8080/update \
  -H "Content-Type: application/json" \
  -H "X-User-Cert: $(base64 -w 0 /root/digital-identity/network/temp/enrollments/org1/users/org1user111/msp/signcerts/cert.pem)" \
  -H "X-User-Key: $(base64 -w 0 /root/digital-identity/network/temp/enrollments/org1/users/org1user111/msp/keystore/key.pem)" \
  -H "X-User-MSPID: Org1MSP" \
  -d '{
    "id": "org1-111",
    "phone": "01777188559",
    "email": "sagor@example.com",
    "dob": "1990-01-01",
    "presentAddress": "123 Main St",
    "permanentAddress": "456 Oak St",
    "gender": "Male"
  }'  
``` 
### delete identity
```curl  
curl -X POST http://localhost:8080/delete \
  -H "Content-Type: application/json" \
  -H "X-User-Cert: $(base64 -w 0 /root/digital-identity/network/temp/enrollments/org1/users/org1user111/msp/signcerts/cert.pem)" \
  -H "X-User-Key: $(base64 -w 0 /root/digital-identity/network/temp/enrollments/org1/users/org1user111/msp/keystore/key.pem)" \
  -H "X-User-MSPID: Org1MSP" \
  -d '{"id": "org1-111"}'  
```
### get identity
```curl
curl -X GET http://localhost:8080/get/org1-111 \
  -H "X-User-Cert: $(base64 -w 0 /root/digital-identity/network/temp/enrollments/org1/users/org1user111/msp/signcerts/cert.pem)" \
  -H "X-User-Key: $(base64 -w 0 /root/digital-identity/network/temp/enrollments/org1/users/org1user111/msp/keystore/key.pem)" \
  -H "X-User-MSPID: Org1MSP"
```