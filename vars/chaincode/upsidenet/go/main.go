// UpsideNet Chaincode
// --------------------
// Este chaincode define dois tipos de ativos:
//   1. DimensionalEnergyMeasurement → ativo PRIVADO (armazenado em uma coleção privada entre organizações específicas)
//   2. EnvironmentalControlData → ativo PÚBLICO (armazenado no ledger público, acessível a todas as organizações)
// O objetivo é demonstrar o uso de dados privados e públicos dentro da mesma rede Fabric.

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

//=========================================================================================
// Estrutura principal do chaincode
//=========================================================================================

// UpsideChaincode é o tipo principal que implementa a interface Chaincode.
// Não armazena estado — serve apenas como “container” para os métodos Init e Invoke.
type UpsideChaincode struct{}

//=========================================================================================
// Estruturas de dados (ativos)
//=========================================================================================

// DimensionalEnergyMeasurement representa uma medição de energia dimensional.
// Este ativo é privado e será armazenado em uma coleção privada entre Hawkins e Montauk.
type DimensionalEnergyMeasurement struct {
	DocType     string  `json:"docType"`     // Tipo do documento (identificação do ativo)
	ID          string  `json:"id"`          // Identificador único
	Location    string  `json:"location"`    // Local onde a medição foi realizada
	EnergyLevel float64 `json:"energyLevel"` // Nível de energia
	Frequency   float64 `json:"frequency"`   // Frequência associada à medição
	Timestamp   string  `json:"timestamp"`   // Momento do registro
	RecordedBy  string  `json:"recordedBy"`  // Usuário ou sensor responsável pela leitura
}

// EnvironmentalControlData representa dados de controle ambiental (públicos).
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

//=========================================================================================
// Métodos obrigatórios da interface Chaincode
//=========================================================================================

// Init é executado apenas uma vez, quando o chaincode é instanciado.
// Pode ser usado para inicializar variáveis ou preparar o ledger.
func (t *UpsideChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("UpsideNet Chaincode Init")
	return shim.Success(nil)
}

// Invoke é chamado a cada transação e determina qual função do chaincode será executada.
// A escolha é feita pelo nome da função passado como primeiro parâmetro.
func (t *UpsideChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
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
		return shim.Error("Função desconhecida: " + function)
	}
}

//=========================================================================================
// Funções relacionadas a DimensionalEnergyMeasurement (dados privados)
//=========================================================================================

// initDimensionalEnergy cria uma nova medição dimensional e salva na coleção privada.
// Espera 4 argumentos: localização, nível de energia, frequência e responsável.
func (t *UpsideChaincode) initDimensionalEnergy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Esperado 4 argumentos: location, energyLevel, frequency, recordedBy")
	}

	// Conversão dos parâmetros recebidos
	location := args[0]
	energyLevel, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error("Energy level deve ser numérico")
	}
	frequency, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return shim.Error("Frequency deve ser numérico")
	}
	recordedBy := args[3]

	// Gera ID único com base no timestamp atual
	id := fmt.Sprintf("energy_%d", time.Now().UnixNano())

	// Monta o objeto a ser armazenado
	dimEnergy := DimensionalEnergyMeasurement{
		DocType:     "dimensionalEnergy",
		ID:          id,
		Location:    location,
		EnergyLevel: energyLevel,
		Frequency:   frequency,
		Timestamp:   time.Now().Format(time.RFC3339),
		RecordedBy:  recordedBy,
	}

	dimEnergyBytes, _ := json.Marshal(dimEnergy)

	// Grava na coleção privada "collectionDimensionalEnergy"
	// Esta coleção deve estar configurada para permitir apenas organizações específicas.
	err = stub.PutPrivateData("collectionDimensionalEnergy", id, dimEnergyBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Retorna o ID gerado
	return shim.Success([]byte(id))
}

// readDimensionalEnergy lê uma medição privada da coleção.
// Só funciona se a organização do peer tiver acesso à coleção.
func (t *UpsideChaincode) readDimensionalEnergy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Esperado 1 argumento: ID")
	}

	id := args[0]
	valAsbytes, err := stub.GetPrivateData("collectionDimensionalEnergy", id)
	if err != nil {
		return shim.Error("Falha ao obter dados privados: " + err.Error())
	}
	if valAsbytes == nil {
		return shim.Error("Nenhum registro encontrado para ID: " + id)
	}

	return shim.Success(valAsbytes)
}

// getAllDimensionalEnergy retorna todos os registros privados acessíveis ao peer.
func (t *UpsideChaincode) getAllDimensionalEnergy(stub shim.ChaincodeStubInterface) pb.Response {
	resultsIterator, err := stub.GetPrivateDataByRange("collectionDimensionalEnergy", "", "")
	if err != nil {
		return shim.Error("Falha ao consultar a coleção privada: " + err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	first := true

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("Erro ao iterar resultados: " + err.Error())
		}
		if !first {
			buffer.WriteString(",")
		}
		buffer.WriteString(fmt.Sprintf("{\"Key\":\"%s\",\"Record\":%s}", queryResponse.Key, string(queryResponse.Value)))
		first = false
	}

	buffer.WriteString("]")
	return shim.Success(buffer.Bytes())
}

//=========================================================================================
// Funções relacionadas a EnvironmentalControlData (dados públicos)
//=========================================================================================

// initEnvironmentalControl cria um novo registro público de controle ambiental.
func (t *UpsideChaincode) initEnvironmentalControl(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 5 {
		return shim.Error("Esperado 5 argumentos: location, temperature, humidity, pressure, recordedBy")
	}

	location := args[0]
	temperature, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error("Temperature deve ser numérico")
	}
	humidity, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return shim.Error("Humidity deve ser numérico")
	}
	pressure, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return shim.Error("Pressure deve ser numérico")
	}
	recordedBy := args[4]

	id := fmt.Sprintf("env_%d", time.Now().UnixNano())

	// Cria e serializa o ativo
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
	envControlBytes, _ := json.Marshal(envControl)

	// Armazena no ledger público
	err = stub.PutState(id, envControlBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(id))
}

// readEnvironmentalControl lê um registro público do ledger.
func (t *UpsideChaincode) readEnvironmentalControl(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Esperado 1 argumento: ID")
	}

	id := args[0]
	valAsbytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Falha ao obter estado: " + err.Error())
	}
	if valAsbytes == nil {
		return shim.Error("Nenhum registro encontrado para ID: " + id)
	}

	return shim.Success(valAsbytes)
}

// getAllEnvironmentalControl lista todos os registros públicos no ledger.
func (t *UpsideChaincode) getAllEnvironmentalControl(stub shim.ChaincodeStubInterface) pb.Response {
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
		// Filtra apenas registros do tipo "environmentalControl"
		var data map[string]interface{}
		json.Unmarshal(queryResponse.Value, &data)
		if docType, ok := data["docType"].(string); ok && docType == "environmentalControl" {
			if bArrayMemberAlreadyWritten {
				buffer.WriteString(",")
			}
			buffer.WriteString(fmt.Sprintf("{\"Key\":\"%s\",\"Record\":%s}", queryResponse.Key, string(queryResponse.Value)))
			bArrayMemberAlreadyWritten = true
		}
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

//=========================================================================================
// Função principal - ponto de entrada do chaincode
//=========================================================================================

func main() {
	err := shim.Start(&UpsideChaincode{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao iniciar o chaincode UpsideNet: %s", err)
		os.Exit(2)
	}
}
