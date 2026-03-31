package mocks

import (
	"context"

	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/data"
)

// ProxyStub -
type ProxyStub struct {
	GetNetworkConfigCalled               func(ctx context.Context) (*data.NetworkConfig, error)
	GetDefaultTransactionArgumentsCalled func(
		ctx context.Context,
		address core.AddressHandler,
		networkConfigs *data.NetworkConfig,
	) (transaction.FrontendTransaction, string, error)
	SendTransactionCalled func(ctx context.Context, tx *transaction.FrontendTransaction) (string, error)
	GetAccountCalled      func(ctx context.Context, address core.AddressHandler) (*data.Account, error)
}

// GetNetworkConfig -
func (ps *ProxyStub) GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error) {
	if ps.GetNetworkConfigCalled != nil {
		return ps.GetNetworkConfigCalled(ctx)
	}

	return nil, nil
}

// GetDefaultTransactionArguments -
func (ps *ProxyStub) GetDefaultTransactionArguments(
	ctx context.Context,
	address core.AddressHandler,
	networkConfigs *data.NetworkConfig,
) (transaction.FrontendTransaction, string, error) {
	if ps.GetDefaultTransactionArgumentsCalled != nil {
		return ps.GetDefaultTransactionArgumentsCalled(ctx, address, networkConfigs)
	}

	return transaction.FrontendTransaction{}, "", nil
}

// SendTransaction -
func (ps *ProxyStub) SendTransaction(ctx context.Context, tx *transaction.FrontendTransaction) (string, error) {
	if ps.SendTransactionCalled != nil {
		return ps.SendTransactionCalled(ctx, tx)
	}

	return "", nil
}

// GetAccount -
func (ps *ProxyStub) GetAccount(ctx context.Context, address core.AddressHandler) (*data.Account, error) {
	if ps.GetAccountCalled != nil {
		return ps.GetAccountCalled(ctx, address)
	}

	return nil, nil
}
