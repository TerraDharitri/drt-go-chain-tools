package main

import (
	"context"
	"fmt"
	"math"
	"os"

	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/dcdt"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/common/errChan"
	"github.com/TerraDharitri/drt-go-chain/state"
	"github.com/TerraDharitri/drt-go-chain/vm/systemSmartContracts"
	logger "github.com/TerraDharitri/drt-go-chain-logger"
	"github.com/TerraDharitri/drt-go-chain-tools/trieTools/trieToolsCommon"
	"github.com/urfave/cli"
)

var log = logger.GetOrCreate("trie")

const (
	logFilePrefix  = "trie"
	rootHashLength = 32
	addressLength  = 32
)

func main() {
	app := cli.NewApp()
	app.Name = "Trie stats CLI app"
	app.Usage = "This is the entry point for the tool that prints stats about the state"
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

	log.Info("execution finished successfully")
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

	shardsState, err := LoadStateForAllShards(flagsConfig)
	if err != nil {
		return err
	}

	systemDcdtAddress := "drt1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls6prdez"
	systemDcdtAccount, err := GetAccountFromBech32String(systemDcdtAddress, shardsState[Meta])
	if err != nil {
		return err
	}

	tokensMap, err := getAllDCDTsFromSystemDcdtAccount(systemDcdtAccount)
	if err != nil {
		return err
	}

	systemAccounts, err := getSystemAccountsForShards(shardsState)
	if err != nil {
		return err
	}

	checkUpdateTokenTypeCalled(systemAccounts, tokensMap)
	for shardId, systemAccount := range systemAccounts {
		log.Info("shard", "shardId", shardId, "rootHash", systemAccount.GetRootHash())
		getNumTokensFromSystemAccount(systemAccount)
	}

	return nil
}

func getNumTokensFromSystemAccount(account state.UserAccountHandler) {
	iteratorChannels := &common.TrieIteratorChannels{
		LeavesChan: make(chan core.KeyValueHolder, common.TrieLeavesChannelDefaultCapacity),
		ErrChan:    errChan.NewErrChanWrapper(),
	}
	err := account.GetAllLeaves(iteratorChannels, context.Background())
	if err != nil {
		log.Error("can not get all leaves", "error", err)
		return
	}
	numNonMetaData := 0
	tokenTypes := make(map[string]uint64)
	for leaf := range iteratorChannels.LeavesChan {
		metaData := &dcdt.DCDigitalToken{}
		errUnmarshal := trieToolsCommon.Marshaller.Unmarshal(metaData, leaf.Value())
		if errUnmarshal != nil {
			numNonMetaData++
			continue
		}
		tokenTypes[core.DCDTType(metaData.Type).String()]++
	}
	err = iteratorChannels.ErrChan.ReadFromChanNonBlocking()
	if err != nil {
		log.Error("can not read from error channel", "error", err)
	}

	log.Info("non metaData keys", "num", numNonMetaData)
	for tokenType, numTokens := range tokenTypes {
		log.Info("tokenType", "type", tokenType, "numTokens", numTokens)
	}
}

func checkUpdateTokenTypeCalled(systemAccounts map[ShardID]state.UserAccountHandler, tokensMap map[string][][]byte) {
	dcdtPrefix := []byte("NUMBATdcdt")
	for tokenType, tokenIds := range tokensMap {
		if tokenType == core.FungibleDCDT {
			continue
		}

		tokenTypeAsInt, err := core.ConvertDCDTTypeToUint32(tokenType)
		if err != nil {
			log.Error("can not convert token type to int", "tokenType", tokenType, "error", err)
			continue
		}

		for _, tokenId := range tokenIds {
			key := append(dcdtPrefix, tokenId...)

			checkTokenInAllShards(key, tokenTypeAsInt, systemAccounts)
		}
		log.Info("tokens", "tokenType", tokenType, "numTokens", len(tokenIds))
	}
}

