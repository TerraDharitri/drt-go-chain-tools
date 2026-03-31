package config

import "github.com/TerraDharitri/drt-go-chain-tools/trieTools/trieToolsCommon"

// ContextFlagsConfigAddr the configuration for flags
type ContextFlagsConfigAddr struct {
	trieToolsCommon.ContextFlagsConfig
	Address string
}
