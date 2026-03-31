package mocks

import (
	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-sdk/core"
)

// TransactionInteractorStub -
type TransactionInteractorStub struct {
	ApplyUserSignatureCalled func(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error
}

// ApplyUserSignature -
func (tis *TransactionInteractorStub) ApplyUserSignature(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error {
	if tis.ApplyUserSignatureCalled != nil {
		return tis.ApplyUserSignatureCalled(cryptoHolder, tx)
	}

	return nil
}
