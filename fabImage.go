/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */


package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Image struct {
	imageName   string `json:"imageName"`
	imageSize  string `json:"imageSize"`
	Owner string `json:"Owner"`
}

/*
 * The Init method is called when the Smart Contract "fabImage" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabImage"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryImage" {
		return s.queryImage(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "UploadImage" {
		return s.UploadImage(APIstub, args)
	} else if function == "queryAllImgs" {
		return s.queryAllImgs(APIstub)
	} else if function == "transferImage" {
		return s.transferImage(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryImage(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	imageAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(imageAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	imgs := []Image{
		Image{imageName: "Typhoid", imageSize: " 1 MB ", Owner: "Tomoko"},
		Image{imageName: "Pnemonia", imageSize: " 1 MB", Owner: "Jin"},
		Image{imageName: "Anemia", imageSize: "1 MB", Owner: "Max"},
	}

	i := 0
	for i < len(imgs) {
		fmt.Println("i is ", i)
		imageAsBytes, _ := json.Marshal(imgs[i])
		APIstub.PutState("IMG"+strconv.Itoa(i), imageAsBytes)
		fmt.Println("Added", imgs[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) UploadImage(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var img = Image{imageName: args[1], imageSize: args[2], Owner: args[3]}

	imageAsBytes, _ := json.Marshal(img)
	APIstub.PutState(args[0], imageAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllImgs(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "IMG0"
	endKey := "IMG999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllImgs:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) transferImage(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	imageAsBytes, _ := APIstub.GetState(args[0])
	img := Image{}

	json.Unmarshal(imageAsBytes, &img)
	img.Owner = args[1]
	img.imageSize = args[2]

	imageAsBytes, _ = json.Marshal(img)
	APIstub.PutState(args[0], imageAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
