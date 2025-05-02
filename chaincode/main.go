package main

import (
	identity "github.com/hyperledger/digital-identity/chaincode/smart-contact"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"log"
)

func main() {
	smartContract, err := contractapi.NewChaincode(&identity.SmartContract{})
	if err != nil {
		log.Panicf("Error creating identity chaincode: %v", err)
	}

	err = smartContract.Start()
	if err != nil {
		log.Panicf("Error strting identity chaincode: %v", err)
	}
}
