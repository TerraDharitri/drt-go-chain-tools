package export

import (
	"strings"

	"github.com/TerraDharitri/drt-go-chain-core/core/pubkeyConverter"
	"github.com/TerraDharitri/drt-go-chain-tools/trieTools/trieToolsCommon"
)

const (
	FormatterNamePlainText   = "plain-text"
	FormatterNamePlainJson   = "plain-json"
	FormatterNameRosettaJson = "rosetta-json"
	fourSpaces               = "    "
	addressLength            = 32
)

var (
	AllFormattersNames  = strings.Join([]string{FormatterNamePlainText, FormatterNamePlainText, FormatterNameRosettaJson}, ", ")
	addressConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(addressLength, trieToolsCommon.WalletHRP)
)

type formatterArgs struct {
	currency         string
	currencyDecimals uint
}
