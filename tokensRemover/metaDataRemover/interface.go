package main

import (
	"context"

	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/data"
)

type proxyProvider interface {
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	GetDefaultTransactionArguments(
		ctx context.Context,
		address core.AddressHandler,
		networkConfigs *data.NetworkConfig,
	) (transaction.FrontendTransaction, string, error)
}

type transactionInteractor interface {
	ApplyUserSignature(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error
}

type pemProvider interface {
	getPrivateKeyAndAddress(pemFile string) (*skAddress, error)
}
