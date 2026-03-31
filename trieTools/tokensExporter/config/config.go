package config

import "github.com/TerraDharitri/drt-go-chain-tools/trieTools/trieToolsCommon"

// ContextFlagsTokensExporter is the flags config for tokens exporter
type ContextFlagsTokensExporter struct {
	trieToolsCommon.ContextFlagsConfig
	Outfile string
}
