package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/nramin/stock-price/structs"
	"gopkg.in/yaml.v3"
)

func main() {
	var result structs.StockPrices
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
		result.Stocks = append(result.Stocks, structs.Stock{
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

	marshaledResult, err := json.Marshal(result)
	if err != nil {
		printError(&result, err.Error())
		os.Exit(0)
	}

	fmt.Println(string(marshaledResult))
	os.Exit(0)
}

func readYamlFile(filePath string, result *structs.StockPrices) YamlConfig {
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

func printError(result *structs.StockPrices, error string) {
	success := new(bool)
	*success = false

	result.Success = success
	result.Error = error
	marshaledResult, _ := json.Marshal(result)
	fmt.Println(string(marshaledResult))
}
