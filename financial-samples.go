/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//participant data
type Participant struct {
	Name    string `json:"name"`
	Balance int    `json:"balanace"`
}

type InterestRateSwap struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("financial-samples Init")

	// _, args := stub.GetFunctionAndParameters()
	// var A, B string    // Entities
	// var Aval, Bval int // Asset holdings
	// var err error

	// if len(args) != 4 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 4")
	// }

	// // Initialize the chaincode
	// A = args[0]
	// Aval, err = strconv.Atoi(args[1])
	// if err != nil {
	// 	return shim.Error("Expecting integer value for asset holding")
	// }
	// B = args[2]
	// Bval, err = strconv.Atoi(args[3])
	// if err != nil {
	// 	return shim.Error("Expecting integer value for asset holding")
	// }
	// fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// // Write the state to the ledger
	// err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	// err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("financial-samples Invoke")
	function, args := stub.GetFunctionAndParameters()
	fmt.Print("function:", function)

	if function == "transfer" {
		// Make payment of X units from A to B
		return t.invoke(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"transfer\" \"delete\" \"query\"")
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A, B string // Entity IDs
	var X int       // Transaction value
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3: payer's key, payee's key, and amount to transfer")
	}

	A = args[0]
	B = args[1]

	AInfoBytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state for %s", A)
	}
	if AInfoBytes == nil {
		return shim.Error("%s not found", A)
	}
	AInfo := Participant{}
	err = json.Unmarshal(AInfoBytes, &AInfo)

	BInfoBytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state for %s", B)
	}
	if BInfoBytes == nil {
		return shim.Error("%s entity not found", B)
	}
	BInfo := Participant{}
	err = json.Unmarshal(BInfoBytes, &BInfo)

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	AInfo.Balance = AInfo.Balance - X
	BInfo.Balance = BInfo.Balance + X
	fmt.Printf("%s's balance = %d, %s's balance' = %d\n", A, AInfo.Balance, B, BInfo.Balance)

	// Write the state back to the ledger
	AInfoUpdatedBytes, _ := JSON.Marshal(AInfo)
	err = stub.PutState(A, AInfoUpdatedBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	BInfoUpdatedBytes, _ := JSON.Marshal(BInfo)
	err = stub.PutState(B, BInfoUpdatedBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: key to delete")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: key to query")
	}

	A = args[0]

	// Get the state from the ledger
	AInfoBytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if AInfoBytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(AInfoBytes) + "\"}"
	fmt.Printf("Query Response: %s\n", jsonResp)
	return shim.Success(AInfoBytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