func checkTokenInAllShards(tokenKey []byte, tokenType uint32, systemAccounts map[ShardID]state.UserAccountHandler) {
	for shardId, systemAccount := range systemAccounts {
		value, _, err := systemAccount.RetrieveValue(tokenKey)
		if err != nil {
			log.Error("can not get token data", "tokenId", tokenKey, "error", err)
			continue
		}

		if len(value) != 2 {
			log.Error("token data has wrong length", "tokenId", tokenKey, "value", value)
			continue
		}

		retrievedTokenType := value[1]
		dcdtTokenType, err := convertToDCDTTokenType(uint32(retrievedTokenType))
		if err != nil {
			log.Error("can not convert token type to int", "tokenType", tokenType, "error", err)
			continue
		}

		if dcdtTokenType != tokenType {
			log.Warn("token type is not the same", "tokenId", tokenKey, "retrievedTokenType", retrievedTokenType, "tokenType", tokenType, "shard", shardId)
		}
	}
}

func getAllDCDTsFromSystemDcdtAccount(systemDcdtAccount state.UserAccountHandler) (map[string][][]byte, error) {
	tokens := make(map[string][][]byte)

	iteratorChannels := &common.TrieIteratorChannels{
		LeavesChan: make(chan core.KeyValueHolder, common.TrieLeavesChannelDefaultCapacity),
		ErrChan:    errChan.NewErrChanWrapper(),
	}
	err := systemDcdtAccount.GetAllLeaves(iteratorChannels, context.Background())
	if err != nil {
		return nil, err
	}

	for leaf := range iteratorChannels.LeavesChan {
		data := &systemSmartContracts.DCDTDataV2{}
		errUnmarshal := trieToolsCommon.Marshaller.Unmarshal(data, leaf.Value())
		if errUnmarshal == nil {
			tokens[string(data.TokenType)] = append(tokens[string(data.TokenType)], leaf.Key())
			continue
		}

		dataV1 := &systemSmartContracts.DCDTDataV1{}
		errUnmarshal = trieToolsCommon.Marshaller.Unmarshal(dataV1, leaf.Value())
		if errUnmarshal == nil {
			tokens[string(data.TokenType)] = append(tokens[string(data.TokenType)], leaf.Key())
			continue
		}

		log.Warn("can not unmarshall data", "key", string(leaf.Key()), "value", leaf.Value())
	}

	return tokens, nil
}

func getSystemAccountsForShards(shardsState map[ShardID]state.AccountsAdapter) (map[ShardID]state.UserAccountHandler, error) {
	globalSettingsShard0Address := "drt1llllllllllllllllllllllllllllllllllllllllllllllllllls9258a4"
	globalSettingsShard0Account, err := GetAccountFromBech32String(globalSettingsShard0Address, shardsState[Shard0])
	if err != nil {
		return nil, err
	}
	globalSettingsShard1Account, err := GetAccountFromBech32String(globalSettingsShard0Address, shardsState[Shard1])
	if err != nil {
		return nil, err
	}
	globalSettingsShard2Account, err := GetAccountFromBech32String(globalSettingsShard0Address, shardsState[Shard2])
	if err != nil {
		return nil, err
	}

	systemAccounts := make(map[ShardID]state.UserAccountHandler)
	systemAccounts[Shard0] = globalSettingsShard0Account
	systemAccounts[Shard1] = globalSettingsShard1Account
	systemAccounts[Shard2] = globalSettingsShard2Account
	return systemAccounts, nil
}

func convertToDCDTTokenType(dcdtType uint32) (uint32, error) {
	switch dcdtType {
	case 0:
		return 0, fmt.Errorf("token type not set inside global settings handler")
	case 1:
		return uint32(core.Fungible), nil
	case 2:
		return uint32(core.NonFungible), nil
	case 3:
		return uint32(core.NonFungibleV2), nil
	case 4:
		return uint32(core.MetaFungible), nil
	case 5:
		return uint32(core.SemiFungible), nil
	case 6:
		return uint32(core.DynamicNFT), nil
	case 7:
		return uint32(core.DynamicSFT), nil
	case 8:
		return uint32(core.DynamicMeta), nil
	default:
		return math.MaxUint32, fmt.Errorf("invalid dcdt type: %d", dcdtType)
	}
}
