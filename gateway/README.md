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

### Extra curl with base64 encoded cert and key
```shell
curl -X POST http://localhost:8080/create \
  -H "Content-Type: application/json" \
  -H "X-User-Cert: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN2VENDQW1TZ0F3SUJBZ0lVY0gwbkdlMysySkVRMlFjZ2J2QlRBbC80eWFNd0NnWUlLb1pJemowRUF3SXcKWWpFTE1Ba0dBMVVFQmhNQ1ZWTXhGekFWQmdOVkJBZ1REazV2Y25Sb0lFTmhjbTlzYVc1aE1SUXdFZ1lEVlFRSwpFd3RJZVhCbGNteGxaR2RsY2pFUE1BMEdBMVVFQ3hNR1JtRmljbWxqTVJNd0VRWURWUVFERXdwdmNtY3hMV05oCkxXTmhNQjRYRFRJMU1EVXlNREF4TkRJd01Gb1hEVE0xTURVeE9EQXpNVEl3TUZvd2JERUxNQWtHQTFVRUJoTUMKVlZNeEZ6QVZCZ05WQkFnVERrNXZjblJvSUVOaGNtOXNhVzVoTVJRd0VnWURWUVFLRXd0SWVYQmxjbXhsWkdkbApjakVjTUFzR0ExVUVDeE1FYjNKbk1UQU5CZ05WQkFzVEJtTnNhV1Z1ZERFUU1BNEdBMVVFQXhNSGRYTmxjakV4Ck5qQlpNQk1HQnlxR1NNNDlBZ0VHQ0NxR1NNNDlBd0VIQTBJQUJPZVJxWksydTJNcTNva3JWcUFwdjQwNFFVWm4KVUN5T0NqbmJ4UENqLzBack5aSllPOW1Zd0FWckZRcFZ3Tm5CWkFqbk1nanlCS3JuOHF5WTNZdDkxSldqZ2UwdwpnZW93RGdZRFZSMFBBUUgvQkFRREFnZUFNQXdHQTFVZEV3RUIvd1FDTUFBd0hRWURWUjBPQkJZRUZHNHd6ZlJXCmtMVGlvdm16L3cxYlVtcWNOWXFmTUI4R0ExVWRJd1FZTUJhQUZMVkRRdmh5NkdrYmlXbHk4YzB5YVJqa1padUsKTUJFR0ExVWRFUVFLTUFpQ0JtUmxZbWxoYmpCM0JnZ3FBd1FGQmdjSUFRUnJleUpoZEhSeWN5STZleUpvWmk1QgpabVpwYkdsaGRHbHZiaUk2SW05eVp6RWlMQ0pvWmk1RmJuSnZiR3h0Wlc1MFNVUWlPaUoxYzJWeU1URTJJaXdpCmFHWXVWSGx3WlNJNkltTnNhV1Z1ZENJc0ltbGtaVzUwYVhSNUxtbGtJam9pYjNKbk1TMHhNVFlpZlgwd0NnWUkKS29aSXpqMEVBd0lEUndBd1JBSWdMN3lINTlNMEJ4VldybFh6UGpoSStYMk9FYkF1YlU5dXM3ZHJLWFpNZlpRQwpJQVpaRi9pbFYyY2cweVpQSnk3dDFUa0t1SnN2SjlETzQyMldrMTNXblhKWQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==" \
  -H "X-User-Key: LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JR0hBZ0VBTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEJHMHdhd0lCQVFRZzB5dUJWWE4xMnZuMWp2MWwKLzNVVUtubm5SWHZLT0F5RmRiNHhJd1V5UDJTaFJBTkNBQVRua2FtU3RydGpLdDZKSzFhZ0tiK05PRUZHWjFBcwpqZ281MjhUd28vOUdheldTV0R2Wm1NQUZheFVLVmNEWndXUUk1eklJOGdTcTUvS3NtTjJMZmRTVgotLS0tLUVORCBQUklWQVRFIEtFWS0tLS0tCg==" \
  -H "X-User-MSPID: Org1MSP" \
  -d '{
    "id": "org1-116",
    "firstName": "sagor",
    "lastName": "azad",
    "phone": "01777188559",
    "nationalID": "101010101"
  }'
 

curl -X GET http://localhost:8080/get/org1-116 \
  -H "X-User-Cert: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN2VENDQW1TZ0F3SUJBZ0lVY0gwbkdlMysySkVRMlFjZ2J2QlRBbC80eWFNd0NnWUlLb1pJemowRUF3SXcKWWpFTE1Ba0dBMVVFQmhNQ1ZWTXhGekFWQmdOVkJBZ1REazV2Y25Sb0lFTmhjbTlzYVc1aE1SUXdFZ1lEVlFRSwpFd3RJZVhCbGNteGxaR2RsY2pFUE1BMEdBMVVFQ3hNR1JtRmljbWxqTVJNd0VRWURWUVFERXdwdmNtY3hMV05oCkxXTmhNQjRYRFRJMU1EVXlNREF4TkRJd01Gb1hEVE0xTURVeE9EQXpNVEl3TUZvd2JERUxNQWtHQTFVRUJoTUMKVlZNeEZ6QVZCZ05WQkFnVERrNXZjblJvSUVOaGNtOXNhVzVoTVJRd0VnWURWUVFLRXd0SWVYQmxjbXhsWkdkbApjakVjTUFzR0ExVUVDeE1FYjNKbk1UQU5CZ05WQkFzVEJtTnNhV1Z1ZERFUU1BNEdBMVVFQXhNSGRYTmxjakV4Ck5qQlpNQk1HQnlxR1NNNDlBZ0VHQ0NxR1NNNDlBd0VIQTBJQUJPZVJxWksydTJNcTNva3JWcUFwdjQwNFFVWm4KVUN5T0NqbmJ4UENqLzBack5aSllPOW1Zd0FWckZRcFZ3Tm5CWkFqbk1nanlCS3JuOHF5WTNZdDkxSldqZ2UwdwpnZW93RGdZRFZSMFBBUUgvQkFRREFnZUFNQXdHQTFVZEV3RUIvd1FDTUFBd0hRWURWUjBPQkJZRUZHNHd6ZlJXCmtMVGlvdm16L3cxYlVtcWNOWXFmTUI4R0ExVWRJd1FZTUJhQUZMVkRRdmh5NkdrYmlXbHk4YzB5YVJqa1padUsKTUJFR0ExVWRFUVFLTUFpQ0JtUmxZbWxoYmpCM0JnZ3FBd1FGQmdjSUFRUnJleUpoZEhSeWN5STZleUpvWmk1QgpabVpwYkdsaGRHbHZiaUk2SW05eVp6RWlMQ0pvWmk1RmJuSnZiR3h0Wlc1MFNVUWlPaUoxYzJWeU1URTJJaXdpCmFHWXVWSGx3WlNJNkltTnNhV1Z1ZENJc0ltbGtaVzUwYVhSNUxtbGtJam9pYjNKbk1TMHhNVFlpZlgwd0NnWUkKS29aSXpqMEVBd0lEUndBd1JBSWdMN3lINTlNMEJ4VldybFh6UGpoSStYMk9FYkF1YlU5dXM3ZHJLWFpNZlpRQwpJQVpaRi9pbFYyY2cweVpQSnk3dDFUa0t1SnN2SjlETzQyMldrMTNXblhKWQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==" \
  -H "X-User-Key: LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JR0hBZ0VBTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEJHMHdhd0lCQVFRZzB5dUJWWE4xMnZuMWp2MWwKLzNVVUtubm5SWHZLT0F5RmRiNHhJd1V5UDJTaFJBTkNBQVRua2FtU3RydGpLdDZKSzFhZ0tiK05PRUZHWjFBcwpqZ281MjhUd28vOUdheldTV0R2Wm1NQUZheFVLVmNEWndXUUk1eklJOGdTcTUvS3NtTjJMZmRTVgotLS0tLUVORCBQUklWQVRFIEtFWS0tLS0tCg==" \
  -H "X-User-MSPID: Org1MSP"
```