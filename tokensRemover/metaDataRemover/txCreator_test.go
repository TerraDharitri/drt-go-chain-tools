package main

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/data/transaction"
	"github.com/TerraDharitri/drt-go-chain-tools/tokensRemover/metaDataRemover/mocks"
	"github.com/TerraDharitri/drt-go-sdk/core"
	"github.com/TerraDharitri/drt-go-sdk/data"
	"github.com/stretchr/testify/require"
)

func TestTxCreator_CreateTxs(t *testing.T) {
	t.Parallel()

	addr, err := data.NewAddressFromBech32String("drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf")
	require.Nil(t, err)
	sk, err := hex.DecodeString("413f42575f7f26fad3317a778771212fdb80245850981e48b58a4f25e344e8f9")
	require.Nil(t, err)
	pemData := &skAddress{
		secretKey: sk,
		address:   addr,
	}

	txData1 := []byte("txData1")
	txData2 := []byte("txData2")
	txsData := [][]byte{txData1, txData2}

	nonce := uint64(4)
	additionalGas := uint64(500)

	networkCfg := &data.NetworkConfig{
		ChainID:        "1",
		MinGasPrice:    100,
		MinGasLimit:    500,
		GasPerDataByte: 15,
	}

	addrStr, err := addr.AddressAsBech32String()
	require.Nil(t, err)

	proxy := &mocks.ProxyStub{
		GetNetworkConfigCalled: func(ctx context.Context) (*data.NetworkConfig, error) {
			return networkCfg, nil
		},

		GetDefaultTransactionArgumentsCalled: func(ctx context.Context, address core.AddressHandler, networkConfigs *data.NetworkConfig) (transaction.FrontendTransaction, string, error) {
			require.Equal(t, networkCfg, networkConfigs)
			require.Equal(t, addr, address)

			return transaction.FrontendTransaction{
				Nonce:    nonce,
				Sender:   addrStr,
				ChainID:  networkCfg.ChainID,
				GasPrice: networkCfg.MinGasPrice,
			}, "", nil
		},
	}

	txIdx := 0
	txInteractor := &mocks.TransactionInteractorStub{
		ApplyUserSignatureCalled: func(cryptoHolder core.CryptoComponentsHolder, tx *transaction.FrontendTransaction) error {
			require.Equal(t, &transaction.FrontendTransaction{
				Nonce:    nonce,
				Value:    "0",
				Data:     txsData[txIdx],
				ChainID:  networkCfg.ChainID,
				GasPrice: networkCfg.MinGasPrice,
				GasLimit: 1105,
				Sender:   addrStr,
				Receiver: addrStr,
			}, tx)

			defer func() {
				nonce++
				txIdx++
			}()

			return nil
		},
	}

	txc, err := newTxCreator(proxy, txInteractor)
	require.Nil(t, err)
	signedTxs, err := txc.createTxs(pemData, txsData, additionalGas)
	require.Nil(t, err)
	require.Len(t, signedTxs, 2)
}
