package main

import (
	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/interactors"
)

type skAddress struct {
	secretKey []byte
	address   core.AddressHandler
}

type pemDataProvider struct {
}

func (pdp *pemDataProvider) getPrivateKeyAndAddress(pemFile string) (*skAddress, error) {
	w := interactors.NewWallet()
	privateKey, err := w.LoadPrivateKeyFromPemFile(pemFile)
	if err != nil {
		return nil, err
	}

	address, err := w.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		return nil, err

	}

	return &skAddress{
		secretKey: privateKey,
		address:   address,
	}, nil
}
