package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Estrutura para medições de energia extra dimensional
type MedicaoEnergia struct {
	ID          string  `json:"id"`
	Laboratorio string  `json:"laboratorio"`
	Potencia    float64 `json:"potencia"`
	Dimensao    string  `json:"dimensao"`
}

// Estrutura para dados de controle ambiental
type DadosAmbientais struct {
	ID          string  `json:"id"`
	Laboratorio string  `json:"laboratorio"`
	Temperatura float64 `json:"temperatura"`
	Umidade     float64 `json:"umidade"`
}

// Contrato inteligente principal
type ContratoUpsideNet struct {
	contractapi.Contract
}

// Criar Medição de Energia
	// Função para registrar uma nova medição de energia (dados privados)
func (c *ContratoUpsideNet) CriarMedicaoEnergia(ctx contractapi.TransactionContextInterface, id string, jsonDados string) error {
	return ctx.GetStub().PutPrivateData("EnergyMeasurementsCollection", id, []byte(jsonDados))
}

// Ler Medição de Energia
	// Função para consultar uma medição de energia (dados privados)
func (c *ContratoUpsideNet) LerMedicaoEnergia(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	dados, err := ctx.GetStub().GetPrivateData("EnergyMeasurementsCollection", id)
	if err != nil {
		return "", err
	}
	if dados == nil {
		return "", fmt.Errorf("nenhuma medição encontrada com o ID %s", id)
	}
	return string(dados), nil
}

// Registrar Dados Ambientais
	// Função para registrar dados ambientais (dados privados)
func (c *ContratoUpsideNet) CriarDadosAmbientais(ctx contractapi.TransactionContextInterface, id string, jsonDados string) error {
	return ctx.GetStub().PutPrivateData("EnvironmentalDataCollection", id, []byte(jsonDados))
}

// Ler Dados Ambientais
	// Função para consultar dados ambientais (dados privados)
func (c *ContratoUpsideNet) LerDadosAmbientais(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	dados, err := ctx.GetStub().GetPrivateData("EnvironmentalDataCollection", id)
	if err != nil {
		return "", err
	}
	if dados == nil {
		return "", fmt.Errorf("nenhum dado ambiental encontrado com o ID %s", id)
	}
	return string(dados), nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(ContratoUpsideNet))
	if err != nil {
		panic(fmt.Sprintf("Erro ao criar o chaincode UpsideNet: %v", err))
	}

	if err := chaincode.Start(); err != nil {
		panic(fmt.Sprintf("Erro ao iniciar o chaincode UpsideNet: %v", err))
	}
}