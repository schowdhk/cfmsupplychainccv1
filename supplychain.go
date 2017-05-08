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
func (t *CFMSupplyChainChainCode) getAllRecordsList(stub shim.ChaincodeStubInterface) ([]string, error) {
	var recordList []string
	recBytes, _ := stub.GetState(ALL_ELEMENENTS)

	err := json.Unmarshal(recBytes, &recordList)
	if err != nil {
		return nil, errors.New("Failed to unmarshal getAllRecordsList ")
	}

	return recordList, nil
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

func (t *CFMSupplyChainChainCode) updateRecord(existingRecord map[string]string, fieldsToUpdate map[string]string) (string, error) {
	for _, key := range fieldsToUpdate {
		value := fieldsToUpdate[key]
		existingRecord[key] = value
	}
	outputMapBytes, _ := json.Marshal(existingRecord)
	logger.Info("updateRecord: Final json after update " + string(outputMapBytes))
	return string(outputMapBytes), nil
}

// Update and existing shipment record
func (t *CFMSupplyChainChainCode) updateShipment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var existingRecMap map[string]string
	var updatedFields map[string]string

	logger.Info("updateShipment called ")
	ext := CFMSupplyChainChainCode{}

	shipmentNumber := args[0]
	payload := args[1]
	logger.Info("updateShipment payload passed " + payload)

	//who :=args[2]
	recBytes, _ := stub.GetState(shipmentNumber)

	json.Unmarshal(recBytes, &existingRecMap)
	json.Unmarshal([]byte(payload), &updatedFields)

	updatedReord, _ := ext.updateRecord(existingRecMap, updatedFields)
	//Store the records
	stub.PutState(shipmentNumber, []byte(updatedReord))
	return nil, nil
}
func (t *CFMSupplyChainChainCode) getShipmentWithStatus(stub shim.ChaincodeStubInterface, status string) ([]byte, error) {
	logger.Info("getShipmentWithStatus called")
	ext := CFMSupplyChainChainCode{}
	recordsList, err := ext.getAllRecordsList(stub)
	if err != nil {
		return nil, errors.New("Unable to get all the records ")
	}
	var outputRecords []TrackingRecord
	outputRecords = make([]TrackingRecord, 0)
	for _, shipmentNumber := range recordsList {
		logger.Info("getShipmentWithStatus: Processing record " + shipmentNumber)
		recBytes, _ := stub.GetState(shipmentNumber)
		var record TrackingRecord
		json.Unmarshal(recBytes, &record)
		if status == record.Status {
			outputRecords = append(outputRecords, record)
		}
	}
	outputBytes, _ := json.Marshal(outputRecords)
	logger.Info("Returning records from getShipmentWithStatus " + string(outputBytes))
	return outputBytes, nil
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

	if function == "createShipment" {
		ext.createShipment(stub, args)
	}
	if function == "updateShipment" {
		ext.updateShipment(stub, args)
	}

	return nil, nil
}

// Query the rcords form the  smart contracts
func (t *CFMSupplyChainChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Info("Query called")
	ext := CFMSupplyChainChainCode{}
	if function == "getAllRecordsByStatus" {
		return ext.getShipmentWithStatus(stub, args[0])
	}
	return nil, nil
}

//Main method
func main() {
	logger.SetLevel(shim.LogInfo)
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(CFMSupplyChainChainCode))
	if err != nil {
		fmt.Printf("Error starting ABC: %s", err)
	}
}
