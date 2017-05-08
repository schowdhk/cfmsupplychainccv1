package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/crypto/primitives"
)

var logger = shim.NewLogger("CFMSupplyChainChainCode")

const ALL_ELEMENENTS = "ALL_RECS"

// Chaincode default interface
type CFMSupplyChainChainCode struct {
}
type TrackingRecord struct {
	ShipmentNumber string `json:"shipmentNumber"`
	Date           string `json:"date"`
	ShipmentWt     string `json:"shipmentWt"`
	ShipType       string `json:"shipType"`
	ShipSrc        string `json:"shipSrc"`
	ShipDest       string `json:"shipDest"`
	ShipingComp    string `json:"shipingComp"`
	VehicleID      string `json:"vehicleId"`
	SealNumber     string `json:"sealNumber"`
	ContractNumber string `json:"contractNumber"`
	ExpYield       string `json:"expYield"`
	Type           string `json:"type"`
	OreType        string `json:"oreType"`
	SrcMine        string `json:"srcMine"`
	ShipperRecvdWt string `json:"shipperRecvdWt"`
	DestRecvdWt    string `json:"destRecvdWt"`
	Status         string `json:"status"`
	Catetogy       string `json:"catetogy"`
}

// Init initializes the smart contracts
func (t *CFMSupplyChainChainCode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Info("Init called")
	//Place an empty arry
	stub.PutState(ALL_ELEMENENTS, []byte("[]"))
	return nil, nil
}

// Invoke entry point
func (t *CFMSupplyChainChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Info("Invoke called")
	ext := CFMSupplyChainChainCode{}

	ext.createShipment(stub, args)
	return nil, nil
}
func (t *CFMSupplyChainChainCode) updateMasterRecords(stub shim.ChaincodeStubInterface, shipmentNumber string) error {
	var recordList []string
	recBytes, _ := stub.GetState(ALL_ELEMENENTS)

	err := json.Unmarshal(recBytes, &recordList)
	if err != nil {
		return errors.New("Failed to unmarshal updateMasterReords ")
	}
	recordList = append(recordList, shipmentNumber)
	bytesToStore, _ := json.Marshal(recordList)
	logger.Info("After addition" + string(bytesToStore))
	stub.PutState(ALL_ELEMENENTS, bytesToStore)
	return nil
}

// Creating a new shipment
func (t *CFMSupplyChainChainCode) createShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Info("createShipment called")

	shipmentNumber := args[0]
	payload := args[1]
	var shipmentRecord TrackingRecord
	err := json.Unmarshal([]byte(payload), &shipmentRecord)
	if err != nil {
		return nil, errors.New("Failed to unmarshal createShipment ")
	}
	stub.PutState(shipmentNumber, []byte(payload))
	ext := CFMSupplyChainChainCode{}
	ext.updateMasterRecords(stub, shipmentNumber)
	logger.Info("Received and unmarshaed the payload : " + payload)

	return nil, nil
}

// Init initializes the smart contracts
func (t *CFMSupplyChainChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Info("Query called")
	return nil, nil
}
func main() {
	logger.SetLevel(shim.LogInfo)
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(CFMSupplyChainChainCode))
	if err != nil {
		fmt.Printf("Error starting ABC: %s", err)
	}
}
