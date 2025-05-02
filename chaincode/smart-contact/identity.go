package identity

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Identity struct {
	Id               string `json:"id"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	Dob              string `json:"dob"`
	PresentAddress   string `json:"presentAddress"`
	PermanentAddress string `json:"permanentAddress"`
	Gender           string `json:"gender"`
	NationalID       string `json:"nationalID"`
	Owner            string `json:"owner"`
}

// CreateIdentity issues a new identity to the world state with given details.
func (s *SmartContract) CreateIdentity(ctx contractapi.TransactionContextInterface, identity Identity) error {
	if isEmptyField(identity.Id) {
		return errors.New("identity id is not provided")
	}

	if isEmptyField(identity.FirstName) {
		return errors.New("identity firstName is not provided")
	}

	if isEmptyField(identity.Phone) {
		return errors.New("identity phone is not provided")
	}

	if isEmptyField(identity.NationalID) {
		return errors.New("identity national id is not provided")
	}

	err := ctx.GetClientIdentity().AssertAttributeValue("identity.id", identity.Id)
	if err != nil {
		return errors.New("submitting identity is not authorized to create, does not have identity.id or not valid identity")
	}

	//idnty, err := s.ReadIdentity(ctx, identity.Id)
	//if err != nil {
	//	return err
	//}
	//if idnty.Id == identity.Id {
	//	return fmt.Errorf("the identity %v already exists", identity.Id)
	//}

	// Get ID of submitting client identity
	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	// set clientID to Owner
	identity.Owner = clientID

	jsonByte, err := json.Marshal(identity)
	if err != nil {
		return err
	}

	fmt.Printf("put identity data to world state")
	return ctx.GetStub().PutState(identity.Id, jsonByte)
}

// ReadIdentity returns the identity stored in the world state with given id.
func (s *SmartContract) ReadIdentity(ctx contractapi.TransactionContextInterface, id string) (*Identity, error) {
	if isEmptyField(id) {
		return nil, errors.New("identity id is not provided")
	}

	fmt.Printf("get identity data from world state")
	jsonByte, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if jsonByte == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var idnty Identity
	err = json.Unmarshal(jsonByte, &idnty)
	if err != nil {
		return nil, err
	}

	return &idnty, nil
}

// DeleteIdentity deletes a given asset from the world state.
func (s *SmartContract) DeleteIdentity(ctx contractapi.TransactionContextInterface, id string) error {
	if isEmptyField(id) {
		return errors.New("identity id is not provided")
	}

	idnty, err := s.ReadIdentity(ctx, id)
	if err != nil {
		return err
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	if clientID != idnty.Owner {
		return fmt.Errorf("submitting client not authorized to delete identity, does not own identity")
	}

	fmt.Printf("delete identity data fron world state")
	return ctx.GetStub().DelState(id)
}

// UpdateIdentity updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateIdentity(ctx contractapi.TransactionContextInterface, id string, update Identity) error {
	if isEmptyField(id) {
		return errors.New("identity id is not provided")
	}

	idnty, err := s.ReadIdentity(ctx, id)
	if err != nil {
		return err
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	if clientID != idnty.Owner {
		return fmt.Errorf("submitting client not authorized to update identity, does not own identity")
	}

	if !isEmptyField(update.FirstName) {
		idnty.LastName = update.LastName
	}
	if !isEmptyField(update.FirstName) {
		idnty.FirstName = update.FirstName
	}
	if !isEmptyField(update.FirstName) {
		idnty.Email = update.Email
	}
	if !isEmptyField(update.FirstName) {
		idnty.Phone = update.Phone
	}
	if !isEmptyField(update.FirstName) {
		idnty.Dob = update.Dob
	}
	if !isEmptyField(update.FirstName) {
		idnty.PresentAddress = update.PresentAddress
	}
	if !isEmptyField(update.FirstName) {
		idnty.PermanentAddress = update.PermanentAddress
	}

	jsonByte, err := json.Marshal(idnty)
	if err != nil {
		return err
	}

	fmt.Printf("update identity data to world state")
	return ctx.GetStub().PutState(id, jsonByte)
}

// GetSubmittingClientIdentity returns the name and issuer of the identity that
// invokes the smart contract. This function base64 decodes the identity string
// before returning the value to the client or smart contract.
func (s *SmartContract) GetSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}

func isEmptyField(f string) bool {
	if f == "" {
		return true
	}
	return false
}
