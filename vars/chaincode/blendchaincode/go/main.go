/*
UpsideNet Chaincode - Proof of Concept
Manages dimensional energy measurements and environmental control data
with private collections for data privacy
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type BlendChaincode struct {
}

// DimensionalEnergyMeasurement - Asset for storing dimensional energy data (PRIVATE - only Hawkins and Montauk)
type DimensionalEnergyMeasurement struct {
	DocType     string  `json:"docType"`
	ID          string  `json:"id"`
	Location    string  `json:"location"`
	EnergyLevel float64 `json:"energyLevel"`
	Frequency   float64 `json:"frequency"`
	Timestamp   string  `json:"timestamp"`
	RecordedBy  string  `json:"recordedBy"`
}

// EnvironmentalControlData - Asset for storing environmental control data (PUBLIC - all organizations)
type EnvironmentalControlData struct {
	DocType        string  `json:"docType"`
	ID             string  `json:"id"`
	Location       string  `json:"location"`
	Temperature    float64 `json:"temperature"`
	Humidity       float64 `json:"humidity"`
	Pressure       float64 `json:"pressure"`
	RadiationLevel float64 `json:"radiationLevel"`
	Timestamp      string  `json:"timestamp"`
	RecordedBy     string  `json:"recordedBy"`
}

func (t *BlendChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("UpsideNet Chaincode Init")
	return shim.Success(nil)
}

func (t *BlendChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "initDimensionalEnergy":
		return t.initDimensionalEnergy(stub, args)
	case "initEnvironmentalControl":
		return t.initEnvironmentalControl(stub, args)
	case "readDimensionalEnergy":
		return t.readDimensionalEnergy(stub, args)
	case "readEnvironmentalControl":
		return t.readEnvironmentalControl(stub, args)
	case "getAllDimensionalEnergy":
		return t.getAllDimensionalEnergy(stub)
	case "getAllEnvironmentalControl":
		return t.getAllEnvironmentalControl(stub)
	default:
		return shim.Error("Unknown function: " + function)
	}
}

// initDimensionalEnergy - Create dimensional energy measurement (PRIVATE collection)
// Only Hawkins and Montauk can access this data
func (t *BlendChaincode) initDimensionalEnergy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5 (id, location, energyLevel, frequency, recordedBy)")
	}

	id := args[0]
	location := args[1]
	energyLevel, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return shim.Error("Energy level must be a number")
	}
	frequency, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return shim.Error("Frequency must be a number")
	}
	recordedBy := args[4]

	dimEnergy := DimensionalEnergyMeasurement{
		DocType:     "dimensionalEnergy",
		ID:          id,
		Location:    location,
		EnergyLevel: energyLevel,
		Frequency:   frequency,
		Timestamp:   time.Now().Format(time.RFC3339),
		RecordedBy:  recordedBy,
	}

	dimEnergyBytes, err := json.Marshal(dimEnergy)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Save to private collection - only accessible by Hawkins and Montauk
	err = stub.PutPrivateData("collectionDimensionalEnergy", id, dimEnergyBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// initEnvironmentalControl - Create environmental control data (PUBLIC - all can access)
func (t *BlendChaincode) initEnvironmentalControl(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6 (id, location, temperature, humidity, pressure, recordedBy)")
	}

	id := args[0]
	location := args[1]
	temperature, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return shim.Error("Temperature must be a number")
	}
	humidity, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return shim.Error("Humidity must be a number")
	}
	pressure, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		return shim.Error("Pressure must be a number")
	}
	recordedBy := args[5]

	envControl := EnvironmentalControlData{
		DocType:        "environmentalControl",
		ID:             id,
		Location:       location,
		Temperature:    temperature,
		Humidity:       humidity,
		Pressure:       pressure,
		RadiationLevel: 0.0,
		Timestamp:      time.Now().Format(time.RFC3339),
		RecordedBy:     recordedBy,
	}

	envControlBytes, err := json.Marshal(envControl)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Save to public state - accessible by all organizations
	err = stub.PutState(id, envControlBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// readDimensionalEnergy - Read dimensional energy measurement from private collection
func (t *BlendChaincode) readDimensionalEnergy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting ID")
	}

	id := args[0]
	valAsbytes, err := stub.GetPrivateData("collectionDimensionalEnergy", id)
	if err != nil {
		return shim.Error("Failed to get dimensional energy: " + err.Error())
	}
	if valAsbytes == nil {
		return shim.Error("Dimensional energy does not exist: " + id)
	}

	return shim.Success(valAsbytes)
}

// readEnvironmentalControl - Read environmental control data
func (t *BlendChaincode) readEnvironmentalControl(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting ID")
	}

	id := args[0]
	valAsbytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get environmental control: " + err.Error())
	}
	if valAsbytes == nil {
		return shim.Error("Environmental control does not exist: " + id)
	}

	return shim.Success(valAsbytes)
}

// getAllDimensionalEnergy - Get all dimensional energy measurements
func (t *BlendChaincode) getAllDimensionalEnergy(stub shim.ChaincodeStubInterface) pb.Response {
	resultsIterator, err := stub.GetPrivateDataByRange("collectionDimensionalEnergy", "", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\",\"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

// getAllEnvironmentalControl - Get all environmental control data
func (t *BlendChaincode) getAllEnvironmentalControl(stub shim.ChaincodeStubInterface) pb.Response {
	resultsIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Filter only environmental control data
		value := string(queryResponse.Value)
		if value != "" && len(value) > 0 {
			var data map[string]interface{}
			json.Unmarshal([]byte(value), &data)
			if docType, ok := data["docType"].(string); ok && docType == "environmentalControl" {
				if bArrayMemberAlreadyWritten == true {
					buffer.WriteString(",")
				}
				buffer.WriteString("{\"Key\":\"")
				buffer.WriteString(queryResponse.Key)
				buffer.WriteString("\",\"Record\":")
				buffer.WriteString(string(queryResponse.Value))
				buffer.WriteString("}")
				bArrayMemberAlreadyWritten = true
			}
		}
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(&BlendChaincode{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Exiting UpsideNet chaincode: %s", err)
		os.Exit(2)
	}
}
