package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
)

func readTxsInput(inputFile string) ([]*transaction.FrontendTransaction, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(workingDir, inputFile)
	jsonFile, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	bytesFromJson, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	txs := make([]*transaction.FrontendTransaction, 0)
	err = json.Unmarshal(bytesFromJson, &txs)
	if err != nil {
		return nil, err
	}

	log.Info("read from input", "file", inputFile, "num of txs", len(txs))
	return txs, nil
}
