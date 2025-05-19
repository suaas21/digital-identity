## Example rest curl command:

## port-forward the gateway rest deployment with 8080 port then do the followings:

### create identity
```curl
curl -X POST http://localhost:8080/create \
  -H "Content-Type: application/json" \
  -H "X-User-Cert: $(cat /root/digital-identity/network/temp/enrollments/org1/users/rcaadmin/msp/signcerts/cert.pem)" \
  -H "X-User-Key: $(cat /root/digital-identity/network/temp/enrollments/org1/users/rcaadmin/msp/signcerts/9a8cfa60817c0ea1b09be2d8cffe51c0bfc7075037878031d06a957d387a3926_sk)" \
  -H "X-User-MSPID: Org1MSP" \
  -d '{
    "id": "org1-102",
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
  -H "X-User-Cert: $(cat /root/digital-identity/network/temp/enrollments/org1/users/rcaadmin/msp/signcerts/cert.pem)" \
  -H "X-User-Key: $(cat /root/digital-identity/network/temp/enrollments/org1/users/rcaadmin/msp/signcerts/9a8cfa60817c0ea1b09be2d8cffe51c0bfc7075037878031d06a957d387a3926_sk)" \
  -H "X-User-MSPID: Org1MSP" \
  -d '{
    "id": "org1-102",
    "phone": "01777188559",
    "email": "sagor@example.com",
    "dob": "1990-01-01",
    "presentAddress": "123 Main St",
    "permanentAddress": "456 Oak St",
    "gender": "Male",
  }'  
``` 
### delete identity
```curl  
curl -X POST http://localhost:8080/delete \
  -H "Content-Type: application/json" \
  -H "X-User-Cert: $(cat /root/digital-identity/network/temp/enrollments/org1/users/rcaadmin/msp/signcerts/cert.pem)" \
  -H "X-User-Key: $(cat /root/digital-identity/network/temp/enrollments/org1/users/rcaadmin/msp/signcerts/9a8cfa60817c0ea1b09be2d8cffe51c0bfc7075037878031d06a957d387a3926_sk)" \
  -H "X-User-MSPID: Org1MSP" \
  -d '{"id": "org1-102"}'  
```
### get identity
```curl
curl localhost:8080/get/org1-102  

curl -X GET http://localhost:8080/get/org1-102 \
  -H "X-User-Cert: $(cat /root/digital-identity/network/temp/enrollments/org1/users/rcaadmin/msp/signcerts/cert.pem)" \
  -H "X-User-Key: $(cat /root/digital-identity/network/temp/enrollments/org1/users/rcaadmin/msp/signcerts/9a8cfa60817c0ea1b09be2d8cffe51c0bfc7075037878031d06a957d387a3926_sk)" \
  -H "X-User-MSPID: Org1MSP"
```