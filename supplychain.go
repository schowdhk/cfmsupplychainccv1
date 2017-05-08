package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/crypto/primitives"
)

var logger = shim.NewLogger("CFMSupplyChainChainCode")

// Chaincode default interface
type CFMSupplyChainChainCode struct {
}

// Init initializes the smart contracts
func (t *CFMSupplyChainChainCode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Info("Init called")

	return nil, nil
}

// Invoke entry point
func (t *CFMSupplyChainChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Info("Invoke called")

	return nil, nil
}

// Init initializes the smart contracts
func (t *CFMSupplyChainChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Info("Query called")
	return nil, nil
}
func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(CFMSupplyChainChainCode))
	if err != nil {
		fmt.Printf("Error starting ABC: %s", err)
	}
}
