package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"gopkg.in/yaml.v3"
)

func main() {
	var result Result
	yamlConfig := readYamlFile("settings.yaml", &result)
	alpacaConfig := yamlConfig.Alpaca
	symbolDetails := yamlConfig.SymbolDetails
	client := marketdata.NewClient(marketdata.ClientOpts{
		APIKey:    alpacaConfig["api-key"],
		APISecret: alpacaConfig["api-secret"],
		BaseURL:   "https://data.alpaca.markets",
	})

	symbols := slices.Collect(maps.Keys(symbolDetails))
	quotes, err := client.GetLatestQuotes(symbols, marketdata.GetLatestQuoteRequest{
		Feed:     marketdata.IEX,
		Currency: "USD",
	})
	if err != nil {
		printError(&result, err.Error())
		os.Exit(0)
	}

	for symbol, quote := range quotes {
		value := quote.BidPrice * symbolDetails[symbol]
		result.Stocks = append(result.Stocks, Stock{
			Symbol:   symbol,
			Price:    quote.BidPrice,
			Quantity: symbolDetails[symbol],
			Value:    value,
		})
		result.TotalValue += value
	}

	success := new(bool)
	*success = true
	result.Success = success

	marshaledResult, _ := json.Marshal(result)
	fmt.Println(string(marshaledResult))
	os.Exit(0)
}

type Result struct {
	Stocks     []Stock `json:"stocks,omitempty"`
	TotalValue float64 `json:"totalValue,omitempty"`
	Success    *bool   `json:"success,omitempty"`
	Error      string  `json:"error,omitempty"`
}

type Stock struct {
	Symbol   string  `json:"symbol,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Quantity float64 `json:"quantity,omitempty"`
	Value    float64 `json:"balance,omitempty"`
}

func readYamlFile(filePath string, result *Result) YamlConfig {
	b, err := os.ReadFile(filePath)
	if err != nil {
		printError(result, "Unable to read input file "+filePath)
		os.Exit(0)
	}
	var yamlConfig YamlConfig

	err = yaml.Unmarshal([]byte(b), &yamlConfig)
	if err != nil {
		printError(result, err.Error())
		os.Exit(0)
	}

	return yamlConfig
}

type YamlConfig struct {
	Alpaca        map[string]string  `yaml:"alpaca"`
	SymbolDetails map[string]float64 `yaml:"symbol-details"`
}

func printError(result *Result, error string) {
	success := new(bool)
	*success = false

	result.Success = success
	result.Error = error
	marshaledResult, _ := json.Marshal(result)
	fmt.Println(string(marshaledResult))
}
