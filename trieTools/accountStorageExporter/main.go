package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/core/check"
	"github.com/TerraDharitri/drt-go-chain-core/core/pubkeyConverter"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/errChan"
	"github.com/TerraDharitri/drt-go-chain/common/holders"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/testscommon/enableEpochsHandlerMock"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/TerraDharitri/drt-go-chain-tools/trieTools/accountStorageExporter/config"
	"github.com/TerraDharitri/drt-go-chain-tools/trieTools/trieToolsCommon"
	"github.com/urfave/cli"
)

const (
	logFilePrefix  = "account-storage-exporter"
	rootHashLength = 32
	addressLength  = 32
)

func main() {
	app := cli.NewApp()
	app.Name = "Accounts Storage Exporter CLI app"
	app.Usage = "This is the entry point for the tool that exports the storage of a given account"
	app.Flags = getFlags()
	app.Authors = []cli.Author{
		{
			Name:  "Team Dharitri",
			Email: "contact@dharitri.org",
		},
	}

	app.Action = func(c *cli.Context) error {
		return startProcess(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
		return
	}

	log.Info("finished exporting the storage")
}

func startProcess(c *cli.Context) error {
	flagsConfig := getFlagsConfig(c)

	_, errLogger := trieToolsCommon.AttachFileLogger(log, logFilePrefix, flagsConfig.ContextFlagsConfig)
	if errLogger != nil {
		return errLogger
	}

	log.Info("sanity checks...")

	err := logger.SetLogLevel(flagsConfig.LogLevel)
	if err != nil {
		return err
	}

	rootHash, err := hex.DecodeString(flagsConfig.HexRootHash)
	if err != nil {
		return fmt.Errorf("%w when decoding the provided hex root hash", err)
	}
	if len(rootHash) != rootHashLength {
		return fmt.Errorf("wrong root hash length: expected %d, got %d", rootHashLength, len(rootHash))
	}

	maxDBValue, err := trieToolsCommon.GetMaxDBValue(filepath.Join(flagsConfig.WorkingDir, flagsConfig.DbDir), log)
	if err != nil {
		return err
	}

	log.Info("starting exporting storage", "pid", os.Getpid())

	return exportStorage(flagsConfig.Address, flagsConfig, rootHash, maxDBValue)
}

func exportStorage(address string, flags config.ContextFlagsConfigAddr, mainRootHash []byte, maxDBValue int) error {
	addressConverter, err := pubkeyConverter.NewBech32PubkeyConverter(addressLength, trieToolsCommon.WalletHRP)
	if err != nil {
		return err
	}

	db, err := trieToolsCommon.CreatePruningStorer(flags.ContextFlagsConfig, maxDBValue)
	if err != nil {
		return err
	}

	enableEpochsHandler := &enableEpochsHandlerMock.EnableEpochsHandlerStub{
		IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
			return flag == common.AutoBalanceDataTriesFlag ||
				flag == common.DynamicDCDTFlag
		},
	}

	tr, err := trieToolsCommon.CreateTrie(db, enableEpochsHandler)
	if err != nil {
		return err
	}

	defer func() {
		errNotCritical := tr.Close()
		log.LogIfError(errNotCritical)
	}()

	accDb, err := trieToolsCommon.NewAccountsAdapter(tr, enableEpochsHandler)
	if err != nil {
		return err
	}

	rootHashHolder := holders.NewDefaultRootHashesHolder(mainRootHash)
	err = accDb.RecreateTrie(rootHashHolder)
	if err != nil {
		return err
	}

	addressBytes, err := addressConverter.Decode(address)
	if err != nil {
		return err
	}

	account, err := accDb.GetExistingAccount(addressBytes)
	if err != nil {
		return err
	}

	userAccount, ok := account.(state.UserAccountHandler)
	if !ok {
		return fmt.Errorf("cannot cast AccountHandler to UserAccountHandler")
	}

	if check.IfNil(userAccount.DataTrie()) {
		return fmt.Errorf("the provided address doesn't have a data trie")
	}

	iteratorChannels := &common.TrieIteratorChannels{
		LeavesChan: make(chan core.KeyValueHolder, common.TrieLeavesChannelDefaultCapacity),
		ErrChan:    errChan.NewErrChanWrapper(),
	}
	err = userAccount.GetAllLeaves(iteratorChannels, context.Background())
	if err != nil {
		return err
	}

	log.Info("data trie for account", "address", address, "root hash", userAccount.GetRootHash())
	for leaf := range iteratorChannels.LeavesChan {
		log.Info("key-value pair", "key", leaf.Key(), "value", leaf.Value())
	}

	return iteratorChannels.ErrChan.ReadFromChanNonBlocking()
}
