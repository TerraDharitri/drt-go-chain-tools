package mocks

import (
	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/data"
)

// TransactionInteractorStub -
type TransactionInteractorStub struct {
	ApplySignatureAndGenerateTxCalled func(cryptoHolder core.CryptoComponentsHolder, arg data.ArgCreateTransaction) (*data.Transaction, error)
}

// ApplySignatureAndGenerateTx -
func (tis *TransactionInteractorStub) ApplySignatureAndGenerateTx(cryptoHolder core.CryptoComponentsHolder, arg data.ArgCreateTransaction) (*data.Transaction, error) {
	if tis.ApplySignatureAndGenerateTxCalled != nil {
		return tis.ApplySignatureAndGenerateTxCalled(cryptoHolder, arg)
	}

	return nil, nil
}
